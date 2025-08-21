package RayWidgets

import (
	"fmt"
	"scratch/RayGui"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderImage struct {
	Label       string
	Layout      *RayGui.Layout
	Bounds      rl.Rectangle
	Visible     bool
	Texture     rl.Texture2D
	TextColor   rl.Color
	AspectRatio float32
	Padding     float32
}

func NewRenderImage(label string, width, height int32) *RenderImage {
	// Create empty texture with black pixels
	emptyPixels := make([]rl.Color, width*height)
	for i := range emptyPixels {

		emptyPixels[i] = rl.Black
	}

	texture := rl.LoadTextureFromImage(&rl.Image{
		Data:    unsafe.Pointer(&emptyPixels[0]),
		Width:   width,
		Height:  height,
		Mipmaps: 1,
		Format:  rl.PixelFormat(7),
	})

	return &RenderImage{
		Label:       label,
		Visible:     true,
		TextColor:   rl.White,
		Bounds:      rl.NewRectangle(0, 0, float32(width), float32(height)),
		Texture:     texture,
		AspectRatio: float32(width) / float32(height),
		Padding:     5.0,
	}
}

func (r *RenderImage) Update(pixels []rl.Color, width, height int32) {
	if r.Texture.Width != width || r.Texture.Height != height {
		// Unload old texture if dimensions changed
		rl.UnloadTexture(r.Texture)

		// Create new texture with correct dimensions
		r.Texture = rl.LoadTextureFromImage(&rl.Image{
			Data:    unsafe.Pointer(&pixels[0]),
			Width:   width,
			Height:  height,
			Mipmaps: 1,
			Format:  rl.PixelFormat(7),
		})
		r.AspectRatio = float32(width) / float32(height)
	} else {
		// Update existing texture
		rl.UpdateTexture(r.Texture, pixels)
	}
}

// [Rest of the methods remain the same as previous implementation]
func (r *RenderImage) SetBounds(bounds rl.Rectangle) {
	r.Bounds = bounds
}

func (r *RenderImage) GetBounds() rl.Rectangle { return r.Bounds }
func (r *RenderImage) GetVisibility() bool     { return r.Visible }
func (r *RenderImage) GetBgColor() rl.Color    { return rl.Blank }
func (r *RenderImage) GetTextFont() rl.Font {

	if r.Layout != nil {
		return r.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}
func (r *RenderImage) GetTextColor() rl.Color { return r.TextColor }

func (r *RenderImage) SetLayout(layout *RayGui.Layout) {
	r.Layout = layout
	textSize := rl.MeasureTextEx(
		r.GetTextFont(),
		r.Label,
		float32(RayGui.Default_Body_Font_Size),
		0,
	)
	r.Bounds.Width = textSize.X + 20
	r.Bounds.Height = textSize.Y + 12
}

// Example fallback function
func (r *RenderImage) TriggerFunc() {
	print := fmt.Sprintf("%v Button clicked...", r.Label)
	fmt.Println(print)
}
