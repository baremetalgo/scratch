package RayGui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LayoutHorizontal = 0
	LayoutVertical   = 1
	LayoutGrid       = 2
)

type MainWidget interface {
	Draw()
	GetBounds() rl.Rectangle
	GetVisibility() bool
	GetBgColor() rl.Color
	GetTextFont() rl.Font
	GetTextColor() rl.Color
}

type BoundsSetter interface {
	SetBounds(bounds rl.Rectangle)
}

type ChildWidget interface {
	MainWidget
	SetLayout(layout *Layout)
	BoundsSetter
}

type Layout struct {
	Type     int
	Widget   MainWidget
	Children []ChildWidget
	Layouts  []*Layout
	Spacing  int
	Padding  rl.Vector2
	Visible  bool
	Bounds   rl.Rectangle
}

func NewLayout(widget MainWidget) *Layout {
	if widget == nil {
		// Create a default invisible widget if none provided
		widget = &BaseWidget{
			Visible: false,
			Bounds:  rl.NewRectangle(0, 0, 0, 0),
		}
	}
	return &Layout{
		Type:    LayoutVertical,
		Widget:  widget,
		Spacing: 5,
		Padding: rl.NewVector2(5, 5),
		Layouts: make([]*Layout, 0),
		Visible: true,
		Bounds:  rl.NewRectangle(0, 0, 0, 0),
	}
}

func (l *Layout) AddChild(child ChildWidget) {
	if child == nil {
		return
	}
	child.SetLayout(l)
	l.Children = append(l.Children, child)
}

func (l *Layout) AddLayout(layout *Layout) {
	if layout == nil {
		return
	}

	// Set default padding if not specified
	if layout.Padding.X == 0 && layout.Padding.Y == 0 {
		layout.Padding = rl.NewVector2(10, 10) // More visible default padding
	}

	layout.SetLayout(l)
	l.Layouts = append(l.Layouts, layout)

	// Calculate initial bounds based on parent
	if l.Widget != nil {
		parentBounds := l.Widget.GetBounds()
		layout.Bounds = rl.NewRectangle(
			parentBounds.X+l.Padding.X,
			parentBounds.Y+l.Padding.Y,
			parentBounds.Width-l.Padding.X*2,
			parentBounds.Height-l.Padding.Y*2,
		)
	}

	// Force initial layout calculation
	l.recalculateLayout()
}

func (l *Layout) recalculateLayout() {
	minSize := l.GetMinSize()

	// Update widget bounds if we have a widget
	if l.Widget != nil {
		if baseWidget, ok := l.Widget.(*BaseWidget); ok {
			if baseWidget.Bounds.Width < minSize.X {
				baseWidget.Bounds.Width = minSize.X
			}
			if baseWidget.Bounds.Height < minSize.Y {
				baseWidget.Bounds.Height = minSize.Y
			}
		}
	}

	// Update our own bounds
	if l.Bounds.Width < minSize.X {
		l.Bounds.Width = minSize.X
	}
	if l.Bounds.Height < minSize.Y {
		l.Bounds.Height = minSize.Y
	}
}

func (l *Layout) RemoveLayout(layout *Layout) {
	for i, lay := range l.Layouts {
		if lay == layout {
			l.Layouts = append(l.Layouts[:i], l.Layouts[i+1:]...)
			break
		}
	}
}

func (l *Layout) GetBounds() rl.Rectangle {
	return l.Bounds
}

func (l *Layout) GetVisibility() bool {
	return l.Visible
}

func (l *Layout) GetBgColor() rl.Color {
	if l.Widget != nil {
		return l.Widget.GetBgColor()
	}
	return rl.Blank
}

func (l *Layout) GetTextFont() rl.Font {
	if l.Widget != nil {
		return l.Widget.GetTextFont()
	}
	return Default_Widget_Body_Text_Font
}

func (l *Layout) GetTextColor() rl.Color {
	if l.Widget != nil {
		return l.Widget.GetTextColor()
	}
	return Default_Text_Color
}

func (l *Layout) SetLayout(parentLayout *Layout) {
	// No need to store parent layout
}

func (l *Layout) SetBounds(bounds rl.Rectangle) {
	l.Bounds = bounds
}

func (l *Layout) Draw() {
	if !l.Visible || (l.Widget != nil && !l.Widget.GetVisibility()) {
		return
	}
	// Calculate available space
	availableWidth := l.Bounds.Width - l.Padding.X*2
	availableHeight := l.Bounds.Height - l.Padding.Y*2
	// Calculate available space accounting for title bar
	var titleBarHeight float32 = 0
	if baseWidget, ok := l.Widget.(*BaseWidget); ok && baseWidget.TitleBar {
		titleBarHeight = baseWidget.TitleBarHeight
	}

	contentStartY := l.Bounds.Y + titleBarHeight
	availableHeight = l.Bounds.Height - titleBarHeight - l.Padding.Y*2

	// Position tracking (start below title bar)
	var x, y float32 = l.Bounds.X + l.Padding.X, contentStartY + l.Padding.Y
	var maxY float32 = y + availableHeight

	// Draw child widgets
	for _, child := range l.Children {
		if !child.GetVisibility() {
			continue
		}

		childBounds := child.GetBounds()
		childBounds.X = x
		childBounds.Y = y

		// Adjust dimensions to available space
		if l.Type == LayoutVertical {
			childBounds.Width = availableWidth
			// Constrain height if it would exceed available space
			if childBounds.Height > (maxY - y) {
				childBounds.Height = maxY - y
			}
		} else {
			// For horizontal layout, constrain height to available height
			childBounds.Height = availableHeight
		}

		child.SetBounds(childBounds)
		child.Draw()

		if l.Type == LayoutVertical {
			y += childBounds.Height + float32(l.Spacing)
			if y > maxY {
				break // Stop drawing if we've exceeded available height
			}
		} else {
			x += childBounds.Width + float32(l.Spacing)
			if x > (l.Bounds.X + l.Bounds.Width - l.Padding.X) {
				break // Stop drawing if we've exceeded available width
			}
		}
	}

	// Draw nested layouts
	for _, layout := range l.Layouts {
		if !layout.GetVisibility() {
			continue
		}

		layoutBounds := layout.GetBounds()
		layoutBounds.X = x
		layoutBounds.Y = y

		// Set dimensions based on layout type
		if l.Type == LayoutVertical {
			layoutBounds.Width = availableWidth
			// Constrain height if it would exceed available space
			if layoutBounds.Height > (maxY - y) {
				layoutBounds.Height = maxY - y
			}
		} else {
			layoutBounds.Height = availableHeight
		}

		layout.SetBounds(layoutBounds)
		layout.Draw()

		if l.Type == LayoutVertical {
			y += layoutBounds.Height + float32(l.Spacing)
			if y > maxY {
				break
			}
		} else {
			x += layoutBounds.Width + float32(l.Spacing)
			if x > (l.Bounds.X + l.Bounds.Width - l.Padding.X) {
				break
			}
		}
	}
}
func (l *Layout) GetMinSize() rl.Vector2 {
	var minWidth, minHeight float32 = 0, 0

	var titleBarHeight float32 = 0
	if baseWidget, ok := l.Widget.(*BaseWidget); ok && baseWidget.TitleBar {
		titleBarHeight = baseWidget.TitleBarHeight
	}

	switch l.Type {
	case LayoutVertical:
		// Calculate from children
		for _, child := range l.Children {
			if !child.GetVisibility() {
				continue
			}
			childBounds := child.GetBounds()
			if childBounds.Width > minWidth {
				minWidth = childBounds.Width
			}
			minHeight += childBounds.Height + float32(l.Spacing)
		}

		// Calculate from nested layouts
		for _, layout := range l.Layouts {
			if !layout.GetVisibility() {
				continue
			}
			layoutSize := layout.GetMinSize()
			if layoutSize.X > minWidth {
				minWidth = layoutSize.X
			}
			minHeight += layoutSize.Y + float32(l.Spacing)
		}

		minHeight += titleBarHeight + l.Padding.Y*2
		minWidth += l.Padding.X * 2

	case LayoutHorizontal:
		// Calculate from children
		for _, child := range l.Children {
			if !child.GetVisibility() {
				continue
			}
			childBounds := child.GetBounds()
			minWidth += childBounds.Width + float32(l.Spacing)
			if childBounds.Height > minHeight {
				minHeight = childBounds.Height
			}
		}

		// Calculate from nested layouts
		for _, layout := range l.Layouts {
			if !layout.GetVisibility() {
				continue
			}
			layoutSize := layout.GetMinSize()
			minWidth += layoutSize.X + float32(l.Spacing)
			if layoutSize.Y > minHeight {
				minHeight = layoutSize.Y
			}
		}

		minHeight += titleBarHeight + l.Padding.Y*2
		minWidth += l.Padding.X * 2
	}

	return rl.NewVector2(minWidth, minHeight)
}
