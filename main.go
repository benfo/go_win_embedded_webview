package main

import (
	"syscall"
	"unsafe"

	webview "github.com/webview/webview_go"
)

const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_VISIBLE          = 0x10000000
)

var (
	user32             = syscall.NewLazyDLL("user32.dll")
	procCreateWindowEx = user32.NewProc("CreateWindowExW")
)

func createWindow() syscall.Handle {
	className, _ := syscall.UTF16PtrFromString("MyWindowClass")
	windowTitle, _ := syscall.UTF16PtrFromString("WebView App")
	hwnd, _, _ := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowTitle)),
		WS_OVERLAPPEDWINDOW|WS_VISIBLE,
		100, 100, 800, 600, 0, 0, 0, 0,
	)

	return syscall.Handle(hwnd)
}

func main() {
	hwnd := createWindow()

	debug := true
	w := webview.NewWindow(debug, unsafe.Pointer(hwnd))
	defer w.Destroy()

	w.SetTitle("Embedded WebView")
	w.SetSize(800, 600, webview.HintNone)

	w.SetHtml("<div><div>Input below</div><div><input id='input' /></div></div>")
	w.Eval(`document.getElementById('input').focus();`)
	w.Run()
}
