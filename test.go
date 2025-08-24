//go:build windows
// +build windows

package main

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

/*
Minimal Win32 + OpenGL 1.1 (no cgo)
- Window creation via user32
- Pixel format via gdi32
- OpenGL context via opengl32 (WGL)
- Upload RGBA buffer to texture and draw full-screen quad
- ESC to quit
*/
var (
	procTextOutW  = gdi32.NewProc("TextOutW")
	procSetBkMode = gdi32.NewProc("SetBkMode")
)

var (
	// DLLs
	user32   = syscall.NewLazyDLL("user32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	opengl32 = syscall.NewLazyDLL("opengl32.dll")

	// user32 procs
	procRegisterClassW     = user32.NewProc("RegisterClassW")
	procCreateWindowExW    = user32.NewProc("CreateWindowExW")
	procDefWindowProcW     = user32.NewProc("DefWindowProcW")
	procShowWindow         = user32.NewProc("ShowWindow")
	procUpdateWindow       = user32.NewProc("UpdateWindow")
	procGetMessageW        = user32.NewProc("GetMessageW")
	procTranslateMessage   = user32.NewProc("TranslateMessage")
	procDispatchMessageW   = user32.NewProc("DispatchMessageW")
	procPostQuitMessage    = user32.NewProc("PostQuitMessage")
	procLoadCursorW        = user32.NewProc("LoadCursorW")
	procLoadIconW          = user32.NewProc("LoadIconW")
	procGetModuleHandleW   = kernel32.NewProc("GetModuleHandleW")
	procAdjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
	procPeekMessageW       = user32.NewProc("PeekMessageW")
	procSleep              = kernel32.NewProc("Sleep")
	procDestroyWindow      = user32.NewProc("DestroyWindow")

	// gdi32 procs
	procChoosePixelFormat = gdi32.NewProc("ChoosePixelFormat")
	procSetPixelFormat    = gdi32.NewProc("SetPixelFormat")
	procSwapBuffers       = gdi32.NewProc("SwapBuffers")

	// opengl32 / wgl
	procWglCreateContext = opengl32.NewProc("wglCreateContext")
	procWglMakeCurrent   = opengl32.NewProc("wglMakeCurrent")
	procWglDeleteContext = opengl32.NewProc("wglDeleteContext")

	// OpenGL 1.1 functions
	glViewport       = opengl32.NewProc("glViewport")
	glClearColor     = opengl32.NewProc("glClearColor")
	glClear          = opengl32.NewProc("glClear")
	glEnable         = opengl32.NewProc("glEnable")
	glDisable        = opengl32.NewProc("glDisable")
	glMatrixMode     = opengl32.NewProc("glMatrixMode")
	glLoadIdentity   = opengl32.NewProc("glLoadIdentity")
	glOrtho          = opengl32.NewProc("glOrtho")
	glBegin          = opengl32.NewProc("glBegin")
	glEnd            = opengl32.NewProc("glEnd")
	glTexCoord2f     = opengl32.NewProc("glTexCoord2f")
	glVertex2f       = opengl32.NewProc("glVertex2f")
	glBindTexture    = opengl32.NewProc("glBindTexture")
	glTexParameteri  = opengl32.NewProc("glTexParameteri")
	glTexImage2D     = opengl32.NewProc("glTexImage2D")
	glTexSubImage2D  = opengl32.NewProc("glTexSubImage2D")
	glGenTextures    = opengl32.NewProc("glGenTextures")
	glDeleteTextures = opengl32.NewProc("glDeleteTextures")
	glPixelStorei    = opengl32.NewProc("glPixelStorei")
)

// --- Win32 types/structs ---

type (
	HINSTANCE = uintptr
	HWND      = uintptr
	HDC       = uintptr
	HGLRC     = uintptr
	HBRUSH    = uintptr
	HICON     = uintptr
	HCURSOR   = uintptr
	ATOM      = uint16
)

type WNDCLASS struct {
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       HCURSOR
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type PIXELFORMATDESCRIPTOR struct {
	NSize           uint16
	NVersion        uint16
	DwFlags         uint32
	IPixelType      byte
	CColorBits      byte
	CRedBits        byte
	CRedShift       byte
	CGreenBits      byte
	CGreenShift     byte
	CBlueBits       byte
	CBlueShift      byte
	CAlphaBits      byte
	CAlphaShift     byte
	CAccumBits      byte
	CAccumRedBits   byte
	CAccumGreenBits byte
	CAccumBlueBits  byte
	CAccumAlphaBits byte
	CDepthBits      byte
	CStencilBits    byte
	CAuxBuffers     byte
	ILayerType      byte
	BReserved       byte
	DwLayerMask     uint32
	DwVisibleMask   uint32
	DwDamageMask    uint32
}

// --- Win32 & GL constants ---

const (
	CS_OWNDC   = 0x0020
	CS_HREDRAW = 0x0002
	CS_VREDRAW = 0x0001

	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000

	CW_USEDEFAULT = 0x80000000

	SW_SHOW = 5

	WM_DESTROY = 0x0002
	WM_SIZE    = 0x0005
	WM_KEYDOWN = 0x0100
	WM_CLOSE   = 0x0010

	VK_ESCAPE = 0x1B

	PFD_DRAW_TO_WINDOW = 0x00000004
	PFD_SUPPORT_OPENGL = 0x00000020
	PFD_DOUBLEBUFFER   = 0x00000001
	PFD_TYPE_RGBA      = 0
	PFD_MAIN_PLANE     = 0

	PM_REMOVE = 0x0001

	// GL enums
	GL_COLOR_BUFFER_BIT   = 0x00004000
	GL_PROJECTION         = 0x1701
	GL_MODELVIEW          = 0x1700
	GL_TEXTURE_2D         = 0x0DE1
	GL_QUADS              = 0x0007
	GL_NEAREST            = 0x2600
	GL_CLAMP              = 0x2900
	GL_UNPACK_ALIGNMENT   = 0x0CF5
	GL_RGBA               = 0x1908
	GL_UNSIGNED_BYTE      = 0x1401
	GL_TEXTURE_MIN_FILTER = 0x2801
	GL_TEXTURE_MAG_FILTER = 0x2800
	GL_TEXTURE_WRAP_S     = 0x2802
	GL_TEXTURE_WRAP_T     = 0x2803
)

var (
	winWidth  = int32(800)
	winHeight = int32(600)

	className = syscall.StringToUTF16Ptr("GoGLWindowClass")
	title     = syscall.StringToUTF16Ptr("Go OpenGL (no cgo) - ESC to Quit")

	hwnd  HWND
	hdc   HDC
	hglrc HGLRC

	texID  uint32
	buffer []byte
)

func loword(l uintptr) uint16 { return uint16(l & 0xFFFF) }
func hiword(l uintptr) uint16 { return uint16((l >> 16) & 0xFFFF) }

func main_() {
	runtime.LockOSThread() // all Win32 + GL must stay on one OS thread

	inst := HINSTANCE(getModuleHandle())
	registerClass(inst)

	hwnd = createWindow(inst, winWidth, winHeight)
	if hwnd == 0 {
		panic("CreateWindowExW failed")
	}

	hdc = getDC(hwnd)
	setupPixelFormat(hdc)
	hglrc = createGLContext(hdc)
	makeCurrent(hdc, hglrc)

	initGL(int(winWidth), int(winHeight))
	initTexture(int(winWidth), int(winHeight))
	makeGradient(int(winWidth), int(winHeight)) // fill initial buffer

	// Show window
	showWindow(hwnd, SW_SHOW)
	updateWindow(hwnd)

	// Main loop: pump messages; if none, draw a frame.
	var msg MSG
loop:
	for {
		// PeekMessage (non-blocking)
		hasMsg := peekMessage(&msg, 0, 0, 0, PM_REMOVE)
		if hasMsg {
			switch msg.Message {
			case WM_DESTROY:
				break loop
			default:
				translateMessage(&msg)
				dispatchMessage(&msg)
			}
		} else {
			// Draw frame
			drawFrame()
			swapBuffers(hdc)

			fps := calculateFPS()
			drawText(uintptr(hdc), 10, 20, fmt.Sprintf("FPS: %d", fps))

			// ~60 FPS
			sleep(16)
		}
	}

	// Cleanup
	deleteTexture()
	makeCurrent(0, 0)
	deleteContext(hglrc)
	destroyWindow(hwnd)
}

// ---------------- Window & GL setup ----------------

func registerClass(hinst HINSTANCE) {
	var wc WNDCLASS
	wc.Style = CS_HREDRAW | CS_VREDRAW | CS_OWNDC
	wc.LpfnWndProc = syscall.NewCallback(wndProc)
	wc.HInstance = hinst
	wc.HCursor = loadCursor(0, 32512) // IDC_ARROW
	wc.HIcon = loadIcon(0, 32512)     // IDI_APPLICATION
	wc.HbrBackground = 0
	wc.LpszClassName = className

	r, _, err := procRegisterClassW.Call(uintptr(unsafe.Pointer(&wc)))
	if r == 0 {
		panic(fmt.Sprintf("RegisterClassW failed: %v", err))
	}
}

func createWindow(hinst HINSTANCE, width, height int32) HWND {
	// Adjust window rect so client area matches requested size
	rect := RECT{0, 0, width, height}
	procAdjustWindowRectEx.Call(uintptr(unsafe.Pointer(&rect)),
		uintptr(WS_OVERLAPPEDWINDOW|WS_VISIBLE), 0, 0)

	w := rect.Right - rect.Left
	h := rect.Bottom - rect.Top

	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(title)),
		uintptr(WS_OVERLAPPEDWINDOW|WS_VISIBLE),
		uintptr(CW_USEDEFAULT), uintptr(CW_USEDEFAULT),
		uintptr(w), uintptr(h),
		0, 0, uintptr(hinst), 0,
	)
	return HWND(hwnd)
}

func wndProc(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_SIZE:
		w := int32(loword(lparam))
		h := int32(hiword(lparam))
		if w == 0 || h == 0 {
			break
		}
		winWidth, winHeight = w, h
		call(glViewport, 0, 0, int32ToUintptr(w), int32ToUintptr(h))
		call(glMatrixMode, uintptr(GL_PROJECTION))
		call(glLoadIdentity)
		call(glOrtho, float64(0), float64(1), float64(1), float64(0), float64(-1), float64(1))
		call(glMatrixMode, uintptr(GL_MODELVIEW))
		call(glLoadIdentity)

	case WM_KEYDOWN:
		if wparam == VK_ESCAPE {
			procDestroyWindow.Call(uintptr(hwnd))
		}

	case WM_CLOSE:
		procDestroyWindow.Call(uintptr(hwnd))

	case WM_DESTROY:
		postQuitMessage(0)

	default:
		r, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
		return r
	}
	return 0
}

// ---------------- GL helpers ----------------

func initGL(w, h int) {
	call(glViewport, 0, 0, int32ToUintptr(int32(w)), int32ToUintptr(int32(h)))
	call(glClearColor, float32ToUintptr(0.1), float32ToUintptr(0.1), float32ToUintptr(0.12), float32ToUintptr(1.0))
	call(glDisable, uintptr(0x0B71)) // GL_DEPTH_TEST
	// Simple 2D ortho 0..1
	call(glMatrixMode, uintptr(GL_PROJECTION))
	call(glLoadIdentity)
	call(glOrtho, 0, 1, 1, 0, -1, 1)
	call(glMatrixMode, uintptr(GL_MODELVIEW))
	call(glLoadIdentity)
	call(glEnable, uintptr(GL_TEXTURE_2D))
}

func initTexture(w, h int) {
	// Generate 1 texture
	var id uint32
	call(glGenTextures, 1, uintptr(unsafe.Pointer(&id)))
	texID = id
	call(glBindTexture, uintptr(GL_TEXTURE_2D), uintptr(texID))

	// Parameters
	call(glTexParameteri, uintptr(GL_TEXTURE_2D), uintptr(GL_TEXTURE_MIN_FILTER), uintptr(GL_NEAREST))
	call(glTexParameteri, uintptr(GL_TEXTURE_2D), uintptr(GL_TEXTURE_MAG_FILTER), uintptr(GL_NEAREST))
	call(glTexParameteri, uintptr(GL_TEXTURE_2D), uintptr(GL_TEXTURE_WRAP_S), uintptr(GL_CLAMP))
	call(glTexParameteri, uintptr(GL_TEXTURE_2D), uintptr(GL_TEXTURE_WRAP_T), uintptr(GL_CLAMP))
	call(glPixelStorei, uintptr(GL_UNPACK_ALIGNMENT), 1)

	// Allocate storage
	buffer = make([]byte, w*h*4)
	call(glTexImage2D,
		uintptr(GL_TEXTURE_2D), 0, uintptr(GL_RGBA),
		uintptr(w), uintptr(h), 0,
		uintptr(GL_RGBA), uintptr(GL_UNSIGNED_BYTE),
		uintptr(unsafe.Pointer(&buffer[0])),
	)
}

func deleteTexture() {
	if texID != 0 {
		call(glDeleteTextures, 1, uintptr(unsafe.Pointer(&texID)))
		texID = 0
	}
}

func drawFrame() {
	// Update CPU buffer (animate the gradient a bit)
	t := float32(time.Now().UnixNano()%2_000_000_000) / 2_000_000_000
	updateGradient(int(winWidth), int(winHeight), t)

	// Upload subimage
	call(glBindTexture, uintptr(GL_TEXTURE_2D), uintptr(texID))
	call(glTexSubImage2D,
		uintptr(GL_TEXTURE_2D), 0, 0, 0,
		uintptr(winWidth), uintptr(winHeight),
		uintptr(GL_RGBA), uintptr(GL_UNSIGNED_BYTE),
		uintptr(unsafe.Pointer(&buffer[0])),
	)

	// Clear and draw textured quad covering 0..1 in X/Y
	call(glClear, uintptr(GL_COLOR_BUFFER_BIT))
	call(glBegin, uintptr(GL_QUADS))

	call(glTexCoord2f, float32ToUintptr(0), float32ToUintptr(1))
	call(glVertex2f, float32ToUintptr(0), float32ToUintptr(float32(winHeight)))

	call(glTexCoord2f, float32ToUintptr(1), float32ToUintptr(1))
	call(glVertex2f, float32ToUintptr(float32(winWidth)), float32ToUintptr(float32(winHeight)))

	call(glTexCoord2f, float32ToUintptr(1), float32ToUintptr(0))
	call(glVertex2f, float32ToUintptr(float32(winWidth)), float32ToUintptr(0))

	call(glTexCoord2f, float32ToUintptr(0), float32ToUintptr(0))
	call(glVertex2f, float32ToUintptr(0), float32ToUintptr(0))

	call(glEnd)

}

// --------------- CPU gradient buffer ----------------

func makeGradient(w, h int) {
	buffer = make([]byte, w*h*4)
	updateGradient(w, h, 0)
}

func updateGradient(w, h int, t float32) {
	if len(buffer) < w*h*4 {
		buffer = make([]byte, w*h*4)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 4
			r := byte((x * 255) / max(1, w-1))
			g := byte((y * 255) / max(1, h-1))
			// animate blue channel with time
			b := byte((int((float32(x+y)/float32(w+h))*255) + int(t*255)) & 0xFF)
			buffer[i+0] = r
			buffer[i+1] = g
			buffer[i+2] = b
			buffer[i+3] = 255
		}
	}
}

// --------------- Thin syscall wrappers ----------------

func getModuleHandle() HINSTANCE {
	h, _, _ := procGetModuleHandleW.Call(0)
	return HINSTANCE(h)
}

func loadCursor(hinst HINSTANCE, id uintptr) HCURSOR {
	h, _, _ := procLoadCursorW.Call(uintptr(hinst), id)
	return HCURSOR(h)
}

func loadIcon(hinst HINSTANCE, id uintptr) HICON {
	h, _, _ := procLoadIconW.Call(uintptr(hinst), id)
	return HICON(h)
}

func showWindow(hwnd HWND, nCmdShow int32) {
	procShowWindow.Call(uintptr(hwnd), uintptr(nCmdShow))
}
func updateWindow(hwnd HWND) {
	procUpdateWindow.Call(uintptr(hwnd))
}
func destroyWindow(hwnd HWND) {
	procDestroyWindow.Call(uintptr(hwnd))
}

func postQuitMessage(code int32) {
	procPostQuitMessage.Call(uintptr(code))
}

func getDC(hwnd HWND) HDC {
	// GetDC is user32, but we avoid keeping a separate proc. Use GetDC via user32.
	getDC := user32.NewProc("GetDC")
	h, _, _ := getDC.Call(uintptr(hwnd))
	return HDC(h)
}

func setupPixelFormat(hdc HDC) {
	var pfd PIXELFORMATDESCRIPTOR
	pfd.NSize = uint16(unsafe.Sizeof(pfd))
	pfd.NVersion = 1
	pfd.DwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER
	pfd.IPixelType = PFD_TYPE_RGBA
	pfd.CColorBits = 32
	pfd.CDepthBits = 24
	pfd.CAlphaBits = 8
	pfd.ILayerType = PFD_MAIN_PLANE

	i, _, _ := procChoosePixelFormat.Call(uintptr(hdc), uintptr(unsafe.Pointer(&pfd)))
	if i == 0 {
		panic("ChoosePixelFormat failed")
	}
	res, _, _ := procSetPixelFormat.Call(uintptr(hdc), i, uintptr(unsafe.Pointer(&pfd)))
	if res == 0 {
		panic("SetPixelFormat failed")
	}
}

func createGLContext(hdc HDC) HGLRC {
	rc, _, _ := procWglCreateContext.Call(uintptr(hdc))
	if rc == 0 {
		panic("wglCreateContext failed")
	}
	return HGLRC(rc)
}

func makeCurrent(hdc HDC, hglrc HGLRC) {
	ok, _, _ := procWglMakeCurrent.Call(uintptr(hdc), uintptr(hglrc))
	if ok == 0 {
		panic("wglMakeCurrent failed")
	}
}
func deleteContext(hglrc HGLRC) {
	if hglrc != 0 {
		procWglDeleteContext.Call(uintptr(hglrc))
	}
}

func swapBuffers(hdc HDC) {
	procSwapBuffers.Call(uintptr(hdc))
}

func translateMessage(msg *MSG) {
	procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
}
func dispatchMessage(msg *MSG) {
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
}
func getMessage(msg *MSG, hwnd HWND, min, max uint32) bool {
	r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(min), uintptr(max))
	return int32(r) > 0
}
func peekMessage(msg *MSG, hwnd HWND, min, max uint32, remove uint32) bool {
	r, _, _ := procPeekMessageW.Call(uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(min), uintptr(max), uintptr(remove))
	return r != 0
}
func sleep(ms int) {
	procSleep.Call(uintptr(ms))
}

// --------------- tiny call helpers & converters ---------------

// replaces: return syscall.SyscallN(p.Addr(), u...)
func call(p *syscall.LazyProc, args ...interface{}) uintptr {
	var u []uintptr
	for _, a := range args {
		switch v := a.(type) {
		case uintptr:
			u = append(u, v)
		case int:
			u = append(u, uintptr(v))
		case int32:
			u = append(u, uintptr(uint32(v)))
		case uint32:
			u = append(u, uintptr(v))
		case float32:
			u = append(u, *(*uintptr)(unsafe.Pointer(&v)))
		case float64:
			u = append(u, *(*uintptr)(unsafe.Pointer(&v)))
		default:
			panic("unsupported arg type in call()")
		}
	}
	r1, _, _ := p.Call(u...) // <-- take only r1
	return r1
}

func float32ToUintptr(f float32) uintptr {
	return *(*uintptr)(unsafe.Pointer(&f))
}
func int32ToUintptr(i int32) uintptr { return uintptr(uint32(i)) }

// --------------- misc helpers ---------------

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func drawText(hdc uintptr, x, y int, s string) {
	utf16, _ := syscall.UTF16PtrFromString(s)
	procSetBkMode.Call(hdc, uintptr(1)) // TRANSPARENT background
	procTextOutW.Call(hdc, uintptr(x), uintptr(y), uintptr(unsafe.Pointer(utf16)), uintptr(len(s)))
}

var (
	lastTime   = time.Now()
	frameCount = 0
	fps        = 0
)

func calculateFPS() int {
	frameCount++
	now := time.Now()
	if now.Sub(lastTime) >= time.Second {
		fps = frameCount
		frameCount = 0
		lastTime = now
	}
	return fps
}
