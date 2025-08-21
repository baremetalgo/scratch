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
	GetTitleBar() bool
	GetLayout() *Layout
}

type BoundsSetter interface {
	SetBounds(bounds rl.Rectangle)
}

type Layout struct {
	Name        string
	Type        int
	Widget      MainWidget
	Children    []MainWidget
	Layouts     []*Layout
	Parent      *Layout
	Spacing     int
	Padding     rl.Vector2
	Visible     bool
	Bounds      rl.Rectangle
	FixedHeight float32
	FixedWidth  float32
	DebugDraw   bool
}

func NewLayout() *Layout {

	return &Layout{
		Type:      LayoutVertical,
		Spacing:   10,
		Padding:   rl.NewVector2(5, 5),
		Layouts:   make([]*Layout, 0),
		Visible:   true,
		Bounds:    rl.NewRectangle(0, 0, 0, 0),
		DebugDraw: true,
	}
}

func (l *Layout) AddChild(child MainWidget) {
	l.Children = append(l.Children, child)
	l.AddLayout(child.GetLayout())
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

func (l *Layout) SetBounds(bounds rl.Rectangle) {
	l.Bounds = bounds
}

func (l *Layout) Update() {

	if l.Widget != nil {
		widget_bounds := l.Widget.GetBounds()
		widget_titlebar_bounds := l.Widget.GetTitleBarBound()

		l.Bounds.X = widget_bounds.X + l.Padding.X + float32(l.Spacing)
		l.Bounds.Y = widget_bounds.Y + l.Padding.Y + float32(l.Spacing)
		if l.Widget.GetTitleBar() {
			l.Bounds.Y = widget_bounds.Y + widget_titlebar_bounds.Height + l.Padding.Y + float32(l.Spacing)
		}

		l.Bounds.Width = widget_bounds.Width - l.Padding.X*2 - float32(l.Spacing)
		l.Bounds.Height = widget_bounds.Height - widget_titlebar_bounds.Height - l.Padding.Y*2 - float32(l.Spacing)
		// Use fixed width/height if specified, otherwise calculate from widget

		if l.FixedWidth > 1 {
			l.Bounds.Width = l.FixedWidth
		}

		if l.FixedHeight > 1 {
			l.Bounds.Height = l.FixedHeight
		}
	} else if l.Parent != nil {
	}

	// Adjust final height based on height of All the children
	l.UpdateChildLayouts()
	l.UpdateChildWidgets()
	rl.DrawRectangleLinesEx(l.Bounds, 1, rl.Pink)
	rl.DrawTextEx(Default_Widget_Header_Font, l.Name, rl.NewVector2(l.Bounds.X, l.Bounds.Y), 15, 0, rl.Black)

}

func (l *Layout) UpdateChildLayouts() {
	if len(l.Layouts) == 0 {
		return
	}

	no_of_children := len(l.Layouts)

	for i, child_layout := range l.Layouts {

		switch l.Type {
		case LayoutHorizontal:
			// Calculate available width (subtract spacing between items)
			availableWidth := l.Bounds.Width - float32(l.Spacing*(no_of_children-1)) - float32(child_layout.Spacing)

			// Use fixed width if specified, otherwise distribute evenly
			width := availableWidth/float32(no_of_children) - float32(l.Spacing)
			if child_layout.Parent.Widget != nil {

				if child_layout.Parent.Widget.GetLayout().FixedWidth > 0 {

					child_layout.Bounds.Width = child_layout.Parent.Widget.GetLayout().FixedWidth - float32(child_layout.Spacing) - float32(l.Spacing)
				}

			}

			if child_layout.FixedWidth > 0 {

				width = child_layout.FixedWidth - float32(child_layout.Spacing) - float32(l.Spacing)
			}

			// Use fixed height if specified, otherwise use parent height
			height := l.Bounds.Height - float32(child_layout.Spacing) - float32(l.Spacing)
			if child_layout.Parent.Widget != nil {

				if child_layout.Parent.Widget.GetLayout().FixedHeight > 0 {
					child_layout.Bounds.Height = child_layout.Parent.Widget.GetLayout().FixedHeight - float32(child_layout.Spacing)
				}

			}
			final_height := height
			if child_layout.FixedHeight > 0 {
				final_height = child_layout.FixedHeight - float32(child_layout.Spacing) - float32(l.Spacing)

			}
			final_width := width
			if child_layout.FixedWidth > 0 {
				final_width = child_layout.FixedWidth - float32(child_layout.Spacing) - float32(l.Spacing)
			}

			xpos := l.Bounds.X + (width+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, final_width, final_height))

			if i > 0 {

				last_layout := l.Layouts[i-1]
				child_layout.Bounds.Y = last_layout.Bounds.Y
				child_layout.Bounds.X = last_layout.Bounds.X + child_layout.Bounds.Width + float32(l.Spacing)
			}

		case LayoutVertical:
			// Use fixed width if specified, otherwise use parent width
			width := l.Bounds.Width - float32(child_layout.Spacing) - float32(l.Spacing)

			if child_layout.Parent.Widget != nil {
				if child_layout.Parent.Widget.GetLayout().FixedWidth > 0 {

					child_layout.Bounds.Width = child_layout.Parent.Widget.GetLayout().FixedWidth - float32(child_layout.Spacing)
				}

			}
			if child_layout.FixedWidth > 0 {
				width = child_layout.FixedWidth - float32(child_layout.Spacing) - float32(l.Spacing)
			}

			// Calculate available height (subtract spacing between items)
			availableHeight := l.Bounds.Height - float32(l.Spacing*(no_of_children-1)) - float32(child_layout.Spacing)

			// Use fixed height if specified, otherwise distribute evenly
			height := availableHeight/float32(no_of_children) - float32(child_layout.Spacing)

			if child_layout.Parent.Widget != nil {
				if child_layout.Parent.Widget.GetLayout().FixedHeight > 0 {
					child_layout.Bounds.Height = child_layout.Parent.Widget.GetLayout().FixedHeight - float32(child_layout.Spacing)
				}

			}

			if child_layout.FixedHeight > 0 {
				height = child_layout.FixedHeight - float32(child_layout.Spacing)
			}

			xpos := l.Bounds.X + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + (height+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))
			if i > 0 {

				last_layout := l.Layouts[i-1]
				child_layout.Bounds.Y = last_layout.Bounds.Y + last_layout.Bounds.Height + float32(l.Spacing)
				child_layout.Bounds.X = last_layout.Bounds.X
			}

		case LayoutGrid:
			cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
			if cols < 1 {
				cols = 1
			}
			rows := (no_of_children + cols - 1) / cols

			availableW := l.Bounds.Width - float32(l.Spacing*(cols+1)) - float32(child_layout.Spacing)
			availableH := l.Bounds.Height - float32(l.Spacing*(rows+1)) - float32(child_layout.Spacing)

			cellW := availableW / float32(cols)
			cellH := availableH / float32(rows)

			// Override with fixed dimensions if specified
			if child_layout.FixedWidth > 0 {
				cellW = child_layout.FixedWidth - float32(l.Spacing)
			}
			if child_layout.FixedHeight > 0 {
				cellH = child_layout.FixedHeight - float32(l.Spacing)
			}

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing)) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing)) + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, cellW, cellH))
		}

		child_layout.UpdateChildLayouts()
		child_layout.UpdateChildWidgets()

		// accumulating height and width  if all children
		if child_layout.FixedHeight > 0 {
			total_height := child_layout.Bounds.Height + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_height += child.Bounds.Height

			}
			l.Bounds.Height = total_height
		}
		if child_layout.FixedWidth > 0 {
			total_width := child_layout.Bounds.Width + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_width += child.Bounds.Width
			}
			l.Bounds.Width = total_width
		}

		if child_layout.DebugDraw {
			rl.DrawRectangleLinesEx(child_layout.Bounds, 1, rl.Pink)
			rl.DrawTextEx(Default_Widget_Header_Font, child_layout.Name, rl.NewVector2(child_layout.Bounds.X, child_layout.Bounds.Y), 15, 0, rl.Black)
		}
	}
}

func (l *Layout) UpdateChildWidgets() {
	if len(l.Children) == 0 {
		return
	}

	no_of_children := len(l.Children)

	for i, child_widget := range l.Children {

		switch l.Type {
		case LayoutHorizontal:

			width := (l.Bounds.Width / float32(no_of_children)) - float32(l.Spacing)
			height := l.Bounds.Height - float32(l.Spacing)
			xpos := l.Bounds.X + (width * float32(i)) + float32(l.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing)

			child_widget.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutVertical:
			width := l.Bounds.Width - float32(l.Spacing) - float32(i*no_of_children)
			height := (l.Bounds.Height / float32(no_of_children)) - float32(l.Spacing)

			xpos := l.Bounds.X + float32(l.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + (height * float32(i)) + float32(l.Spacing)

			child_widget.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

		case LayoutGrid:
			cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
			if cols < 1 {
				cols = 1
			}
			rows := (no_of_children + cols - 1) / cols

			availableW := l.Bounds.Width - float32(l.Spacing*(cols+1))
			availableH := l.Bounds.Height - float32(l.Spacing*(rows+1))*float32(rows)

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

func (l *Layout) Draw() {
	if !l.Visible {
		return
	}

	for _, child_layout := range l.Layouts {
		child_layout.Draw()
	}
	for _, child_widget := range l.Children {
		child_widget.Draw()
	}
}
