# SynapSeq Hub - UI Specification

## Overview

This document defines the complete UI specification for SynapSeq Hub, including color palettes, typography, spacing, components, and interaction patterns. This specification enables consistent design implementation across all SynapSeq projects.

---

## Color Palette

### Base Colors

```css
--color-white: #ffffff;
--color-black: #000000;
```

### Gray Scale

```css
--color-gray-50: #f8fafc;
--color-gray-100: #f1f5f9;
--color-gray-200: #e2e8f0;
--color-gray-300: #cbd5e1;
--color-gray-400: #94a3b8;
--color-gray-500: #64748b;
--color-gray-600: #475569;
--color-gray-700: #334155;
--color-gray-800: #1e293b;
--color-gray-850: #172033;
--color-gray-900: #0f172a;
--color-gray-950: #020617;
```

### Blue Colors (Primary Brand)

```css
--color-blue-50: #eff6ff;
--color-blue-100: #dbeafe;
--color-blue-200: #bfdbfe;
--color-blue-300: #93c5fd;
--color-blue-400: #60a5fa;
--color-blue-500: #3b82f6;
--color-blue-600: #2563eb;
--color-blue-700: #1d4ed8;
--color-blue-800: #1e40af;
--color-blue-900: #1e3a8a;
--color-blue-950: #172554;
```

### Cyan Colors (Secondary Brand)

```css
--color-cyan-300: #67e8f9;
--color-cyan-400: #22d3ee;
--color-cyan-500: #06b6d4;
--color-cyan-600: #0891b2;
--color-cyan-700: #0e7490;
```

### Green Colors (Success/Action)

```css
--color-green-100: #dcfce7;
--color-green-200: #bbf7d0;
--color-green-300: #86efac;
--color-green-400: #4ade80;
--color-green-600: #16a34a;
--color-green-700: #15803d;
--color-green-800: #166534;
--color-green-900: #14532d;
--color-green-950: #052e16;
```

---

## Theme Application

### Light Mode

- **Background**: Linear gradient from `gray-50` via `blue-50/30` to `cyan-50/30`
- **Cards**: `white/80` with backdrop blur
- **Borders**: `gray-200`
- **Text Primary**: `gray-900`
- **Text Secondary**: `gray-600`
- **Text Tertiary**: `gray-500`
- **Accents**: `blue-600` to `cyan-500` gradient

### Dark Mode

Applied via `.dark` class on `<html>` element.

- **Background**: Linear gradient from `gray-950` via `gray-950` to `gray-900`
- **Cards**: `gray-900/80` with backdrop blur
- **Borders**: `gray-700` to `gray-800`
- **Text Primary**: `white` or `gray-100`
- **Text Secondary**: `gray-300`
- **Text Tertiary**: `gray-400`
- **Accents**: `blue-400` to `cyan-400` gradient

---

## Typography

### Font Families

- **Primary**: System font stack (default)
- **Monospace**: `'Monaco', 'Menlo', 'Ubuntu Mono', monospace`

### Font Sizes

- **Hero Title**: `3rem` (48px) on mobile, `3.75rem` (60px) on desktop
- **H1**: `1.875rem` (30px) to `2rem` (32px)
- **H2**: `1.5rem` (24px)
- **H3**: `1.25rem` (20px)
- **Body**: `1rem` (16px)
- **Small**: `0.875rem` (14px)
- **Extra Small**: `0.75rem` (12px)

### Font Weights

- **Regular**: 400
- **Medium**: 500
- **Semibold**: 600
- **Bold**: 700

### Line Heights

- **Tight**: 1.2
- **Normal**: 1.5
- **Relaxed**: 1.7

---

## Spacing Scale

Using Tailwind's default spacing scale (0.25rem increments):

- `0`: 0
- `1`: 0.25rem (4px)
- `2`: 0.5rem (8px)
- `3`: 0.75rem (12px)
- `4`: 1rem (16px)
- `5`: 1.25rem (20px)
- `6`: 1.5rem (24px)
- `8`: 2rem (32px)
- `12`: 3rem (48px)
- `16`: 4rem (64px)

---

## Border Radius

- **Small**: `0.5rem` (8px) - `.rounded-lg`
- **Medium**: `0.75rem` (12px) - `.rounded-xl`
- **Large**: `1rem` (16px) - `.rounded-2xl`
- **Full**: `9999px` - `.rounded-full`

---

## Shadows

### Light Mode

- **Small**: `0 1px 3px rgb(0 0 0 / 0.1)`
- **Medium**: `0 4px 6px rgb(0 0 0 / 0.1)`
- **Large**: `0 20px 25px -5px rgb(0 0 0 / 0.1)`

### Dark Mode

- **Small**: `0 1px 3px rgb(0 0 0 / 0.3)`
- **Medium**: `0 4px 6px rgb(0 0 0 / 0.3)`
- **Large**: `0 20px 25px -5px rgb(0 0 0 / 0.5)`

### Colored Shadows (Hover States)

- **Blue**: `0 6px 10px -1px rgba(59, 130, 246, 0.1)`
- **Cyan**: `0 6px 10px -1px rgba(34, 211, 238, 0.1)`

---

## Gradients

### Brand Gradient (Primary)

**Light Mode:**

```css
background: linear-gradient(to right, #2563eb, #06b6d4, #2563eb);
```

**Dark Mode:**

```css
background: linear-gradient(to right, #60a5fa, #22d3ee, #60a5fa);
```

### Background Gradients

**Light Mode:**

```css
background: linear-gradient(
	to bottom right,
	#f8fafc,
	rgba(239, 246, 255, 0.3),
	rgba(207, 250, 254, 0.3)
);
```

**Dark Mode:**

```css
background: linear-gradient(to bottom right, #020617, #020617, #0f172a);
```

### Button Gradients

**Primary (Green):**

```css
background: linear-gradient(135deg, #22c55e, #16a34a);
```

**Download (Blue):**

```css
background: linear-gradient(135deg, #3b82f6, #2563eb);
```

---

## Component Specifications

### Cards

**Base Card:**

- Background: `white/80` (light), `gray-900/80` (dark)
- Border: 1px solid `gray-200` (light), `gray-700` (dark)
- Border radius: `1rem` (16px)
- Padding: `2rem` (32px)
- Backdrop filter: `blur-sm`
- Box shadow: `0 1px 3px rgb(0 0 0 / 0.1)`

**Hover State:**

- Border color: `blue-300` (light), `cyan-600` (dark)
- Shadow: `0 6px 10px rgba(59, 130, 246, 0.1)`
- Transform: `translateY(-2px)`

### Buttons

**Primary Button:**

- Padding: `0.75rem 1.5rem`
- Background: Green gradient
- Color: White
- Border radius: `0.75rem`
- Font weight: 600
- Shadow: `0 4px 6px rgb(0 0 0 / 0.1)`
- Transition: `all 0.2s`

**Primary Button Hover:**

- Background: Darker green gradient
- Shadow: `0 6px 8px rgb(0 0 0 / 0.15)`
- Transform: `translateY(-2px)`

**Secondary Button:**

- Padding: `0.5rem 0.875rem`
- Background: `blue-500`
- Color: White
- Border radius: `0.5rem`
- Font weight: 500

**Icon Button:**

- Width/Height: `2.5rem`
- Border: 1px solid `gray-200` (light), `gray-700` (dark)
- Border radius: `0.5rem`
- Background: `gray-50` (light), `gray-800` (dark)

### Badges

**Category Badge:**

- Padding: `0.375rem 0.75rem`
- Background: `blue-100` (light), `blue-900/30` (dark)
- Color: `blue-700` (light), `blue-300` (dark)
- Border: 1px solid `blue-200` (light), `blue-700` (dark)
- Border radius: `0.5rem`
- Font size: `0.875rem`
- Font weight: 600

**Author Badge (Official):**

- Background: `blue-100` (light), `blue-900/40` (dark)
- Color: `blue-700` (light), `blue-300` (dark)
- Border: 1px solid `blue-300` (light), `blue-700` (dark)

**Author Badge (Community):**

- Background: `green-100` (light), `green-900/40` (dark)
- Color: `green-700` (light), `green-200` (dark)
- Border: 1px solid `green-300` (light), `green-700` (dark)

### Input Fields

**Text Input:**

- Padding: `0.875rem 1rem` (vertical), `3rem` (horizontal with icon)
- Background: `white` (light), `gray-900/50` (dark)
- Border: 1px solid `gray-200` (light), `gray-700` (dark)
- Border radius: `0.5rem`
- Font size: `0.875rem`

**Focus State:**

- Outline: None
- Ring: 2px solid `blue-500` (light), `cyan-500` (dark)
- Border: Transparent

**Select/Dropdown:**

- Same as text input
- Additional chevron icon on right

### Tables

**Header:**

- Background: `gray-50` (light), `gray-900/50` (dark)
- Padding: `1rem 1.5rem`
- Font size: `0.75rem`
- Font weight: 700
- Color: `blue-600` (light), `cyan-400` (dark)
- Text transform: Uppercase
- Letter spacing: Wide

**Row:**

- Background: `white` (light), `gray-950` (dark)
- Border bottom: 1px solid `gray-200` (light), `gray-800` (dark)
- Padding: `1rem 1.5rem`

**Row Hover:**

- Background: `blue-50/50` (light), `gray-800/30` (dark)
- Border left: 4px solid `blue-500` (light), `cyan-500` (dark)
- Shadow: `0 1px 2px rgb(0 0 0 / 0.05)`
- Cursor: Pointer

### Code Blocks

**Background:**

- Light: `gray-50` or `gray-900` (for terminal)
- Dark: `gray-800` or `black` (for terminal)

**Text Color:**

- Light: `blue-700`
- Dark: `blue-300` or `gray-100`

**Border:**

- 1px solid `gray-200` (light), `gray-700` (dark)

**Border Radius:** `0.5rem` to `0.75rem`

**Font:**

- Family: `'Monaco', 'Menlo', 'Ubuntu Mono', monospace`
- Size: `0.875rem`

### Modals

**Backdrop:**

- Background: `rgba(0, 0, 0, 0.7)`
- Z-index: 1000

**Modal Container:**

- Background: `white` (light), `gray-800` (dark)
- Border radius: `1rem`
- Padding: `2rem`
- Max width: `42rem`
- Shadow: `0 20px 25px -5px rgb(0 0 0 / 0.1)`

**Close Button:**

- Size: `2.5rem × 2.5rem`
- Background: `gray-100` (light), `gray-700` (dark)
- Border radius: `0.5rem`
- Position: Absolute top-right

### Stats Cards

**Container:**

- Background: `white` (light), `gray-900` (dark)
- Border: 1px solid `gray-200` (light), `gray-800` (dark)
- Border radius: `0.75rem`
- Padding: `1.25rem`

**Hover:**

- Border color: `blue-400` (light), `blue-600` (dark)
- Transform: `translateY(-4px)`
- Shadow: `0 10px 15px rgba(59, 130, 246, 0.1)`

**Icon Container:**

- Padding: `0.625rem`
- Background: `gray-100` (light), `gray-800` (dark)
- Border radius: `0.5rem`
- Color: `blue-600` (light), `cyan-400` (dark)

---

## Animations & Transitions

### Standard Transition

```css
transition: all 0.2s ease;
```

### Slide In From Top

```css
@keyframes slide-in-from-top {
	from {
		opacity: 0;
		transform: translateY(-0.5rem);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
```

**Usage:**

- Duration: `0.2s`
- Fill mode: `both`

### Hover Transforms

- **Lift**: `translateY(-2px)` to `translateY(-4px)`
- **Scale**: `scale(1.05)`

### Loading Spinner

```css
@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}
```

**Spinner Styles:**

- Size: `3rem × 3rem`
- Border: `4px solid`
- Border color: `gray-200` (light), `gray-700` (dark)
- Border top color: `blue-500` (light), `blue-300` (dark)
- Border radius: Full
- Animation: `spin 0.8s linear infinite`

---

## Scrollbar Styling

### Width/Height

- `8px`

### Track

- Background: `gray-100` (light), `gray-900` (dark)

### Thumb

- Background: `gray-300` (light), `gray-700` (dark)
- Border radius: `0.5rem`

### Thumb Hover

- Background: `gray-400` (light), `gray-600` (dark)

---

## Layout Patterns

### Container

```css
max-width: 1152px; /* 72rem */
margin: 0 auto;
padding: 1rem; /* 16px on mobile */
```

### Grid Layouts

**Two Columns (Desktop):**

```css
display: grid;
grid-template-columns: repeat(2, 1fr);
gap: 1rem;
```

**Three Columns (Desktop):**

```css
display: grid;
grid-template-columns: repeat(3, 1fr);
gap: 1.5rem;
```

**Responsive Breakpoints:**

- Mobile: `< 640px` - Single column
- Tablet: `640px - 768px` - Two columns
- Desktop: `> 768px` - Three columns

### Header (Fixed)

```css
position: fixed;
top: 0;
left: 0;
right: 0;
z-index: 50;
padding: 1rem;
```

**Spacer Below Header:**

```css
height: 5rem; /* 80px */
```

---

## Icon Specifications

**Library:** `lucide-svelte`

### Common Sizes

- Small: `w-4 h-4` (16px)
- Medium: `w-5 h-5` (20px)
- Large: `w-6 h-6` (24px)
- Extra Large: `w-16 h-16` (64px - hero icons)

### Common Icons

- **Navigation**: `Home`, `List`, `Menu`, `X`, `ArrowLeft`, `ArrowRight`
- **Actions**: `Download`, `Copy`, `Search`, `Filter`, `ChevronDown`
- **Status**: `Clock`, `Github`, `Moon`, `Lightbulb`
- **Branding**: `Brain`, `Zap`

---

## Responsive Design Rules

### Mobile First Approach

All base styles target mobile. Desktop enhancements use `@media (min-width: ...)`.

### Breakpoints

- `sm`: `640px`
- `md`: `768px`
- `lg`: `1024px`
- `xl`: `1280px`

### Common Responsive Patterns

**Hidden on Mobile, Visible on Desktop:**

```html
<div class="hidden md:block">...</div>
```

**Visible on Mobile, Hidden on Desktop:**

```html
<div class="md:hidden">...</div>
```

**Stacked on Mobile, Grid on Desktop:**

```html
<div class="space-y-4 md:grid md:grid-cols-2 md:gap-4">...</div>
```

---

## Accessibility Guidelines

### Color Contrast

- All text meets WCAG AA standards (4.5:1 for normal text, 3:1 for large text)
- Interactive elements have clear focus states (2px ring)

### Focus States

```css
focus:outline-none
focus:ring-2
focus:ring-blue-500 /* light */
dark:focus:ring-cyan-500 /* dark */
focus:border-transparent
```

### Interactive Elements

- All buttons and links have visible hover states
- Keyboard navigation supported (tab index, enter/space keys)
- ARIA labels for icon-only buttons

### Text Readability

- Minimum font size: `0.875rem` (14px)
- Maximum line length: `65-75 characters` (prose)
- Adequate line height: `1.5` minimum

---

## State Variations

### Loading States

- Display spinner with message
- Disable interactive elements
- Reduce opacity to `0.7`

### Error States

- Red border: `border-red-500`
- Red text: `text-red-600` (light), `text-red-400` (dark)

### Success States

- Green border: `border-green-500`
- Green text: `text-green-600` (light), `text-green-400` (dark)

### Active/Selected States

- Background: `blue-100` (light), `blue-950` (dark)
- Border: `blue-300` (light), `blue-700` (dark)
- Text: `blue-700` (light), `blue-300` (dark)

---

## Special Components

### Header Navigation

**Desktop:**

- Fixed top bar with logo, nav links, and theme toggle
- Glassmorphism effect (backdrop blur + opacity)
- Border: 1px solid with opacity
- Border radius: `1rem`
- Shadow: Large

**Mobile:**

- Hamburger menu
- Collapsible navigation drawer
- Full-width menu items

### Search Bar

- Icon: `Search` on left (absolute positioned)
- Padding left: `3rem` to accommodate icon
- Clear button: `X` icon on right (conditional)
- Active indicator badge when search term present

### Filter Panel

- Collapsible with `ChevronDown` rotation
- Toggle button shows active filter count
- Dropdown selects for category and author
- Clear all button at bottom

### Sequence Cards

- Gradient title with brand colors
- Category badge top-right
- Author name with metadata
- Relative time display with `Clock` icon
- Hover state with border accent and shadow

### Download Modal

- Step-by-step instructions
- Numbered circles for steps
- Code blocks with copy buttons
- Countdown timer for auto-download
- Platform-specific instructions (Windows/Unix)

---

## Brand Voice & Tone

### Visual Language

- **Modern**: Clean, minimal interfaces with subtle gradients
- **Technical**: Monospace fonts for code, precise spacing
- **Trustworthy**: Consistent color usage, clear hierarchy
- **Approachable**: Soft shadows, rounded corners, smooth transitions

### Copy Guidelines

- **Concise**: Short, clear sentences
- **Helpful**: Provide context and guidance
- **Professional**: Technical but not intimidating
- **Friendly**: Welcoming to contributors

---

## Implementation Notes

### Tailwind Configuration

The design system is implemented using Tailwind CSS v4 with custom theme configuration in `@theme` block within `app.css`.

### Dark Mode

Dark mode is controlled via the `.dark` class on the `<html>` element. Theme preference is stored in `localStorage` and respects system preference as fallback.

### Component Structure

All components follow the Svelte 5 runes syntax:

- `$state` for reactive state
- `$derived` for computed values
- `$effect` for side effects
- `$props` for component props

### File Organization

```
src/
  lib/
    components/    # Reusable UI components
    store.ts       # Global state management
    types.ts       # TypeScript interfaces
  routes/          # Page components
  app.css          # Global styles + Tailwind
```

---

## Version History

- **v1.0.0** - Initial design system specification

---

## Maintenance

This design specification should be updated whenever:

- New colors are added to the palette
- Component patterns are established
- Interaction behaviors change
- Accessibility improvements are made

For questions or suggestions, please open an issue in the [SynapSeq repository](https://github.com/ruanklein/synapseq).
