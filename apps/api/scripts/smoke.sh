#!/usr/bin/env bash
# ============================================================================
# Smoke Test Script — War Tiket API
# ============================================================================
# HITS every critical endpoint and verifies response status.
# Distinguishes chi 404-page (plain text, route missing) from JSON 404-handler.
# Exits 0 iff all tests PASS.
#
# Usage:
#   bash apps/api/scripts/smoke.sh
#
# Override defaults:
#   API_BASE=http://localhost:8080 bash apps/api/scripts/smoke.sh
#   ADMIN_EMAIL=admin@example.com ADMIN_PASSWORD=admin123 ...
# ============================================================================

set -euo pipefail

API_BASE="${API_BASE:-http://localhost:8080}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@example.com}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin123}"
USER_EMAIL="${USER_EMAIL:-smoke-$(date +%s)@test.local}"
USER_PASSWORD="${USER_PASSWORD:-smoketest123}"

PASS=0
FAIL=0
FAILURES=()

# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------
green()  { printf "\033[32m%s\033[0m\n" "$1"; }
red()    { printf "\033[31m%s\033[0m\n" "$1"; }
bold()   { printf "\033[1m%s\033[0m\n" "$1"; }

check_404_page() {
  local body="$1"
  # chi 404-page: plain text "404 page not found"
  # JSON handler: starts with "{" 
  if [[ "$body" == *"404 page not found"* ]]; then
    return 0
  fi
  return 1
}

run_test() {
  local method="$1"
  local path="$2"
  local expected="$3"
  local label="$4"
  local token="$5"
  local extra_args="$6"

  local url="${API_BASE}${path}"
  local args=(-s -S -w "\n%{http_code}" --max-time 10)

  if [ -n "$token" ]; then
    args+=(-H "Authorization: Bearer $token")
  fi
  if [ -n "$extra_args" ]; then
    # shellcheck disable=SC2206
    args+=($extra_args)
  fi

  local response
  response=$(curl "${args[@]}" -X "$method" "$url" 2>/dev/null || true)

  local http_code
  http_code=$(echo "$response" | tail -1)
  local response_body
  response_body=$(echo "$response" | sed '$d')

  if [ -z "$http_code" ]; then
    red "  FAIL — connection error (is API running at $API_BASE?)"
    FAIL=$((FAIL + 1))
    FAILURES+=("$label — connection error")
    return
  fi

  # Detect chi 404-page (plain text) vs JSON 404-handler
  if [ "$http_code" = "404" ]; then
    if check_404_page "$response_body"; then
      red "  FAIL — chi 404 page not found (route missing): $method $path"
      FAIL=$((FAIL + 1))
      FAILURES+=("$label — chi 404 page")
      return
    fi
  fi

  # Check expected status — accept both single code and comma-separated list
  local ok=0
  IFS=',' read -ra expected_codes <<< "$expected"
  for code in "${expected_codes[@]}"; do
    code="$(echo "$code" | xargs)"
    if [ "$http_code" = "$code" ]; then
      ok=1
      break
    fi
  done

  if [ "$ok" = "1" ]; then
    green "  PASS — $http_code ($label)"
    PASS=$((PASS + 1))
  else
    red "  FAIL — expected $expected, got $http_code ($label)"
    FAIL=$((FAIL + 1))
    FAILURES+=("$label — expected $expected, got $http_code")
  fi
}

# ---------------------------------------------------------------------------
# 1. Login as admin → get admin token
# ---------------------------------------------------------------------------
bold "=== Admin Login ==="
ADMIN_TOKEN=$(curl -s -S -X POST "${API_BASE}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}" \
  --max-time 10 | grep -o '"token":"[^"]*"' | cut -d'"' -f4 || true)

if [ -z "$ADMIN_TOKEN" ]; then
  red "FAIL — could not obtain admin token (is API running and seeded?)"
  exit 1
fi
green "  Admin token obtained"

# ---------------------------------------------------------------------------
# 2. Register a new user → get user token
# ---------------------------------------------------------------------------
bold "=== User Registration ==="
USER_TOKEN=$(
  curl -s -S -X POST "${API_BASE}/api/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${USER_EMAIL}\",\"password\":\"${USER_PASSWORD}\",\"full_name\":\"Smoke Tester\",\"phone\":\"08123456789\"}" \
    --max-time 10 | grep -o '"token":"[^"]*"' | cut -d'"' -f4 || true
)
if [ -z "$USER_TOKEN" ]; then
  red "FAIL — could not register user"
  exit 1
fi
green "  User registered: $USER_EMAIL"

# ---------------------------------------------------------------------------
# 3. Create event + category for live tests (reuse across test cases)
# ---------------------------------------------------------------------------
bold "=== Setup: Create event & category ==="
EVENT_ID=$(
  curl -s -S -X POST "${API_BASE}/api/admin/events" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{\"title\":\"Smoke Test Event $(date +%s)\",\"venue\":\"Smoke Venue\",\"description\":\"Auto-created by smoke.sh\",\"start_date\":\"$(date -u +%Y-%m-%d -d '+30 days' 2>/dev/null || echo '2099-12-31')\",\"end_date\":\"$(date -u +%Y-%m-%d -d '+31 days' 2>/dev/null || echo '2099-12-31')\"}" \
    --max-time 10 | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2 || true
)
if [ -z "$EVENT_ID" ]; then
  red "FAIL — could not create event"
  exit 1
fi
green "  Event created: id=$EVENT_ID"

CATEGORY_ID=$(
  curl -s -S -X POST "${API_BASE}/api/admin/events/${EVENT_ID}/categories" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{"name":"Smoke Category","price":50000,"total_stock":100,"max_per_user":4}' \
    --max-time 10 | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2 || true
)
if [ -z "$CATEGORY_ID" ]; then
  red "FAIL — could not create category (chi fix regression?)"
  exit 1
fi
green "  Category created: id=$CATEGORY_ID"

# Publish event
PUBLISH_STATUS=$(
  curl -s -o /dev/null -w "%{http_code}" -X POST \
    "${API_BASE}/api/admin/events/${EVENT_ID}/publish" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    --max-time 10 || true
)
if [ "$PUBLISH_STATUS" != "200" ]; then
  red "FAIL — could not publish event (got $PUBLISH_STATUS)"
  exit 1
fi
green "  Event published"

# ---------------------------------------------------------------------------
# 4. Run test cases
# ---------------------------------------------------------------------------
bold "=== Smoke Tests ==="

# --- Public endpoints ---
run_test GET  /api/healthz                 200           "healthz"                    ""      ""
run_test GET  /api/events                  200           "public event list"          ""      ""
run_test GET  "/api/events/${EVENT_ID}"    200,404       "public event detail"        ""      ""
run_test GET  "/api/events/${EVENT_ID}/categories" 200  "public event categories"    ""      ""

# --- Auth endpoints ---
run_test POST /api/auth/login              200           "admin login (token check)"  ""      "-H \"Content-Type: application/json\" -d \"{\\\"email\\\":\\\"${ADMIN_EMAIL}\\\",\\\"password\\\":\\\"${ADMIN_PASSWORD}\\\"}\""
run_test GET  /api/auth/me                 200           "auth me (admin)"            "$ADMIN_TOKEN" ""
run_test GET  /api/auth/me                 200           "auth me (user)"             "$USER_TOKEN"  ""

# --- Admin endpoints ---
run_test GET  /api/admin/stats             200           "admin stats"                "$ADMIN_TOKEN" ""
run_test GET  /api/admin/events            200           "admin event list"           "$ADMIN_TOKEN" ""
run_test POST /api/admin/events            201           "admin create event"         "$ADMIN_TOKEN" "-H \"Content-Type: application/json\" -d '{\"title\":\"Smoke Aux $(date +%s)\",\"venue\":\"Aux\",\"start_date\":\"2099-12-31\",\"end_date\":\"2100-01-01\"}'"
run_test POST "/api/admin/events/${EVENT_ID}/categories" 201 "admin add category (regression check)" "$ADMIN_TOKEN" "-H \"Content-Type: application/json\" -d '{\"name\":\"Aux Cat\",\"price\":25000,\"total_stock\":50,\"max_per_user\":2}'"
run_test POST "/api/admin/events/${EVENT_ID}/publish"  200  "admin publish event"      "$ADMIN_TOKEN" ""
run_test GET  "/api/admin/bookings/${EVENT_ID}" 200,404 "admin booking detail"       "$ADMIN_TOKEN" ""
run_test GET  /api/admin/bookings         200           "admin booking list"         "$ADMIN_TOKEN" ""
run_test GET  /api/admin/payments         200           "admin payment list"         "$ADMIN_TOKEN" ""
run_test GET  /api/admin/users            200           "admin user list"            "$ADMIN_TOKEN" ""

# --- War queue endpoints ---
run_test POST "/api/war/join?event_id=${EVENT_ID}" 200 "war join"                    "$USER_TOKEN" ""
run_test GET  "/api/war/status?event_id=${EVENT_ID}" 200 "war status"                "$USER_TOKEN" ""

# --- Booking endpoints ---
run_test GET  /api/bookings/me            200           "my bookings"                "$USER_TOKEN" ""

# --- Payment endpoints ---
run_test GET  /api/payments/me            200           "my payments"                "$USER_TOKEN" ""

# --- 404-handler test (route exists, resource missing) ---
run_test GET  /api/admin/bookings/999999   404           "admin booking 404 (JSON)"   "$ADMIN_TOKEN" ""

# --- Auth failure test ---
run_test GET  /api/auth/me                401           "auth me (no token)"         ""      ""

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
bold "=============================="
bold "        Smoke Test Results"
bold "=============================="
echo "  PASS: $PASS"
echo "  FAIL: $FAIL"

if [ "$FAIL" -gt 0 ]; then
  echo ""
  bold "Failed tests:"
  for f in "${FAILURES[@]}"; do
    red "  - $f"
  done
  echo ""
  red "Some tests FAILED. Check API status and try again."
  exit 1
fi

green "All tests PASSED!"
exit 0
