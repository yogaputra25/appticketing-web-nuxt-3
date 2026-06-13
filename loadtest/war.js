import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const API_BASE = __ENV.API_URL || 'http://localhost:8080';

const warJoinDuration = new Trend('war_join_duration');
const warStatusDuration = new Trend('war_status_duration');
const bookingReserveDuration = new Trend('booking_reserve_duration');
const warJoinErrors = new Rate('war_join_errors');
const bookingErrors = new Rate('booking_errors');

export const options = {
  stages: [
    { duration: '10s', target: 50 },
    { duration: '20s', target: 100 },
    { duration: '30s', target: 200 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.05'],
    war_join_errors: ['rate<0.1'],
    booking_errors: ['rate<0.1'],
  },
};

// Shared: register a test user and return JWT
function registerUser(idx) {
  const email = `waruser${idx}@test.com`;
  const res = http.post(`${API_BASE}/api/auth/register`, JSON.stringify({
    email,
    password: 'test123',
    full_name: `War User ${idx}`,
  }), { headers: { 'Content-Type': 'application/json' } });
  if (res.status === 200 || res.status === 201) {
    return res.json().token || res.json().data?.token;
  }
  return null;
}

// Login and return JWT
function loginUser(idx) {
  const email = `waruser${idx}@test.com`;
  const res = http.post(`${API_BASE}/api/auth/login`, JSON.stringify({
    email,
    password: 'test123',
  }), { headers: { 'Content-Type': 'application/json' } });
  if (res.status === 200) {
    return res.json().token;
  }
  return null;
}

// Fetch published events, return first event ID
function fetchFirstEvent(token) {
  const res = http.get(`${API_BASE}/api/events?limit=1`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (res.status === 200) {
    const body = res.json();
    if (body.data && body.data.length > 0) {
      return body.data[0].id;
    }
  }
  return null;
}

export default function () {
  const idx = __VU; // virtual user index

  // 1. Login
  let token = loginUser(idx);
  if (!token) {
    token = registerUser(idx);
    if (!token) {
      warJoinErrors.add(1);
      return;
    }
  }

  // 2. Fetch event to war on
  const eventId = fetchFirstEvent(token);
  if (!eventId) {
    warJoinErrors.add(1);
    return;
  }

  // 3. Join war queue
  {
    const start = Date.now();
    const res = http.post(`${API_BASE}/api/war/join`, JSON.stringify({ event_id: eventId }), {
      headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
    });
    warJoinDuration.add(Date.now() - start);
    if (res.status !== 200 && res.status !== 429) {
      warJoinErrors.add(1);
    }
  }

  // 4. Poll queue status for a bit
  for (let i = 0; i < 5; i++) {
    const start = Date.now();
    const res = http.get(`${API_BASE}/api/war/status?event_id=${eventId}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    warStatusDuration.add(Date.now() - start);
    if (res.status === 200) {
      const body = res.json();
      if (body.is_ready) {
        break;
      }
    }
    sleep(1);
  }

  // Simulate a booking attempt after queue
  {
    const start = Date.now();
    const res = http.post(`${API_BASE}/api/bookings/reserve`, JSON.stringify({
      event_id: eventId,
      category_id: null, // will fail gracefully
      quantity: 1,
    }), {
      headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
    });
    bookingReserveDuration.add(Date.now() - start);
    if (res.status >= 500) {
      bookingErrors.add(1);
    }
  }

  sleep(1);
}
