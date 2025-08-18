package main

import (
	"scratch/RayGui"
	"scratch/RayWidgets"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Proper Layout Test")

	// Initialize fonts
	RayGui.InitializeFonts()

	// Create main widget
	mainWidget := RayGui.NewBaseWidget("Main Window")
	mainWidget.BgColor = rl.NewColor(60, 60, 60, 255)
	mainWidget.Layout.Type = RayGui.LayoutVertical
	mainWidget.IsMainWindow = true

	// Set bounds accounting for title bar
	mainWidget.Bounds = rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
	mainWidget.Layout.Bounds = rl.NewRectangle(
		0,
		mainWidget.TitleBarHeight, // Start below title bar
		float32(rl.GetScreenWidth()),
		float32(rl.GetScreenHeight())-mainWidget.TitleBarHeight,
	)

	// Add label - will now appear below title bar
	label := RayWidgets.NewRayLabel("PROPERLY POSITIONED LABEL")
	label.TextColor = rl.Yellow
	label.FontSize = 24
	mainWidget.Layout.AddChild(label)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Handle window resize
		if rl.IsWindowResized() {
			mainWidget.Bounds = rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
			mainWidget.Layout.Bounds = rl.NewRectangle(
				0,
				mainWidget.TitleBarHeight,
				float32(rl.GetScreenWidth()),
				float32(rl.GetScreenHeight())-mainWidget.TitleBarHeight,
			)
		}

		mainWidget.Update()
		mainWidget.Draw()

		rl.EndDrawing()
	}

	mainWidget.Unload()
	rl.CloseWindow()
}
