package RayGui

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LayoutHorizontal = 0
	LayoutVertical   = 1
	LayoutGrid       = 2
)

type MainWidget interface {
	Update()
	Draw()
	SetBounds(rl.Rectangle)
	GetBounds() rl.Rectangle
	GetVisibility() bool
	GetBgColor() rl.Color
	GetTextFont() rl.Font
	GetTextColor() rl.Color
	GetName() string
	GetTitleBarBound() rl.Rectangle
	GetLayout() *Layout
}

type BoundsSetter interface {
	SetBounds(bounds rl.Rectangle)
}

type Layout struct {
	Name     string
	Type     int
	Widget   MainWidget
	Children []MainWidget
	Layouts  []*Layout
	Parent   *Layout
	Spacing  int
	Padding  rl.Vector2
	Visible  bool
	Bounds   rl.Rectangle
}

func NewLayout() *Layout {

	return &Layout{
		Type:    LayoutVertical,
		Spacing: 5,
		Padding: rl.NewVector2(5, 5),
		Layouts: make([]*Layout, 0),
		Visible: true,
		Bounds:  rl.NewRectangle(0, 0, 0, 0),
	}
}

func (l *Layout) AddChild(child MainWidget) {
	l.Children = append(l.Children, child)
	l.Layouts = append(l.Layouts, child.GetLayout())

}

func (l *Layout) AddLayout(layout *Layout) {
	l.Layouts = append(l.Layouts, layout)
	layout.Parent = l
	layout.Update()

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

func (l *Layout) Update() {

	if l.Widget != nil {
		widget_bounds := l.Widget.GetBounds()
		widget_titlebar_bounds := l.Widget.GetTitleBarBound()
		l.Bounds.X = widget_bounds.X + float32(l.Spacing)
		l.Bounds.Y = widget_bounds.Y + widget_titlebar_bounds.Height
		l.Bounds.Width = widget_bounds.Width - float32(l.Spacing)
		l.Bounds.Height = widget_bounds.Height - widget_titlebar_bounds.Height - float32(l.Spacing)
	} else {
		parent_bounds := l.Parent.Bounds
		l.Bounds.X = parent_bounds.X + float32(l.Spacing)
		l.Bounds.Y = parent_bounds.Y + float32(l.Spacing)
		l.Bounds.Width = parent_bounds.Width - float32(l.Spacing)
		l.Bounds.Height = parent_bounds.Height - float32(l.Spacing)

	}
	rl.DrawRectangleLinesEx(l.Bounds, 1, rl.Pink)
	l.UpdateChildLayouts()

}

func (l *Layout) UpdateChildLayouts() {
	no_of_children := len(l.Layouts)
	for i, child_layout := range l.Layouts {
		switch l.Type {
		case LayoutHorizontal:
			// Calculate available width (subtract spacing between items)
			availableWidth := l.Bounds.Width - float32(l.Spacing*(no_of_children-1))
			width := availableWidth / float32(no_of_children)
			height := l.Bounds.Height - float32(l.Spacing)

			xpos := l.Bounds.X + (width+float32(l.Spacing))*float32(i)
			ypos := l.Bounds.Y + float32(l.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutVertical:
			width := l.Bounds.Width - float32(l.Spacing)*2
			height := (l.Bounds.Height)/float32(no_of_children) - float32(l.Spacing)

			xpos := l.Bounds.X + float32(l.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + (height+float32(l.Spacing))*float32(i)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutGrid:
			cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
			if cols < 1 {
				cols = 1
			}
			rows := (no_of_children + cols - 1) / cols

			availableW := l.Bounds.Width - float32(l.Spacing*(cols+1))
			availableH := l.Bounds.Height - float32(l.Spacing*(rows+1)) - float32(l.Spacing)

			cellW := availableW / float32(cols)
			cellH := availableH / float32(rows)

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing))
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing))

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, cellW, cellH))
		}
		rl.DrawRectangleLinesEx(child_layout.GetBounds(), 1, rl.Pink)
		child_layout.UpdateChildLayouts()
		child_layout.UpdateChildWidgets()
	}
}

func (l *Layout) Draw() {
	if !l.Visible {
		return
	}
	for _, child_widget := range l.Children {
		child_widget.Draw()
	}
	for _, child_layout := range l.Layouts {
		child_layout.Draw()
	}
}

func (l *Layout) UpdateChildWidgets() {
	no_of_children := len(l.Children)
	if no_of_children == 0 {
		return
	}

	for i, child_widget := range l.Children {
		var titleBarHeight float32 = 0
		if baseWidget, ok := child_widget.(*BaseWidget); ok && baseWidget.TitleBar {
			titleBarHeight = baseWidget.TitleBarHeight
		}

		switch l.Type {
		case LayoutHorizontal:
			// Calculate available width (subtract spacing between items)
			availableWidth := l.Bounds.Width - float32(l.Spacing*(no_of_children-1))
			width := availableWidth / float32(no_of_children)
			height := l.Bounds.Height - titleBarHeight

			xpos := l.Bounds.X + (width+float32(l.Spacing))*float32(i)
			ypos := l.Bounds.Y + titleBarHeight

			child_widget.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutVertical:
			width := l.Bounds.Width - float32(l.Spacing)*2
			height := (l.Bounds.Height)/float32(no_of_children) - float32(l.Spacing)

			xpos := l.Bounds.X + float32(l.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + (height+float32(l.Spacing))*float32(i)

			child_widget.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutGrid:
			cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
			if cols < 1 {
				cols = 1
			}
			rows := (no_of_children + cols - 1) / cols

			availableW := l.Bounds.Width - float32(l.Spacing*(cols+1))
			availableH := l.Bounds.Height - float32(l.Spacing*(rows+1)) - titleBarHeight

			cellW := availableW / float32(cols)
			cellH := availableH / float32(rows)

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing))
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing))

			child_widget.SetBounds(rl.NewRectangle(xpos, ypos, cellW, cellH))
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
