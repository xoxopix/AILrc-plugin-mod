package main

import (
	"context"
	"path/filepath"
	"sync/atomic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MusicInfo struct {
	FileName    string `json:"FileName"`
	TrackNumber int    `json:"TrackNumber"`
}

type PlayerState struct {
	Position int `json:"Position"`
	State    int `json:"State"`
}

type App struct {
	ctx             context.Context
	inSettings      bool
	pendingSettings atomic.Bool
	savedPosX       int
	savedPosY       int
	hasSavedPos     bool
	savedLocked     bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	setupIPC(ctx, a)

	hwnd := getWindowHandle()
	hideFromTaskbar(hwnd)
	initTray(hwnd, a)

	config := LoadAppConfig()
	applyWindowStyle(hwnd, config.BackdropType, config.BorderRadius, config.BgOpacity)

	runtime.WindowSetMinSize(ctx, 280, 64)
	runtime.WindowSetMaxSize(ctx, 4000, 1000)

	if config.WindowWidth >= 280 && config.WindowWidth <= 3840 {
		runtime.WindowSetSize(ctx, config.WindowWidth, 120)
	}

	if config.WindowPositionSaved && (isPointOnScreen(config.WindowX, config.WindowY) || isPointOnScreen(config.WindowX+120, config.WindowY+40)) {
		runtime.WindowSetPosition(ctx, config.WindowX, config.WindowY)
	}

	if config.IsLocked {
		SetClickThrough(true, nil)
	}
}

func (a *App) beforeClose(ctx context.Context) bool {
	removeTray()

	hwnd := getWindowHandle()
	currentLocked := (getWindowLongPtr(hwnd, GWL_EXSTYLE) & WS_EX_TRANSPARENT) != 0
	if a.inSettings {
		currentLocked = a.savedLocked
	}

	if a.inSettings && a.hasSavedPos {
		config := LoadAppConfig()
		config.WindowX = a.savedPosX
		config.WindowY = a.savedPosY
		config.WindowPositionSaved = true
		config.IsLocked = currentLocked
		_ = SaveAppConfig(config)
		return false
	}

	config := LoadAppConfig()
	config.IsLocked = currentLocked
	if !runtime.WindowIsMinimised(ctx) {
		x, y := runtime.WindowGetPosition(ctx)
		if isPointOnScreen(x, y) || isPointOnScreen(x+120, y+40) {
			config.WindowX = x
			config.WindowY = y
			config.WindowPositionSaved = true
		}
	}
	_ = SaveAppConfig(config)

	return false
}

func (a *App) FetchMusicInfo() MusicInfo {
	internal := GetMusicInfo()
	return MusicInfo{
		FileName:    internal.FileName,
		TrackNumber: internal.TrackNumber,
	}
}

func (a *App) FetchPlayerState() PlayerState {
	state := GetRealtimeState()
	actualState := getAimpState()
	if actualState >= 0 {
		state.State = actualState
	}
	return state
}

func (a *App) FetchLyrics(filename string) []LyricLine {
	internal := GetMusicInfo()

	if filepath.Base(internal.FullPath) == filename {
		return LoadLyrics(internal.FullPath)
	}
	return LoadLyrics(filename)
}

func (a *App) SetWindowClickThrough(enabled bool, rect *DOMRect) {
	SetClickThrough(enabled, rect)
}

func (a *App) GetConfig() AppConfig {
	return LoadAppConfig()
}

func (a *App) SaveConfig(config AppConfig) {
	if a.inSettings {
		config.IsLocked = a.savedLocked
	} else {
		config.IsLocked = (getWindowLongPtr(getWindowHandle(), GWL_EXSTYLE) & WS_EX_TRANSPARENT) != 0
	}
	_ = SaveAppConfig(config)
	applyWindowStyle(getWindowHandle(), config.BackdropType, config.BorderRadius, config.BgOpacity)
}

func (a *App) ApplyWindowStyle(backdropMode string, cornerRadius int, bgOpacity float64) {
	applyWindowStyle(getWindowHandle(), backdropMode, cornerRadius, bgOpacity)
}

func (a *App) GetInstalledFonts() []string {
	return enumInstalledFonts()
}

func (a *App) ConsumePendingSettings() bool {
	return a.pendingSettings.Swap(false)
}

func (a *App) EnterSettingsView(width, height int) {
	a.pendingSettings.Store(false)
	if !a.inSettings {
		x, y := runtime.WindowGetPosition(a.ctx)
		if isPointOnScreen(x, y) || isPointOnScreen(x+120, y+40) {
			a.savedPosX = x
			a.savedPosY = y
			a.hasSavedPos = true
		}
		a.savedLocked = (getWindowLongPtr(getWindowHandle(), GWL_EXSTYLE) & WS_EX_TRANSPARENT) != 0
		a.inSettings = true
	}

	SetClickThrough(false, nil)
	runtime.WindowSetMinSize(a.ctx, width, height)
	runtime.WindowSetMaxSize(a.ctx, width, height)
	runtime.WindowSetSize(a.ctx, width, height)
	runtime.WindowCenter(a.ctx)
}

func (a *App) ExitSettingsView(width, height int) {
	a.inSettings = false

	runtime.WindowSetMinSize(a.ctx, 280, 64)
	runtime.WindowSetMaxSize(a.ctx, 4000, 1000)
	runtime.WindowSetSize(a.ctx, width, height)

	if a.hasSavedPos && (isPointOnScreen(a.savedPosX, a.savedPosY) || isPointOnScreen(a.savedPosX+120, a.savedPosY+40)) {
		runtime.WindowSetPosition(a.ctx, a.savedPosX, a.savedPosY)
	}
	if a.savedLocked {
		SetClickThrough(true, nil)
		runtime.EventsEmit(a.ctx, "lock_state_change", true)
	}
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

func (a *App) ResizeWindow(width, height int) {
	if a.inSettings {
		return
	}
	if width < 280 {
		width = 280
	}
	if height < 64 {
		height = 64
	}
	runtime.WindowSetSize(a.ctx, width, height)
}
