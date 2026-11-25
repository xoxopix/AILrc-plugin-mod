<script setup lang="ts">
import { ref, computed } from 'vue';
import { hexToRgba } from "../utils/color";
import type { AppConfig, LyricLine } from "../types";

interface Props {
    mainText: string;
    subText: string;
    currentLine: LyricLine | null;
    nextLine: LyricLine | null;
    position?: number;
    config: AppConfig;
}

interface Token {
    text: string;
    start: number;
    end: number;
    progress: number;
    active: boolean;
}

const props = defineProps<Props>();
const lyricRef = ref<HTMLDivElement>();

const fontStack = computed(() => {
    const family = props.config.fontFamily?.trim() || 'Segoe UI';
    return `"${family}", "Segoe UI", "Microsoft YaHei", sans-serif`;
});

const baseTextColor = computed(() => hexToRgba(props.config.fontColor || '#FFFFFF', props.config.textOpacity));
const inactiveTextColor = computed(() => hexToRgba(props.config.fontColor || '#FFFFFF', props.config.textOpacity * 0.42));
const karaokeTextColor = computed(() => hexToRgba(props.config.karaokeColor || '#FFD500', props.config.textOpacity));
const glowColor = computed(() => hexToRgba(props.config.strokeColor || '#000000', props.config.textOpacity * 0.65));
const deepShadow = computed(() => hexToRgba(props.config.strokeColor || '#000000', props.config.textOpacity * 0.92));

const effectGlowTint = computed(() =>
    hexToRgba(
        props.config.karaokeMode === 'off' ? (props.config.fontColor || '#FFFFFF') : (props.config.karaokeColor || '#FFD500'),
        0.75
    )
);

const rootStyle = computed(() => ({
    '--fx-strength': String(props.config.effectStrength ?? 1),
    '--fx-speed': String(props.config.effectSpeed ?? 1),
    '--fx-glow-color': effectGlowTint.value,
    textAlign: props.config.textAlign || 'center',
}));

const effectClass = computed(() => `fx-${props.config.lyricEffect || 'none'}`);

const baseShadow = computed(() => {
    const blur = props.config.glowIntensity ?? 10;
    return `0 0 1px ${deepShadow.value}, 0 2px ${Math.max(2, Math.round(blur * 0.45))}px ${glowColor.value}, 0 0 ${blur}px ${glowColor.value}`;
});

const alignClass = computed(() => {
    switch (props.config.textAlign) {
        case 'left': return 'text-left items-start';
        case 'right': return 'text-right items-end';
        default: return 'text-center items-center';
    }
});

const mainTextStyle = computed(() => ({
    fontSize: `${props.config.fontSize}px`,
    fontWeight: props.config.fontWeight || 700,
    fontFamily: fontStack.value,
    letterSpacing: `${props.config.letterSpacing ?? 0}px`,
    lineHeight: '1.22',
    width: '100%',
    textShadow: baseShadow.value,
    color: props.config.karaokeMode === 'off' ? baseTextColor.value : inactiveTextColor.value,
}));

const subTextStyle = computed(() => ({
    fontSize: `${Math.max(13, Math.round(props.config.fontSize * 0.54))}px`,
    fontWeight: Math.max(400, (props.config.fontWeight || 700) - 100),
    color: baseTextColor.value,
    textShadow: `0 1px 2px ${deepShadow.value}`,
    fontFamily: fontStack.value,
    letterSpacing: `${Math.max(-1, (props.config.letterSpacing ?? 0) * 0.5)}px`,
    lineHeight: '1.3',
    width: '100%',
    opacity: 0.85,
}));

const nextLineStyle = computed(() => ({
    fontSize: `${Math.max(12, Math.round(props.config.fontSize * 0.5))}px`,
    fontWeight: 500,
    color: hexToRgba(props.config.fontColor || '#FFFFFF', props.config.textOpacity * 0.45),
    textShadow: `0 1px 2px ${deepShadow.value}`,
    fontFamily: fontStack.value,
    lineHeight: '1.25',
    width: '100%',
}));

const tokens = computed<Token[]>(() => {
    const text = props.mainText || '';
    if (!text) return [];

    const pos = props.position ?? -1;
    const line = props.currentLine;

    if (!line || props.config.karaokeMode === 'off') {
        const rawParts = text.match(/\S+\s*/g) || [text];
        return rawParts.map((part) => ({
            text: part,
            start: 0,
            end: 0,
            progress: 1,
            active: true,
        }));
    }

    if (line.words && line.words.length > 0) {
        return line.words.map((w, i) => {
            const start = w.time;
            const nextStart = i + 1 < line.words!.length ? line.words![i + 1].time : line.endTime;
            const end = Math.max(start + 120, nextStart);
            const progress = pos < 0 ? 0 : Math.min(1, Math.max(0, (pos - start) / (end - start)));
            return {
                text: w.text,
                start,
                end,
                progress,
                active: pos >= start,
            };
        });
    }

    const rawParts = /\s/.test(text)
        ? (text.match(/\S+\s*/g) || [text])
        : Array.from(text);

    const lineStart = line.time;
    const rawDuration = line.endTime > lineStart ? line.endTime - lineStart : 3500;
    const duration = Math.min(7000, Math.max(800, rawDuration));
    const totalChars = Math.max(1, rawParts.reduce((sum, p) => sum + Math.max(1, p.trim().length), 0));

    let cursor = lineStart;
    return rawParts.map((part) => {
        const weight = Math.max(1, part.trim().length) / totalChars;
        const sliceMs = Math.max(80, duration * weight);
        const start = cursor;
        const end = cursor + sliceMs;
        cursor = end;
        const progress = pos < 0 ? 0 : Math.min(1, Math.max(0, (pos - start) / (end - start)));
        return {
            text: part,
            start,
            end,
            progress,
            active: pos >= start,
        };
    });
});

defineExpose({ lyricRef });
</script>

<template>
    <div
        ref="lyricRef"
        :class="['pointer-events-none flex flex-col justify-center p-1.5 w-full h-auto', alignClass, effectClass]"
        :style="rootStyle"
    >
        <div
            class="fx-active-line leading-tight wrap-break-word whitespace-pre-wrap"
            :style="mainTextStyle"
        >
            <template v-if="config.karaokeMode === 'sweep' && currentLine">
                <span
                    v-for="(tok, idx) in tokens"
                    :key="`${currentLine.time}-${idx}`"
                    :class="['fx-tok relative', tok.active ? 'fx-sung-motion' : '']"
                    :style="{ '--tok-idx': idx }"
                >
                    <span>{{ tok.text }}</span>
                    <span
                        v-if="tok.progress > 0.001"
                        class="absolute inset-0 pointer-events-none fx-sung-glow"
                        :style="{
                            color: karaokeTextColor,
                            textShadow: baseShadow,
                            clipPath: tok.progress < 0.999 ? `inset(-1.5em ${((1 - tok.progress) * 100).toFixed(1)}% -1.5em -1.5em)` : undefined,
                        }"
                    >{{ tok.text }}</span>
                </span>
            </template>
            <template v-else-if="config.karaokeMode === 'word' && currentLine">
                <span
                    v-for="(tok, idx) in tokens"
                    :key="`${currentLine.time}-${idx}`"
                    :class="['fx-tok', tok.active ? 'fx-sung-motion fx-sung-glow' : '']"
                    :style="{ '--tok-idx': idx }"
                >
                    <span
                        class="inline-block fx-tok-word"
                        :style="{
                            color: tok.active ? karaokeTextColor : inactiveTextColor,
                            transform: tok.active && tok.progress < 0.85 ? 'translate3d(0, -1px, 0) scale(1.05)' : 'none',
                        }"
                    >{{ tok.text }}</span>
                </span>
            </template>
            <template v-else>
                <span
                    v-for="(tok, idx) in tokens"
                    :key="idx"
                    class="fx-tok fx-sung-motion fx-sung-glow"
                    :style="{ '--tok-idx': idx }"
                >{{ tok.text }}</span>
            </template>
        </div>

        <div
            v-if="subText"
            class="leading-tight wrap-break-word whitespace-pre-wrap mt-1"
            :style="subTextStyle"
        >
            {{ subText }}
        </div>

        <div
            v-if="config.showNextLine && nextLine?.mainText"
            class="truncate max-w-full mt-1"
            :style="nextLineStyle"
        >
            {{ nextLine.mainText }}
        </div>
    </div>
</template>
