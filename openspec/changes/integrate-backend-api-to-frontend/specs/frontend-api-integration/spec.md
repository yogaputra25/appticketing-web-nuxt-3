## ADDED Requirements

### Requirement: Consistent API Path Prefix
All frontend HTTP calls to the backend SHALL use the `/api` prefix and the paths defined in `apps/api/internal/router/router.go`.

#### Scenario: Public pages call correct endpoints
- **WHEN** a user browses public pages (home, events list, event detail)
- **THEN** the frontend calls `GET /api/events` and `GET /api/events/{id}` exactly as registered in the router

#### Scenario: My bookings page calls correct endpoint
- **WHEN** an authenticated user opens `/my/bookings`
- **THEN** the frontend calls `GET /api/bookings/me` (not `/bookings` or `/api/me/bookings`)

#### Scenario: My booking detail page calls correct endpoint
- **WHEN** an authenticated user opens `/my/bookings/{id}`
- **THEN** the frontend calls `GET /api/bookings/{id}` (not `/bookings/{id}`)

#### Scenario: Profile page calls correct endpoint
- **WHEN** an authenticated user opens `/profile`
- **THEN** the frontend calls `PUT /api/auth/me` (not `/auth/profile`)

#### Scenario: Admin pages call correct endpoints
- **WHEN** admin uses any admin page
- **THEN** all calls use the `/api/admin/...` paths and HTTP methods registered in the router (no invented endpoints)

### Requirement: Field Name Alignment with Backend Models
Frontend TypeScript types and component templates SHALL use field names that exactly match the JSON returned by the Go backend (snake_case from `apps/api/internal/model/model.go`).

#### Scenario: Booking list uses backend field names
- **WHEN** the user booking list is rendered
- **THEN** the template uses `booking.booking_code`, `booking.total_amount`, and `booking.e_ticket_codes.length` for ticket count — not `b.code`, `b.total_price`, or `b.ticket_count`

#### Scenario: Booking ID is rendered as code
- **WHEN** a booking row is displayed
- **THEN** the ID column shows `booking.booking_code` (e.g. "BK-a1b2c3") and `String(booking.id)` is not sliced (since `id` is a `number`)

#### Scenario: Queue status uses backend field names
- **WHEN** the queue polling page renders
- **THEN** it reads `position`, `total_in_queue`, `estimated_wait`, and `is_ready` from the response — and uses `session_token` (not `token`) when `is_ready` is true

#### Scenario: Reservation request body matches backend schema
- **WHEN** the user submits a booking
- **THEN** the request body sent to `POST /api/bookings/reserve` is `{ event_id, session_token, items: [{ category_id, quantity }] }` exactly as the backend handler `BookingHandler.Reserve` expects

### Requirement: Pagination Response Shape
Frontend pagination logic SHALL use the response shape `{ data, total, page, limit }` returned by the backend (not `last_page` or `total_pages`).

#### Scenario: Admin list pagination uses total + limit
- **WHEN** an admin page lists events/bookings/payments/users
- **THEN** the page computes `totalPages = Math.ceil(total / limit)` from the response, not from `res.last_page`

### Requirement: useApi Composable Used Consistently
All authenticated API calls from pages and components SHALL go through the `useApi()` composable so that the JWT token, baseURL, and 401 handling are applied uniformly.

#### Scenario: Auth token is sent on protected requests
- **WHEN** a logged-in user makes any request to a protected endpoint
- **THEN** the `Authorization: Bearer <token>` header is attached by `useApi()`

#### Scenario: 401 response triggers re-login flow
- **WHEN** the backend returns 401 for an authenticated request
- **THEN** the auth store is cleared and the user is redirected to `/login` by the `useApi()` interceptor

### Requirement: nuxt.config apiBase Configured
The `nuxt.config.ts` SHALL define `runtimeConfig.public.apiBase` pointing to the backend default `http://localhost:8080`, and `.env.example` SHALL document `NUXT_PUBLIC_API_BASE` for overrides.

#### Scenario: Default apiBase points to local backend
- **WHEN** the Nuxt app starts with no env override
- **THEN** `useRuntimeConfig().public.apiBase` resolves to `http://localhost:8080`

#### Scenario: Env override works
- **WHEN** the user sets `NUXT_PUBLIC_API_BASE=http://my-api:8080` in `.env`
- **THEN** all `$fetch` calls from the app use that base URL

### Requirement: CORS Allows Web Origin
The Go backend SHALL allow the Nuxt dev origin (default `http://localhost:3000`) via the existing `cors.Handler` configuration so that the web can call the API in development.

#### Scenario: Web can call API in dev
- **WHEN** the Nuxt dev server at `http://localhost:3000` makes a request to the API
- **THEN** the request is not blocked by CORS preflight

### Requirement: API Errors Surface in UI
The frontend SHALL display backend error messages (from `err.data.message` or `err.message`) to the user in form error regions and toast/alert components, not just log to console.

#### Scenario: Login error is shown to user
- **WHEN** the login API returns 401 with message "Email atau password salah"
- **THEN** the login form displays that message in the red error region

#### Scenario: Booking reservation error is shown
- **WHEN** the booking reserve API returns 409 with message "Not enough stock"
- **THEN** the booking page shows that error inline near the "Lanjut Bayar" button
