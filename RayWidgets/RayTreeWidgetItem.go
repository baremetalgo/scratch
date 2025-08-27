package RayWidgets

import (
	"fmt"
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TreeWidgetItem struct {
	RayGui.BaseWidget
	Parent     *TreeWidgetItem
	Children   []*TreeWidgetItem
	isExpanded bool
	toggleRect rl.Rectangle
}

func NewTreeWidgetItem(name string) *TreeWidgetItem {
	item := TreeWidgetItem{}
	item.Visible = true
	item.Name = name
	item.BorderColor = RayGui.Default_Border_Color

	item.Layout = RayGui.NewLayout()
	item.Layout.Name = fmt.Sprintf("%v_layout", name)
	item.Layout.Type = RayGui.LayoutHorizontal
	item.Layout.Widget = &item
	item.Parent = nil
	item.HeaderFont = RayGui.Default_Widget_Body_Text_Font
	item.Children = make([]*TreeWidgetItem, 0)
	item.isExpanded = true

	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, &item)
	return &item
}

func (item *TreeWidgetItem) ClearChildren() {
	new_children_list := make([]*TreeWidgetItem, 0)
	for _, child := range item.Children {
		child.Parent = nil
	}

	item.Children = new_children_list

}

func (item *TreeWidgetItem) SetParent(parent *TreeWidgetItem) {
	item.Parent = parent
}

func (item *TreeWidgetItem) AddChildItem(child_item *TreeWidgetItem) {
	for _, child := range item.Children {
		if child.Name == child_item.Name {
			return
		}
	}

	item.Children = append(item.Children, child_item)
	item.Layout.AddLayout(child_item.Layout)
	child_item.SetParent(item)
}

func (item *TreeWidgetItem) GetAllChildrenRecusively() []*TreeWidgetItem {
	children := make([]*TreeWidgetItem, 0)
	if len(item.Children) > 0 {

		for _, child := range item.Children {
			children = append(children, child)
			children = append(children, child.GetAllChildrenRecusively()...)

		}
	}

	return children
}

func (item *TreeWidgetItem) RemoveChildren(child_item *TreeWidgetItem) {
	new_children_list := make([]*TreeWidgetItem, 0)
	for _, child := range item.Children {
		if child.Name != child_item.Name {
			new_children_list = append(new_children_list, child)
		} else {
			child_item.Parent = nil
		}
	}

	item.Children = new_children_list
}

func (item *TreeWidgetItem) Update() {
	// compute toggle rect (needs to happen every frame)
	posx := item.Layout.Bounds.X + float32(item.Layout.Spacing)
	textSizeVec := rl.MeasureTextEx(item.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0)
	posy := item.Layout.Bounds.Y + textSizeVec.Y + float32(item.Layout.Spacing)

	var sign string
	if len(item.Children) == 0 {
		sign = "o"
	} else if item.isExpanded {
		sign = "-"
	} else {
		sign = "+"
	}

	toggleSize := rl.MeasureTextEx(item.HeaderFont, sign, float32(RayGui.Default_Body_Font_Size), 0)
	item.toggleRect = rl.NewRectangle(
		posx, posy-toggleSize.Y,
		toggleSize.X+4, toggleSize.Y+4,
	)

	// handle click
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()
		if rl.CheckCollisionPointRec(mousePos, item.toggleRect) {
			if len(item.Children) > 0 {
				item.isExpanded = !item.isExpanded
			}
		}
	}

	// recurse into children if expanded
	if item.isExpanded {
		for _, child := range item.Children {
			child.Update()
		}
	}
}

func (item *TreeWidgetItem) Draw() {
	posx := item.Layout.Bounds.X + float32(item.Layout.Spacing)
	textSizeVec := rl.MeasureTextEx(item.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0)
	posy := item.Layout.Bounds.Y + textSizeVec.Y + float32(item.Layout.Spacing)

	// choose symbol based on expand state
	var sign string
	if len(item.Children) == 0 {
		sign = "o"
	} else if item.isExpanded {
		sign = "-"
	} else {
		sign = "+"
	}

	// draw toggle
	rl.DrawTextEx(
		item.HeaderFont,
		sign,
		rl.NewVector2(posx, posy),
		float32(RayGui.Default_Body_Font_Size),
		0,
		RayGui.Default_Text_Color,
	)

	// draw node name next to toggle
	toggleSize := rl.MeasureTextEx(item.HeaderFont, sign, float32(RayGui.Default_Body_Font_Size), 0)
	rl.DrawTextEx(
		item.HeaderFont,
		item.Name,
		rl.NewVector2(posx+toggleSize.X+6, posy),
		float32(RayGui.Default_Body_Font_Size),
		0,
		RayGui.Default_Text_Color,
	)

	// draw children if expanded
	if item.isExpanded {
		for _, child := range item.Children {
			child.Draw()
		}
	}
}
