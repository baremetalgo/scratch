package RayGui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var app_icon_path = "icons/puzzle.png"
var Default_Widget_Header_Font rl.Font
var Default_Widget_Body_Text_Font rl.Font
var Default_Bg_Color rl.Color = rl.DarkGray
var Default_Text_Color rl.Color = rl.White
var Default_Border_Color rl.Color = rl.Gray
var Default_Titlebar_Color rl.Color = rl.NewColor(54, 54, 54, 255)
var Default_ResizerHandler_Color rl.Color = rl.NewColor(54, 54, 54, 255)
var Default_Titlebar_Height float32 = 25.0
var Default_Header_Font_Size int32 = 14
var Default_Body_Font_Size int32 = 14

func InitializeFonts() {
	Default_Widget_Header_Font = rl.LoadFontEx("fonts/CALIBRIB.TTF", Default_Header_Font_Size, nil, 0)
	Default_Widget_Body_Text_Font = rl.LoadFontEx("fonts/CALIBRI.TTF", Default_Body_Font_Size, nil, 0)

	icon := rl.LoadImage(app_icon_path)
	rl.SetWindowIcon(*icon)
	rl.UnloadImage(icon)

}
