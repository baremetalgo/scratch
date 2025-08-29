package RayWidgets

import (
	"github.com/baremetalgo/scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ContextMenu struct {
	RayGui.BaseWidget
	Bounds      rl.Rectangle
	ActionItems []*ActionMenuItem
	isClicked   bool
	isVisible   bool // Add this flag to control visibility
}

func NewContextMenu(name string) *ContextMenu {
	cmenu := &ContextMenu{}
	cmenu.Name = name
	cmenu.Visible = true

	cmenu.TitleBar = false
	cmenu.DrawBackground = false
	cmenu.DrawWidgetBorder = false
	cmenu.BgColor = RayGui.Default_Titlebar_Color
	cmenu.BorderColor = RayGui.Default_Border_Color

	// Initialize the layout properly
	cmenu.SetLayout(RayGui.LayoutVertical)
	cmenu.Bounds = rl.NewRectangle(cmenu.Layout.Bounds.X, cmenu.Layout.Bounds.Y, 200, 500)
	cmenu.HeaderFont = RayGui.Default_Widget_Header_Font
	cmenu.TextColor = rl.White
	cmenu.isClicked = false
	cmenu.isVisible = false // Start with menu hidden
	cmenu.SetZIndex(10000)
	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, cmenu)
	return cmenu
}

// Add these methods to control visibility
func (cmenu *ContextMenu) Show() {
	cmenu.isVisible = true
}

func (cmenu *ContextMenu) Hide() {
	cmenu.isVisible = false
}

func (cmenu *ContextMenu) Toggle() {
	cmenu.isVisible = !cmenu.isVisible
}

func (cmenu *ContextMenu) IsVisible() bool {
	return cmenu.isVisible
}

func (cmenu *ContextMenu) AddAction(actionItem *ActionMenuItem) {
	for _, action := range cmenu.ActionItems {
		if action.Name == actionItem.Name {
			return
		}

	}
	cmenu.ActionItems = append(cmenu.ActionItems, actionItem)

}

func (cmenu *ContextMenu) RemoveAction(actionItemName string) {
	new_action_items := make([]*ActionMenuItem, 0)
	for _, item := range cmenu.ActionItems {
		if item.Name != actionItemName {
			new_action_items = append(new_action_items, item)
		} else {
			item = nil
		}
	}
	cmenu.ActionItems = new_action_items
}

func (cmenu *ContextMenu) Draw() {
	// Only draw if visible and has items
	if !cmenu.isVisible || !cmenu.Visible || len(cmenu.ActionItems) == 0 {
		return
	}

	// Rest of your existing Draw method remains the same...
	const (
		padding     = float32(12)
		itemHeight  = float32(28)
		minWidth    = float32(120)
		borderWidth = float32(1)
	)

	// Calculate maximum text width
	maxTextWidth := float32(0)
	for _, item := range cmenu.ActionItems {
		textWidth := rl.MeasureTextEx(item.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0).X
		if textWidth > maxTextWidth {
			maxTextWidth = textWidth
		}
	}

	// Calculate final dimensions
	menuWidth := maxTextWidth + padding*2
	if menuWidth < minWidth {
		menuWidth = minWidth
	}
	menuHeight := itemHeight * float32(len(cmenu.ActionItems))

	// Set bounds with border consideration
	cmenu.Bounds = rl.NewRectangle(
		cmenu.Bounds.X,
		cmenu.Bounds.Y,
		menuWidth+borderWidth*2,
		menuHeight+borderWidth*2,
	)

	// Draw background with border
	rl.DrawRectangleRec(cmenu.Bounds, cmenu.BgColor)
	rl.DrawRectangleLinesEx(cmenu.Bounds, borderWidth, cmenu.BorderColor)

	// Draw each menu item
	for i, item := range cmenu.ActionItems {
		itemRect := rl.NewRectangle(
			cmenu.Bounds.X+borderWidth,
			cmenu.Bounds.Y+borderWidth+float32(i)*itemHeight,
			menuWidth,
			itemHeight,
		)

		// Hover effect
		if cmenu.IsItemHovered(i) {
			rl.DrawRectangleRec(itemRect, cmenu.BorderColor)
		}

		// Draw text (properly aligned)
		textSizeVec := rl.MeasureTextEx(item.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0)
		textX := itemRect.X + padding
		textY := itemRect.Y + (itemHeight-textSizeVec.Y)/2

		rl.DrawTextEx(
			item.HeaderFont,
			item.Name,
			rl.NewVector2(textX, textY),
			float32(RayGui.Default_Header_Font_Size),
			0,
			cmenu.TextColor,
		)

		// Draw separator (except for last item)
		if i < len(cmenu.ActionItems)-1 {
			separatorY := itemRect.Y + itemHeight
			rl.DrawLine(
				int32(itemRect.X+padding/2),
				int32(separatorY),
				int32(itemRect.X+itemRect.Width-padding/2),
				int32(separatorY),
				cmenu.BorderColor,
			)
		}
	}
}

// Add a method to handle clicks on menu items
func (cmenu *ContextMenu) HandleClick() bool {
	if !cmenu.isVisible {
		return false
	}

	mousePos := rl.GetMousePosition()
	if !rl.CheckCollisionPointRec(mousePos, cmenu.Bounds) {
		return false
	}

	// Check if any menu item was clicked
	for i, item := range cmenu.ActionItems {
		if cmenu.IsItemHovered(i) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if item.OnTrigger != nil {
				item.OnTrigger()
			}
			cmenu.Hide() // Hide menu after selection
			return true
		}
	}

	return false
}

func (cmenu *ContextMenu) IsItemHovered(index int) bool {
	mousePos := rl.GetMousePosition()
	itemHeight := float32(25)
	itemRect := rl.NewRectangle(
		cmenu.Bounds.X,
		cmenu.Bounds.Y+float32(index)*itemHeight,
		cmenu.Bounds.Width,
		itemHeight,
	)
	return rl.CheckCollisionPointRec(mousePos, itemRect)
}
