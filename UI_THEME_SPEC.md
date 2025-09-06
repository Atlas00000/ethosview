## EthosView UI Theme Spec — "Ethos Cloud"

Status
- Planning-only. No frontend implementation begins until explicitly requested.

Theme Summary
- Cloudy, glassy, milky aesthetic with vibrant blue→teal gradients.
- FAANG-level polish; tasteful motion; no heavy effects.
- Color system blends bold blues/teals with milky white and deep navy.

Palette
- Primary Blue A/B: `#1E6AE1` / `#3986FF` (actions, gradients)
- Forest/Sea Green A/B: `#1D9A6C` / `#2AB3A6` (ESG accents, gradients)
- Milky White: `#F5F7FA` (ambient base)
- Deep Navy: `#0B2545` (primary text)
- Greys (cool): `#111827`, `#374151`, `#6B7280`, `#D1D5DB`, `#E5E7EB`, `#F3F4F6`

Gradients
- BlueGlass: `linear-gradient(135deg, #1E6AE1 0%, #2AB3A6 50%, #3986FF 100%)`
- EcoSky: `linear-gradient(180deg, #F5F7FA 0%, #E6F7F2 100%)`
- TextGradient: `linear-gradient(135deg, #1E6AE1, #2AB3A6)`

Surfaces (Glassy/Milky)
- Glass card
  - Background: `rgba(255,255,255,0.6)`
  - Border: `1px solid rgba(255,255,255,0.35)`
  - Shadow: `0 8px 30px rgba(16,24,40,0.08)`
  - Backdrop: `backdrop-filter: blur(16px) saturate(120%)`
  - Radius: `16px`
 - Elevation shadow: `0 10px 22px rgba(16,24,40,0.12), 0 2px 6px rgba(16,24,40,0.06)`
- Elevation scale
  - L1: blur 8px, low shadow, radius 12px
  - L2: blur 16px, medium shadow, radius 16px
  - L3: blur 24px, high shadow, radius 20px (modals)

Typography
- Primary: Inter
- Headlines alt: Plus Jakarta Sans (or Satoshi)
- Numbers/KPIs: Tabular lining (`font-variant-numeric: tabular-nums;`)
- Weights: Hero 600; Headings 600; Body 400/500

Motion & Interaction
- Durations: tap 120–160ms; hover 180–220ms; sections 280–360ms; ambient 900–1200ms
- Easing: enter `cubic-bezier(0.2, 0.8, 0.2, 1)`, exit `cubic-bezier(0.4, 0, 0.2, 1)`
- Micro-interactions
  - Card hover: lift 2–4px, subtle tilt, +10% shadow
  - KPI count-up: 600ms ease-out
  - Sparkline draw: 400ms on first reveal; crossfade on refresh
  - Shimmer skeleton: 1200ms diagonal shimmer
- Accessibility: respect `prefers-reduced-motion`; disable nonessential animations

Backgrounds
- App shell: animated EcoSky with radial blue/teal glows
- Cloudy field: ultra-slow gradient shift (EcoSky) + 1–2% noise
- Hero ambient: soft cloud contour; static on mobile

Homepage Components (Visual Intent)
- Hero + KPI band: headline, CTA; KPI pills in glass; last-updated timestamp
- Market snapshot bar: BlueGlass gradient; index/breadth; mini sparkline
- Featured companies carousel: glass slides, masked edge fade, ESG badge
- ESG highlights: sector tiles/mini-bars; hover reveals delta
- Ticker & status: non-blocking ticker; WS status dot (green/amber/red)

Charts
- Duotone lines/areas, rounded joins, 1.5px strokes; low-opacity grids
- Colors: Blue for price, Green for ESG, Amber for alerts, Red `#E25555` for negatives
- Empty/error states: gentle cloud icon + clear message

Accessibility
- Contrast target 4.5:1; avoid translucent text over glass without darkening
- Focus ring: `0 0 0 3px rgba(30,106,225,0.35)`
- Color cues paired with icons/labels (green vs red safe)

Performance Guardrails
- Prefer opacity/transform; limit heavy blurs on mobile
- Reuse shadow tokens; avoid large drop-shadow radii
- Lazy-render offscreen carousels; cap ambient animation GPU usage

Backend-Hugging UX (when implemented)
- Data freshness badges from `/api/v1/dashboard` timestamps
- Skeletons while fetching; crossfade on refresh (no layout shift)
- WS status indicator from `/api/v1/ws/status`; fallback polling of `/alerts`
- Respect endpoint TTLs; replace long spinners with skeletons

Semantic Tokens (CSS variables)
- Brand: `--brand-blue:#1E6AE1; --brand-green:#1D9A6C; --brand-milk:#F5F7FA; --blue-a:#1E6AE1; --blue-b:#3986FF; --teal-a:#2AB3A6; --teal-b:#1DAA8E;`
- Text: `--text-primary:#0B2545; --text-secondary:#374151;`
- Surfaces: `--glass-bg:rgba(255,255,255,0.6); --glass-border:rgba(255,255,255,0.35); --glass-blur:16px;`
- States: `--ok:#1D9A6C; --warn:#FFB547; --bad:#E25555;`
- Utilities: `.text-gradient`, `.animated-gradient`, `.hover-lift`, `.animate-fade-in-up`, `.edge-fade`, `.skeleton`, `.app-shell`, `.blob`, `.shadow-elevated`

Enforcement
- Use `.text-gradient` for section headings and brand wordmark.
- Use glass cards for KPI/summary surfaces; avoid plain borders.
- Prefer `animated-gradient` for bars and accent surfaces.
- Use `.blob` accents to avoid boxy sections; keep opacity ≤ 0.2.

Notes
- Keep animations subtle; respect reduced-motion.

