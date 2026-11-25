package main

import (
	"syscall"
	"testing"
	"unsafe"
)

func TestLyricsAndScreenSafety(t *testing.T) {
	expectedNidSize := uintptr(976)
	if unsafe.Sizeof(uintptr(0)) == 4 {
		expectedNidSize = 956
	}
	if unsafe.Sizeof(NOTIFYICONDATAW{}) != expectedNidSize {
		t.Fatalf("expected NOTIFYICONDATAW size %d, got %d", expectedNidSize, unsafe.Sizeof(NOTIFYICONDATAW{}))
	}
	lrc := "\xef\xbb\xbf[00:01.50]<00:01.50>Hello <00:02.00>world\n[00:04.00]Next line\n"
	lines := parseLrc(lrc)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].MainText != "Hello world" {
		t.Fatalf("expected stripped text 'Hello world', got %q", lines[0].MainText)
	}
	if len(lines[0].Words) != 2 || lines[0].Words[1].Time != 2000 {
		t.Fatalf("unexpected enhanced words: %+v", lines[0].Words)
	}
	if lines[0].EndTime != 4000 {
		t.Fatalf("expected EndTime 4000, got %d", lines[0].EndTime)
	}
	if isPointOnScreen(-32000, -32000) {
		t.Fatalf("minimized sentinel (-32000,-32000) must not be treated as on-screen")
	}
	// Verify ITaskbarList COM initialization + DeleteTab(0) and system font enumeration.
	if !hideFromTaskbar(0) {
		t.Fatalf("expected ITaskbarList COM initialization to succeed")
	}
	if fonts := enumInstalledFonts(); len(fonts) == 0 {
		t.Fatalf("expected system fonts to be enumerated")
	}
	h1, ok1 := acquireSingleInstanceMutex("Local\\AILrc_UnitTest_Mutex_Check")
	if !ok1 || h1 == 0 {
		t.Fatalf("expected first mutex acquisition to succeed")
	}
	_, ok2 := acquireSingleInstanceMutex("Local\\AILrc_UnitTest_Mutex_Check")
	if ok2 {
		t.Fatalf("expected second mutex acquisition to fail with ERROR_ALREADY_EXISTS")
	}
	_ = syscall.CloseHandle(syscall.Handle(h1))
}
