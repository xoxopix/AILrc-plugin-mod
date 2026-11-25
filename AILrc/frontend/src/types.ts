export interface MusicInfo {
    FileName: string;
    TrackNumber: number;
}

export interface PlayerState {
    Position: number;
    State: number;
}

export interface LyricWord {
    time: number;
    text: string;
}

export interface LyricLine {
    time: number;
    endTime: number;
    mainText: string;
    subText: string;
    words?: LyricWord[];
}

export type KaraokeMode = 'off' | 'sweep' | 'word';
export type TextAlign = 'left' | 'center' | 'right';
export type BackdropType = 'none' | 'acrylic';
export type LyricEffect =
    | 'none'
    | 'glow'
    | 'wave'
    | 'bounce'
    | 'float'
    | 'sparkle'
    | 'glitch'
    | 'echo'
    | 'pop'
    | 'pulse'
    | 'sway'
    | 'flicker'
    | 'whisper'
    | 'blur'
    | 'tremble'
    | 'adlib';

export interface AppConfig {
    fontSize: number;
    fontFamily: string;
    fontWeight: number;
    fontColor: string;
    karaokeColor: string;
    strokeColor: string;
    textOpacity: number;
    bgOpacity: number;
    karaokeMode: KaraokeMode;
    lyricEffect: LyricEffect;
    effectStrength: number;
    effectSpeed: number;
    textAlign: TextAlign;
    letterSpacing: number;
    glowIntensity: number;
    borderRadius: number;
    backdropType: BackdropType;
    showNextLine: boolean;
    isLocked: boolean;
    windowWidth: number;
}
