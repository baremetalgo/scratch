package main

import (
	"github.com/baremetalgo/scratch/RayGui"
	"github.com/baremetalgo/scratch/RayWidgets"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func create_menu_bar(menubarLayout *RayGui.Layout) *RayWidgets.MenuBar {
	// MenuBar Widget
	menubar := RayWidgets.NewMenubar("Menubar")
	menubar.TitleBar = false
	menubar.Layout.SetFixedHeight(50) // Make sure layout exists before calling
	menubarLayout.AddChild(menubar)

	file_menu := RayWidgets.NewContextMenu("File")
	save_menu_action := RayWidgets.NewActionMenuItem("Save")
	save_as_menu_action := RayWidgets.NewActionMenuItem("Save As..")
	open_asset_action := RayWidgets.NewActionMenuItem("Open Asset")
	open_level_action := RayWidgets.NewActionMenuItem("Open Level")
	exit_action := RayWidgets.NewActionMenuItem("Exit")

	file_menu.AddAction(open_level_action)
	file_menu.AddAction(open_asset_action)
	file_menu.AddAction(save_menu_action)
	file_menu.AddAction(save_as_menu_action)
	file_menu.AddAction(exit_action)
	menubar.AddContextMenu(file_menu)

	Edit_menu := RayWidgets.NewContextMenu("Edit")
	asset_action_item := RayWidgets.NewActionMenuItem("Create Assembly")
	Edit_menu.AddAction(asset_action_item)
	menubar.AddContextMenu(Edit_menu)

	about_menu := RayWidgets.NewContextMenu("About")
	menubar.AddContextMenu(about_menu)

	return menubar
}

func create_scratch_window() *RayGui.BaseWidget {
	// Initialize fonts
	RayGui.InitializeFonts()

	// Create main widget (fills entire window)
	mainWidget := RayGui.NewBaseWidget("MainWindow")
	mainWidget.IsMainWindow = true
	mainWidget.TitleBar = true
	mainWidget.Layout.Type = RayGui.LayoutVertical
	mainWidget.Layout.Padding = rl.NewVector2(5, 5)
	mainWidget.Layout.Spacing = 5

	// Layouts - Initialize them properly
	menubarLayout := RayGui.NewLayout()
	menubarLayout.Name = "MenuBarLayout"
	menubarLayout.Type = RayGui.LayoutHorizontal
	mainWidget.Layout.AddLayout(menubarLayout)

	midPanelLayout := RayGui.NewLayout()
	midPanelLayout.Name = "MidPanelLayout"
	midPanelLayout.Type = RayGui.LayoutHorizontal
	midPanelLayout.SetFixedHeight(640) // This was causing the panic
	mainWidget.Layout.AddLayout(midPanelLayout)

	lowerPanelLayout := RayGui.NewLayout()
	lowerPanelLayout.Name = "LowerPanelLayout"
	lowerPanelLayout.Type = RayGui.LayoutVertical
	mainWidget.Layout.AddLayout(lowerPanelLayout)

	// Level Explorer
	levelExplorer := RayWidgets.NewTreeWidget("Level Explorer")
	levelExplorer.Layout.SetFixedWidth(350)
	midPanelLayout.AddChild(levelExplorer)
	light_item := RayWidgets.NewTreeWidgetItem("Lights")
	levelExplorer.AddItem(light_item)
	renderer_item := RayWidgets.NewTreeWidgetItem("Renderer")
	levelExplorer.AddItem(renderer_item)
	shadows := RayWidgets.NewTreeWidgetItem("Shadows")
	renderer_item.AddChildItem(shadows)

	render_image := RayWidgets.NewRayImage("E:/GitHub/scratch/sources/splash_screen.png", 1280, 720)
	midPanelLayout.AddChild(render_image)

	// PropertiesPanel
	propertiesPanel := RayGui.NewBaseWidget("Properties")
	propertiesPanel.Layout.Name = "PropertiesWidgetLayout"
	propertiesPanel.Layout.SetFixedWidth(350)
	midPanelLayout.AddChild(propertiesPanel)

	// Asset Browser
	assetBrowser := RayGui.NewBaseWidget("Asset Browser")
	lowerPanelLayout.AddChild(assetBrowser)

	// Menubar
	create_menu_bar(menubarLayout)
	return mainWidget
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowTopmost)
	rl.InitWindow(1024, 720, "Scratch GUI Framework")

	mainWidget := create_scratch_window()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(RayGui.Default_Bg_Color)

		mainWidget.Update()
		mainWidget.Draw()
		rl.DrawFPS(500, 50)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
