# Nord Design System

Shared by **thomaskraeuter.com**, **smoothtime.io** and **status.smoothtime.io**. Live reference: https://thomaskraeuter.com/design (hidden, `noindex`).

The single source of truth is `tokens.css` (identical copy in all three repos:
`thomaskraeuter-website/src/styles/tokens.css`, `smoothtime/website/tokens.css`, `smoothtime-status/styles/tokens.css`).
It is the **only** file allowed to contain raw colour values; everything else uses `var(--…)`.
When changing a token, change it in all three repos. Exceptions that must hold literals (kept in sync by hand): favicons/logos (SVG files), `theme-color` meta tags, web manifests, and the Mantine colour tuple in `smoothtime-status/util/nord.ts`.

## Principles
Calm, technical, timeless. Classic Nord palette (nordtheme.com), dark Polar Night base,
Frost for interaction, Aurora sparingly for status and LCARS accents. Motion is slow and
purposeful (LCARS bars draw in, orbit/topo drift), and fully disabled under `prefers-reduced-motion`.

## Palette
| Group | Tokens |
|---|---|
| Polar Night | `--nord0` #2E3440 · `--nord1` #3B4252 · `--nord2` #434C5E · `--nord3` #4C566A |
| Snow Storm | `--nord4` #D8DEE9 · `--nord5` #E5E9F0 · `--nord6` #ECEFF4 |
| Frost | `--nord7` #8FBCBB · `--nord8` #88C0D0 · `--nord9` #81A1C1 · `--nord10` #5E81AC |
| Aurora | `--nord11` #BF616A · `--nord12` #D08770 · `--nord13` #EBCB8B · `--nord14` #A3BE8C · `--nord15` #B48EAD |
| Derived (AA) | `--nord-muted` #A7B1C2 · `--nord-red-text` #D9A0A6 · `--nord-sunken` #292E39 |

## Semantic colour
| Token | Value | Use |
|---|---|---|
| `--color-bg` | nord0 | page background |
| `--color-bg-sunken` | nord-sunken | inputs, code wells |
| `--color-surface` | nord1 | cards, panels |
| `--color-surface-raised` | nord2 | hover / raised cards |
| `--color-border` / `-strong` | nord1 / nord3 | hairlines / card + control borders |
| `--color-text` | nord4 | body |
| `--color-text-strong` | nord6 | headings |
| `--color-text-muted` | nord-muted | secondary text |
| `--color-text-faint` | nord3 | decoration only (never text) |
| `--color-accent` / `-hover` / `-strong` | nord8 / nord7 / nord9 | links, primary buttons |
| `--color-focus`, `--focus-ring` | nord8 | keyboard focus |
| `--color-success/warning/danger/info` | nord14 / nord13 / nord11 / nord8 | status (`--color-danger-text` for red text) |
| `--lcars-elbow,-1,-2,-3,-4,-rail` | nord8, nord13, nord12, nord9, nord15, nord1 | LCARS bar segments |

## Contrast (WCAG 2.x, computed)
| Pair | Ratio | |
|---|---|---|
| text nord4 / bg | 9.25 | AAA |
| text-strong nord6 / bg | 10.84 | AAA |
| text-muted / bg · / surface | 5.78 · 4.65 | AA |
| accent nord8 / bg · / surface | 6.24 · 5.03 | AA |
| bg on accent (primary button) · on accent-hover | 6.24 · 5.99 | AA |
| success / bg · warning / bg | 6.13 · 8.00 | AA / AAA |
| danger-text / bg · / surface | 5.67 · 4.57 | AA |
| nord11 / bg | 3.05 | UI / large only |
| nord10 / bg | 3.10 | UI only |
| nord3 / bg | 1.69 | decoration only |

## Typography
Sans `--font-sans` (Inter), mono `--font-mono` (JetBrains Mono).
Sizes `--text-2xs` .66rem · `xs` .75 · `sm` .875 · `base` 1 · `md` 1.125 · `lg` 1.375 · `xl` clamp(1.6-2.4rem) · `2xl` clamp(2.2-3.6rem) · `3xl` clamp(3-6.4rem).
Weights 400/500/600/650. Leading 1.05/1.25/1.6/1.75. Tracking tight −.035em, caps .1em (mono labels).

## Spacing, radii, elevation
Spacing (4px base): `--space-1` .25rem · 2 .5 · 3 .75 · 4 1 · 5 1.5 · 6 2 · 7 3 · 8 4 · 9 6rem. Container 1180px, gutter clamp(1.25-2.5rem).
Radii: none 0 · sm 4 · md 6 · lg 10 · pill 999px.
Shadows: `--shadow-1/2/3`, glows `--glow-accent`, `--glow-success`.

## Motion
`--dur-fast` 150ms · `base` 300 · `slow` 600 · `slower` 1200 · `ambient` 90s; `--ease-standard` (.2,.6,.2,1), `--ease-emphasized` (.2,.7,.2,1), `--ease-in-out` (.4,0,.2,1); `--stagger` 70ms.
`prefers-reduced-motion: reduce` sets all durations to 0 and disables animations/transitions globally.

## Breakpoints
`--bp-sm` 420 · `--bp-md` 640 · `--bp-lg` 860 · `--bp-xl` 1180 (reference values; media queries use the literals).

## Components (`ds-` classes in tokens.css)
- Buttons: `.ds-btn` + `--primary` / `--secondary` / `--ghost`
- Links: `.ds-link`
- Cards: `.ds-card`, `.ds-card--raised`
- Badges/tags: `.ds-badge` (+ `--success|warning|danger|info`), `.ds-tag`
- Inputs: `.ds-input`, `.ds-label`
- LCARS bar: `.ds-lcars` with `.elbow .s1 .s2 .s3 .rail .s4`
- Section label: `.ds-section-label` (`[LOG]` prefix)
- Code: `.ds-code`, `.ds-pre`

Site-specific classes (e.g. `.cta`, `.btn`, `.card`) are thin aliases built only from these tokens.

## Marks, logos & favicons
All marks are a nord0 rounded square (rx 13/64) with a 2.5/64 Frost (nord8) frame.
- **thomaskraeuter.com**: geometric "TK" in nord6 above a three-segment LCARS bar (nord13 · nord12 · nord9).
- **smoothtime.io**: the smooth-return clock glyph in nord8.
- **status.smoothtime.io**: the smoothtime mark plus a nord14 status dot.
Each site ships `favicon.svg` (also used as the header logo), `favicon.ico` (16/32/48), `favicon-32.png`,
`apple-touch-icon.png` (180, full-bleed nord0), `icon-192.png`, `icon-512.png` and `site.webmanifest`.

## Status page (smoothtime-status)
UptimeFlare stays upstream-mergeable: styling lives in `styles/tokens.css` + `styles/globals.css` (aliases `--st-*` onto tokens),
status colours map to `--color-success/warning/danger(-text)`, Mantine uses the `nord` primary tuple, and the header adds the LCARS rail.
