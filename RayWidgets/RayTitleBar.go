package RayWidgets

import (
	"github.com/baremetalgo/scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type RayTitleBar struct {
	Layout      *RayGui.Layout
	EngineLogo  string
	LogoTexture rl.Texture2D
	LogoBounds  rl.Rectangle
}

func NewRayTitleBar(layout RayGui.Layout) *RayTitleBar {
	t := &RayTitleBar{
		EngineLogo: "E:/GitHub/GopherEngine/icons/go_logo_small.png",
	}
	t.LogoTexture = rl.LoadTexture(t.EngineLogo)
	t.LogoBounds = rl.NewRectangle(0, 0, 200.0, 200.0)
	return t
}

func (t RayTitleBar) Draw() {
	t.Update()

	rl.DrawTexturePro(
		t.LogoTexture,
		t.LogoBounds,
		t.LogoBounds,
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

}

func (t *RayTitleBar) Update() {

}
