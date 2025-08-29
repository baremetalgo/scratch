package RayWidgets

import (
	"fmt"

	"github.com/baremetalgo/scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaySlider struct {
	Label     string
	Value     float32
	Min, Max  float32
	Dragging  bool
	Layout    *RayGui.Layout
	Bounds    rl.Rectangle
	Visible   bool
	TextColor rl.Color
	OnChange  func(float32)
}

func NewRaySlider(label string, value, min, max float32) *RaySlider {
	rs := &RaySlider{
		Label:     label,
		Value:     value,
		Min:       min,
		Max:       max,
		Dragging:  false,
		Visible:   true,
		TextColor: RayGui.Default_Text_Color,
		Bounds:    rl.NewRectangle(0, 0, 150, 20), // default width for slider
	}
	rs.OnChange = rs.TriggerFunc
	return rs
}

func (rs *RaySlider) Draw() {
	if !rs.Visible {
		return
	}

	rs.Update()

	// --- Calculate knob position ---
	percent := (rs.Value - rs.Min) / (rs.Max - rs.Min)
	knobX := rs.Bounds.X + percent*rs.Bounds.Width
	knobRect := rl.NewRectangle(knobX-5, rs.Bounds.Y, 10, rs.Bounds.Height)

	barY := rs.Bounds.Y + rs.Bounds.Height/2 - 3
	barHeight := float32(6)

	// --- Draw left (filled) area in green ---
	rl.DrawRectangle(int32(rs.Bounds.X), int32(barY),
		int32(knobX-rs.Bounds.X), int32(barHeight), rl.Green)

	// --- Draw right (unfilled) area in dark gray ---
	rl.DrawRectangle(int32(knobX), int32(barY),
		int32(rs.Bounds.Width-(knobX-rs.Bounds.X)), int32(barHeight), rl.Black)

	// --- Draw knob as rectangle ---
	rl.DrawRectangleRec(knobRect, rl.LightGray)
	rl.DrawRectangleLines(int32(knobRect.X), int32(knobRect.Y),
		int32(knobRect.Width), int32(knobRect.Height), rl.Black)

	// --- Draw label and value on right side ---
	labelText := fmt.Sprintf("%s: %.2f", rs.Label, rs.Value)
	textSize := rl.MeasureTextEx(rs.GetTextFont(), labelText, float32(RayGui.Default_Body_Font_Size), 0)

	textX := rs.Bounds.X + rs.Bounds.Width + 10
	textY := rs.Bounds.Y + (rs.Bounds.Height-textSize.Y)/2

	rl.DrawTextEx(rs.GetTextFont(), labelText,
		rl.NewVector2(textX, textY),
		float32(RayGui.Default_Body_Font_Size), 0, rs.TextColor)
}

func (rs *RaySlider) Update() {
	mouse := rl.GetMousePosition()
	knobWidth := float32(10)
	percent := (rs.Value - rs.Min) / (rs.Max - rs.Min)
	knobX := rs.Bounds.X + percent*rs.Bounds.Width
	knobRect := rl.NewRectangle(knobX-5, rs.Bounds.Y, knobWidth, rs.Bounds.Height)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mouse, knobRect) {
		rs.Dragging = true
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		rs.Dragging = false
	}
	if rs.Dragging {
		newPercent := (mouse.X - rs.Bounds.X) / rs.Bounds.Width
		if newPercent < 0 {
			newPercent = 0
		}
		if newPercent > 1 {
			newPercent = 1
		}
		newValue := rs.Min + newPercent*(rs.Max-rs.Min)
		if newValue != rs.Value {
			rs.Value = newValue
			if rs.OnChange != nil {
				rs.OnChange(rs.Value)
			}
		}
	}
}

// Implement bounds setter
func (rs *RaySlider) SetBounds(bounds rl.Rectangle) {
	rs.Bounds = bounds
}

// Implement MainWidget interface methods
func (rs *RaySlider) GetBounds() rl.Rectangle { return rs.Bounds }
func (rs *RaySlider) GetVisibility() bool     { return rs.Visible }
func (rs *RaySlider) GetBgColor() rl.Color    { return rl.Blank }
func (rs *RaySlider) GetTextFont() rl.Font {
	if rs.Layout != nil {
		return rs.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}

func (rs *RaySlider) GetTextColor() rl.Color {
	return rs.TextColor
}

func (rs *RaySlider) SetLayout(layout *RayGui.Layout) {
	rs.Layout = layout
}

func (rs *RaySlider) TriggerFunc(value float32) {
	fmt.Printf("%v Slider value changed: %.2f\n", rs.Label, value)
}
