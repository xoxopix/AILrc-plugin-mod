package main

import (
	"math"
	"sort"
	"syscall"
	"unsafe"
)

const (
	GWL_EXSTYLE       = -20
	WS_EX_LAYERED     = 0x80000
	WS_EX_TRANSPARENT = 0x20
	MB_ICONERROR      = 0x10

	SW_HIDE   = 0
	SW_SHOWNA = 8

	WM_GETICON      = 0x007F
	ICON_SMALL      = 0
	IDI_APPLICATION = 32512

	DWMWA_USE_IMMERSIVE_DARK_MODE  = 20
	DWMWA_WINDOW_CORNER_PREFERENCE = 33
	DWMWA_SYSTEMBACKDROP_TYPE      = 38

	DWMWCP_DONOTROUND = 1
	DWMWCP_ROUND      = 2
	DWMWCP_ROUNDSMALL = 3

	WCA_ACCENT_POLICY               = 19
	ACCENT_DISABLED                 = 0
	ACCENT_ENABLE_ACRYLICBLURBEHIND = 4

	COINIT_APARTMENTTHREADED = 0x2
	CLSCTX_INPROC_SERVER     = 0x1

	DEFAULT_CHARSET = 1

	// AIMP Remote API Constants
	WM_USER                       = 0x0400
	WM_AIMP_PROPERTY              = WM_USER + 0x77
	AIMP_RA_PROPERTY_PLAYER_STATE = 0x40
	AIMP_RA_PROPVALUE_GET         = 0
)

type POINT struct {
	X int32
	Y int32
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type ACCENT_POLICY struct {
	AccentState   uint32
	AccentFlags   uint32
	GradientColor uint32
	AnimationId   uint32
}

type WINDOWCOMPOSITIONATTRIBDATA struct {
	Attrib uint32
	PvData unsafe.Pointer
	CbData uint32
}

type guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

type LOGFONTW struct {
	LfHeight         int32
	LfWidth          int32
	LfEscapement     int32
	LfOrientation    int32
	LfWeight         int32
	LfItalic         byte
	LfUnderline      byte
	LfStrikeOut      byte
	LfCharSet        byte
	LfOutPrecision   byte
	LfClipPrecision  byte
	LfQuality        byte
	LfPitchAndFamily byte
	LfFaceName       [32]uint16
}

var (
	user32                           = syscall.NewLazyDLL("user32.dll")
	gdi32                            = syscall.NewLazyDLL("gdi32.dll")
	ole32                            = syscall.NewLazyDLL("ole32.dll")
	dwmapi                           = syscall.NewLazyDLL("dwmapi.dll")
	procFindWindowW                  = user32.NewProc("FindWindowW")
	procSetWindowLongPtrW            = user32.NewProc("SetWindowLongPtrW")
	procGetWindowLongPtrW            = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongW               = user32.NewProc("SetWindowLongW")
	procGetWindowLongW               = user32.NewProc("GetWindowLongW")
	procCallWindowProcW              = user32.NewProc("CallWindowProcW")
	procMessageBoxW                  = user32.NewProc("MessageBoxW")
	procGetCursorPos                 = user32.NewProc("GetCursorPos")
	procScreenToClient               = user32.NewProc("ScreenToClient")
	procSendMessageW                 = user32.NewProc("SendMessageW")
	procMonitorFromPoint              = user32.NewProc("MonitorFromPoint")
	procShowWindow                    = user32.NewProc("ShowWindow")
	procIsWindowVisible               = user32.NewProc("IsWindowVisible")
	procSetForegroundWindow           = user32.NewProc("SetForegroundWindow")
	procLoadIconW                     = user32.NewProc("LoadIconW")
	procGetDC                         = user32.NewProc("GetDC")
	procReleaseDC                     = user32.NewProc("ReleaseDC")
	procSetWindowCompositionAttribute = user32.NewProc("SetWindowCompositionAttribute")
	procEnumFontFamiliesExW           = gdi32.NewProc("EnumFontFamiliesExW")
	procCoInitializeEx                = ole32.NewProc("CoInitializeEx")
	procCoUninitialize                = ole32.NewProc("CoUninitialize")
	procCoCreateInstance              = ole32.NewProc("CoCreateInstance")
	procDwmSetWindowAttribute         = dwmapi.NewProc("DwmSetWindowAttribute")

	clsidTaskbarList = guid{0x56FDF344, 0xFD6D, 0x11D0, [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90}}
	iidTaskbarList   = guid{0x56FDF342, 0xFD6D, 0x11D0, [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90}}
)

const monitorDefaultToNull = 0

// isPointOnScreen checks via MonitorFromPoint whether (x, y) falls on any active monitor,
// rejecting minimized sentinel coordinates like (-32000, -32000).
func isPointOnScreen(x, y int) bool {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		ret, _, _ := procMonitorFromPoint.Call(uintptr(int32(x)), uintptr(int32(y)), monitorDefaultToNull)
		return ret != 0
	}
	// Win64 ABI packs 8-byte POINT{int32 X, int32 Y} into a single uintptr argument.
	pt := uint64(uint32(int32(x))) | (uint64(uint32(int32(y))) << 32)
	ret, _, _ := procMonitorFromPoint.Call(uintptr(pt), monitorDefaultToNull)
	return ret != 0
}

func getWindowHandle() uintptr {
	title, _ := syscall.UTF16PtrFromString("AILrc")
	ret, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(title)))
	return ret
}

func setWindowLongPtr(hwnd uintptr, index int, newLong uintptr) uintptr {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		ret, _, _ := procSetWindowLongW.Call(hwnd, uintptr(index), newLong)
		return ret
	}
	ret, _, _ := procSetWindowLongPtrW.Call(hwnd, uintptr(index), newLong)
	return ret
}

func getWindowLongPtr(hwnd uintptr, index int) uintptr {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		ret, _, _ := procGetWindowLongW.Call(hwnd, uintptr(index))
		return ret
	}
	ret, _, _ := procGetWindowLongPtrW.Call(hwnd, uintptr(index))
	return ret
}

func callWindowProc(prevWndProc uintptr, hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procCallWindowProcW.Call(prevWndProc, hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func showErrorMessageBox(title, message string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), MB_ICONERROR)
}

func getCursorPos(pt *POINT) bool {
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(pt)))
	return ret != 0
}

func screenToClient(hwnd uintptr, pt *POINT) bool {
	ret, _, _ := procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(pt)))
	return ret != 0
}

func showWindowCmd(hwnd uintptr, cmd int) {
	if hwnd != 0 {
		procShowWindow.Call(hwnd, uintptr(cmd))
	}
}

func isWindowVisible(hwnd uintptr) bool {
	if hwnd == 0 {
		return false
	}
	ret, _, _ := procIsWindowVisible.Call(hwnd)
	return ret != 0
}

func setForegroundWindow(hwnd uintptr) {
	if hwnd != 0 {
		procSetForegroundWindow.Call(hwnd)
	}
}

// hideFromTaskbar removes the window's taskbar button via COM ITaskbarList::DeleteTab
// without setting WS_EX_TOOLWINDOW (which causes Windows to draw a native toolwindow frame/close button).
func hideFromTaskbar(hwnd uintptr) bool {
	hrInit, _, _ := procCoInitializeEx.Call(0, COINIT_APARTMENTTHREADED)
	if hrInit == 0 || hrInit == 1 {
		defer procCoUninitialize.Call()
	}

	var tbl uintptr
	hr, _, _ := procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidTaskbarList)),
		0,
		CLSCTX_INPROC_SERVER,
		uintptr(unsafe.Pointer(&iidTaskbarList)),
		uintptr(unsafe.Pointer(&tbl)),
	)
	if hr != 0 || tbl == 0 {
		return false
	}

	vtbl := *(*[8]uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(tbl))))
	defer syscall.SyscallN(vtbl[2], tbl) // Release

	if hwnd == 0 {
		return true
	}
	if r, _, _ := syscall.SyscallN(vtbl[3], tbl); r == 0 { // HrInit
		syscall.SyscallN(vtbl[5], tbl, hwnd) // DeleteTab
	}
	return true
}

func getWindowIcon(hwnd uintptr) uintptr {
	if hwnd != 0 {
		ico, _, _ := procSendMessageW.Call(hwnd, WM_GETICON, ICON_SMALL, 0)
		if ico != 0 {
			return ico
		}
	}
	ico, _, _ := procLoadIconW.Call(0, uintptr(IDI_APPLICATION))
	return ico
}

func applyWindowStyle(hwnd uintptr, backdropMode string, cornerRadius int, bgOpacity float64) {
	if hwnd == 0 {
		return
	}

	var dark int32 = 1
	procDwmSetWindowAttribute.Call(hwnd, DWMWA_USE_IMMERSIVE_DARK_MODE, uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))

	var disableBackdrop int32 = 1
	procDwmSetWindowAttribute.Call(hwnd, DWMWA_SYSTEMBACKDROP_TYPE, uintptr(unsafe.Pointer(&disableBackdrop)), unsafe.Sizeof(disableBackdrop))

	var corner int32 = DWMWCP_DONOTROUND
	if backdropMode == "acrylic" {
		if cornerRadius <= 0 {
			corner = DWMWCP_DONOTROUND
		} else if cornerRadius <= 6 {
			corner = DWMWCP_ROUNDSMALL
		} else {
			corner = DWMWCP_ROUND
		}
	}
	procDwmSetWindowAttribute.Call(hwnd, DWMWA_WINDOW_CORNER_PREFERENCE, uintptr(unsafe.Pointer(&corner)), unsafe.Sizeof(corner))

	// ACCENT_ENABLE_ACRYLICBLURBEHIND keeps blur active even when the lyric overlay is unfocused.
	if procSetWindowCompositionAttribute.Find() == nil {
		policy := ACCENT_POLICY{AccentState: ACCENT_DISABLED}
		if backdropMode == "acrylic" {
			alpha := uint32(math.Round(bgOpacity * 210))
			if alpha < 1 {
				alpha = 1
			}
			if alpha > 235 {
				alpha = 235
			}
			policy = ACCENT_POLICY{
				AccentState:   ACCENT_ENABLE_ACRYLICBLURBEHIND,
				AccentFlags:   2,
				GradientColor: (alpha << 24) | 0x00140E0E, // ABGR dark tint
			}
		}
		data := WINDOWCOMPOSITIONATTRIBDATA{
			Attrib: WCA_ACCENT_POLICY,
			PvData: unsafe.Pointer(&policy),
			CbData: uint32(unsafe.Sizeof(policy)),
		}
		procSetWindowCompositionAttribute.Call(hwnd, uintptr(unsafe.Pointer(&data)))
	}
}

func enumInstalledFonts() []string {
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return nil
	}
	defer procReleaseDC.Call(0, hdc)

	seen := make(map[string]struct{})
	cb := syscall.NewCallback(func(lpelfe, lpntme uintptr, fontType uint32, lParam uintptr) uintptr {
		lf := (*LOGFONTW)(unsafe.Pointer(lpelfe))
		name := syscall.UTF16ToString(lf.LfFaceName[:])
		if name != "" && name[0] != '@' {
			seen[name] = struct{}{}
		}
		return 1
	})

	var lf LOGFONTW
	lf.LfCharSet = DEFAULT_CHARSET
	procEnumFontFamiliesExW.Call(hdc, uintptr(unsafe.Pointer(&lf)), cb, 0, 0)

	fonts := make([]string, 0, len(seen))
	for name := range seen {
		fonts = append(fonts, name)
	}
	sort.Strings(fonts)
	return fonts
}

var aimpClassNamePtr *uint16

func init() {
	aimpClassNamePtr, _ = syscall.UTF16PtrFromString("AIMP2_RemoteInfo")
}

func getAimpState() int {
	hwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(aimpClassNamePtr)), 0)
	if hwnd != 0 {
		ret, _, _ := procSendMessageW.Call(hwnd, uintptr(WM_AIMP_PROPERTY), uintptr(AIMP_RA_PROPERTY_PLAYER_STATE|AIMP_RA_PROPVALUE_GET), 0)
		return int(ret)
	}
	return -1
}
