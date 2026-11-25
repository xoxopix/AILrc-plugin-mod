package main

import (
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	WM_APP          = 0x8000
	WM_TRAYICON     = WM_APP + 1
	WM_COMMAND      = 0x0111
	WM_LBUTTONUP    = 0x0202
	WM_LBUTTONDBL   = 0x0203
	WM_RBUTTONUP    = 0x0205
	WM_CONTEXTMENU  = 0x007B

	NIM_ADD    = 0x00000000
	NIM_DELETE = 0x00000002

	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004

	MF_STRING    = 0x00000000
	MF_SEPARATOR = 0x00000800
	TPM_RIGHTBUTTON = 0x0002
	TPM_BOTTOMALIGN = 0x0020

	IDM_SHOW_HIDE   = 1001
	IDM_TOGGLE_LOCK = 1002
	IDM_SETTINGS    = 1003
	IDM_EXIT        = 1004
)

type NOTIFYICONDATAW struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

var (
	shell32              = syscall.NewLazyDLL("shell32.dll")
	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")
	procCreatePopupMenu  = user32.NewProc("CreatePopupMenu")
	procAppendMenuW      = user32.NewProc("AppendMenuW")
	procTrackPopupMenu   = user32.NewProc("TrackPopupMenu")
	procDestroyMenu      = user32.NewProc("DestroyMenu")

	trayHwnd uintptr
	trayApp  *App
)

func initTray(hwnd uintptr, app *App) {
	if hwnd == 0 {
		return
	}
	trayHwnd = hwnd
	trayApp = app

	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = hwnd
	nid.UID = 1
	nid.UFlags = NIF_MESSAGE | NIF_ICON | NIF_TIP
	nid.UCallbackMessage = WM_TRAYICON
	nid.HIcon = getWindowIcon(hwnd)

	tip, _ := syscall.UTF16FromString("AILrc - Desktop Lyrics")
	copy(nid.SzTip[:], tip)

	procShellNotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))
}

func removeTray() {
	if trayHwnd == 0 {
		return
	}
	var nid NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = trayHwnd
	nid.UID = 1
	procShellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
	trayHwnd = 0
}

func appendTrayItem(menu uintptr, flags uint32, id uintptr, label string) {
	var ptr uintptr
	if label != "" {
		utf16Label, _ := syscall.UTF16PtrFromString(label)
		ptr = uintptr(unsafe.Pointer(utf16Label))
	}
	procAppendMenuW.Call(menu, uintptr(flags), id, ptr)
}

func showTrayMenu(hwnd uintptr) {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)

	showLabel := "Show Lyrics"
	if isWindowVisible(hwnd) {
		showLabel = "Hide Lyrics"
	}

	isLocked := (getWindowLongPtr(hwnd, GWL_EXSTYLE) & WS_EX_TRANSPARENT) != 0
	lockLabel := "Lock Window"
	if isLocked {
		lockLabel = "Unlock Window"
	}

	appendTrayItem(menu, MF_STRING, IDM_SHOW_HIDE, showLabel)
	appendTrayItem(menu, MF_STRING, IDM_TOGGLE_LOCK, lockLabel)
	appendTrayItem(menu, MF_STRING, IDM_SETTINGS, "Settings")
	appendTrayItem(menu, MF_SEPARATOR, 0, "")
	appendTrayItem(menu, MF_STRING, IDM_EXIT, "Quit")

	var pt POINT
	getCursorPos(&pt)
	setForegroundWindow(hwnd)
	procTrackPopupMenu.Call(menu, TPM_RIGHTBUTTON|TPM_BOTTOMALIGN, uintptr(pt.X), uintptr(pt.Y), 0, hwnd, 0)
}

func handleTrayMessage(hwnd uintptr, lParam uintptr) {
	switch uint32(lParam) {
	case WM_LBUTTONUP, WM_LBUTTONDBL:
		if isWindowVisible(hwnd) {
			showWindowCmd(hwnd, SW_HIDE)
		} else {
			showWindowCmd(hwnd, SW_SHOWNA)
			hideFromTaskbar(hwnd)
			setForegroundWindow(hwnd)
		}
	case WM_RBUTTONUP, WM_CONTEXTMENU:
		showTrayMenu(hwnd)
	}
}

func handleTrayCommand(hwnd uintptr, cmdID uint16) bool {
	switch cmdID {
	case IDM_SHOW_HIDE:
		if isWindowVisible(hwnd) {
			showWindowCmd(hwnd, SW_HIDE)
		} else {
			showWindowCmd(hwnd, SW_SHOWNA)
			hideFromTaskbar(hwnd)
			setForegroundWindow(hwnd)
		}
		return true
	case IDM_TOGGLE_LOCK:
		isLocked := (getWindowLongPtr(hwnd, GWL_EXSTYLE) & WS_EX_TRANSPARENT) != 0
		nextLocked := !isLocked
		SetClickThrough(nextLocked, nil)
		SaveLockState(nextLocked)
		if trayApp != nil && trayApp.ctx != nil {
			go runtime.EventsEmit(trayApp.ctx, "lock_state_change", nextLocked)
		}
		return true
	case IDM_SETTINGS:
		if !isWindowVisible(hwnd) {
			showWindowCmd(hwnd, SW_SHOWNA)
			hideFromTaskbar(hwnd)
		}
		setForegroundWindow(hwnd)
		if trayApp != nil && trayApp.ctx != nil {
			go runtime.EventsEmit(trayApp.ctx, "open_settings")
		}
		return true
	case IDM_EXIT:
		if trayApp != nil && trayApp.ctx != nil {
			go runtime.Quit(trayApp.ctx)
		}
		return true
	}
	return false
}
