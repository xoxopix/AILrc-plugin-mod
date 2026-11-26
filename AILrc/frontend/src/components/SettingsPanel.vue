<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { X, Search, AlignLeft, AlignCenter, AlignRight } from 'lucide-vue-next';
import type { AppConfig, BackdropType, KaraokeMode, LyricEffect, TextAlign } from '../types';
import CustomSlider from './CustomSlider.vue';
import packageJson from '../../package.json';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { ApplyWindowStyle, GetInstalledFonts } from '../../wailsjs/go/main/App';

interface Props {
    config: AppConfig;
}

const props = defineProps<Props>();

const emit = defineEmits<{
    change: [config: AppConfig];
    close: [];
}>();

type TabId = 'fx' | 'type' | 'window';
const activeTab = ref<TabId>('fx');

const systemFonts = ref<string[]>([]);
const fontQuery = ref('');

onMounted(async () => {
    try {
        const list = await GetInstalledFonts();
        if (Array.isArray(list)) {
            systemFonts.value = list;
        }
    } catch {
        systemFonts.value = [];
    }
});

const filteredFonts = computed(() => {
    const q = fontQuery.value.trim().toLowerCase();
    const list = q
        ? systemFonts.value.filter((f) => f.toLowerCase().includes(q))
        : systemFonts.value;
    return list.slice(0, 60);
});

const karaokeModes: { label: string; value: KaraokeMode }[] = [
    { label: 'Off', value: 'off' },
    { label: 'Sweep', value: 'sweep' },
    { label: 'Word Pop', value: 'word' },
];

const lyricEffects: { label: string; value: LyricEffect }[] = [
    { label: 'None', value: 'none' },
    { label: 'Glow', value: 'glow' },
    { label: 'Wave', value: 'wave' },
    { label: 'Bounce', value: 'bounce' },
    { label: 'Float', value: 'float' },
    { label: 'Sparkle', value: 'sparkle' },
    { label: 'Glitch', value: 'glitch' },
    { label: 'Echo', value: 'echo' },
    { label: 'Pop', value: 'pop' },
    { label: 'Pulse', value: 'pulse' },
    { label: 'Sway', value: 'sway' },
    { label: 'Flicker', value: 'flicker' },
    { label: 'Whisper', value: 'whisper' },
    { label: 'Blur', value: 'blur' },
    { label: 'Tremble', value: 'tremble' },
    { label: 'Adlib', value: 'adlib' },
];

const backdropOptions: { label: string; value: BackdropType }[] = [
    { label: 'None', value: 'none' },
    { label: 'Acrylic', value: 'acrylic' },
];

const handleChange = (key: keyof AppConfig, value: any) => {
    if (typeof value === 'number' && !['fontWeight', 'letterSpacing', 'glowIntensity', 'borderRadius', 'lyricOffsetMs'].includes(key)) {
        value = parseFloat(value.toFixed(2));
    }
    const nextConfig = { ...props.config, [key]: value };
    if (key === 'backdropType' || key === 'borderRadius' || key === 'bgOpacity') {
        ApplyWindowStyle(nextConfig.backdropType, nextConfig.borderRadius, nextConfig.bgOpacity);
    }
    emit('change', nextConfig);
};

const formatLyricOffset = (ms: number) => {
    const v = Math.round(ms || 0);
    if (v === 0) return '0 ms (Default)';
    const sign = v > 0 ? '+' : '-';
    const abs = Math.abs(v);
    const val = abs >= 1000 ? `${(abs / 1000).toFixed(2)} s` : `${abs} ms`;
    return `${sign}${val} (${v > 0 ? 'ahead' : 'behind'})`;
};

const nudgeLyricOffset = (delta: number) => {
    const cur = Math.round(props.config.lyricOffsetMs ?? 0);
    handleChange('lyricOffsetMs', Math.max(-5000, Math.min(5000, cur + delta)));
};

const previewBg = computed(() => {
    if (props.config.backdropType === 'acrylic') {
        return `rgba(14, 14, 20, ${(props.config.bgOpacity * 0.72).toFixed(2)})`;
    }
    return `rgba(0, 0, 0, ${props.config.bgOpacity})`;
});
</script>

<template>
    <div class="flex flex-col h-full text-white/90 px-5 py-4 select-none">
        <!-- Header -->
        <div class="flex justify-between items-center mb-3 border-b border-white/10 pb-3 shrink-0">
            <div class="flex items-center gap-2">
                <h2 class="text-xs font-black tracking-[0.25em] text-pink-400 uppercase">Settings</h2>
                <button
                    @click="BrowserOpenURL('https://github.com/Leovikii/AILrc-plugin')"
                    class="text-[10px] font-bold px-1.5 py-0.5 rounded-sm tracking-widest text-pink-300 bg-pink-500/15 hover:bg-pink-500/25 border border-pink-500/25 transition-colors cursor-pointer"
                    :style="{ '--wails-draggable': 'none' }"
                >
                    v{{ packageJson.version }}
                </button>
            </div>
            <button @click="emit('close')" class="p-1.5 -mr-1 text-white/50 hover:text-white transition-colors cursor-pointer" :style="{ '--wails-draggable': 'none' }">
                <X :size="18" :stroke-width="2.5" />
            </button>
        </div>

        <!-- Tabs -->
        <div class="grid grid-cols-3 gap-1.5 bg-white/5 p-1 rounded-lg mb-4 shrink-0" :style="{ '--wails-draggable': 'none' }">
            <button
                v-for="tab in ([['fx', 'Karaoke & FX'], ['type', 'Typography'], ['window', 'Window']] as const)"
                :key="tab[0]"
                @click="activeTab = tab[0]"
                :class="[
                    'py-1.5 rounded-md text-[11px] font-bold transition-colors cursor-pointer',
                    activeTab === tab[0] ? 'bg-pink-600 text-white shadow-sm' : 'text-white/55 hover:text-white'
                ]"
            >
                {{ tab[1] }}
            </button>
        </div>

        <!-- Body -->
        <div class="flex-1 space-y-5 overflow-y-auto pr-1 scrollbar-hide min-h-0" :style="{ '--wails-draggable': 'none' }">
            <!-- Tab 1: Karaoke & FX -->
            <template v-if="activeTab === 'fx'">
                <div class="space-y-2">
                    <span class="text-[11px] font-bold text-white/50 uppercase tracking-wider">Karaoke Mode</span>
                    <div class="grid grid-cols-3 gap-2">
                        <button
                            v-for="mode in karaokeModes"
                            :key="mode.value"
                            @click="handleChange('karaokeMode', mode.value)"
                            :class="[
                                'py-2 rounded-lg text-[11px] font-bold border transition-colors cursor-pointer',
                                config.karaokeMode === mode.value
                                    ? 'bg-pink-600/80 border-pink-400/60 text-white'
                                    : 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white'
                            ]"
                        >
                            {{ mode.label }}
                        </button>
                    </div>
                </div>

                <div class="space-y-2">
                    <div class="flex justify-between items-center">
                        <span class="text-[11px] font-bold text-white/50 uppercase tracking-wider">ivLyrics Effects</span>
                        <span class="text-[10px] font-mono text-pink-300 uppercase">{{ config.lyricEffect }}</span>
                    </div>
                    <div class="grid grid-cols-4 gap-1.5">
                        <button
                            v-for="fx in lyricEffects"
                            :key="fx.value"
                            @click="handleChange('lyricEffect', fx.value)"
                            :class="[
                                'py-1.5 rounded-md text-[11px] font-semibold border transition-colors cursor-pointer',
                                config.lyricEffect === fx.value
                                    ? 'bg-pink-600/80 border-pink-400/60 text-white'
                                    : 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white'
                            ]"
                        >
                            {{ fx.label }}
                        </button>
                    </div>
                </div>

                <CustomSlider
                    label="Effect Intensity"
                    :value="config.effectStrength ?? 1"
                    :min="0.3"
                    :max="1.5"
                    :step="0.1"
                    :display-fn="(v) => `${Math.round(v * 100)}%`"
                    @change="(v) => handleChange('effectStrength', v)"
                />

                <CustomSlider
                    label="Effect Speed"
                    :value="config.effectSpeed ?? 1"
                    :min="0.5"
                    :max="2.0"
                    :step="0.1"
                    :display-fn="(v) => `${v.toFixed(1)}x`"
                    @change="(v) => handleChange('effectSpeed', v)"
                />

                <CustomSlider
                    label="Lyric Sync"
                    :value="config.lyricOffsetMs ?? 0"
                    :min="-5000"
                    :max="5000"
                    :step="50"
                    :centered="true"
                    :inverted="true"
                    :display-fn="formatLyricOffset"
                    @change="(v) => handleChange('lyricOffsetMs', v)"
                >
                    <template #between>
                        <div class="flex items-center justify-between px-1 text-[11px] font-bold">
                            <button
                                @click="nudgeLyricOffset(50)"
                                title="Advance +50 ms"
                                class="w-7 h-6 rounded-md bg-white/5 hover:bg-white/15 border border-white/10 text-white/80 hover:text-white transition-colors cursor-pointer flex items-center justify-center text-sm leading-none"
                            >+</button>
                            <button
                                @click="handleChange('lyricOffsetMs', 0)"
                                :class="[
                                    'px-2.5 py-0.5 rounded-md border text-[10px] uppercase tracking-wider transition-colors cursor-pointer',
                                    (config.lyricOffsetMs ?? 0) !== 0
                                        ? 'bg-pink-600/30 border-pink-400/50 text-pink-200 hover:bg-pink-600/50'
                                        : 'bg-white/5 border-white/10 text-white/45 hover:text-white/75'
                                ]"
                            >Reset</button>
                            <button
                                @click="nudgeLyricOffset(-50)"
                                title="Delay -50 ms"
                                class="w-7 h-6 rounded-md bg-white/5 hover:bg-white/15 border border-white/10 text-white/80 hover:text-white transition-colors cursor-pointer flex items-center justify-center text-sm leading-none"
                            >−</button>
                        </div>
                    </template>
                </CustomSlider>

                <div class="flex items-center justify-between py-1">
                    <span class="text-[11px] font-bold text-white/50 uppercase tracking-wider">Next Line Preview</span>
                    <button
                        @click="handleChange('showNextLine', !config.showNextLine)"
                        :class="[
                            'px-3 py-1 rounded-md text-[11px] font-bold border transition-colors cursor-pointer',
                            config.showNextLine
                                ? 'bg-pink-600/80 border-pink-400/60 text-white'
                                : 'bg-white/5 border-white/10 text-white/40'
                        ]"
                    >
                        {{ config.showNextLine ? 'ON' : 'OFF' }}
                    </button>
                </div>
            </template>

            <!-- Tab 2: Typography & Colors -->
            <template v-else-if="activeTab === 'type'">
                <div class="space-y-2">
                    <div class="flex justify-between items-center text-[11px] font-bold text-white/50 uppercase tracking-wider">
                        <span>Font Search</span>
                        <span class="text-[10px] font-mono text-pink-300 truncate max-w-44">{{ config.fontFamily }}</span>
                    </div>
                    <div class="relative">
                        <Search :size="14" class="absolute left-3 top-1/2 -translate-y-1/2 text-white/35 pointer-events-none" />
                        <input
                            v-model="fontQuery"
                            type="text"
                            placeholder="Search installed system fonts..."
                            class="w-full bg-white/5 border border-white/10 rounded-lg pl-8 pr-3 py-1.5 text-xs text-white/90 placeholder-white/30 outline-none focus:border-pink-500/60"
                        />
                    </div>
                    <div class="max-h-32 overflow-y-auto rounded-lg border border-white/10 bg-black/25 divide-y divide-white/5">
                        <button
                            v-for="font in filteredFonts"
                            :key="font"
                            @click="handleChange('fontFamily', font)"
                            :class="[
                                'w-full text-left px-3 py-1.5 text-xs transition-colors cursor-pointer flex items-center justify-between',
                                config.fontFamily === font ? 'bg-pink-600/75 text-white' : 'text-white/75 hover:bg-white/10'
                            ]"
                            :style="{ fontFamily: font }"
                        >
                            <span class="truncate">{{ font }}</span>
                            <span v-if="config.fontFamily === font" class="text-[10px] font-mono uppercase">Active</span>
                        </button>
                        <div v-if="filteredFonts.length === 0" class="px-3 py-2 text-[11px] text-white/40 text-center">
                            No matching fonts found
                        </div>
                    </div>
                </div>

                <div class="flex items-center justify-between">
                    <span class="text-[11px] font-bold text-white/50 uppercase tracking-wider">Alignment</span>
                    <div class="flex gap-1 bg-white/5 p-1 rounded-lg border border-white/10">
                        <button
                            v-for="align in (['left', 'center', 'right'] as TextAlign[])"
                            :key="align"
                            @click="handleChange('textAlign', align)"
                            :class="[
                                'p-1.5 rounded transition-colors cursor-pointer',
                                config.textAlign === align ? 'bg-pink-600 text-white' : 'text-white/50 hover:text-white'
                            ]"
                        >
                            <AlignLeft v-if="align === 'left'" :size="14" />
                            <AlignCenter v-else-if="align === 'center'" :size="14" />
                            <AlignRight v-else :size="14" />
                        </button>
                    </div>
                </div>

                <CustomSlider label="Font Size" :value="config.fontSize" :min="12" :max="120" :step="1" unit=" PT" @change="(v) => handleChange('fontSize', v)" />
                <CustomSlider label="Font Weight" :value="config.fontWeight || 700" :min="400" :max="900" :step="100" @change="(v) => handleChange('fontWeight', v)" />
                <CustomSlider label="Letter Spacing" :value="config.letterSpacing ?? 0" :min="-2" :max="8" :step="1" unit=" PX" @change="(v) => handleChange('letterSpacing', v)" />
                <CustomSlider label="Shadow / Glow Blur" :value="config.glowIntensity ?? 10" :min="0" :max="24" :step="1" unit=" PX" @change="(v) => handleChange('glowIntensity', v)" />
                <CustomSlider label="Text Opacity" :value="config.textOpacity" :min="0.1" :max="1" :step="0.05" :display-fn="(v) => `${Math.round(v * 100)}%`" @change="(v) => handleChange('textOpacity', v)" />

                <div class="grid grid-cols-3 gap-3 pt-1">
                    <div class="bg-white/5 rounded-lg p-2.5 border border-white/10 flex flex-col gap-2">
                        <label class="text-[10px] font-bold text-white/50 uppercase">Base</label>
                        <div class="relative h-7">
                            <input type="color" :value="config.fontColor" @input="(e) => handleChange('fontColor', (e.target as HTMLInputElement).value)" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" />
                            <div class="w-full h-full rounded ring-1 ring-white/15" :style="{ backgroundColor: config.fontColor }" />
                        </div>
                    </div>
                    <div class="bg-white/5 rounded-lg p-2.5 border border-white/10 flex flex-col gap-2">
                        <label class="text-[10px] font-bold text-white/50 uppercase">Karaoke</label>
                        <div class="relative h-7">
                            <input type="color" :value="config.karaokeColor" @input="(e) => handleChange('karaokeColor', (e.target as HTMLInputElement).value)" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" />
                            <div class="w-full h-full rounded ring-1 ring-white/15" :style="{ backgroundColor: config.karaokeColor }" />
                        </div>
                    </div>
                    <div class="bg-white/5 rounded-lg p-2.5 border border-white/10 flex flex-col gap-2">
                        <label class="text-[10px] font-bold text-white/50 uppercase">Shadow</label>
                        <div class="relative h-7">
                            <input type="color" :value="config.strokeColor" @input="(e) => handleChange('strokeColor', (e.target as HTMLInputElement).value)" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" />
                            <div class="w-full h-full rounded ring-1 ring-white/15" :style="{ backgroundColor: config.strokeColor }" />
                        </div>
                    </div>
                </div>
            </template>

            <!-- Tab 3: Window -->
            <template v-else>
                <div class="space-y-2">
                    <span class="text-[11px] font-bold text-white/50 uppercase tracking-wider">Background Style</span>
                    <div class="grid grid-cols-2 gap-2">
                        <button
                            v-for="bd in backdropOptions"
                            :key="bd.value"
                            @click="handleChange('backdropType', bd.value)"
                            :class="[
                                'py-2 rounded-lg text-[11px] font-bold border transition-colors cursor-pointer',
                                config.backdropType === bd.value
                                    ? 'bg-pink-600/80 border-pink-400/60 text-white'
                                    : 'bg-white/5 border-white/10 text-white/60 hover:bg-white/10 hover:text-white'
                            ]"
                        >
                            {{ bd.label }}
                        </button>
                    </div>
                </div>

                <CustomSlider
                    label="Background Opacity"
                    :value="config.bgOpacity"
                    :min="0"
                    :max="1"
                    :step="0.05"
                    :display-fn="(v) => `${Math.round(v * 100)}%`"
                    @change="(v) => handleChange('bgOpacity', v)"
                />

                <CustomSlider
                    label="Corner Radius"
                    :value="config.borderRadius ?? 12"
                    :min="0"
                    :max="24"
                    :step="2"
                    unit=" PX"
                    @change="(v) => handleChange('borderRadius', v)"
                />

                <!-- Live window box preview so opacity & corner radius changes are visible immediately -->
                <div class="p-3 rounded-lg bg-[linear-gradient(45deg,#1e1e26_25%,#262632_25%,#262632_50%,#1e1e26_50%,#1e1e26_75%,#262632_75%)] bg-[length:16px_16px] border border-white/10">
                    <div
                        class="w-full py-3 px-4 text-center text-xs font-semibold border border-white/15 transition-colors"
                        :style="{
                            backgroundColor: previewBg,
                            borderRadius: `${config.backdropType === 'acrylic' ? ((config.borderRadius ?? 12) <= 0 ? 0 : (config.borderRadius ?? 12) <= 6 ? 4 : 8) : (config.borderRadius ?? 12)}px`,
                            color: config.karaokeColor
                        }"
                    >
                        Preview ({{ config.backdropType === 'acrylic' ? 'Acrylic' : 'None' }} · {{ Math.round(config.bgOpacity * 100) }}% · {{ config.borderRadius ?? 12 }}px)
                    </div>
                </div>
            </template>
        </div>

        <!-- Footer -->
        <div class="mt-4 pt-3 border-t border-white/10 shrink-0">
            <button @click="emit('close')" class="w-full py-2.5 bg-pink-600 hover:bg-pink-500 text-white text-[11px] font-bold rounded-lg shadow-lg shadow-pink-900/40 transition-colors active:scale-[0.98] cursor-pointer" :style="{ '--wails-draggable': 'none' }">
                APPLY &amp; CLOSE
            </button>
        </div>
    </div>
</template>
