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
	GetVisibility() bool
	GetBgColor() rl.Color
	GetTextFont() rl.Font
	GetTextColor() rl.Color
	GetName() string
	GetTitleBar() bool
	GetLayout() *Layout
	GetZIndex() int
	MainWindow() bool
}

type BoundsSetter interface {
	SetBounds(bounds rl.Rectangle)
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
		DebugDraw:     true,
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
	l.minimumHeight = 0
	l.Bounds.Height = height
}

func (l *Layout) GetFixedWidth() float32 {
	return l.fixedWidth
}

func (l *Layout) SetFixedWidth(width float32) {
	l.fixedWidth = width
	l.minimumWidth = 0
	l.maximumWidth = 0
	l.Bounds.Width = width
}

func (l *Layout) GetMinimumWidth() float32 {
	return l.minimumWidth
}

func (l *Layout) GetMaximumWidth() float32 {
	return l.maximumWidth
}

func (l *Layout) GetMinimumHeight() float32 {
	return l.minimumHeight
}

func (l *Layout) GetMaximumHeight() float32 {
	return l.maximumheight
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
	l.Bounds.Height = l.getSanitizedHeight(bounds.Height)
	l.Bounds.Width = l.getSanitizedWidth(bounds.Width)
}

func (l *Layout) SetMinimumWidth(minimum_width float32) {
	if l.maximumWidth > 0 {
		if minimum_width > l.maximumWidth {
			panic_error := fmt.Sprintf(
				"ERROR while seting minimumWidth on layout %v!, minimum width (%v)"+
					" cannot be greater than maximum width (%v)",
				l.Name, minimum_width, l.maximumWidth)
			panic(panic_error)
		}

	}

	l.minimumWidth = minimum_width
	l.fixedWidth = 0
}

func (l *Layout) SetMaximumWidth(maximum_width float32) {
	if l.minimumWidth > 0 {
		if maximum_width < l.minimumWidth {
			panic_error := fmt.Sprintf(
				"ERROR while seting maximumWidth of Layout %v!"+
					" maximum width (%v) cannot be less than minimum width (%v)",
				l.Name, maximum_width, l.minimumWidth)
			panic(panic_error)
		}

	}

	l.maximumWidth = maximum_width
	l.fixedWidth = 0
}

func (l *Layout) SetMinimumHeight(minimum_height float32) {
	if l.maximumheight > 0 {
		if minimum_height > l.maximumheight {
			panic_error := fmt.Sprintf(
				"ERROR while seting MinumumWidth! on layout %v,"+
					" minimum width (%v) cannot be greater than maximum width (%v)",
				l.Name, minimum_height, l.maximumheight)
			panic(panic_error)
		}
	}
	l.minimumHeight = minimum_height
	l.fixedHeight = 0
}

func (l *Layout) SetMaximumHeight(maximum_height float32) {
	if l.minimumHeight > 0 {
		if maximum_height < l.minimumHeight {
			panic_error := fmt.Sprintf(
				"ERROR while seting maximumWidth! on layout %v, maximum width (%v)"+
					" cannot be less than minimum width (%v)",
				l.Name, maximum_height, l.maximumheight)
			panic(panic_error)
		}
	}
	l.maximumheight = maximum_height
	l.fixedHeight = 0
}

func (l *Layout) Update() {

	l.UpdateChildLayouts()

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
	no_of_non_fixed_width_children := float32(0)
	no_of_non_fixed_height_children := float32(0)

	fixed_width := float32(0)
	fixed_height := float32(0)

	for _, child_layout := range l.Layouts {
		if child_layout.fixedWidth > 0 {
			fixed_width += child_layout.fixedWidth
		} else {
			no_of_non_fixed_width_children += 1
		}
		if child_layout.fixedHeight > 0 {
			fixed_height += child_layout.fixedHeight
		} else {
			no_of_non_fixed_height_children += 1
		}
	}

	defered_width := l.Bounds.Width - fixed_width

	for i, child_layout := range l.Layouts {

		switch l.Type {
		case LayoutHorizontal:

			available_width := l.Bounds.Width / float32(no_of_children)
			if child_layout.fixedWidth == 0 { // stretch to fit if not fixed width
				available_width = defered_width / no_of_non_fixed_width_children
			}
			width := child_layout.getSanitizedWidth(available_width)
			height := child_layout.getSanitizedHeight(l.Bounds.Height)
			xpos := l.Bounds.X + (width+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))

			if i > 0 {
				last_layout := l.Layouts[i-1]
				child_layout.Bounds.X = last_layout.Bounds.X + last_layout.Bounds.Width + float32(l.Spacing)
			}

		case LayoutVertical:
			width := child_layout.getSanitizedWidth(l.Bounds.Width)
			availableHeight := l.Bounds.Height
			if child_layout.fixedHeight == 0 { // stretch to fit if not fixed width
				availableHeight = defered_width / no_of_non_fixed_height_children
			}
			height := child_layout.getSanitizedHeight(availableHeight)

			xpos := l.Bounds.X + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + (height+float32(l.Spacing))*float32(i) + float32(child_layout.Spacing)
			if i > 0 {
				ypos = l.Bounds.Y + availableHeight
			}
			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))
			if i > 0 {

				last_layout := l.Layouts[i-1]
				child_layout.Bounds.Y = last_layout.Bounds.Y + last_layout.Bounds.Height + float32(l.Spacing)
			}

		case LayoutGrid:
			cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
			if cols < 1 {
				cols = 1
			}
			rows := (no_of_children + cols - 1) / cols

			availableW := l.Bounds.Width - float32(l.Spacing*(cols+1)) - float32(child_layout.Spacing)
			availableH := l.Bounds.Height - float32(l.Spacing*(rows+1)) - float32(child_layout.Spacing)

			cellW := child_layout.getSanitizedWidth(availableW / float32(cols))
			cellH := child_layout.getSanitizedHeight(availableH / float32(rows))

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing)) + float32(child_layout.Spacing)
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing)) + float32(child_layout.Spacing)

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, cellW, cellH))

		}

		child_layout.UpdateChildLayouts()

		// accumulating height and width  if all children

		if child_layout.fixedHeight > 0 {
			total_height := child_layout.Bounds.Height + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_height += child.Bounds.Height

			}
			l.Bounds.Height = l.getSanitizedHeight(total_height)
		}

		if child_layout.fixedWidth > 0 {
			total_width := child_layout.Bounds.Width + float32(l.Spacing) + float32(child_layout.Spacing)

			for _, child := range child_layout.Layouts {
				total_width += child.Bounds.Width
			}
			l.Bounds.Width = l.getSanitizedWidth(total_width)
		}

	}

}

func (l *Layout) getSanitizedWidth(width float32) float32 {
	sanitized_width := width - float32(l.Spacing)

	if l.fixedWidth > 1 {

		sanitized_width = l.fixedWidth - float32(l.Spacing)

	}
	if l.maximumWidth > 1 {
		if width > l.maximumWidth {
			sanitized_width = l.maximumWidth - float32(l.Spacing)
		}
	}
	if l.minimumWidth > 1 {
		if width < l.minimumWidth {
			sanitized_width = l.minimumWidth - float32(l.Spacing)
		}
	}

	return sanitized_width

}

func (l *Layout) getSanitizedHeight(height float32) float32 {
	sanitized_height := height - float32(l.Spacing)
	if l.fixedHeight > 1 {
		sanitized_height = l.fixedHeight - float32(l.Spacing)

	}
	if l.maximumheight > 1 {
		if height > l.maximumheight {
			sanitized_height = l.maximumheight - float32(l.Spacing)
		}
	}
	if l.minimumHeight > 1 {
		if height < l.minimumHeight {
			sanitized_height = l.minimumHeight - float32(l.Spacing)
		}
	}

	return sanitized_height

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
	if l.DebugDraw {
		rl.DrawRectangleLinesEx(l.Bounds, 1, rl.Pink)
		rl.DrawTextEx(Default_Widget_Header_Font, l.Name, rl.NewVector2(l.Bounds.X, l.Bounds.Y), 15, 0, rl.Black)
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
