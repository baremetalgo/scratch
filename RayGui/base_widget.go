package RayGui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BaseWidget struct {
	Name                  string
	Visible               bool
	IsMainWindow          bool
	Bounds                rl.Rectangle
	Parent                *BaseWidget
	BgColor               rl.Color
	TextColor             rl.Color
	BorderColor           rl.Color
	TitleBar              bool
	TitleBarHeight        float32
	TitleBarBounds        rl.Rectangle
	TitleBarColor         rl.Color
	Layout                Layout
	Children              []*BaseWidget
	dragging              bool
	dragOffset            rl.Vector2
	minimized             bool
	resizeHandler         rl.Rectangle
	resizeHandlerColor    rl.Color
	resizehandlerDragging bool
	last_position         rl.Vector2
	lastHeight            float32
	Closed                bool
	HeaderFont            rl.Font
	TextFont              rl.Font
	childPositions        map[*BaseWidget]rl.Vector2
	childMinimizedStates  map[*BaseWidget]bool
	CloseButton           rl.Rectangle
}

func NewBaseWidget(name string) *BaseWidget {
	b := &BaseWidget{
		Name:                  name,
		Visible:               true,
		IsMainWindow:          false,
		Bounds:                rl.NewRectangle(0, 0, 200, 200),
		Parent:                nil,
		BgColor:               Default_Bg_Color,
		TextColor:             Default_Text_Color,
		BorderColor:           rl.Gray,
		TitleBar:              true,
		TitleBarHeight:        Default_Titlebar_Height,
		TitleBarColor:         Default_Titlebar_Color,
		Children:              nil,
		last_position:         rl.NewVector2(0, 0),
		minimized:             false,
		resizehandlerDragging: false,
		resizeHandlerColor:    Default_ResizerHandler_Color,
		childPositions:        make(map[*BaseWidget]rl.Vector2),
		childMinimizedStates:  make(map[*BaseWidget]bool),
		Closed:                false,
	}
	b.Layout = *NewLayout(b)
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		b.Bounds.Width,
		b.TitleBarHeight,
	)
	b.last_position = rl.NewVector2(b.Bounds.X, b.Bounds.Y)
	InitializeFonts()
	b.HeaderFont = Default_Widget_Header_Font
	b.TextFont = Default_Widget_Body_Text_Font
	return b
}

func (b *BaseWidget) GetBounds() rl.Rectangle {
	return b.Bounds
}

func (b *BaseWidget) SetBounds(bounds rl.Rectangle) {
	// Calculate movement delta
	deltaX := bounds.X - b.Bounds.X
	deltaY := bounds.Y - b.Bounds.Y

	// Update bounds
	b.Bounds = bounds

	// Update title bar bounds
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		b.Bounds.Width,
		b.TitleBarHeight,
	)

	// Move children by the same delta
	for _, child := range b.Children {
		child.Bounds.X += deltaX
		child.Bounds.Y += deltaY
	}
}

func (b *BaseWidget) GetVisibility() bool {
	return b.Visible
}

func (b *BaseWidget) GetBgColor() rl.Color {
	return b.BgColor
}

func (b *BaseWidget) GetTextFont() rl.Font {
	return b.TextFont
}

func (b *BaseWidget) GetTextColor() rl.Color {
	return b.TextColor
}

func (b *BaseWidget) SetLayout(layout *Layout) {
	b.Layout = *layout
}

func (b *BaseWidget) AddChildWidget(child *BaseWidget) {
	child.Parent = b
	b.Children = append(b.Children, child)
	b.childMinimizedStates[child] = child.minimized

	// Position child widget relative to parent with padding
	child.Bounds.X = b.Bounds.X + 10
	child.Bounds.Y = b.Bounds.Y + b.TitleBarHeight + 10

	// Store initial relative position
	b.childPositions[child] = rl.NewVector2(
		child.Bounds.X-b.Bounds.X,
		child.Bounds.Y-b.Bounds.Y,
	)

	// Ensure child stays within parent bounds
	if child.Bounds.X < b.Bounds.X {
		child.Bounds.X = b.Bounds.X
	}
	if child.Bounds.Y < b.Bounds.Y+b.TitleBarHeight {
		child.Bounds.Y = b.Bounds.Y + b.TitleBarHeight
	}
	if child.Bounds.X+child.Bounds.Width > b.Bounds.X+b.Bounds.Width {
		child.Bounds.X = b.Bounds.X + b.Bounds.Width - child.Bounds.Width
	}
	if child.Bounds.Y+child.Bounds.Height > b.Bounds.Y+b.Bounds.Height {
		child.Bounds.Y = b.Bounds.Y + b.Bounds.Height - child.Bounds.Height
	}

	// Set initial size if zero
	if child.Bounds.Width <= 0 {
		child.Bounds.Width = b.Bounds.Width - 20
	}
	if child.Bounds.Height <= 0 {
		child.Bounds.Height = 30 // Default height
	}
}

func (b *BaseWidget) buttonRects() (minBtn, maxBtn, closeBtn rl.Rectangle, minSize, maxSize, closeSize int32) {
	size := int32(b.TitleBarHeight - 10)
	if size < 8 {
		size = 8
	}

	y := b.Bounds.ToInt32().Y + (int32(b.TitleBarHeight) / 4)

	closeX := int32(b.Bounds.X+b.Bounds.Width) - int32(b.TitleBarHeight) - 2
	maxX := int32(b.Bounds.X+b.Bounds.Width) - int32(b.TitleBarHeight)*2 + 2
	minX := int32(b.Bounds.X+b.Bounds.Width) - int32(b.TitleBarHeight)*3 + 6

	minBtn = rl.NewRectangle(float32(minX), float32(y), float32(size), float32(size))
	maxBtn = rl.NewRectangle(float32(maxX), float32(y), float32(size), float32(size))
	closeBtn = rl.NewRectangle(float32(closeX), float32(y), float32(size), float32(size))
	b.CloseButton = closeBtn
	return minBtn, maxBtn, closeBtn, size, size, size
}

func (b *BaseWidget) Draw() {
	if !b.Visible || b.Closed {
		return
	}

	// Update bounds to match window size
	b.Bounds.Width = float32(rl.GetScreenWidth())
	b.Bounds.Height = float32(rl.GetScreenHeight())

	// Body
	if !b.minimized {
		b.Layout.Draw()

		// Draw border
		rl.DrawRectangleLinesEx(b.Bounds, 1, b.BorderColor)

		// Draw child widgets that aren't minimized
		for _, child := range b.Children {
			if !child.minimized {
				child.Draw()
			}
		}
	}

	// Title bar
	if b.TitleBar {
		rl.DrawRectangleRec(b.TitleBarBounds, b.TitleBarColor)
		rl.DrawRectangleLinesEx(b.TitleBarBounds, 1, b.BorderColor)
	}

	// Title text
	rl.DrawTextEx(
		Default_Widget_Header_Font,
		b.Name,
		rl.NewVector2(b.Bounds.X+7, b.Bounds.Y+7),
		14,
		0,
		rl.White,
	)

	// Buttons
	minBtn, maxBtn, closeBtn, minSize, maxSize, closeSize := b.buttonRects()

	// Minimize button
	rl.DrawRectangle(int32(minBtn.X), int32(minBtn.Y), minSize, minSize, rl.LightGray)
	if b.minimized {
		rl.DrawText("-", int32(minBtn.X)+3, int32(minBtn.Y)-2, 20, rl.Black)
	} else {
		rl.DrawText("_", int32(minBtn.X)+3, int32(minBtn.Y)-2, 20, rl.Black)
	}

	// Maximize button
	rl.DrawRectangle(int32(maxBtn.X), int32(maxBtn.Y), maxSize, maxSize, rl.LightGray)
	if b.IsMainWindow && rl.IsWindowMaximized() {
		rl.DrawText("❐", int32(maxBtn.X)+3, int32(maxBtn.Y)-2, 20, rl.Black)
	} else {
		rl.DrawText("□", int32(maxBtn.X)+3, int32(maxBtn.Y)-2, 20, rl.Black)
	}

	// Close button
	rl.DrawRectangle(int32(closeBtn.X), int32(closeBtn.Y), closeSize, closeSize, rl.LightGray)
	rl.DrawText("x", int32(closeBtn.X)+2, int32(closeBtn.Y)-2, 20, rl.Black)

	// Resize handle
	if !b.minimized {
		handle_xpos := b.Bounds.X + b.Bounds.Width - 18
		handle_ypos := b.Bounds.Y + b.Bounds.Height - 18
		handleRect := rl.NewRectangle(handle_xpos, handle_ypos, 15, 15)

		points := []rl.Vector2{
			{handleRect.X, handleRect.Y + handleRect.Height},
			{handleRect.X + handleRect.Width, handleRect.Y + handleRect.Height},
			{handleRect.X + handleRect.Width, handleRect.Y},
		}

		rl.DrawTriangle(points[0], points[1], points[2], b.resizeHandlerColor)
		b.resizeHandler = handleRect
	}

	b.last_position = rl.NewVector2(b.Bounds.X, b.Bounds.Y)
}

func (b *BaseWidget) Update() {
	if !b.Visible || b.Closed {
		return
	}

	// Always sync main window bounds with actual window size
	if b.IsMainWindow {
		b.SyncWithWindow()
	}

	// Update child widgets first
	for _, child := range b.Children {
		child.Update()
	}
	mouse := rl.GetMousePosition()

	// Handle window operations differently for main window
	if b.IsMainWindow {
		// Handle close button
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
			rl.CheckCollisionPointRec(mouse, b.CloseButton) {
			b.Closed = true
			return
		}

		// Resize handler
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
			rl.CheckCollisionPointRec(mouse, b.resizeHandler) {
			b.resizeHandlerColor = rl.LightGray
			b.resizehandlerDragging = true
		}

		if b.resizehandlerDragging {
			newWidth := int(mouse.X)
			newHeight := int(mouse.Y)

			// Set minimum window size
			minSize := b.Layout.GetMinSize()
			if newWidth < int(minSize.X) {
				newWidth = int(minSize.X)
			}
			if newHeight < int(minSize.Y) {
				newHeight = int(minSize.Y)
			}

			// Set new window size
			rl.SetWindowSize(newWidth, newHeight)
			b.SyncWithWindow()
		}
	} else {
		// Regular widget resize logic
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
			rl.CheckCollisionPointRec(mouse, b.resizeHandler) {
			b.resizeHandlerColor = rl.LightGray
			b.resizehandlerDragging = true
		}

		if b.resizehandlerDragging {
			newWidth := mouse.X - b.Bounds.X
			newHeight := mouse.Y - b.Bounds.Y

			minSize := b.Layout.GetMinSize()
			if newWidth < minSize.X {
				newWidth = minSize.X
			}
			if newHeight < minSize.Y {
				newHeight = minSize.Y
			}

			b.Bounds.Width = newWidth
			b.Bounds.Height = newHeight
		}
	}

	// Buttons
	minBtn, maxBtn, closeBtn, _, _, _ := b.buttonRects()

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
		rl.CheckCollisionPointRec(mouse, minBtn) {
		if b.IsMainWindow {
			rl.MinimizeWindow()
		} else {
			b.ToggleMinimize()
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
		rl.CheckCollisionPointRec(mouse, maxBtn) {
		if b.IsMainWindow {
			if rl.IsWindowMaximized() {
				rl.RestoreWindow()
			} else {
				rl.MaximizeWindow()
			}
		}
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
		rl.CheckCollisionPointRec(mouse, closeBtn) {
		if b.IsMainWindow {
			b.Closed = true
		} else {
			b.Close()
		}
		return
	}

	// Dragging - move the window or widget
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
		(rl.CheckCollisionPointRec(mouse, b.TitleBarBounds)) {

		// Only allow dragging if not clicking on a button
		if !rl.CheckCollisionPointRec(mouse, minBtn) &&
			!rl.CheckCollisionPointRec(mouse, maxBtn) &&
			!rl.CheckCollisionPointRec(mouse, closeBtn) {
			b.dragging = true
			b.dragOffset = rl.NewVector2(mouse.X-b.Bounds.X, mouse.Y-b.Bounds.Y)
		}
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		b.dragging = false
		b.resizehandlerDragging = false
		b.resizeHandlerColor = b.TitleBarColor
	}

	if b.dragging {
		if b.IsMainWindow {
			// Get current window position
			windowPos := rl.GetWindowPosition()

			// Calculate new position based on mouse and drag offset
			newX := int(mouse.X - b.dragOffset.X)
			newY := int(mouse.Y - b.dragOffset.Y)

			// Only move if position changed
			if newX != int(windowPos.X) || newY != int(windowPos.Y) {
				rl.SetWindowPosition(newX, newY)
			}
		} else {
			// Move the widget
			b.Bounds.X = mouse.X - b.dragOffset.X
			b.Bounds.Y = mouse.Y - b.dragOffset.Y

			// Clamp position to screen bounds
			screenW := float32(rl.GetScreenWidth())
			screenH := float32(rl.GetScreenHeight())

			if b.Bounds.X < 0 {
				b.Bounds.X = 0
			}
			if b.Bounds.Y < 0 {
				b.Bounds.Y = 0
			}
			if b.Bounds.X+b.Bounds.Width > screenW {
				b.Bounds.X = screenW - b.Bounds.Width
			}
			if b.Bounds.Y+b.Bounds.Height > screenH {
				b.Bounds.Y = screenH - b.Bounds.Height
			}
		}
	}

	// Update title bar bounds
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		b.Bounds.Width,
		b.TitleBarHeight,
	)

	b.EnforceMinWidth()
}

func (b *BaseWidget) ToggleMinimize() {
	if b.minimized {
		// Restore previous height
		b.Bounds.Height = b.lastHeight
		b.minimized = false

		// Restore child widgets to their previous state
		for child, wasMinimized := range b.childMinimizedStates {
			if !wasMinimized {
				child.minimized = false
				child.Bounds.Height = child.lastHeight
			}
		}
	} else {
		// Save current height before minimizing
		b.lastHeight = b.Bounds.Height

		// Store child minimized states and minimize them
		for _, child := range b.Children {
			b.childMinimizedStates[child] = child.minimized
			if !child.minimized {
				child.lastHeight = child.Bounds.Height
				child.Bounds.Height = child.TitleBarHeight
				child.minimized = true
			}
		}

		b.Bounds.Height = b.TitleBarHeight
		b.minimized = true
	}
}

func (b *BaseWidget) EnforceMinWidth() {
	titleMin := b.getTitleBarMinWidth()
	layoutMin := b.Layout.GetMinSize().X

	minWidth := titleMin
	if layoutMin > minWidth {
		minWidth = layoutMin
	}

	if b.Bounds.Width < minWidth {
		b.Bounds.Width = minWidth
	}
}

func (b *BaseWidget) getTitleBarMinWidth() float32 {
	textSize := rl.MeasureTextEx(
		Default_Widget_Header_Font,
		b.Name,
		14,
		0,
	)
	buttonPad := float32(b.TitleBarHeight*2 + 12)
	return textSize.X + 20 + buttonPad
}

func (b *BaseWidget) Close() {
	if b.IsMainWindow {
		// For main window, we'll let the main loop handle the actual closing
		b.Closed = true
	} else {
		b.Closed = true
		b.Visible = false
	}
}

// Interface implementation check
var _ ChildWidget = (*BaseWidget)(nil)

func (b *BaseWidget) SyncWithWindow() {
	if b.IsMainWindow {
		b.Bounds = rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight()))
	}
}

func (b *BaseWidget) Unload() {
	// Unload any resources here if needed
	b.Closed = true
	b.Visible = false

	// Unload child widgets
	for _, child := range b.Children {
		child.Unload()
	}

	// Clear references
	b.Children = nil
	b.Layout.Children = nil
}
