package RayWidgets

import (
	"fmt"
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayCheckBox struct {
	Label     string
	Layout    *RayGui.Layout
	Bounds    rl.Rectangle
	Visible   bool
	TextColor rl.Color
	IsChecked bool
	OnToggle  func(bool)
}

func NewRayCheckBox(label string) *RayCheckBox {
	cb := &RayCheckBox{
		Label:     label,
		Visible:   true,
		TextColor: RayGui.Default_Text_Color,
		Bounds:    rl.NewRectangle(0, 0, 0, 0),
		IsChecked: false,
	}
	cb.OnToggle = cb.TriggerFunc
	return cb
}

func (cb *RayCheckBox) Draw() {
	if !cb.Visible {
		return
	}
	cb.Update()
	textSize := rl.MeasureTextEx(
		cb.GetTextFont(),
		cb.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)

	// Center text vertically within bounds
	textY := cb.Bounds.Y + (cb.Bounds.Height-textSize.Y)/2

	if cb.IsChecked {
		rl.DrawRectangle(int32(cb.Bounds.X+5), int32(textY), 20, 20, rl.Green)
	} else {
		rl.DrawRectangle(int32(cb.Bounds.X+5), int32(textY), 20, 20, rl.DarkGray)
	}
	rl.DrawRectangleLines(int32(cb.Bounds.X+5), int32(textY), 20, 20, rl.Black)

	rl.DrawTextEx(
		cb.GetTextFont(),
		cb.Label,
		rl.NewVector2(cb.Bounds.X+30, textY+4),
		float32(RayGui.Default_Body_Font_Size),
		0,
		cb.TextColor,
	)
}

// Implement bounds setter
func (cb *RayCheckBox) SetBounds(bounds rl.Rectangle) {
	cb.Bounds = bounds
}

// Implement MainWidget interface methods
func (cb *RayCheckBox) GetBounds() rl.Rectangle { return cb.Bounds }
func (cb *RayCheckBox) GetVisibility() bool     { return cb.Visible }
func (cb *RayCheckBox) GetBgColor() rl.Color    { return rl.Blank }
func (cb *RayCheckBox) GetTextFont() rl.Font {
	if cb.Layout != nil {
		return cb.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}

func (cb *RayCheckBox) GetTextColor() rl.Color {
	return cb.TextColor
}

func (cb *RayCheckBox) SetLayout(layout *RayGui.Layout) {
	cb.Layout = layout
	// Auto-size the label based on text
	textSize := rl.MeasureTextEx(
		cb.GetTextFont(),
		cb.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)
	cb.Bounds.Width = textSize.X + 10 // Add some padding
	cb.Bounds.Height = textSize.Y + 10
}

func (cb *RayCheckBox) Update() {
	mousePos := rl.GetMousePosition()

	// Check if the mouse is pressed and within the bounds of the checkbox
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mousePos, cb.Bounds) {
		cb.IsChecked = !cb.IsChecked
		if cb.OnToggle != nil {
			cb.OnToggle(cb.IsChecked)
		}
	}

}

func (cb *RayCheckBox) TriggerFunc(value bool) {
	print := fmt.Sprintf("%v CheckBox clicked...", cb.Label)
	fmt.Println(print)
}
