## ADDED Requirements

### Requirement: Auth middleware must not redirect during SSR

Nuxt middleware (`auth.ts` and `admin.ts`) SHALL skip authentication checks during server-side rendering. The `plugins/auth.client.ts` plugin handles token restoration from localStorage on client hydration, after which the middleware runs correctly on client-side navigation.

#### Scenario: User refreshes an authenticated page
- **WHEN** an authenticated user refreshes `/admin/events`
- **THEN** the server renders the page without redirecting to `/login`
- **AND** the client restores auth state from localStorage and revalidates via `fetchMe()`

#### Scenario: Unauthenticated user accesses protected route (client-side)
- **WHEN** an unauthenticated user navigates to `/my/bookings` via client-side navigation
- **THEN** the middleware redirects to `/login?redirect=/my/bookings` (unchanged behavior)

### Requirement: Login page auto-redirects authenticated users

The `/login` page SHALL redirect authenticated users to the home page (or redirect query param) on mount, preventing logged-in users from seeing the login form.

#### Scenario: Authenticated user visits /login
- **WHEN** an authenticated user navigates to `/login`
- **THEN** the page redirects to `/` immediately

#### Scenario: Authenticated user visits /login?redirect=/admin/events
- **WHEN** an authenticated user navigates to `/login?redirect=/admin/events`
- **THEN** the page redirects to `/admin/events`

### Requirement: Booking detail includes event and user data

The `GET /api/bookings/{id}` endpoint SHALL return the `Booking` resource with `Event` and `User` associations preloaded, so the frontend can display `event.title` and other nested fields.

#### Scenario: Admin or owner views booking detail
- **WHEN** a user requests `GET /api/bookings/1`
- **THEN** the response includes `event: { id, title, venue, start_date, end_date }` and `user: { id, email, full_name }`
