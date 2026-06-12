# Responsive Web UI

## Purpose

This spec defines the mobile-first responsive design requirements for the web application. _Initial baseline established from make-web-responsive-for-mobile change._

## Requirements

### Requirement: Mobile-First Breakpoint System
The web application SHALL follow a mobile-first responsive design using Tailwind CSS default breakpoints (`sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`, `2xl: 1536px`). All pages MUST be usable and visually correct at viewport widths from 320px up.

#### Scenario: Page renders correctly on small mobile
- **WHEN** a user opens any public page (home, events list, event detail, login, register) on a device with viewport width 320px
- **THEN** all content fits within the viewport without horizontal scroll, all text is readable, and primary actions are tappable

#### Scenario: Layout adapts at tablet breakpoint
- **WHEN** viewport width reaches 768px
- **THEN** multi-column grids expand to show 2 columns where applicable and navigation transitions from hamburger to horizontal menu

#### Scenario: Desktop layout uses full width
- **WHEN** viewport width is 1024px or greater
- **THEN** content is constrained to a max-width container (e.g. `max-w-7xl`) with appropriate gutters and 3+ column grids where applicable

### Requirement: Responsive Navigation
The site MUST provide a navigation pattern that works on mobile (hamburger menu) and desktop (horizontal menu) without losing access to any menu item.

#### Scenario: Mobile hamburger opens drawer
- **WHEN** a user on viewport ≤ 640px taps the hamburger icon in the navbar
- **THEN** a slide-in drawer (or full-screen overlay) opens showing all navigation links and the user avatar/login link

#### Scenario: Hamburger closes on link click
- **WHEN** a user taps any navigation link inside the mobile drawer
- **THEN** the drawer closes and the user is navigated to the selected page

#### Scenario: Desktop shows inline navigation
- **WHEN** viewport width is ≥ 768px
- **THEN** all primary navigation links are visible inline in the navbar and no hamburger icon is shown

### Requirement: Touch-Friendly Interactive Elements
All clickable elements (buttons, links, form inputs, card actions) MUST have a minimum tap target size of 44×44px on mobile viewports to comply with WCAG 2.5.5 and Apple/Google touch guidelines.

#### Scenario: Buttons have minimum tap target on mobile
- **WHEN** a user views any page on a viewport ≤ 640px
- **THEN** all buttons and links have a height and width of at least 44px (via padding or `min-h-[44px]`)

#### Scenario: Form inputs are easily tappable
- **WHEN** a user opens login, register, booking, or payment forms on mobile
- **THEN** all input fields have height ≥ 44px and adequate spacing between fields (≥ 12px)

### Requirement: Responsive Tables on Admin Pages
Admin tables (events, bookings, payments, users) MUST be usable on mobile viewports either via horizontal scroll within a constrained container, a card-based mobile view, or a column-prioritised display.

#### Scenario: Table scrollable on small screens
- **WHEN** an admin user views the events, bookings, payments, or users list on viewport ≤ 640px
- **THEN** the table container allows horizontal scroll without breaking the page layout, and the first column (e.g. ID or name) remains visible

#### Scenario: Filter controls collapse on mobile
- **WHEN** admin pages include filters (date range, status, search)
- **THEN** those filters collapse into a toggleable accordion or "Filter" sheet on viewport ≤ 640px and expand inline at ≥ 768px

### Requirement: Responsive War Tiket Flow
The critical war tiket pages (war landing, queue waiting, booking form) MUST be optimized for one-handed mobile use with sticky CTAs and large countdown timers.

#### Scenario: Sticky CTA on booking page
- **WHEN** a user is on `/events/[id]/booking` on mobile
- **THEN** a sticky bottom action bar (e.g. "Lanjut Bayar") is visible and does not overlap form content

#### Scenario: Countdown timer legible on mobile
- **WHEN** a user is on the queue or war page on mobile
- **THEN** the countdown timer occupies a prominent area with font size ≥ 32px and updates without layout shift

#### Scenario: Queue position is glanceable
- **WHEN** a user is on `/events/[id]/queue` on mobile
- **THEN** the user's current position and total queue length are visible above the fold without scrolling

### Requirement: Safe Area and Viewport Configuration
The application MUST set a proper viewport meta tag and respect device safe areas (notches, home indicators) on all pages.

#### Scenario: Viewport meta tag is correct
- **WHEN** any page is loaded
- **THEN** the document head includes `<meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">`

#### Scenario: Sticky elements respect safe area
- **WHEN** a user views pages with sticky header, sticky bottom bar, or full-height content on a device with a notch or home indicator
- **THEN** those elements have appropriate padding (e.g. `pt-safe`, `pb-safe` or `env(safe-area-inset-*)`) so content is not obscured

### Requirement: Responsive Typography and Spacing
Text and spacing MUST scale appropriately across breakpoints. Body text SHALL be at least 16px on mobile to prevent iOS auto-zoom on input focus.

#### Scenario: Body text size on mobile
- **WHEN** a user views any page on viewport ≤ 640px
- **THEN** body text has font size ≥ 16px and headings use a mobile-friendly scale (e.g. h1 ≈ 28–32px, h2 ≈ 22–24px)

#### Scenario: Form input text prevents iOS zoom
- **WHEN** a user focuses an input on iOS Safari
- **THEN** the input font size is ≥ 16px so the browser does not auto-zoom the viewport

### Requirement: Shared Mobile UI Components
The application MUST provide shared components for common mobile UI patterns to avoid duplication.

#### Scenario: Bottom sheet component available
- **WHEN** a developer needs a modal-style UI on mobile
- **THEN** a `MobileBottomSheet` component is available that slides up from the bottom on mobile and renders as a centered modal on desktop

#### Scenario: Tabs component is responsive
- **WHEN** a developer uses the `MobileTabs` component
- **THEN** it renders as horizontal scrollable tabs on mobile and centered tabs on desktop

#### Scenario: Viewport composable available
- **WHEN** a developer needs to branch logic by viewport size
- **THEN** a `useViewport()` composable returns reactive booleans (`isMobile`, `isTablet`, `isDesktop`) based on current viewport

### Requirement: Image and Media Responsiveness
All images, banners, and media MUST scale appropriately and use `srcset` or Tailwind responsive classes to avoid oversized downloads on mobile.

#### Scenario: Event banner scales to container
- **WHEN** a user views an event detail page on any viewport
- **THEN** the event banner fills the container width, maintains aspect ratio, and does not overflow

#### Scenario: Event card thumbnail optimized
- **WHEN** a user browses the events list on mobile
- **THEN** card thumbnails are served at a smaller size and use `loading="lazy"` to defer off-screen images

### Requirement: Performance Budget on Mobile
Pages on mobile MUST meet a basic performance budget: no layout shift (CLS) larger than 0.1, and interactive elements respond within 100ms of tap.

#### Scenario: No layout shift on image load
- **WHEN** a page loads on mobile
- **THEN** image containers have explicit width/height or aspect-ratio classes so images do not cause layout shift

#### Scenario: Tap feedback is immediate
- **WHEN** a user taps any button or link on mobile
- **THEN** visual feedback (active state, ripple, or color change) appears within 100ms
