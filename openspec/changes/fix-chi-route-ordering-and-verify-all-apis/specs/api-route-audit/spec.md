## ADDED Requirements

### Requirement: Chi Router No-Fallthrough Rule
The Go API router declared in `apps/api/internal/router/router.go` SHALL follow the rule that **event-scoped admin routes** (paths that begin with `/admin/events/{id}/...`) MUST be declared inside the `r.Route("/admin/events", ...)` subroute block, NOT in the outer `r.Route("/admin", ...)` block. This rule prevents chi from returning 404 on routes that are syntactically declared but matched under a sibling subroute.

#### Scenario: Admin adds a ticket category to an event
- **WHEN** admin POSTs to `/api/admin/events/{eventId}/categories` with valid body `{name, price, total_stock, max_per_user}` and a valid admin token
- **THEN** API returns 201 Created with the new category object (id, event_id, name, price, total_stock, available_stock, max_per_user, timestamps)

#### Scenario: Admin publishes an event
- **WHEN** admin POSTs to `/api/admin/events/{id}/publish` for a draft event that has at least one category
- **THEN** API returns 200 OK with the event object whose `status` is now `published`

#### Scenario: Admin updates a category
- **WHEN** admin PUTs to `/api/admin/categories/{id}` with partial body
- **THEN** API returns 200 OK with the updated category

#### Scenario: Router fallback rule is documented
- **WHEN** a developer opens `apps/api/internal/router/router.go`
- **THEN** a comment block above the `/admin/events` subroute explains the no-fallthrough rule and instructs future maintainers to add event-scoped routes inside the `/admin/events` subroute

### Requirement: API Route Smoke Test Script
A shell script SHALL exist at `apps/api/scripts/smoke.sh` (or similar canonical location) that hits every route declared in `router.go` and prints PASS/FAIL based on the expected status code. This script SHALL be runnable manually (`bash apps/api/scripts/smoke.sh`) and from CI to detect regressions in route registration.

#### Scenario: Script lists all public, authenticated, and admin routes
- **WHEN** the script runs against a running API at `$API_BASE` (default `http://localhost:8080`)
- **THEN** it tests at minimum the following routes:
  - `GET /api/healthz` → 200
  - `GET /api/events` → 200
  - `GET /api/events/1` → 200 (or 404 if id=1 missing)
  - `GET /api/events/1/categories` → 200
  - `POST /api/auth/register` → 201 (uses random email to avoid conflict)
  - `POST /api/auth/login` → 200 (uses seeded admin)
  - `GET /api/auth/me` (with token) → 200
  - `GET /api/admin/stats` (with admin token) → 200
  - `GET /api/admin/events` (with admin token) → 200
  - `POST /api/admin/events` (with admin token) → 201
  - `POST /api/admin/events/{id}/categories` (with admin token) → 201 ← the previously broken route
  - `GET /api/war/status?event_id=1` (with any token) → 200

#### Scenario: Script distinguishes 404-page from 404-handler
- **WHEN** the script checks an endpoint and receives a 404 response
- **THEN** it inspects the response body: a plain-text body `404 page not found` (chi default) is treated as FAIL (route missing), while a JSON body `{"error":"not_found","message":"..."}` is treated as PASS (route exists, resource missing)

#### Scenario: Script exits non-zero on any failure
- **WHEN** at least one endpoint check fails
- **THEN** the script exits with a non-zero status code and prints a summary of which checks failed

### Requirement: Critical Routes Audit Snapshot
The following route inventory SHALL be considered the "audit baseline" for the project. Every route in this list MUST be present in `router.go` and respond as expected after the chi router fix.

#### Scenario: Public routes are accessible without auth
- **WHEN** an unauthenticated client requests:
  - `GET /api/events`
  - `GET /api/events/1`
  - `GET /api/events/1/categories`
- **THEN** each returns a 2xx response (or 404 for missing event id), proving the route is registered and chi routes correctly

#### Scenario: Authenticated user routes are accessible with a valid user token
- **WHEN** a logged-in user requests:
  - `GET /api/auth/me`
  - `GET /api/war/status?event_id=1`
  - `POST /api/war/join?event_id=1`
  - `GET /api/bookings/me`
- **THEN** each returns a 2xx response (or 401 if token invalid)

#### Scenario: Admin routes are accessible with a valid admin token
- **WHEN** an admin requests:
  - `GET /api/admin/stats`
  - `GET /api/admin/events`
  - `GET /api/admin/bookings`
  - `GET /api/admin/payments`
  - `GET /api/admin/users`
  - `POST /api/admin/events` (create)
  - `POST /api/admin/events/{id}/categories` (add category)
  - `POST /api/admin/events/{id}/publish` (publish)
  - `GET /api/admin/bookings/{id}` (detail)
  - `PUT /api/admin/categories/{id}` (update)
- **THEN** each returns a 2xx response (or 404 JSON for missing resources)
