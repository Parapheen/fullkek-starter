package stacks

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// Stack describes the composed project blueprint built from modular features.
type Stack struct {
	ID          string
	Name        string
	Description string
	Templates   []Template
	Directories []string
	Tags        []string
	Features    []Feature
}

// HasFeature reports whether the stack includes the feature with the provided ID.
func (s Stack) HasFeature(id string) bool {
	for _, feature := range s.Features {
		if feature.ID == id {
			return true
		}
	}
	return false
}

// Template describes a templated file sourced from the embedded filesystem.
type Template struct {
	// Source is the path inside internal/templates that should be rendered.
	Source string
	// Destination is the relative path to write within the generated project.
	Destination string
	// Mode controls the filesystem permissions for the generated file.
	Mode fsFileMode
}

// FeatureCategory represents a group of compatible modular features.
type FeatureCategory struct {
	ID            string
	Name          string
	Description   string
	Required      bool
	AllowMultiple bool
}

// Feature enumerates a single modular capability that can be composed together.
type Feature struct {
	ID          string
	CategoryID  string
	Name        string
	Description string
	Templates   []Template
	Directories []string
	Tags        []string
}

// Selection captures the chosen feature identifiers per category.
type Selection map[string][]string

// Compose builds a stack from the provided feature selection.
func Compose(sel Selection) (Stack, error) {
	if err := ValidateSelection(sel); err != nil {
		return Stack{}, err
	}

	resolved, err := resolveSelection(sel)
	if err != nil {
		return Stack{}, err
	}

	dirSet := map[string]struct{}{}
	tagSet := map[string]struct{}{}
	tmplSet := map[string]Template{}
	idParts := make([]string, 0, len(resolved))
	nameParts := make([]string, 0, len(resolved))
	descriptionParts := make([]string, 0, len(resolved))

	for _, res := range resolved {
		feature := res.Feature
		idParts = append(idParts, feature.ID)
		nameParts = append(nameParts, feature.Name)
		descriptionParts = append(descriptionParts, fmt.Sprintf("%s: %s", res.Category.Name, feature.Name))

		for _, dir := range feature.Directories {
			dirSet[dir] = struct{}{}
		}
		for _, tag := range feature.Tags {
			tagSet[tag] = struct{}{}
		}
		for _, tmpl := range feature.Templates {
			if existing, ok := tmplSet[tmpl.Destination]; ok {
				return Stack{}, fmt.Errorf("conflicting template destination %q between %s and %s", tmpl.Destination, existing.Source, tmpl.Source)
			}
			tmplSet[tmpl.Destination] = tmpl
		}
	}

	directories := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		directories = append(directories, dir)
	}
	sort.Strings(directories)

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	templates := make([]Template, 0, len(tmplSet))
	for _, tmpl := range tmplSet {
		templates = append(templates, tmpl)
	}
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Destination < templates[j].Destination
	})

	featureList := make([]Feature, 0, len(resolved))
	for _, res := range resolved {
		featureList = append(featureList, res.Feature)
	}

	id := strings.Join(idParts, "+")
	name := strings.Join(nameParts, " + ")
	if id == "" {
		id = "base"
	}
	if name == "" {
		name = "Base"
	}

	description := "Modular stack"
	if len(descriptionParts) > 0 {
		description = fmt.Sprintf("Modular stack composed of %s", strings.Join(descriptionParts, ", "))
	}

	return Stack{
		ID:          id,
		Name:        name,
		Description: description,
		Templates:   templates,
		Directories: directories,
		Tags:        tags,
		Features:    featureList,
	}, nil
}

// ValidateSelection checks that all selected feature IDs are valid and compatible.
func ValidateSelection(sel Selection) error {
	_, err := resolveSelection(sel)
	return err
}

// CloneSelection returns a deep copy of the provided selection to avoid mutation.
func CloneSelection(sel Selection) Selection {
	out := make(Selection, len(sel))
	for key, values := range sel {
		copyValues := make([]string, len(values))
		copy(copyValues, values)
		out[key] = copyValues
	}
	return out
}

type resolvedFeature struct {
	Category FeatureCategory
	Feature  Feature
}

func resolveSelection(sel Selection) ([]resolvedFeature, error) {
	categories := Categories()
	index := featureIndex()
	categoryIndex := categoryByID(categories)
	resolved := make([]resolvedFeature, 0)
	selectedByCategory := make(map[string][]Feature, len(categories))

	for categoryID := range sel {
		if _, ok := categoryIndex[categoryID]; !ok {
			return nil, fmt.Errorf("unknown category %q", categoryID)
		}
	}

	for _, category := range categories {
		ids := sel[category.ID]
		if len(ids) == 0 {
			if category.Required {
				return nil, fmt.Errorf("no selection provided for required category %q", category.Name)
			}
			continue
		}
		if !category.AllowMultiple && len(ids) > 1 {
			return nil, fmt.Errorf("multiple selections provided for single-choice category %q", category.Name)
		}

		availableFeatures := FeaturesForCategory(category.ID)
		availableIDs := featureIDs(availableFeatures)

		for _, id := range ids {
			feature, ok := index[id]
			if !ok {
				suggestion := suggestClosestID(id, availableIDs)
				if suggestion != "" {
					return nil, fmt.Errorf("unknown feature %q for %s; did you mean %q? valid values: %s", id, category.Name, suggestion, strings.Join(availableIDs, ", "))
				}
				return nil, fmt.Errorf("unknown feature %q for %s; valid values: %s", id, category.Name, strings.Join(availableIDs, ", "))
			}
			if feature.CategoryID != category.ID {
				actualCategory := feature.CategoryID
				if actual, ok := categoryIndex[feature.CategoryID]; ok {
					actualCategory = actual.Name
				}
				return nil, fmt.Errorf("feature %q does not belong to category %q (belongs to %q)", id, category.Name, actualCategory)
			}
			resolved = append(resolved, resolvedFeature{Category: category, Feature: feature})
			selectedByCategory[category.ID] = append(selectedByCategory[category.ID], feature)
		}
	}

	for _, selectedFeatures := range selectedByCategory {
		for _, selected := range selectedFeatures {
			requiredIDs := FeatureDependencies(selected.ID)
			for _, requiredID := range requiredIDs {
				requiredFeature, ok := index[requiredID]
				if !ok {
					return nil, fmt.Errorf("feature %q depends on unknown feature %q", selected.ID, requiredID)
				}

				requiredCategory, ok := categoryIndex[requiredFeature.CategoryID]
				if !ok {
					return nil, fmt.Errorf("feature %q depends on %q in unknown category %q", selected.ID, requiredID, requiredFeature.CategoryID)
				}

				categorySelection := selectedByCategory[requiredFeature.CategoryID]
				if len(categorySelection) == 0 {
					return nil, fmt.Errorf("feature %q requires %q in category %q", selected.ID, requiredID, requiredCategory.Name)
				}

				if !containsFeature(categorySelection, requiredID) {
					chosenIDs := make([]string, 0, len(categorySelection))
					for _, chosen := range categorySelection {
						chosenIDs = append(chosenIDs, chosen.ID)
					}
					return nil, fmt.Errorf("feature %q requires %q; selected in %s: %s", selected.ID, requiredID, requiredCategory.Name, strings.Join(chosenIDs, ", "))
				}
			}
		}
	}

	if containsFeature(selectedByCategory[CategoryAuth], "auth-oauth2") && len(selectedByCategory[CategoryOAuthProviders]) == 0 {
		return nil, fmt.Errorf("feature %q requires at least one selection in category %q", "auth-oauth2", categoryIndex[CategoryOAuthProviders].Name)
	}

	return resolved, nil
}

func categoryByID(categories []FeatureCategory) map[string]FeatureCategory {
	index := make(map[string]FeatureCategory, len(categories))
	for _, category := range categories {
		index[category.ID] = category
	}
	return index
}

func featureIDs(features []Feature) []string {
	ids := make([]string, 0, len(features))
	for _, feature := range features {
		ids = append(ids, feature.ID)
	}
	sort.Strings(ids)
	return ids
}

func containsFeature(features []Feature, id string) bool {
	for _, feature := range features {
		if feature.ID == id {
			return true
		}
	}
	return false
}

func suggestClosestID(value string, candidates []string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(candidates) == 0 {
		return ""
	}

	best := ""
	bestDistance := 1 << 30
	for _, candidate := range candidates {
		distance := levenshteinDistance(value, candidate)
		if distance < bestDistance {
			bestDistance = distance
			best = candidate
		}
	}

	threshold := len(value) / 3
	if threshold < 2 {
		threshold = 2
	}
	if bestDistance > threshold {
		return ""
	}

	return best
}

func levenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len(b)
	}
	if b == "" {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr := make([]int, len(b)+1)
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			deletion := prev[j] + 1
			insertion := curr[j-1] + 1
			substitution := prev[j-1] + cost

			curr[j] = deletion
			if insertion < curr[j] {
				curr[j] = insertion
			}
			if substitution < curr[j] {
				curr[j] = substitution
			}
		}
		prev = curr
	}

	return prev[len(b)]
}

// featureIndex returns the immutable feature lookup table.
func featureIndex() map[string]Feature {
	features := allFeatures()
	index := make(map[string]Feature, len(features))
	for _, feature := range features {
		index[feature.ID] = feature
	}
	return index
}

// SelectionFromIDs normalizes per-category input ensuring each value is tracked as a slice.
func SelectionFromIDs(values map[string]string) Selection {
	selection := make(Selection, len(values))
	for key, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		selection[key] = []string{value}
	}
	return selection
}

// MergeSelections merges b into a, overriding category selections present in b.
func MergeSelections(a, b Selection) Selection {
	merged := CloneSelection(a)
	for key, values := range b {
		copyValues := make([]string, len(values))
		copy(copyValues, values)
		merged[key] = copyValues
	}
	return merged
}

// fsFileMode mirrors fs.FileMode without importing io/fs to avoid circular deps in templates.
type fsFileMode = fs.FileMode
