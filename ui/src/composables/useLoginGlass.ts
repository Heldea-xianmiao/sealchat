import { computed, onBeforeUnmount, ref, watch } from 'vue';
import type { ComputedRef, Ref } from 'vue';
import type { LoginBackgroundConfig } from '@/types';

type RGB = { r: number; g: number; b: number };
type SampleResult = { color: RGB; luma: number };

type LoginGlassOptions = {
  imageUrl: Ref<string>;
  config: Ref<LoginBackgroundConfig | null | undefined>;
  enabled?: Ref<boolean> | boolean;
  radius?: string;
};

const DEFAULT_LIGHT_BASE: RGB = { r: 247, g: 248, b: 250 };
const DEFAULT_DARK_BASE: RGB = { r: 16, g: 20, b: 28 };

const sampleCache = new Map<string, SampleResult>();

const clamp = (value: number, min: number, max: number) => Math.min(max, Math.max(min, value));
const lerp = (from: number, to: number, t: number) => from + (to - from) * t;

const mixColors = (a: RGB, b: RGB, t: number): RGB => ({
  r: Math.round(lerp(a.r, b.r, t)),
  g: Math.round(lerp(a.g, b.g, t)),
  b: Math.round(lerp(a.b, b.b, t)),
});

const colorLuma = (color: RGB) => {
  return (0.2126 * color.r + 0.7152 * color.g + 0.0722 * color.b) / 255;
};

const desaturate = (color: RGB, amount: number): RGB => {
  const gray = Math.round((color.r + color.g + color.b) / 3);
  return {
    r: Math.round(lerp(color.r, gray, amount)),
    g: Math.round(lerp(color.g, gray, amount)),
    b: Math.round(lerp(color.b, gray, amount)),
  };
};

const toRgbString = (color: RGB) => `rgb(${color.r}, ${color.g}, ${color.b})`;
const toRgbaString = (color: RGB, alpha: number) => `rgba(${color.r}, ${color.g}, ${color.b}, ${alpha.toFixed(3)})`;

const parseHexColor = (value: string): RGB | null => {
  const raw = value.trim().replace('#', '');
  if (raw.length === 3) {
    const r = Number.parseInt(raw[0] + raw[0], 16);
    const g = Number.parseInt(raw[1] + raw[1], 16);
    const b = Number.parseInt(raw[2] + raw[2], 16);
    if (Number.isNaN(r) || Number.isNaN(g) || Number.isNaN(b)) return null;
    return { r, g, b };
  }
  if (raw.length === 6) {
    const r = Number.parseInt(raw.slice(0, 2), 16);
    const g = Number.parseInt(raw.slice(2, 4), 16);
    const b = Number.parseInt(raw.slice(4, 6), 16);
    if (Number.isNaN(r) || Number.isNaN(g) || Number.isNaN(b)) return null;
    return { r, g, b };
  }
  return null;
};

const parseRgbColor = (value: string): RGB | null => {
  const match = value.match(/^rgba?\((.+)\)$/i);
  if (!match) return null;
  const parts = match[1].split(',').map((part) => Number.parseFloat(part.trim()));
  if (parts.length < 3) return null;
  const [r, g, b] = parts;
  if ([r, g, b].some((v) => Number.isNaN(v))) return null;
  return {
    r: clamp(Math.round(r), 0, 255),
    g: clamp(Math.round(g), 0, 255),
    b: clamp(Math.round(b), 0, 255),
  };
};

const parseColor = (value?: string | null): RGB | null => {
  if (!value) return null;
  if (value.startsWith('#')) return parseHexColor(value);
  if (value.startsWith('rgb')) return parseRgbColor(value);
  return null;
};

const getDisplayPalette = () => {
  if (typeof document === 'undefined') return 'day';
  return document.documentElement?.dataset?.displayPalette === 'night' ? 'night' : 'day';
};

const loadImage = (url: string) => new Promise<HTMLImageElement>((resolve, reject) => {
  const img = new Image();
  img.crossOrigin = 'anonymous';
  img.onload = () => resolve(img);
  img.onerror = (err) => reject(err);
  img.src = url;
});

const sampleImage = async (url: string): Promise<SampleResult | null> => {
  if (!url) return null;
  if (sampleCache.has(url)) return sampleCache.get(url) || null;
  try {
    const img = await loadImage(url);
    const size = 24;
    const canvas = document.createElement('canvas');
    canvas.width = size;
    canvas.height = size;
    const ctx = canvas.getContext('2d', { willReadFrequently: true });
    if (!ctx) return null;
    ctx.drawImage(img, 0, 0, size, size);
    const data = ctx.getImageData(0, 0, size, size).data;
    let r = 0;
    let g = 0;
    let b = 0;
    let total = 0;
    for (let i = 0; i < data.length; i += 4) {
      const alpha = data[i + 3] / 255;
      if (alpha < 0.05) continue;
      r += data[i] * alpha;
      g += data[i + 1] * alpha;
      b += data[i + 2] * alpha;
      total += alpha;
    }
    if (total <= 0) return null;
    const color = {
      r: clamp(Math.round(r / total), 0, 255),
      g: clamp(Math.round(g / total), 0, 255),
      b: clamp(Math.round(b / total), 0, 255),
    };
    const result = { color, luma: colorLuma(color) };
    sampleCache.set(url, result);
    return result;
  } catch {
    return null;
  }
};

const applyBrightness = (color: RGB, brightness: number): RGB => ({
  r: clamp(Math.round(color.r * brightness), 0, 255),
  g: clamp(Math.round(color.g * brightness), 0, 255),
  b: clamp(Math.round(color.b * brightness), 0, 255),
});

const applyOverlay = (color: RGB, overlay: RGB, opacity: number): RGB => {
  return mixColors(color, overlay, clamp(opacity, 0, 1));
};

const adjustLumaRange = (color: RGB, minLuma: number, maxLuma: number): RGB => {
  let luma = colorLuma(color);
  if (luma < minLuma) {
    const mix = (minLuma - luma) / (1 - luma);
    const adjusted = mixColors(color, { r: 255, g: 255, b: 255 }, clamp(mix, 0, 1));
    luma = colorLuma(adjusted);
    return luma < minLuma ? mixColors(adjusted, { r: 255, g: 255, b: 255 }, 0.1) : adjusted;
  }
  if (luma > maxLuma) {
    const mix = (luma - maxLuma) / luma;
    const adjusted = mixColors(color, { r: 0, g: 0, b: 0 }, clamp(mix, 0, 1));
    luma = colorLuma(adjusted);
    return luma > maxLuma ? mixColors(adjusted, { r: 0, g: 0, b: 0 }, 0.1) : adjusted;
  }
  return color;
};

const buildGlassVars = (cfg: LoginBackgroundConfig | null | undefined, sample: SampleResult | null, palette: 'day' | 'night') => {
  const autoTint = cfg?.panelAutoTint ?? true;
  const baseTint = parseColor(cfg?.panelTintColor) ?? (palette === 'night' ? DEFAULT_DARK_BASE : DEFAULT_LIGHT_BASE);
  const baseOpacity = clamp((cfg?.panelTintOpacity ?? 72) / 100, 0.25, 0.98);
  const blur = clamp(cfg?.panelBlur ?? 14, 0, 30);
  const saturate = clamp(cfg?.panelSaturate ?? 120, 80, 180);
  const contrast = clamp(cfg?.panelContrast ?? 105, 90, 140);
  const borderOpacity = clamp((cfg?.panelBorderOpacity ?? 18) / 100, 0, 0.6);
  const shadowStrength = clamp((cfg?.panelShadowStrength ?? 22) / 100, 0, 1);

  let tint = baseTint;
  let sourceLuma = colorLuma(baseTint);

  if (autoTint && sample) {
    let adjusted = sample.color;
    const brightness = clamp((cfg?.brightness ?? 100) / 100, 0.5, 1.5);
    adjusted = applyBrightness(adjusted, brightness);
    if (cfg?.overlayColor && cfg?.overlayOpacity) {
      const overlayColor = parseColor(cfg.overlayColor);
      if (overlayColor) {
        adjusted = applyOverlay(adjusted, overlayColor, clamp(cfg.overlayOpacity / 100, 0, 1));
      }
    }
    const bgOpacity = clamp((cfg?.opacity ?? 30) / 100, 0, 1);
    const neutral = palette === 'night' ? DEFAULT_DARK_BASE : DEFAULT_LIGHT_BASE;
    adjusted = mixColors(neutral, adjusted, bgOpacity);
    sourceLuma = colorLuma(adjusted);

    const desat = desaturate(adjusted, 0.45);
    const autoBase = sourceLuma < 0.52 ? DEFAULT_LIGHT_BASE : DEFAULT_DARK_BASE;
    const baseMix = mixColors(baseTint, autoBase, 0.6);
    tint = mixColors(desat, baseMix, 0.55);
  }

  tint = adjustLumaRange(tint, 0.18, 0.92);

  let alpha = baseOpacity;
  if (autoTint && sample) {
    const adjust = clamp((0.55 - sourceLuma) * 0.28, -0.12, 0.18);
    alpha = clamp(baseOpacity + adjust, 0.5, 0.92);
  }

  const tintLuma = colorLuma(tint);
  const borderBase = tintLuma < 0.5 ? { r: 255, g: 255, b: 255 } : { r: 0, g: 0, b: 0 };
  const borderColor = mixColors(tint, borderBase, 0.5);

  const shadowAlphaA = (0.08 + 0.18 * shadowStrength).toFixed(3);
  const shadowAlphaB = (0.04 + 0.12 * shadowStrength).toFixed(3);
  const shadow = `0 10px 28px rgba(0, 0, 0, ${shadowAlphaA}), 0 4px 12px rgba(0, 0, 0, ${shadowAlphaB})`;

  return {
    '--sc-glass-bg': toRgbaString(tint, alpha),
    '--sc-glass-solid': toRgbString(tint),
    '--sc-glass-border': toRgbaString(borderColor, borderOpacity),
    '--sc-glass-shadow': shadow,
    '--sc-glass-blur': `${blur}px`,
    '--sc-glass-saturate': `${saturate}%`,
    '--sc-glass-contrast': `${contrast}%`,
  } as Record<string, string>;
};

export const useLoginGlass = (options: LoginGlassOptions) => {
  const { imageUrl, config, radius } = options;
  const enabled = options.enabled ?? true;
  const palette = ref<'day' | 'night'>(getDisplayPalette());
  let observer: MutationObserver | null = null;

  if (typeof window !== 'undefined' && typeof MutationObserver !== 'undefined') {
    observer = new MutationObserver(() => {
      palette.value = getDisplayPalette();
    });
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['data-display-palette'] });
  }

  onBeforeUnmount(() => {
    observer?.disconnect();
    observer = null;
  });

  const glassVars = ref<Record<string, string>>({});
  let refreshToken = 0;

  const refresh = async () => {
    const token = ++refreshToken;
    const isEnabled = typeof enabled === 'boolean' ? enabled : enabled.value;
    const url = imageUrl.value;
    let sample: SampleResult | null = null;
    if (isEnabled && url && typeof window !== 'undefined') {
      sample = await sampleImage(url);
    }
    if (token !== refreshToken) return;
    glassVars.value = buildGlassVars(config.value, sample, palette.value);
  };

  watch(imageUrl, refresh, { immediate: true });
  watch(config, refresh, { deep: true });
  watch(palette, refresh);
  watch(() => (typeof enabled === 'boolean' ? enabled : enabled.value), refresh);

  const glassStyle: ComputedRef<Record<string, string>> = computed(() => {
    if (!radius) return glassVars.value;
    return { ...glassVars.value, '--sc-glass-radius': radius };
  });

  return { glassStyle, refresh };
};
