package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	FontSize            int     `json:"fontSize"`
	FontFamily          string  `json:"fontFamily"`
	FontWeight          int     `json:"fontWeight"`
	FontColor           string  `json:"fontColor"`
	KaraokeColor        string  `json:"karaokeColor"`
	StrokeColor         string  `json:"strokeColor"`
	BgOpacity           float64 `json:"bgOpacity"`
	TextOpacity         float64 `json:"textOpacity"`
	KaraokeMode         string  `json:"karaokeMode"`
	LyricEffect         string  `json:"lyricEffect"`
	EffectStrength      float64 `json:"effectStrength"`
	EffectSpeed         float64 `json:"effectSpeed"`
	TextAlign           string  `json:"textAlign"`
	LetterSpacing       int     `json:"letterSpacing"`
	GlowIntensity       int     `json:"glowIntensity"`
	BorderRadius        int     `json:"borderRadius"`
	BackdropType        string  `json:"backdropType"`
	ShowNextLine        bool    `json:"showNextLine"`
	IsLocked            bool    `json:"isLocked"`
	WindowWidth         int     `json:"windowWidth"`
	WindowX             int     `json:"windowX"`
	WindowY             int     `json:"windowY"`
	WindowPositionSaved bool    `json:"windowPositionSaved"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		FontSize:            48,
		FontFamily:          "Segoe UI",
		FontWeight:          700,
		FontColor:           "#FFFFFF",
		KaraokeColor:        "#FFD500",
		StrokeColor:         "#000000",
		BgOpacity:           0.6,
		TextOpacity:         1.0,
		KaraokeMode:         "sweep",
		LyricEffect:         "glow",
		EffectStrength:      1.0,
		EffectSpeed:         1.0,
		TextAlign:           "center",
		LetterSpacing:       0,
		GlowIntensity:       10,
		BorderRadius:        12,
		BackdropType:        "none",
		ShowNextLine:        true,
		WindowWidth:         800,
		WindowX:             -1,
		WindowY:             -1,
		WindowPositionSaved: false,
	}
}

func GetConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "config.json"
	}
	appDir := filepath.Join(configDir, "AILrc")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v", err)
		return "config.json"
	}
	return filepath.Join(appDir, "config.json")
}

func LoadAppConfig() AppConfig {
	config := DefaultConfig()

	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			log.Printf("Failed to parse config file: %v, using defaults", err)
		}
	}

	if config.FontFamily == "" {
		config.FontFamily = "Segoe UI"
	}
	if config.FontWeight < 300 || config.FontWeight > 900 {
		config.FontWeight = 700
	}
	if config.KaraokeColor == "" {
		config.KaraokeColor = "#FFD500"
	}
	if config.KaraokeMode == "" {
		config.KaraokeMode = "sweep"
	}
	if config.LyricEffect == "" {
		config.LyricEffect = "glow"
	}
	if config.EffectStrength < 0.2 || config.EffectStrength > 1.5 {
		config.EffectStrength = 1.0
	}
	if config.EffectSpeed < 0.4 || config.EffectSpeed > 2.5 {
		config.EffectSpeed = 1.0
	}
	if config.TextAlign != "left" && config.TextAlign != "right" {
		config.TextAlign = "center"
	}
	if config.LetterSpacing < -4 || config.LetterSpacing > 12 {
		config.LetterSpacing = 0
	}
	if config.GlowIntensity < 0 || config.GlowIntensity > 32 {
		config.GlowIntensity = 10
	}
	if config.BorderRadius < 0 || config.BorderRadius > 24 {
		config.BorderRadius = 12
	}
	if config.BackdropType != "acrylic" {
		config.BackdropType = "none"
	}
	if config.WindowWidth < 280 || config.WindowWidth > 3840 {
		config.WindowWidth = 800
	}

	return config
}

func SaveAppConfig(config AppConfig) error {
	path := GetConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SaveLockState(locked bool) {
	config := LoadAppConfig()
	if config.IsLocked == locked {
		return
	}
	config.IsLocked = locked
	_ = SaveAppConfig(config)
}
