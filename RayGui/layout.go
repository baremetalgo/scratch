package RayGui

import (
	"fmt"
	"math"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LayoutHorizontal = 0
	LayoutVertical   = 1
	LayoutGrid       = 2
)

const (
	SizePolicyExpanding = 0
	SizePolicyMinimum   = 1
	SizePolicyMaximum   = 2
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
	GetZIndex() int
	MainWindow() bool
}

type BoundsSetter interface {
	SetBounds(bounds rl.Rectangle)
}

type ZIndexable interface {
	GetZIndex() int
	SetZIndex(zIndex int)
}

type Layout struct {
	Name          string
	Type          int
	Widget        MainWidget
	Children      []MainWidget
	Layouts       []*Layout
	Parent        *Layout
	Spacing       int
	Padding       rl.Vector2
	Visible       bool
	Bounds        rl.Rectangle
	fixedHeight   float32
	fixedWidth    float32
	minimumHeight float32
	MinumumWidth  float32
	maximumheight float32
	maximumWidth  float32
	minimumWidth  float32
	SizePolicy    int
	DebugDraw     bool
}

func NewLayout() *Layout {
	return &Layout{
		Type:          LayoutVertical,
		Spacing:       10,
		Padding:       rl.NewVector2(5, 5),
		Layouts:       make([]*Layout, 0),
		Children:      make([]MainWidget, 0),
		Visible:       true,
		Bounds:        rl.NewRectangle(0, 0, 0, 0),
		DebugDraw:     false,
		fixedHeight:   0,
		fixedWidth:    0,
		minimumHeight: 0,
		minimumWidth:  0,
		maximumheight: 0,
		maximumWidth:  0,
		SizePolicy:    SizePolicyExpanding,
	}
}

func (l *Layout) GetFixedHeight() float32 {
	return l.fixedHeight
}

func (l *Layout) SetFixedHeight(height float32) {
	l.fixedHeight = height
	l.maximumheight = 0
}

func (l *Layout) GetFixedWidth() float32 {
	return l.fixedWidth
}

func (l *Layout) SetFixedWidth(width float32) {
	l.fixedWidth = width
	l.maximumWidth = 0
}

func (l *Layout) GetMinimumWidth() float32 {
	return l.minimumWidth
}

func (l *Layout) SetMinimumWidth(width float32) {
	if l.maximumWidth < width {
		panic_error := fmt.Sprintf("ERROR while seting MinumumWidth!, minimum width (%v) cannot be greater than maximum width (%v)", width, l.maximumWidth)
		fmt.Errorf(panic_error)
	}
	l.minimumWidth = width
	l.SetFixedWidth(0)
}

func (l *Layout) GetMaximumWidth() float32 {
	return l.maximumWidth
}

func (l *Layout) SetMaximumWidth(width float32) {
	if l.minimumWidth < width {
		panic_error := fmt.Sprintf("ERROR while seting maximumWidth!, maximum width (%v) cannot be less than minimum width (%v)", width, l.minimumWidth)
		fmt.Errorf(panic_error)
	}
	l.maximumWidth = width
	l.SetFixedWidth(0)
}

func (l *Layout) GetMinimumHeight() float32 {
	return l.minimumHeight
}

func (l *Layout) SetMinimumHeight(height float32) {
	if l.minimumHeight < height {
		panic_error := fmt.Sprintf("ERROR while seting MinumumWidth!, minimum width (%v) cannot be greater than maximum width (%v)", height, l.maximumheight)
		fmt.Errorf(panic_error)
	}
	l.minimumHeight = height
	l.SetFixedHeight(0)
}

func (l *Layout) GetMaximumHeight() float32 {
	return l.maximumheight
}

func (l *Layout) SetMaximumHeight(height float32) {
	if l.maximumheight < height {
		panic_error := fmt.Sprintf("ERROR while seting maximumWidth!, maximum width (%v) cannot be less than minimum width (%v)", height, l.maximumheight)
		fmt.Errorf(panic_error)
	}
	l.maximumheight = height
	l.SetFixedHeight(0)
}

func (l *Layout) AddChild(child MainWidget) {
	l.Children = append(l.Children, child)
	l.AddLayout(child.GetLayout())
	// ALL_WIDGETS = append(ALL_WIDGETS, child)
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

		l.Bounds.Width = l.getSanitizedWidth(widget_bounds.Width - l.Padding.X*2)
		l.Bounds.Height = l.getSanitizedHeight(widget_bounds.Height - widget_titlebar_bounds.Height - l.Padding.Y*2)

	} else if l.Parent != nil {
	}

	// Adjust final height based on height of All the children
	l.UpdateChildLayouts()
	l.UpdateChildWidgets()
	if l.DebugDraw {
		rl.DrawRectangleLinesEx(l.Bounds, 1, rl.Pink)
		rl.DrawTextEx(Default_Widget_Header_Font, l.Name, rl.NewVector2(l.Bounds.X, l.Bounds.Y), 15, 0, rl.Black)
	}
}

func (l *Layout) UpdateChildLayouts() {
	if len(l.Layouts) == 0 {
		return
	}

	no_of_children := len(l.Layouts)

	for i, child_layout := range l.Layouts {

		switch l.Type {
		case LayoutHorizontal:
			availableWidth := l.Bounds.Width - float32(l.Spacing*(no_of_children-1)) - float32(child_layout.Spacing)

			width := availableWidth / float32(no_of_children)

			height := l.Bounds.Height - float32(child_layout.Spacing)

			xpos := l.Bounds.X + (width+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, child_layout.getSanitizedWidth(width), child_layout.getSanitizedHeight(height)))

			if i > 0 {

				last_layout := l.Layouts[i-1]
				child_layout.Bounds.Y = last_layout.Bounds.Y
				child_layout.Bounds.X = last_layout.Bounds.X + child_layout.Bounds.Width + float32(l.Spacing)
			}

		case LayoutVertical:
			width := l.Bounds.Width - float32(child_layout.Spacing)

			availableHeight := l.Bounds.Height - float32(l.Spacing*(no_of_children-1)) - float32(child_layout.Spacing)

			height := availableHeight/float32(no_of_children) - float32(child_layout.Spacing)

			xpos := l.Bounds.X + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + (height+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, child_layout.getSanitizedWidth(width), child_layout.getSanitizedHeight(height)))
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

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing)) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing)) + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, child_layout.getSanitizedWidth(cellW), child_layout.getSanitizedHeight(cellH)))
		}

		child_layout.UpdateChildLayouts()
		child_layout.UpdateChildWidgets()

		// accumulating height and width  if all children
		if child_layout.GetFixedHeight() > 0 {
			total_height := child_layout.Bounds.Height + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_height += child.Bounds.Height

			}
			l.Bounds.Height = l.getSanitizedHeight(total_height)
		}
		if child_layout.GetFixedWidth() > 0 {
			total_width := child_layout.Bounds.Width + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_width += child.Bounds.Width
			}
			l.Bounds.Width = l.getSanitizedWidth(total_width)
		}

		if child_layout.DebugDraw {
			rl.DrawRectangleLinesEx(child_layout.Bounds, 1, rl.Pink)
			rl.DrawTextEx(Default_Widget_Header_Font, child_layout.Name, rl.NewVector2(child_layout.Bounds.X, child_layout.Bounds.Y), 15, 0, rl.Black)
		}
	}
}

func (l *Layout) getSanitizedWidth(width float32) float32 {
	if l.fixedWidth > 1 {
		return l.fixedWidth - float32(l.Spacing)

	}
	if l.maximumWidth > 1 {
		if width > l.maximumWidth {
			return l.maximumWidth - float32(l.Spacing)
		}
	}
	if l.minimumWidth > 1 {
		if width < l.minimumWidth {
			return l.minimumWidth - float32(l.Spacing)
		}
	}
	return width - float32(l.Spacing)

}

func (l *Layout) getSanitizedHeight(height float32) float32 {

	if l.fixedHeight > 1 {
		return l.fixedHeight - float32(l.Spacing)

	}
	if l.maximumheight > 1 {
		if height > l.maximumheight {
			return l.maximumheight - float32(l.Spacing)
		}
	}
	if l.minimumHeight > 1 {
		if height < l.minimumHeight {
			return l.minimumHeight - float32(l.Spacing)
		}
	}
	return height - float32(l.Spacing)

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

func (l *Layout) GetWidgets() []MainWidget {
	list_ := make([]MainWidget, 0)
	if !l.Visible {
		return list_
	}

	list_ = append(list_, l.Children...)

	return list_
}

func (l *Layout) Draw() {
	if !l.Visible {
		return
	}

	for _, child_layout := range l.Layouts {
		child_layout.Draw()
	}

	if l.Widget != nil {
		if l.Widget.MainWindow() {
			l.DrawWidgetsByDepth(ALL_WIDGETS)
		}
	}
}

func (l *Layout) DrawWidgetsByDepth(widgets []MainWidget) {
	// Sort widgets by z-index (lowest first)
	sort.Slice(widgets, func(i, j int) bool {
		return widgets[i].GetZIndex() < widgets[j].GetZIndex()
	})

	// Draw in sorted order
	for _, widget := range widgets {
		if !widget.MainWindow() {

			widget.Draw()
		}
	}
}
