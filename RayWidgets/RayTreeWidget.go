package RayWidgets

import (
	"fmt"

	"github.com/baremetalgo/scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TreeWidget struct {
	RayGui.BaseWidget
	TreeItems map[*TreeWidgetItem][]string
}

func NewTreeWidget(name string) *TreeWidget {
	tree := TreeWidget{}
	tree.Visible = true
	tree.Name = name
	tree.DrawWidgetBorder = true
	tree.BorderColor = RayGui.Default_Border_Color
	tree.DrawBackground = true
	tree.BgColor = rl.NewColor(70, 70, 70, 255)

	tree.Layout = RayGui.NewLayout()
	tree.Layout.Name = fmt.Sprintf("%v_layout", name)
	tree.Layout.Type = RayGui.LayoutVertical
	tree.Layout.Widget = &tree

	tree.TreeItems = make(map[*TreeWidgetItem][]string)

	tree.DrawWidgetBorder = true
	tree.TitleBar = true
	tree.HeaderFont = RayGui.Default_Widget_Header_Font
	tree.DrawPostHook = tree.DrawChildren
	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, &tree)

	return &tree
}

func (tree *TreeWidget) Clear() {
	tree.TreeItems = make(map[*TreeWidgetItem][]string)
}

// Recursively populate the tree with data of the form map[string]any
func (tree *TreeWidget) Populate(data map[*TreeWidgetItem]any) {
	for k, v := range data {
		item := NewTreeWidgetItem(fmt.Sprintf("%v", k))
		tree.AddItem(item)

		// if value is nested map, recurse
		if childMap, ok := v.(map[any]any); ok {
			tree.populateChildren(item, childMap)
		}
	}
}

func (tree *TreeWidget) populateChildren(parent *TreeWidgetItem, data map[any]any) {
	for k, v := range data {
		child := NewTreeWidgetItem(fmt.Sprintf("%v", k))
		parent.Children = append(parent.Children, child)
		tree.TreeItems[parent] = append(tree.TreeItems[parent], child.Name)

		if childMap, ok := v.(map[any]any); ok {
			tree.populateChildren(child, childMap)
		}
	}
}

// Add a single item to the tree (top-level)
func (tree *TreeWidget) AddItem(item *TreeWidgetItem) {
	if item == nil {
		return
	}
	// Add to TreeItems as a root item
	if _, exists := tree.TreeItems[item]; !exists {
		tree.TreeItems[item] = []string{}
	}
	// Add widget’s layout node (so it can be drawn)
	tree.Layout.AddLayout(item.Layout)
}

// Remove an item (recursively clears its children too)
func (tree *TreeWidget) RemoveItem(item *TreeWidgetItem) {
	if item == nil {
		return
	}

	// Remove from parent’s child list if applicable
	if item.Parent != nil {
		newChildren := []*TreeWidgetItem{}
		for _, c := range item.Parent.Children {
			if c != item {
				newChildren = append(newChildren, c)
			}
		}
		item.Parent.Children = newChildren
	}

	// Remove from tree map
	delete(tree.TreeItems, item)

	// Also remove its children recursively
	for _, child := range item.Children {
		tree.RemoveItem(child)
	}

	// Remove from layout
	tree.Layout.RemoveLayout(item.Layout)
}

func (tree *TreeWidget) DrawChildren() {
	for item := range tree.TreeItems {
		item.Draw()
	}

}

func (tree *TreeWidget) DrawConnections(item1 *TreeWidgetItem, item2 *TreeWidget) {
	rl.DrawLine(
		item1.Layout.Bounds.ToInt32().X,
		item1.Layout.Bounds.ToInt32().Y,
		item2.Layout.Bounds.ToInt32().X,
		item2.Layout.Bounds.ToInt32().Y,
		rl.White)

}
