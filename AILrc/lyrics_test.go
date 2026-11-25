package main

import (
	"testing"
	"unsafe"
)

func TestLyricsAndScreenSafety(t *testing.T) {
	if unsafe.Sizeof(NOTIFYICONDATAW{}) != 976 {
		t.Fatalf("expected NOTIFYICONDATAW size 976, got %d", unsafe.Sizeof(NOTIFYICONDATAW{}))
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
}
