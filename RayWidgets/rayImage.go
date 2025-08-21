package RayWidgets

import (
	"fmt"
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayImage struct {
	FilePath  string
	Texture   rl.Texture2D
	Height    int32
	Width     int32
	Layout    *RayGui.Layout
	Bounds    rl.Rectangle
	Visible   bool
	IsChecked bool
	OnToggle  func(bool)
}

// Constructor
func NewRayImage(filepath string, width, height int32) *RayImage {
	texture := rl.LoadTexture(filepath)

	img := &RayImage{
		FilePath:  filepath,
		Texture:   texture,
		Width:     width,
		Height:    height,
		Visible:   true,
		IsChecked: false,
		Bounds:    rl.NewRectangle(0, 0, float32(width), float32(height)),
	}
	img.OnToggle = img.TriggerFunc
	return img
}

// Draw the image
func (r *RayImage) Draw() {
	if !r.Visible {
		return
	}
	r.Update()
	// Scale to bounds
	dest := rl.NewRectangle(r.Bounds.X, r.Bounds.Y, r.Bounds.Width, r.Bounds.Height)
	src := rl.NewRectangle(0, 0, float32(r.Texture.Width), float32(r.Texture.Height))

	rl.DrawTexturePro(
		r.Texture,
		src,
		dest,
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)
}

func (r *RayImage) Update() {
	if !r.Visible {
		return
	}
	mousePos := rl.GetMousePosition()

	// Check if the mouse is pressed and within the bounds of the checkbox
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mousePos, r.Bounds) {
		r.IsChecked = !r.IsChecked
		if r.OnToggle != nil {
			r.OnToggle(r.IsChecked)
		}
	}

}

// Implement bounds setter
func (r *RayImage) SetBounds(bounds rl.Rectangle) {
	r.Bounds = bounds
}

// Implement MainWidget interface methods
func (r *RayImage) GetBounds() rl.Rectangle { return r.Bounds }
func (r *RayImage) GetVisibility() bool     { return r.Visible }
func (r *RayImage) GetBgColor() rl.Color    { return rl.Blank }
func (r *RayImage) GetTextFont() rl.Font {
	if r.Layout != nil {
		return r.Layout.Widget.GetTextFont()
	}
	return RayGui.Default_Widget_Body_Text_Font
}
func (r *RayImage) GetTextColor() rl.Color { return rl.White }

// Implement ChildWidget interface
func (r *RayImage) SetLayout(layout *RayGui.Layout) {
	r.Layout = layout
	// If layout resizes the image, update bounds accordingly
	r.Bounds.Width = float32(r.Width)
	r.Bounds.Height = float32(r.Height)
}

func (r *RayImage) Unload() {
	rl.UnloadTexture(r.Texture)
}

func (r *RayImage) Load(filePath string) {
	rl.UnloadTexture(r.Texture)
	r.FilePath = filePath
	tex := rl.LoadTexture(filePath)
	// defer rl.UnloadTexture(tex)
	r.Texture = tex
}

func (r *RayImage) TriggerFunc(value bool) {
	path := "C:/Users/think/OneDrive/Desktop/wip.png"
	if r.FilePath == path {
		r.Load("C:/Users/think/OneDrive/Desktop/go_engine_ico.png")
	} else {
		r.Load("C:/Users/think/OneDrive/Desktop/wip.png")

	}
	print := fmt.Sprintf("%v Image clicked...", r.FilePath)
	fmt.Println(print)
}
