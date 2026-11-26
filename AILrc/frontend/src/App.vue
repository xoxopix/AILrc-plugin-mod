<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { SetWindowClickThrough, ResizeWindow, EnterSettingsView, ExitSettingsView, ConsumePendingSettings } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { useAimp } from './composables/useAimp';
import { useLyrics } from './composables/useLyrics';
import { useConfig } from './composables/useConfig';
import LyricRenderer from './components/LyricRenderer.vue';
import SettingsPanel from './components/SettingsPanel.vue';

const isLocked = ref(false);
const isSettingsOpen = ref(false);
const isTransitioning = ref(false);

const { config, setConfig, save: saveConfig, isLoaded } = useConfig();
const { musicInfo, playerState } = useAimp();
const { mainText, subText, currentLine, nextLine } = useLyrics(
    computed(() => musicInfo.value?.FileName),
    computed(() => playerState.value?.Position)
);

const lyricRef = ref<InstanceType<typeof LyricRenderer>>();
const resizeTimeoutRef = ref<number | null>(null);
const resizeRafId = ref<number | null>(null);

watch(isLoaded, (loaded) => {
    if (loaded && !isSettingsOpen.value) {
        isLocked.value = Boolean(config.value.isLocked);
    }
}, { immediate: true });

watch(isLocked, (locked) => {
    if (isSettingsOpen.value) {
        return;
    }
    SetWindowClickThrough(locked, { x: 0, y: 0, width: 0, height: 0 });
});

const computeTargetHeight = () => {
    const contentH = lyricRef.value?.lyricRef?.offsetHeight || 0;
    const minH = Math.ceil(config.value.fontSize * (config.value.showNextLine ? 2.4 : 1.7) + 20);
    return Math.max(64, minH, Math.ceil(contentH + 16));
};

const fitLyricHeight = () => {
    if (isSettingsOpen.value || isTransitioning.value) return;
    const targetH = computeTargetHeight();
    if (Math.abs(window.innerHeight - targetH) > 8) {
        const w = Math.max(280, window.innerWidth || config.value.windowWidth || 800);
        ResizeWindow(Math.ceil(w), targetH);
    }
};

const scheduleFitHeight = () => {
    if (resizeRafId.value !== null || isSettingsOpen.value) return;
    resizeRafId.value = requestAnimationFrame(() => {
        fitLyricHeight();
        resizeRafId.value = null;
    });
};

const handleOpenSettings = async () => {
    if (isSettingsOpen.value || isTransitioning.value) return;
    setConfig({ ...config.value, windowWidth: Math.max(280, window.innerWidth) });
    isTransitioning.value = true;
    isLocked.value = false;
    isSettingsOpen.value = true;
    await EnterSettingsView(480, 600);
    await nextTick();
    setTimeout(() => {
        isTransitioning.value = false;
    }, 120);
};

const handleCloseSettings = async () => {
    if (!isSettingsOpen.value || isTransitioning.value) return;
    saveConfig(config.value);
    isTransitioning.value = true;
    isSettingsOpen.value = false;
    await nextTick();
    const w = Math.max(280, config.value.windowWidth || 800);
    const h = computeTargetHeight();
    await ExitSettingsView(Math.ceil(w), h);
    setTimeout(() => {
        isTransitioning.value = false;
        scheduleFitHeight();
    }, 120);
};

onMounted(() => {
    if (!config.value.windowWidth) {
        setConfig({ ...config.value, windowWidth: Math.max(280, window.innerWidth) });
    }

    const cancelLock = EventsOn('lock_state_change', (locked: boolean) => {
        isLocked.value = locked;
        setConfig({ ...config.value, isLocked: locked });
    });
    const cancelOpenSettings = EventsOn('open_settings', () => {
        handleOpenSettings();
    });
    ConsumePendingSettings().then((pending) => {
        if (pending) {
            handleOpenSettings();
        }
    }).catch(() => {});

    const handleResize = () => {
        if (isSettingsOpen.value || isTransitioning.value) return;
        const currentW = window.innerWidth;
        if (isLoaded.value && currentW >= 280 && Math.abs(currentW - config.value.windowWidth) > 4) {
            if (resizeTimeoutRef.value) clearTimeout(resizeTimeoutRef.value);
            resizeTimeoutRef.value = window.setTimeout(() => {
                const newConfig = { ...config.value, windowWidth: currentW };
                setConfig(newConfig);
                saveConfig(newConfig);
            }, 250);
        }
    };

    window.addEventListener('resize', handleResize);
    onUnmounted(() => {
        cancelLock();
        cancelOpenSettings();
        window.removeEventListener('resize', handleResize);
        if (resizeRafId.value !== null) {
            cancelAnimationFrame(resizeRafId.value);
        }
        if (resizeTimeoutRef.value) {
            clearTimeout(resizeTimeoutRef.value);
        }
    });
});

watch(
    [
        isSettingsOpen,
        mainText,
        subText,
        () => nextLine.value?.mainText,
        () => config.value.showNextLine,
        () => config.value.fontSize,
        () => config.value.fontFamily,
    ],
    async () => {
        if (isSettingsOpen.value) return;
        await nextTick();
        scheduleFitHeight();
    },
    { flush: 'post' }
);

const rootBackground = computed(() => {
    const isAcrylic = config.value.backdropType === 'acrylic';
    if (isAcrylic) {
        // SetWindowCompositionAttribute paints Acrylic + alpha tint natively behind WebView2.
        return isSettingsOpen.value ? 'rgba(10, 10, 14, 0.28)' : 'transparent';
    }
    if (isSettingsOpen.value) {
        return `rgba(12, 12, 16, ${Math.max(0.65, config.value.bgOpacity).toFixed(2)})`;
    }
    return `rgba(0, 0, 0, ${config.value.bgOpacity})`;
});

const showOuterBorder = computed(() =>
    isSettingsOpen.value || config.value.backdropType === 'acrylic' || config.value.bgOpacity > 0.05
);

const effectiveBorderRadius = computed(() => {
    const r = config.value.borderRadius ?? 12;
    if (config.value.backdropType === 'acrylic') {
        if (r <= 0) return 0;
        if (r <= 6) return 4;
        return 8;
    }
    return r;
});
</script>

<template>
    <div
        :class="[
            'w-screen h-screen overflow-hidden flex flex-col relative',
            showOuterBorder ? 'border border-white/10' : ''
        ]"
        :style="{
            '--wails-draggable': isLocked ? 'none' : 'drag',
            backgroundColor: rootBackground,
            borderRadius: `${effectiveBorderRadius}px`
        }"
    >
        <div
            :class="[
                'w-full h-full flex flex-col relative transition-opacity duration-150 ease-out',
                isTransitioning ? 'opacity-0' : 'opacity-100'
            ]"
        >
            <SettingsPanel
                v-if="isSettingsOpen"
                :config="config"
                @change="setConfig"
                @close="handleCloseSettings"
            />
            <div v-else class="w-full h-full flex items-center justify-center px-4 py-2">
                <LyricRenderer
                    ref="lyricRef"
                    :main-text="mainText"
                    :sub-text="subText"
                    :current-line="currentLine"
                    :next-line="nextLine"
                    :position="playerState?.Position"
                    :config="config"
                />
            </div>
        </div>
    </div>
</template>
