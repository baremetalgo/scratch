package RayWidgets

import (
	"scratch/RayGui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuBar struct {
	RayGui.BaseWidget
	ContextMenus []*ContextMenu
	activeMenu   *ContextMenu // Track which menu is currently active
}

func NewMenubar(name string) *MenuBar {
	m := &MenuBar{}
	m.Name = name
	m.Visible = true
	m.TitleBar = false
	m.DrawBackground = false
	m.DrawWidgetBorder = false
	m.BgColor = RayGui.Default_Bg_Color
	m.BorderColor = rl.Gray
	m.activeMenu = nil // Initialize as nil

	m.SetLayout(RayGui.LayoutHorizontal)

	m.HeaderFont = RayGui.Default_Widget_Header_Font
	m.TextColor = rl.White
	m.SetZIndex(10000)
	RayGui.ALL_WIDGETS = append(RayGui.ALL_WIDGETS, m)

	return m
}

func (m *MenuBar) AddContextMenu(context_menu *ContextMenu) {
	for _, contextmenu := range m.ContextMenus {
		if contextmenu.Name == context_menu.Name {
			return
		}
	}
	m.ContextMenus = append(m.ContextMenus, context_menu)
}

func (m *MenuBar) RemoveContextMenu(context_menu_name string) {
	new_list := make([]*ContextMenu, 0)
	for _, menu := range m.ContextMenus {
		if menu.Name != context_menu_name {
			new_list = append(new_list, menu)
		} else {
			menu = nil
		}
	}
}

func (m *MenuBar) Draw() {
	if !m.Visible {
		return
	}
	m.Update()
	// Draw background
	rl.DrawRectangleLinesEx(m.Layout.Bounds, 1, m.BorderColor)
	// fmt.Println(m.Bounds)
	// Draw menu items
	xPos := m.Layout.Bounds.X + 10

	for _, item := range m.ContextMenus {
		// Only draw the context menu if it's the active one
		if item == m.activeMenu {
			item.Draw()
		}

		textSize := rl.MeasureTextEx(m.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0)
		rl.DrawTextEx(m.HeaderFont, item.Name, rl.NewVector2(xPos+10, m.Layout.Bounds.Y+15), float32(RayGui.Default_Header_Font_Size), 0, m.TextColor)
		xPos += textSize.X + 40
	}
}

func (m *MenuBar) Update() {
	m.Layout.Update()
	xPos := m.Layout.Bounds.X + 10

	for _, item := range m.ContextMenus {
		textSize := rl.MeasureTextEx(m.HeaderFont, item.Name, float32(RayGui.Default_Header_Font_Size), 0)
		item.Bounds.X = xPos + 10
		item.Bounds.Y = m.Layout.Bounds.Y + 30
		xPos += textSize.X + 40
	}

	// Handle menu bar clicks
	m.HandleClicks()
}

// Add this method to handle menu bar interactions
func (m *MenuBar) HandleClicks() {
	// First check if any context menu item was clicked
	for _, menu := range m.ContextMenus {
		if menu.HandleClick() {
			m.activeMenu = nil // Menu will handle its own hiding
			return
		}
	}

	// If click was outside all menus, hide all
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()
		clickedOnMenuBar := rl.CheckCollisionPointRec(mousePos, m.Layout.Bounds)

		if !clickedOnMenuBar {
			// Click outside menu bar - hide all menus
			for _, menu := range m.ContextMenus {
				menu.Hide()
			}
			m.activeMenu = nil
			return
		}

		// Check if a menu title was clicked
		xPos := m.Layout.Bounds.X + 10
		for _, menu := range m.ContextMenus {
			textSize := rl.MeasureTextEx(m.HeaderFont, menu.Name, float32(RayGui.Default_Header_Font_Size), 0)
			menuRect := rl.NewRectangle(
				xPos,
				m.Layout.Bounds.Y,
				textSize.X+20,
				m.Layout.Bounds.Height,
			)

			if rl.CheckCollisionPointRec(mousePos, menuRect) {
				// Toggle this menu (show if hidden, hide if shown)
				if m.activeMenu == menu {
					menu.Hide()
					m.activeMenu = nil
				} else {
					// Hide all other menus first
					for _, otherMenu := range m.ContextMenus {
						otherMenu.Hide()
					}
					menu.Show()
					m.activeMenu = menu
				}
				return
			}

			xPos += textSize.X + 40
		}
	}
}
