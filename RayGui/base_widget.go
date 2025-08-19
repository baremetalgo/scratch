package RayGui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BaseWidget struct {
	Name                  string
	Visible               bool
	IsMainWindow          bool
	Bounds                rl.Rectangle
	MaxWidth              float32
	MaxHeight             float32
	MinWidth              float32
	MinHeight             float32
	Parent                *BaseWidget
	DrawBackground        bool
	BgColor               rl.Color
	TextColor             rl.Color
	BorderColor           rl.Color
	TitleBar              bool
	TitleBarHeight        float32
	TitleBarBounds        rl.Rectangle
	TitleBarColor         rl.Color
	Layout                *Layout
	Children              []*BaseWidget
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
		Name:         name,
		Visible:      true,
		IsMainWindow: false,
		Bounds:       rl.NewRectangle(0, 0, 800, 600),
		// MaxWidth:             ,
		// MaxHeight:            300,
		MinWidth:              50.0,
		MinHeight:             50,
		Parent:                nil,
		BgColor:               Default_Bg_Color,
		DrawBackground:        false,
		TextColor:             Default_Text_Color,
		BorderColor:           rl.Gray,
		TitleBar:              true,
		TitleBarHeight:        Default_Titlebar_Height,
		TitleBarColor:         Default_Titlebar_Color,
		resizehandlerDragging: false,
		resizeHandlerColor:    Default_ResizerHandler_Color,
		last_position:         rl.NewVector2(0, 0),
		Closed:                false,
	}
	b.Layout = NewLayout()
	b.Layout.Name = fmt.Sprintf("%v_Layout", b.Name)
	b.Layout.Widget = b
	b.Layout.Bounds = b.Bounds
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		b.Bounds.Width,
		b.TitleBarHeight,
	)
	b.last_position = rl.NewVector2(b.Bounds.X, b.Bounds.Y)
	b.HeaderFont = Default_Widget_Header_Font
	b.TextFont = Default_Widget_Body_Text_Font
	return b
}

func (b *BaseWidget) GetBounds() rl.Rectangle {
	return b.Bounds
}

func SetLayout(widget *BaseWidget, layout *Layout) {
	layout.Widget = widget
}

func (b *BaseWidget) GetLayout() *Layout {
	return b.Layout
}

func (b *BaseWidget) SetBounds(bounds rl.Rectangle) {
	b.Bounds = bounds

	// Update title bar bounds
	b.TitleBarBounds = rl.NewRectangle(
		b.Bounds.X,
		b.Bounds.Y,
		b.Bounds.Width,
		b.TitleBarHeight,
	)
	/*
		// Calculate movement delta

		deltaX := bounds.X - b.Bounds.X
		deltaY := bounds.Y - b.Bounds.Y
			// Move children by the same delta
			for _, child := range b.Children {
				child.Bounds.X += deltaX
				child.Bounds.Y += deltaY
			}
	*/
}

func (b *BaseWidget) GetTitleBarBound() rl.Rectangle {
	return b.TitleBarBounds
}

func (b *BaseWidget) GetVisibility() bool {
	return b.Visible
}

func (b *BaseWidget) GetName() string {
	return b.Name
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

	if b.IsMainWindow {
		b.SetBounds(
			rl.NewRectangle(
				0,
				0,
				float32(rl.GetScreenWidth()),
				float32(rl.GetScreenHeight()),
			))
	}

	// Draw background first
	if b.DrawBackground {
		rl.DrawRectangleRec(b.Bounds, b.BgColor)
	}

	// Title bar
	if b.TitleBar {
		rl.DrawRectangleRec(b.TitleBarBounds, b.TitleBarColor)
		rl.DrawRectangleLinesEx(b.TitleBarBounds, 1, b.BorderColor)

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
		rl.DrawText("_", int32(minBtn.X)+3, int32(minBtn.Y)-2, 20, rl.Black)

		// Maximize button
		rl.DrawRectangle(int32(maxBtn.X), int32(maxBtn.Y), maxSize, maxSize, rl.LightGray)
		rl.DrawText("□", int32(maxBtn.X)+3, int32(maxBtn.Y)-2, 20, rl.Black)

		// Close button
		rl.DrawRectangle(int32(closeBtn.X), int32(closeBtn.Y), closeSize, closeSize, rl.LightGray)
		rl.DrawText("x", int32(closeBtn.X)+2, int32(closeBtn.Y)-2, 20, rl.Black)
	}

	b.Layout.Draw()

	// Border last so it's on top
	rl.DrawRectangleLinesEx(b.Bounds, 2, b.BorderColor)
	b.last_position = rl.NewVector2(b.Bounds.X, b.Bounds.Y)

	// drawing resize handle in case of mainwindow
	if b.IsMainWindow {
		handleSize := float32(15)
		handle_xpos := b.Bounds.X + b.Bounds.Width - handleSize
		handle_ypos := b.Bounds.Y + b.Bounds.Height - handleSize
		handleRect := rl.NewRectangle(handle_xpos, handle_ypos, handleSize, handleSize)

		// Draw a more visible resize handle
		rl.DrawRectangleRec(handleRect, b.resizeHandlerColor)
		rl.DrawRectangleLinesEx(handleRect, 1, b.BorderColor)

		// Draw diagonal lines for better visibility
		rl.DrawLineEx(
			rl.NewVector2(handleRect.X, handleRect.Y+handleRect.Height),
			rl.NewVector2(handleRect.X+handleRect.Width, handleRect.Y),
			2,
			b.BorderColor,
		)

		b.resizeHandler = handleRect
	}

}

func (b *BaseWidget) Update() {
	if !b.Visible || b.Closed {
		return
	}

	// Sync with window if main window
	if b.IsMainWindow {
		b.SyncWithWindow()
	}

	mouse := rl.GetMousePosition()
	mousePressed := rl.IsMouseButtonPressed(rl.MouseLeftButton)
	mouseReleased := rl.IsMouseButtonReleased(rl.MouseLeftButton)
	mouseDown := rl.IsMouseButtonDown(rl.MouseLeftButton)

	// Handle resize handler dragging
	if mousePressed && rl.CheckCollisionPointRec(mouse, b.resizeHandler) {
		b.resizeHandlerColor = rl.LightGray
		b.resizehandlerDragging = true
	}

	if mouseReleased {
		b.resizehandlerDragging = false
		b.resizeHandlerColor = Default_ResizerHandler_Color
	}

	if b.resizehandlerDragging && mouseDown {
		if b.IsMainWindow {
			// For main window, resize the actual window
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
	}

	// Update layout bounds to match widget bounds
	b.Layout.Bounds = b.Bounds
	b.Layout.Update()
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
