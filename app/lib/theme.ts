/**
 * Centralized design tokens. Mirrors the CSS variables defined in
 * assets/css/main.css and assets/css/font.css. Use these whenever a value
 * is needed in JS (chart configs, dynamic styles, motion variants, etc.)
 * so that the source of truth stays single.
 */

export const colors = {
  canvas: "var(--jv-canvas)",
  white: "var(--jv-white)",
  slate: "var(--jv-slate)",
  ink: "var(--jv-ink)",
  ink2: "var(--jv-ink-2)",
  muted: "var(--jv-muted)",
  coral: "var(--jv-coral)",
  yellow: "var(--jv-yellow)",
  yellow2: "var(--jv-yellow-2)",
  yellowSoft: "var(--jv-yellow-soft)",
  mint: "var(--jv-mint)",
  mint2: "var(--jv-mint-2)",
  salmon: "var(--jv-salmon)",
  lavender: "var(--jv-lavender)",
  ivory: "var(--jv-ivory)",
} as const;

export const accentColors = {
  green: "var(--jv-accent-green)",
  gold: "var(--jv-accent-gold)",
  red: "var(--jv-accent-red)",
  purple: "var(--jv-accent-purple)",
  orange: "var(--jv-accent-orange)",
} as const;

export const statusColors = {
  success: accentColors.green,
  warning: accentColors.gold,
  danger: accentColors.red,
  info: accentColors.purple,
  notice: accentColors.orange,
} as const;

export const fonts = {
  body: "var(--font-body)",
  headings: "var(--font-headings)",
  feature: "var(--font-feature)",
} as const;

export const shadows = {
  brutalSm: "var(--shadow-brutal-sm)",
  brutal: "var(--shadow-brutal)",
  brutalLg: "var(--shadow-brutal-lg)",
} as const;

export const radius = {
  sm: "calc(var(--radius) - 4px)",
  md: "calc(var(--radius) - 2px)",
  lg: "var(--radius)",
  xl: "calc(var(--radius) + 4px)",
} as const;

/** Resolves a CSS variable token to its actual computed color string. */
export function resolveCssVar(token: string): string {
  if (typeof window === "undefined") return token;
  const name = token.replace(/^var\(|\)$/g, "").trim();
  return getComputedStyle(document.documentElement)
    .getPropertyValue(name)
    .trim();
}

export const theme = {
  colors,
  accentColors,
  statusColors,
  fonts,
  shadows,
  radius,
} as const;

export type Theme = typeof theme;
