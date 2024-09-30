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

type IncrementResult struct {
	Count uint `json:"count"`
}

func main() {
	var count uint = 0

	hwnd := createWindow()

	debug := true
	w := webview.NewWindow(debug, unsafe.Pointer(hwnd))
	defer w.Destroy()

	w.SetTitle("Embedded WebView")
	w.SetSize(800, 600, webview.HintNone)

	w.Bind("increment", func() IncrementResult {
		count++
		return IncrementResult{Count: count}
	})

	w.SetHtml(html)
	w.Run()
}

const html = `
<div>
	<button id="increment">Tap me</button>
	<div>You tapped <span id="count">0</span> time(s).</div>
	<input id="input" />
	<div>
		<label for="cars">Choose a car:</label>

		<select name="cars" id="cars">
  		<option value="volvo">Volvo</option>
  		<option value="saab">Saab</option>
  		<option value="mercedes">Mercedes</option>
  		<option value="audi">Audi</option>
		</select>
	</div>
</div>
<script>
	const [incrementElement, countElement, inputElement] =
    document.querySelectorAll("#increment, #count, #input");
  document.addEventListener("DOMContentLoaded", () => {
		inputElement.focus();

    incrementElement.addEventListener("click", () => {
      window.increment().then(result => {
        countElement.textContent = result.count;
      });
    });
  });
</script>
`
