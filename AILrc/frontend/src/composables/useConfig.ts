import { ref, onMounted } from 'vue';
import { GetConfig, SaveConfig } from "../../wailsjs/go/main/App";
import { AppConfig } from "../types";

const defaultConfig: AppConfig = {
    fontSize: 48,
    fontFamily: "Segoe UI",
    fontWeight: 700,
    fontColor: "#FFFFFF",
    karaokeColor: "#FFD500",
    strokeColor: "#000000",
    bgOpacity: 0.6,
    textOpacity: 1.0,
    karaokeMode: "sweep",
    lyricEffect: "glow",
    effectStrength: 1.0,
    effectSpeed: 1.0,
    textAlign: "center",
    letterSpacing: 0,
    glowIntensity: 10,
    borderRadius: 12,
    backdropType: "none",
    showNextLine: true,
    lyricOffsetMs: 0,
    isLocked: false,
    windowWidth: 800
};

export function useConfig() {
    const config = ref<AppConfig>(defaultConfig);
    const isLoaded = ref(false);

    onMounted(async () => {
        try {
            const cfg = await GetConfig();
            if (cfg) {
                config.value = { ...config.value, ...(cfg as Partial<AppConfig>) };
            }
            isLoaded.value = true;
        } catch {
            isLoaded.value = true;
        }
    });

    const setConfig = (newConfig: AppConfig) => {
        config.value = newConfig;
    };

    const save = (newConfig: AppConfig) => {
        config.value = newConfig;
        SaveConfig(newConfig as any);
    };

    return { config, setConfig, save, isLoaded };
}
