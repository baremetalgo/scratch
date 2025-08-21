package RayWidgets

import (
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayLabel struct {
	Label     string
	Layout    *RayGui.Layout
	Bounds    rl.Rectangle
	Visible   bool
	TextColor rl.Color
	FontSize  float32
}

func NewRayLabel(label string) *RayLabel {
	return &RayLabel{
		Label:     label,
		Visible:   true,
		TextColor: rl.Yellow,                      // Changed to bright yellow
		FontSize:  24,                             // Larger font size
		Bounds:    rl.NewRectangle(0, 0, 300, 40), // Fixed initial size
	}
}

func (l *RayLabel) Draw() {
	if !l.Visible {
		return
	}

	// Draw background for debugging
	rl.DrawRectangleRec(l.Bounds, rl.NewColor(80, 80, 80, 255))

	textSize := rl.MeasureTextEx(
		l.GetTextFont(),
		l.Label,
		l.FontSize,
		1,
	)

	textX := l.Bounds.X + 10
	textY := l.Bounds.Y + (l.Bounds.Height-textSize.Y)/2

	rl.DrawTextEx(
		l.GetTextFont(),
		l.Label,
		rl.NewVector2(textX, textY),
		l.FontSize,
		1,
		l.TextColor,
	)
}

// Implement bounds setter
func (l *RayLabel) SetBounds(bounds rl.Rectangle) {
	l.Bounds = bounds
}

// Implement MainWidget interface methods
func (l *RayLabel) GetBounds() rl.Rectangle { return l.Bounds }
func (l *RayLabel) GetVisibility() bool     { return l.Visible }
func (l *RayLabel) GetBgColor() rl.Color    { return rl.Blank }
func (l *RayLabel) GetTextFont() rl.Font {
	if l.Layout != nil {
		return l.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}
func (l *RayLabel) GetTextColor() rl.Color { return l.TextColor }

// Implement ChildWidget interface
func (l *RayLabel) SetLayout(layout *RayGui.Layout) {
	l.Layout = layout
	// Auto-size the label based on text
	textSize := rl.MeasureTextEx(
		l.GetTextFont(),
		l.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)
	l.Bounds.Width = textSize.X + 10 // Add some padding
	l.Bounds.Height = textSize.Y + 10
}
