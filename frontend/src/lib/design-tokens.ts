// Design tokens for Remove Background App
// Inspired by Linear.app / Vercel Dashboard / Arc Browser

export const tokens = {
  colors: {
    background: "#0A0A0F",
    surface: "#111118",
    surfaceElevated: "#1C1C27",
    border: "#2A2A3D",
    accent: "#7C6EFA",
    accentHover: "#9D8FFF",
    success: "#34D399",
    error: "#F87171",
    textPrimary: "#F0F0FF",
    textMuted: "#8B8BA7",
  },
  font: {
    family: {
      sans: ["Inter", "system-ui", "-apple-system", "sans-serif"],
    },
    weight: {
      light: 300,
      normal: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
    },
  },
  radius: {
    DEFAULT: "12px",
    lg: "16px",
    full: "999px",
  },
  shadow: {
    glow: "0 0 20px rgba(124, 110, 250, 0.25)",
    "glow-lg": "0 0 40px rgba(124, 110, 250, 0.35)",
  },
  transition: {
    fast: "150ms ease",
    normal: "200ms ease",
    slow: "300ms ease",
  },
  spacing: {
    xs: "4px",
    sm: "8px",
    md: "16px",
    lg: "24px",
    xl: "32px",
    "2xl": "48px",
  },
} as const;

export type Tokens = typeof tokens;
export type ColorKey = keyof typeof tokens.colors;
export type RadiusKey = keyof typeof tokens.radius;

// CSS Variables export for runtime use
export const cssVariables = `
  --color-background: ${tokens.colors.background};
  --color-surface: ${tokens.colors.surface};
  --color-surface-elevated: ${tokens.colors.surfaceElevated};
  --color-border: ${tokens.colors.border};
  --color-accent: ${tokens.colors.accent};
  --color-accent-hover: ${tokens.colors.accentHover};
  --color-success: ${tokens.colors.success};
  --color-error: ${tokens.colors.error};
  --color-text-primary: ${tokens.colors.textPrimary};
  --color-text-muted: ${tokens.colors.textMuted};
  --radius-default: ${tokens.radius.DEFAULT};
  --radius-lg: ${tokens.radius.lg};
  --radius-full: ${tokens.radius.full};
  --shadow-glow: ${tokens.shadow.glow};
  --shadow-glow-lg: ${tokens.shadow["glow-lg"]};
  --transition-fast: ${tokens.transition.fast};
  --transition-normal: ${tokens.transition.normal};
  --transition-slow: ${tokens.transition.slow};
`.trim();
