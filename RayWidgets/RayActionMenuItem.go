package RayWidgets

import (
	"fmt"
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ActionMenuItem struct {
	RayGui.BaseWidget
	OnTrigger func()
}

func NewActionMenuItem(name string) *ActionMenuItem {
	item := &ActionMenuItem{}
	item.Name = name
	item.Visible = true
	item.Bounds = rl.NewRectangle(0, 0, 200, 500)
	item.TitleBar = false
	item.DrawBackground = false
	item.DrawWidgetBorder = false
	item.BgColor = RayGui.Default_Bg_Color
	item.BorderColor = rl.Gray

	item.SetLayout(RayGui.LayoutVertical)
	item.HeaderFont = RayGui.Default_Widget_Header_Font
	item.TextColor = rl.White
	item.OnTrigger = item.run_trigger
	item.SetZIndex(10000)
	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, item)
	return item
}

func (item *ActionMenuItem) run_trigger() {
	fmt.Printf("Action menu item %v triggered.\n", item.Name)

}
