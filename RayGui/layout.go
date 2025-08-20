package RayGui

import (
	"fmt"
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
	// Don't add the child's layout to the Layouts slice
	// Layouts slice should only contain pure layout objects, not widget-attached layouts
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

		l.Bounds.X = widget_bounds.X + l.Padding.X
		if l.Widget.GetTitleBar() {
			l.Bounds.Y = widget_bounds.Y + widget_titlebar_bounds.Height + l.Padding.Y
		} else {
			l.Bounds.Y = widget_bounds.Y + l.Padding.Y
		}

		// Use fixed width/height if specified, otherwise calculate from widget
		if l.FixedWidth > 0 {
			l.Bounds.Width = l.FixedWidth
		} else {
			l.Bounds.Width = widget_bounds.Width - l.Padding.X*2
		}

		if l.FixedHeight > 0 {
			l.Bounds.Height = l.FixedHeight
		} else {
			l.Bounds.Height = widget_bounds.Height - widget_titlebar_bounds.Height - l.Padding.Y*2
		}
	} else if l.Parent != nil {
		// For nested layouts, bounds are already set by parent layout in UpdateChildLayouts()
		// We don't need to recalculate them here
	}

	// Update child layouts FIRST to ensure proper bounds calculation
	l.UpdateChildLayouts()
	l.UpdateChildWidgets()

	// Draw debug rectangles after bounds are properly set
	if l.Parent != nil && len(l.Layouts) > 0 {
		for _, child_layout := range l.Layouts {
			rl.DrawRectangleLinesEx(child_layout.Bounds, 1, rl.Pink)
		}
	}
}

func (l *Layout) UpdateChildLayouts() {
	if len(l.Layouts) == 0 {
		return
	}

	no_of_children := len(l.Layouts)

	switch l.Type {
	case LayoutVertical:
		// First pass: calculate total fixed height and count flexible children
		var fixedHeightTotal float32 = 0
		var flexibleChildren int = 0

		for _, child_layout := range l.Layouts {
			if child_layout.FixedHeight > 0 {
				fixedHeightTotal += child_layout.FixedHeight
			} else {
				flexibleChildren++
			}
		}

		// Calculate available height for flexible children
		availableHeight := l.Bounds.Height - float32(l.Spacing*(no_of_children-1))
		availableFlexibleHeight := availableHeight - fixedHeightTotal
		flexibleHeight := float32(0)
		if flexibleChildren > 0 && availableFlexibleHeight > 0 {
			flexibleHeight = availableFlexibleHeight / float32(flexibleChildren)
		}

		// Second pass: set bounds for each child layout
		var currentY float32 = l.Bounds.Y

		for _, child_layout := range l.Layouts {
			width := l.Bounds.Width
			height := float32(0)

			// Use fixed height if specified, otherwise use calculated flexible height
			if child_layout.FixedHeight > 0 {
				height = child_layout.FixedHeight
			} else {
				height = flexibleHeight
			}

			// Ensure height doesn't exceed available space
			if height < 0 {
				height = 0
			}

			// Use fixed width if specified for the child layout
			if child_layout.FixedWidth > 0 {
				width = child_layout.FixedWidth
			}

			child_layout.SetBounds(rl.NewRectangle(l.Bounds.X, currentY, width, height))
			currentY += height + float32(l.Spacing)
		}

	case LayoutGrid:
		cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
		if cols < 1 {
			cols = 1
		}
		rows := (no_of_children + cols - 1) / cols

		availableW := l.Bounds.Width - float32(l.Spacing*(cols+1))
		availableH := l.Bounds.Height - float32(l.Spacing*(rows+1))

		// Default cell dimensions
		cellW := availableW / float32(cols)
		cellH := availableH / float32(rows)

		for i, child_layout := range l.Layouts {
			// Use fixed dimensions if specified, otherwise use calculated cell dimensions
			width := cellW
			if child_layout.FixedWidth > 0 {
				width = child_layout.FixedWidth
			}

			height := cellH
			if child_layout.FixedHeight > 0 {
				height = child_layout.FixedHeight
			}

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing))
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing))

			child_layout.SetBounds(rl.NewRectangle(xpos, ypos, width, height))
		}
	}

	// Update all child layouts
	for _, child_layout := range l.Layouts {
		child_layout.Update()
	}
}

func (l *Layout) UpdateChildWidgets() {
	if len(l.Children) == 0 {
		return
	}

	no_of_children := len(l.Children)

	switch l.Type {
	case LayoutHorizontal:
		// First pass: calculate total fixed width and count flexible children
		var fixedWidthTotal float32 = 0
		var flexibleChildren int = 0

		for _, child := range l.Children {
			childLayout := child.GetLayout()
			if childLayout.FixedWidth > 0 {
				fixedWidthTotal += childLayout.FixedWidth
			} else {
				flexibleChildren++
			}
		}

		// Calculate available width for flexible children
		availableWidth := l.Bounds.Width - float32(l.Spacing*(no_of_children-1))
		availableFlexibleWidth := availableWidth - fixedWidthTotal
		flexibleWidth := float32(0)
		if flexibleChildren > 0 {
			flexibleWidth = availableFlexibleWidth / float32(flexibleChildren)
		}

		// Second pass: set bounds for each child
		var currentX float32 = l.Bounds.X

		for _, child := range l.Children {
			childLayout := child.GetLayout()

			width := float32(0)
			height := l.Bounds.Height

			// Use fixed width if specified, otherwise use calculated flexible width
			if childLayout.FixedWidth > 0 {
				width = childLayout.FixedWidth
			} else {
				width = flexibleWidth
			}

			// Use fixed height if specified for the child
			if childLayout.FixedHeight > 0 {
				height = childLayout.FixedHeight
			}

			child.SetBounds(rl.NewRectangle(currentX, l.Bounds.Y, width, height))
			currentX += width + float32(l.Spacing)
		}

	case LayoutVertical:
		// First pass: calculate total fixed height and count flexible children
		var fixedHeightTotal float32 = 0
		var flexibleChildren int = 0

		for _, child := range l.Children {
			childLayout := child.GetLayout()
			if childLayout.FixedHeight > 0 {
				fixedHeightTotal += childLayout.FixedHeight
			} else {
				flexibleChildren++
			}
		}

		// Calculate available height for flexible children
		availableHeight := l.Bounds.Height - float32(l.Spacing*(no_of_children-1))
		availableFlexibleHeight := availableHeight - fixedHeightTotal
		flexibleHeight := float32(0)
		if flexibleChildren > 0 {
			flexibleHeight = availableFlexibleHeight / float32(flexibleChildren)
		}

		// Second pass: set bounds for each child
		var currentY float32 = l.Bounds.Y

		for _, child := range l.Children {
			childLayout := child.GetLayout()

			width := l.Bounds.Width
			height := float32(0)

			// Use fixed height if specified, otherwise use calculated flexible height
			if childLayout.FixedHeight > 0 {
				height = childLayout.FixedHeight
			} else {
				height = flexibleHeight
			}

			// Use fixed width if specified for the child
			if childLayout.FixedWidth > 0 {
				width = childLayout.FixedWidth
			}

			child.SetBounds(rl.NewRectangle(l.Bounds.X, currentY, width, height))
			currentY += height + float32(l.Spacing)
		}

	case LayoutGrid:
		cols := int(math.Ceil(math.Sqrt(float64(no_of_children))))
		if cols < 1 {
			cols = 1
		}
		rows := (no_of_children + cols - 1) / cols

		availableW := l.Bounds.Width - float32(l.Spacing*(cols+1))
		availableH := l.Bounds.Height - float32(l.Spacing*(rows+1))

		// Default cell dimensions
		cellW := availableW / float32(cols)
		cellH := availableH / float32(rows)

		for i, child := range l.Children {
			childLayout := child.GetLayout()

			// Use fixed dimensions if specified, otherwise use calculated cell dimensions
			width := cellW
			if childLayout.FixedWidth > 0 {
				width = childLayout.FixedWidth
			}

			height := cellH
			if childLayout.FixedHeight > 0 {
				height = childLayout.FixedHeight
			}

			row := i / cols
			col := i % cols

			xpos := l.Bounds.X + float32(l.Spacing) + float32(col)*(cellW+float32(l.Spacing))
			ypos := l.Bounds.Y + float32(l.Spacing) + float32(row)*(cellH+float32(l.Spacing))

			child.SetBounds(rl.NewRectangle(xpos, ypos, width, height))
		}
	}
}

func (l *Layout) GetMinSize() rl.Vector2 {
	var minWidth, minHeight float32 = 0, 0

	titleBarHeight := float32(1)

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

	if l.FixedHeight > 0 {
		l.Bounds.Height = l.FixedHeight
		fmt.Println(l.Bounds)

	}
	if l.FixedWidth > 0 {
		l.Bounds.Width = l.FixedWidth
	}
	if l.Name == "MenuBarLayout" || l.Name == "MidPanelLayout" || l.Name == "LowerPanelLayout" {
		rl.DrawRectangleLinesEx(l.Bounds, 1, rl.Pink)
	}
}
