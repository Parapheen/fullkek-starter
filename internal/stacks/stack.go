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
	resolved := make([]resolvedFeature, 0)

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
		for _, id := range ids {
			feature, ok := index[id]
			if !ok {
				return nil, fmt.Errorf("unknown feature %q", id)
			}
			if feature.CategoryID != category.ID {
				return nil, fmt.Errorf("feature %q does not belong to category %q", id, category.Name)
			}
			resolved = append(resolved, resolvedFeature{Category: category, Feature: feature})
		}
	}

	return resolved, nil
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
