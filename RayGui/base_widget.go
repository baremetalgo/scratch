package RayGui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ALL_WIDGETS []MainWidget

type BaseWidget struct {
	Name                  string
	Visible               bool
	IsMainWindow          bool
	Parent                *BaseWidget
	DrawBackground        bool
	DrawWidgetBorder      bool
	BgColor               rl.Color
	TextColor             rl.Color
	BorderColor           rl.Color
	TitleBar              bool
	Layout                *Layout
	minimized             bool
	resizeHandler         rl.Rectangle
	resizeHandlerColor    rl.Color
	resizehandlerDragging bool
	last_position         rl.Vector2
	Closed                bool
	HeaderFont            rl.Font
	TextFont              rl.Font
	zIndex                int
	lastHeight            float32
	drawMinButton         bool
	drawMaxButton         bool
	drawCloseButton       bool
	DrawPostHook          func()
}

func NewBaseWidget(name string) *BaseWidget {
	b := &BaseWidget{
		Name:                  name,
		Visible:               true,
		IsMainWindow:          false,
		Parent:                nil,
		BgColor:               Default_Bg_Color,
		DrawBackground:        false,
		DrawWidgetBorder:      true,
		TextColor:             Default_Text_Color,
		BorderColor:           rl.Gray,
		TitleBar:              true,
		resizehandlerDragging: false,
		resizeHandlerColor:    Default_ResizerHandler_Color,
		last_position:         rl.NewVector2(0, 0),
		Closed:                false,
		drawMinButton:         false,
		drawMaxButton:         false,
		drawCloseButton:       false,
	}
	b.SetZIndex(1)

	// Initialize layout FIRST
	b.SetLayout(0)

	b.HeaderFont = Default_Widget_Header_Font
	b.TextFont = Default_Widget_Body_Text_Font

	ALL_WIDGETS = append(ALL_WIDGETS, b)
	return b
}

func (b *BaseWidget) SetLayout(layout_type int) {
	b.Layout = NewLayout()
	b.Layout.Type = layout_type
	b.Layout.Name = fmt.Sprintf("%v_Layout", b.Name)
	b.Layout.Widget = b
}

func (b *BaseWidget) GetZIndex() int {
	return b.zIndex
}

func (b *BaseWidget) SetZIndex(zIndex int) {
	b.zIndex = zIndex
}

func (b *BaseWidget) MainWindow() bool {
	return b.IsMainWindow
}

func (b *BaseWidget) GetTitleBar() bool {
	return b.TitleBar
}

func SetLayout(widget *BaseWidget, layout *Layout) {
	layout.Widget = widget
}

func (b *BaseWidget) GetLayout() *Layout {
	return b.Layout
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
	size := int32(Default_Titlebar_Height - 10)
	if size < 8 {
		size = 8
	}
	y := b.Layout.Bounds.ToInt32().Y + (int32(Default_Titlebar_Height) / 4)

	closeX := int32(b.Layout.Bounds.X+b.Layout.Bounds.Width) - int32(Default_Titlebar_Height) - 2
	maxX := int32(b.Layout.Bounds.X+b.Layout.Bounds.Width) - int32(Default_Titlebar_Height)*2 + 2
	minX := int32(b.Layout.Bounds.X+b.Layout.Bounds.Width) - int32(Default_Titlebar_Height)*3 + 6

	minBtn = rl.NewRectangle(float32(minX), float32(y), float32(size), float32(size))
	maxBtn = rl.NewRectangle(float32(maxX), float32(y), float32(size), float32(size))
	closeBtn = rl.NewRectangle(float32(closeX), float32(y), float32(size), float32(size))
	return minBtn, maxBtn, closeBtn, size, size, size
}

func (b *BaseWidget) Draw() {
	if !b.Visible || b.Closed {
		return
	}

	if b.DrawBackground {
		rl.DrawRectangleRec(b.Layout.Bounds, b.BgColor)
	}

	// Title bar - draw for all widgets that have TitleBar true, except main window

	if b.TitleBar && !b.IsMainWindow {

		TitleBarBounds := rl.NewRectangle(
			b.Layout.Bounds.X,
			b.Layout.Bounds.Y,
			b.Layout.Bounds.Width,
			Default_Titlebar_Height,
		)

		rl.DrawRectangleRec(TitleBarBounds, Default_Titlebar_Color)
		if b.DrawWidgetBorder {
			rl.DrawRectangleLinesEx(TitleBarBounds, 1, b.BorderColor)
		}
		// Title text
		rl.DrawTextEx(
			b.HeaderFont,
			b.Name,
			rl.NewVector2(b.Layout.Bounds.X+7, b.Layout.Bounds.Y+7),
			14,
			0,
			rl.White,
		)

		// Buttons
		minBtn, maxBtn, closeBtn, minSize, maxSize, closeSize := b.buttonRects()

		// Minimize button
		if b.drawMinButton {
			rl.DrawRectangle(int32(minBtn.X), int32(minBtn.Y), minSize, minSize, rl.LightGray)
			rl.DrawText("_", int32(minBtn.X)+3, int32(minBtn.Y)-2, 20, rl.Black)
		}

		// Maximize button
		if b.drawMaxButton {
			rl.DrawRectangle(int32(maxBtn.X), int32(maxBtn.Y), maxSize, maxSize, rl.LightGray)
			rl.DrawText("□", int32(maxBtn.X)+3, int32(maxBtn.Y)-2, 20, rl.Black)
		}

		// Close button
		if b.drawCloseButton {
			rl.DrawRectangle(int32(closeBtn.X), int32(closeBtn.Y), closeSize, closeSize, rl.LightGray)
			rl.DrawText("x", int32(closeBtn.X)+2, int32(closeBtn.Y)-2, 20, rl.Black)
		}
	}

	// Border last so it's on top
	if b.DrawWidgetBorder {
		rl.DrawRectangleLinesEx(b.Layout.Bounds, 1, b.BorderColor)
	}

	b.last_position = rl.NewVector2(b.Layout.Bounds.X, b.Layout.Bounds.Y)

	// drawing resize handle in case of mainwindow
	if b.IsMainWindow {
		handleSize := float32(15)
		handle_xpos := b.Layout.Bounds.X + b.Layout.Bounds.Width - handleSize
		handle_ypos := b.Layout.Bounds.Y + b.Layout.Bounds.Height - handleSize
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

	b.Layout.Draw()

	if b.DrawPostHook != nil {
		b.DrawPostHook()
	}

}

func (b *BaseWidget) Update() {
	if !b.Visible || b.Closed {
		return
	}
	b.Layout.Update()

	// For main window, set layout bounds to match window size
	if b.IsMainWindow {
		windowWidth := float32(rl.GetScreenWidth())
		windowHeight := float32(rl.GetScreenHeight())
		b.Layout.Bounds.Width = windowWidth - float32(b.Layout.Spacing)
		b.Layout.Bounds.Height = windowHeight - float32(b.Layout.Spacing)
	}

}

func (b *BaseWidget) ToggleMinimize() {
	if b.minimized {
		// Restore previous height
		b.Layout.Bounds.Height = b.lastHeight
		b.minimized = false

	} else {
		// Save current height before minimizing
		b.lastHeight = b.Layout.Bounds.Height

		b.Layout.Bounds.Height = Default_Titlebar_Height
		b.minimized = true
	}
}
