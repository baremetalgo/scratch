package RayWidgets

import (
	"fmt"

	"github.com/baremetalgo/scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayButton struct {
	Label     string
	Layout    *RayGui.Layout
	Bounds    rl.Rectangle
	Visible   bool
	TextColor rl.Color
	IsPressed bool
	OnClick   func() // <-- callback for clicks
}

func NewRayButton(label string) *RayButton {
	return &RayButton{
		Label:     label,
		Visible:   true,
		TextColor: RayGui.Default_Text_Color,
		Bounds:    rl.NewRectangle(0, 0, 0, 0),
		IsPressed: false,
	}
}

func (b *RayButton) Update() {
	if !b.Visible {
		return
	}

	mousePos := rl.GetMousePosition()
	inside := rl.CheckCollisionPointRec(mousePos, b.Bounds)

	if inside {
		// Button is visually pressed if mouse is down
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			b.IsPressed = true
		} else {
			b.IsPressed = false
		}

		// Fire callback only on release inside button
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			if b.OnClick != nil {
				b.OnClick()
			} else {
				b.TriggerFunc()
			}
		}
	} else {
		// Reset pressed state if mouse is outside
		b.IsPressed = false
	}
}

func (b *RayButton) Draw() {
	if !b.Visible {
		return
	}
	b.Update()

	// Pick color depending on pressed state
	bgColor := rl.DarkGray
	if b.IsPressed {
		bgColor = rl.LightGray // held down color
	}

	// Draw button background
	rl.DrawRectangleRec(b.Bounds, bgColor)
	rl.DrawRectangleLinesEx(b.Bounds, 1, rl.Black)

	// Draw text centered
	textSize := rl.MeasureTextEx(
		b.GetTextFont(),
		b.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)
	textX := b.Bounds.X + (b.Bounds.Width-textSize.X)/2
	textY := b.Bounds.Y + (b.Bounds.Height-textSize.Y)/2

	rl.DrawTextEx(
		b.GetTextFont(),
		b.Label,
		rl.NewVector2(textX, textY),
		float32(RayGui.Default_Body_Font_Size),
		0,
		b.TextColor,
	)
}

// Implement bounds setter
func (b *RayButton) SetBounds(bounds rl.Rectangle) {
	b.Bounds = bounds
}

func (b *RayButton) GetBounds() rl.Rectangle { return b.Bounds }
func (b *RayButton) GetVisibility() bool     { return b.Visible }
func (b *RayButton) GetBgColor() rl.Color    { return rl.Blank }
func (b *RayButton) GetTextFont() rl.Font {
	if b.Layout != nil {
		return b.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}
func (b *RayButton) GetTextColor() rl.Color { return b.TextColor }

func (b *RayButton) SetLayout(layout *RayGui.Layout) {
	b.Layout = layout
	textSize := rl.MeasureTextEx(
		b.GetTextFont(),
		b.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)
	b.Bounds.Width = textSize.X + 20
	b.Bounds.Height = textSize.Y + 12
}

// Example fallback function
func (b *RayButton) TriggerFunc() {
	print := fmt.Sprintf("%v Button clicked...", b.Label)
	fmt.Println(print)
}
