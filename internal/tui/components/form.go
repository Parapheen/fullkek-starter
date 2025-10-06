package components

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultTextWidth   = 40
	defaultToggleWidth = 10
)

var (
	labelStyle  = lipgloss.NewStyle().Bold(true)
	helperStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	textStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	focusColor  = lipgloss.Color("205")
	blurColor   = lipgloss.Color("240")
)

func containerStyle(width int, focused bool) lipgloss.Style {
	if width < 4 {
		width = 4
	}
	style := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).Width(width + 2)
	borderColor := blurColor
	if focused {
		borderColor = focusColor
	}
	return style.BorderForeground(borderColor)
}

// TextField bundles a Bubbles text input with a shared presentation style.
type TextField struct {
	label   string
	hint    string
	width   int
	input   textinput.Model
	focused bool
}

// NewTextField constructs a text input configured with shared styling.
func NewTextField(label, placeholder, initial string) TextField {
	ti := textinput.New()
	ti.CharLimit = 0
	ti.Placeholder = placeholder
	ti.SetValue(initial)
	ti.Prompt = ""
	ti.TextStyle = textStyle
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(blurColor)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(focusColor)
	ti.Cursor.SetMode(cursor.CursorBlink)
	ti.Width = defaultTextWidth

	return TextField{
		label: label,
		hint:  "",
		width: defaultTextWidth,
		input: ti,
	}
}

// SetHint configures optional helper text rendered under the field.
func (f *TextField) SetHint(hint string) {
	f.hint = hint
}

// Focus marks the field as active.
func (f *TextField) Focus() tea.Cmd {
	f.focused = true
	return f.input.Focus()
}

// Blur marks the field as inactive.
func (f *TextField) Blur() {
	f.focused = false
	f.input.Blur()
}

// Update propagates messages to the underlying text input.
func (f *TextField) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	f.input, cmd = f.input.Update(msg)
	return cmd
}

// View renders the field with shared styling.
func (f TextField) View() string {
	container := containerStyle(f.width, f.focused)

	content := container.Render(f.input.View())
	if f.hint != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, helperStyle.Render(f.hint))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(f.label),
		content,
	)
}

// Value returns the current text input value.
func (f TextField) Value() string {
	return f.input.Value()
}

// SetValue mutates the underlying input value.
func (f *TextField) SetValue(v string) {
	f.input.SetValue(v)
}

// SetPlaceholder adjusts the placeholder text.
func (f *TextField) SetPlaceholder(v string) {
	f.input.Placeholder = v
}

// SetWidth updates the field width ensuring layouts remain consistent.
func (f *TextField) SetWidth(width int) {
	if width < 20 {
		width = 20
	}
	f.width = width
	f.input.Width = width
}

// ToggleField displays a boolean input using the same styling as text fields.
type ToggleField struct {
	label   string
	width   int
	input   textinput.Model
	value   bool
	focused bool
	hint    string
}

// NewToggleField constructs a toggle backed by a styled text input.
func NewToggleField(label string, initial bool) ToggleField {
	ti := textinput.New()
	ti.CharLimit = 0
	ti.Prompt = ""
	ti.SetCursorMode(textinput.CursorHide)
	ti.TextStyle = textStyle
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(blurColor)
	ti.Width = defaultToggleWidth

	tf := ToggleField{label: label, width: defaultToggleWidth, input: ti}
	tf.Set(initial)
	return tf
}

// SetHint adds helper text rendered below the toggle.
func (t *ToggleField) SetHint(hint string) {
	t.hint = hint
}

// Value yields the current boolean value.
func (t ToggleField) Value() bool {
	return t.value
}

// Set overrides the toggle's value.
func (t *ToggleField) Set(v bool) {
	t.value = v
	if v {
		t.input.SetValue("Yes")
	} else {
		t.input.SetValue("No")
	}
}

// Focus activates the toggle field.
func (t *ToggleField) Focus() tea.Cmd {
	t.focused = true
	return t.input.Focus()
}

// Blur deactivates the toggle field.
func (t *ToggleField) Blur() {
	t.focused = false
	t.input.Blur()
}

// Toggle flips the boolean value.
func (t *ToggleField) Toggle() {
	t.Set(!t.value)
}

// Update keeps the underlying model responsive (cursor blink, etc.).
func (t *ToggleField) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	t.Set(t.value)
	return cmd
}

// View renders the toggle consistent with other fields.
func (t ToggleField) View() string {
	container := containerStyle(t.width, t.focused)
	content := container.Render(t.input.View())
	if t.hint != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, helperStyle.Render(t.hint))
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render(t.label),
		content,
	)
}

// SetWidth adjusts the toggle's render width.
func (t *ToggleField) SetWidth(width int) {
	if width < 10 {
		width = 10
	}
	t.width = width
	t.input.Width = width
}
