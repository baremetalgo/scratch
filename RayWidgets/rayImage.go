package RayWidgets

import (
	"fmt"
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayImage struct {
	RayGui.BaseWidget
	FilePath    string
	Texture     rl.Texture2D
	IsChecked   bool
	OnToggle    func(bool)
	AspectRatio float32
}

func NewRayImage(filepath string, width, height int32) *RayImage {
	texture := rl.LoadTexture(filepath)

	img := &RayImage{
		FilePath:    filepath,
		Texture:     texture,
		IsChecked:   false,
		AspectRatio: float32(texture.Width) / float32(texture.Height),
	}
	img.Name = filepath
	img.Visible = true
	img.TitleBar = false

	// Initialize layout properly
	img.SetLayout(RayGui.LayoutVertical)
	img.Layout.Bounds.Width = float32(width)
	img.Layout.Bounds.Height = float32(height)
	img.Layout.Widget = img

	img.OnToggle = img.TriggerFunc

	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, img)
	return img
}

func (r *RayImage) getScaledBounds() rl.Rectangle {
	containerWidth := r.Layout.Bounds.Width
	containerHeight := r.Layout.Bounds.Height
	containerAspect := containerWidth / containerHeight

	var scaledWidth, scaledHeight float32
	var x, y float32

	if containerAspect > r.AspectRatio {
		// Container is wider than image - letterbox on sides
		scaledHeight = containerHeight
		scaledWidth = scaledHeight * r.AspectRatio
		x = r.Layout.Bounds.X + (containerWidth-scaledWidth)/2
		y = r.Layout.Bounds.Y
	} else {
		// Container is taller than image - letterbox on top/bottom
		scaledWidth = containerWidth
		scaledHeight = scaledWidth / r.AspectRatio
		x = r.Layout.Bounds.X
		y = r.Layout.Bounds.Y + (containerHeight-scaledHeight)/2
	}

	return rl.NewRectangle(x, y, scaledWidth, scaledHeight)
}

func (r *RayImage) Draw() {
	if !r.Visible {
		return
	}

	scaledBounds := r.getScaledBounds()

	// Draw background in letterbox areas
	rl.DrawRectangleRec(r.Layout.Bounds, rl.Black) // Or your preferred background color

	// Draw the texture with aspect ratio preservation
	rl.DrawTexturePro(
		r.Texture,
		rl.NewRectangle(0, 0, float32(r.Texture.Width), float32(r.Texture.Height)),
		scaledBounds,
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	// Debug: Show aspect ratio info
	rl.DrawText(fmt.Sprintf("AR: %.2f", r.AspectRatio),
		int32(r.Layout.Bounds.X+5),
		int32(r.Layout.Bounds.Y+5),
		12, rl.White)
}

func (r *RayImage) Update() {
	if !r.Visible {
		return
	}

	// Update layout first
	r.Layout.Update()

	// Use scaled bounds for click detection
	scaledBounds := r.getScaledBounds()
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mousePos, scaledBounds) {
		r.IsChecked = !r.IsChecked
		if r.OnToggle != nil {
			r.OnToggle(r.IsChecked)
		}
	}
}

func (r *RayImage) Load(filePath string) {
	rl.UnloadTexture(r.Texture)
	r.FilePath = filePath
	tex := rl.LoadTexture(filePath)
	defer rl.UnloadTexture(tex)
	r.Texture = tex
}

func (r *RayImage) TriggerFunc(value bool) {
	path := "/sources/splash_screen.png"
	if r.FilePath == path {
		r.Load("C:/Users/think/OneDrive/Desktop/go_engine_ico.png")
	} else {
		r.Load("C:/Users/think/OneDrive/Desktop/wip.png")

	}
	print := fmt.Sprintf("%v Image clicked...", r.FilePath)
	fmt.Println(print)
}
