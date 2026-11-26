import { ref, watch, type Ref, type ComputedRef } from 'vue';
import { FetchLyrics } from "../../wailsjs/go/main/App";
import { LyricLine } from "../types";

export function useLyrics(
    fileName: Ref<string | undefined> | ComputedRef<string | undefined>,
    position: Ref<number | undefined> | ComputedRef<number | undefined>
) {
    const lyrics = ref<LyricLine[]>([]);
    const mainText = ref<string>("AILrc Ready");
    const subText = ref<string>("");
    const currentLine = ref<LyricLine | null>(null);
    const nextLine = ref<LyricLine | null>(null);

    let fetchSeq = 0;

    watch(fileName, async (newFileName) => {
        const seq = ++fetchSeq;
        if (!newFileName) {
            lyrics.value = [];
            currentLine.value = null;
            nextLine.value = null;
            mainText.value = "AILrc Ready";
            subText.value = "";
            return;
        }

        const lrc = await FetchLyrics(newFileName);
        if (seq !== fetchSeq) return;
        lyrics.value = lrc || [];
        if (lyrics.value.length === 0) {
            currentLine.value = null;
            nextLine.value = null;
            mainText.value = "♪";
            subText.value = "";
        }
    }, { immediate: true });

    let throttleTimer: number | null = null;
    let lastPosition: number | undefined = undefined;

    watch([position, lyrics], ([newPosition, newLyrics]) => {
        if (newPosition === undefined || newLyrics.length === 0) return;

        if (throttleTimer !== null) {
            lastPosition = newPosition;
            return;
        }

        const updateLyric = (pos: number) => {
            let left = 0;
            let right = newLyrics.length - 1;
            let idx = -1;

            while (left <= right) {
                const mid = Math.floor((left + right) / 2);
                if (newLyrics[mid].time <= pos) {
                    idx = mid;
                    left = mid + 1;
                } else {
                    right = mid - 1;
                }
            }

            if (idx !== -1) {
                currentLine.value = newLyrics[idx];
                nextLine.value = idx + 1 < newLyrics.length ? newLyrics[idx + 1] : null;
                mainText.value = newLyrics[idx].mainText;
                subText.value = newLyrics[idx].subText;
            } else {
                currentLine.value = null;
                nextLine.value = newLyrics[0] || null;
                if (pos < newLyrics[0].time) {
                    mainText.value = "...";
                    subText.value = "";
                }
            }
        };

        updateLyric(newPosition);

        throttleTimer = window.setTimeout(() => {
            throttleTimer = null;
            if (lastPosition !== undefined && lastPosition !== newPosition) {
                updateLyric(lastPosition);
                lastPosition = undefined;
            }
        }, 40);
    });

    return { mainText, subText, currentLine, nextLine };
}
