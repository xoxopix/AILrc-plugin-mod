# AILrc (Mod)

A customized mod of [AILrc-plugin](https://github.com/Leovikii/AILrc-plugin) — a modern desktop lyric renderer for **AIMP** built with Go (Wails), Vue 3, and a native C++ AIMP plugin.

> **Credits & Attribution**
> - **Original Creator:** [Viki / LeoViki](https://github.com/Leovikii) — original author of [Leovikii/AILrc-plugin](https://github.com/Leovikii/AILrc-plugin) (MIT License, Copyright © 2026 Viki).
> - **Visual Effects Inspiration:** [ivLis-Studio/ivLyrics](https://github.com/ivLis-Studio/ivLyrics) — lyric motion and visual style effects (`Glow`, `Wave`, `Bounce`, `Float`, `Sparkle`, `Glitch`, `Echo`, `Pop`, `Pulse`, `Sway`, `Flicker`, `Whisper`, `Blur`, `Tremble`, `Adlib`) adapted for this mod.

---

## What Changed in This Mod

### Bug Fixes
- **Off-screen window recovery:** Validates coordinates with Win32 `MonitorFromPoint` and `WindowIsMinimised` so closing while minimized (`-32000, -32000`) never traps the overlay off-screen.
- **Persistent lock state (`isLocked`):** Remembers click-through lock state across restarts and settings edits.
- **Stale lyrics & encoding fixes:** Clears previous song lyrics when a new track has no lyric file, validates UTF-8 before `chardet`, fixes `EUC-JP` decoding, and strips/parses Enhanced LRC (`<mm:ss.xx>`) inline tags.
- **Cleaner resizing & layout:** Removes vertical height lock fighting cursor resizing and eliminates overlapping/unused overlay controls.

### New Features
- **Karaoke Modes:** `Off`, `Sweep` (smooth pixel-by-pixel highlight with progressive glow/sparkle/echo/glitch), and `Word Pop` (word-by-word highlight and motion).
- **15 `ivLyrics`-Inspired Effects:** `Glow`, `Wave`, `Bounce`, `Float`, `Sparkle`, `Glitch`, `Echo`, `Pop`, `Pulse`, `Sway`, `Flicker`, `Whisper`, `Blur`, `Tremble`, and `Adlib` with configurable **Intensity** and **Speed** (respecting `prefers-reduced-motion`).
- **System Font Search:** Enumerates installed Windows fonts via `EnumFontFamiliesExW` with live search, weight (`400–900`), alignment, letter spacing, glow blur, and colors.
- **Windows 11 Acrylic & Corner Radius:** Live `None` / `Acrylic` backdrop (stays blurred when unfocused via `SetWindowCompositionAttribute`) and corner rounding.
- **System Tray + AIMP Menu Integration:** Hides from the Windows taskbar (`ITaskbarList::DeleteTab`), adds a System Tray icon/menu, opens **Settings** centered on screen (restoring original position on close), and adds **`AILrc Settings`** to AIMP's main Functions menu.

---

## Tech Stack
- **Frontend:** Vue 3, TypeScript, TailwindCSS v4
- **Backend:** Go (Wails v2)
- **Plugin:** C++ using the AIMP SDK (`WM_COPYDATA` IPC)

## Build from Source

### C++ AIMP Plugin (`AILrc_plugin.dll`)
```bash
cd plugin
cmake -S . -B build -G "Visual Studio 17 2022" -A x64
cmake --build build --config Release
```

### Go / Wails Renderer (`AILrc.exe`)
```bash
cd AILrc
wails build -clean -o AILrc.exe
```

Place both `AILrc_plugin.dll` and `AILrc.exe` inside your `AIMP\Plugins\AILrc_plugin\` folder and restart AIMP.

## License
MIT License — see [LICENSE](LICENSE). Original work Copyright (c) 2026 Viki.
