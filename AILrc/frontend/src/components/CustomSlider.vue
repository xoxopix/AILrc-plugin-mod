<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';

interface Props {
    label: string;
    value: number;
    min: number;
    max: number;
    step?: number;
    unit?: string;
    displayFn?: (val: number) => string;
    centered?: boolean;
    inverted?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
    change: [value: number];
}>();

const trackRef = ref<HTMLDivElement>();
const isDragging = ref(false);
let rafId: number | null = null;

const safeValue = computed(() => props.value ?? props.min);
const percentage = computed(() => {
    const raw = props.inverted
        ? ((props.max - safeValue.value) / (props.max - props.min)) * 100
        : ((safeValue.value - props.min) / (props.max - props.min)) * 100;
    return Math.min(100, Math.max(0, raw));
});

const fillStyle = computed(() => {
    if (!props.centered) {
        return { left: '0%', width: `${percentage.value}%` };
    }
    return {
        left: `${Math.min(50, percentage.value)}%`,
        width: `${Math.abs(percentage.value - 50)}%`,
    };
});

const displayValue = computed(() => {
    if (props.displayFn) {
        return props.displayFn(safeValue.value);
    }
    return `${safeValue.value}${props.unit || ''}`;
});

const calculateValue = (clientX: number) => {
    if (!trackRef.value) return;
    const rect = trackRef.value.getBoundingClientRect();

    let percent = (clientX - rect.left) / rect.width;
    percent = Math.max(0, Math.min(1, percent));

    let rawValue = props.inverted
        ? props.max - percent * (props.max - props.min)
        : props.min + percent * (props.max - props.min);

    if (props.step) {
        rawValue = Math.round(rawValue / props.step) * props.step;
    }

    const newValue = parseFloat(rawValue.toFixed(2));
    emit('change', newValue);
};

const handlePointerDown = (e: PointerEvent) => {
    e.preventDefault();
    e.stopPropagation();
    (e.currentTarget as HTMLElement | null)?.setPointerCapture?.(e.pointerId);
    isDragging.value = true;
    calculateValue(e.clientX);
};

const handleGlobalMove = (e: PointerEvent) => {
    if (!isDragging.value) return;

    e.preventDefault();

    if (rafId !== null) return;

    rafId = requestAnimationFrame(() => {
        calculateValue(e.clientX);
        rafId = null;
    });
};

const stopDragging = () => {
    if (!isDragging.value) return;

    isDragging.value = false;

    if (rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
    }
};

const handleGlobalUp = (e: PointerEvent) => {
    if (!isDragging.value) return;
    e.preventDefault();
    stopDragging();
};

onMounted(() => {
    window.addEventListener('pointermove', handleGlobalMove);
    window.addEventListener('pointerup', handleGlobalUp);
    window.addEventListener('pointercancel', handleGlobalUp);
    window.addEventListener('blur', stopDragging);
});

onUnmounted(() => {
    window.removeEventListener('pointermove', handleGlobalMove);
    window.removeEventListener('pointerup', handleGlobalUp);
    window.removeEventListener('pointercancel', handleGlobalUp);
    window.removeEventListener('blur', stopDragging);

    if (rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
    }
});
</script>

<template>
    <section class="space-y-3 select-none">
        <div class="flex justify-between items-center text-[11px] font-bold text-white/50 uppercase tracking-wider">
            <span>{{ label }}</span>
            <span :class="['font-mono transition-colors', isDragging ? 'text-pink-400' : 'text-white/30']">
                {{ displayValue }}
            </span>
        </div>

        <slot name="between" />

        <div
            class="h-6 flex items-center cursor-pointer group touch-none px-3.5"
            @pointerdown="handlePointerDown"
            :style="{ '--wails-draggable': 'none' }"
        >
            <div
                ref="trackRef"
                class="relative w-full h-1.5 bg-white/10 rounded-full overflow-visible transition-colors group-hover:bg-white/20"
            >
                <div
                    v-if="centered"
                    class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-0.5 h-3 bg-white/30 rounded-full pointer-events-none"
                />

                <div
                    class="absolute top-0 h-full bg-pink-500 rounded-full pointer-events-none"
                    :style="fillStyle"
                />

                <div
                    :class="['absolute top-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 bg-white rounded-full shadow-lg shadow-black/50 pointer-events-none transition-transform duration-150 border border-black/10', isDragging ? 'scale-125 bg-pink-100 shadow-pink-500/50' : 'scale-100 group-hover:scale-110']"
                    :style="{ left: `${percentage}%` }"
                />
            </div>
        </div>
    </section>
</template>
