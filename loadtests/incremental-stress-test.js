import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // Warm up - 10 VUs
    { duration: '30s', target: 50 },   // Ramp up to 50 VUs
    { duration: '30s', target: 100 },  // Ramp up to 100 VUs
    { duration: '30s', target: 150 },  // Ramp up to 150 VUs
    { duration: '30s', target: 200 },  // Ramp up to 200 VUs
    { duration: '30s', target: 200 },  // Sustain at 200 VUs
    { duration: '20s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% under 1 second
    http_req_failed: ['rate<0.05'],    // Error rate under 5%
  },
};

// Test data
const USER_IDS = ['user-001', 'user-002', 'user-003', 'user-004', 'user-005'];
const PRODUCT_IDS = ['PROD-001', 'PROD-002', 'PROD-003', 'PROD-004', 'PROD-005'];

export default function () {
  // Test 1: Health check
  let healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
  });

  // Test 2: Get recommendations
  const userId = USER_IDS[Math.floor(Math.random() * USER_IDS.length)];
  const recRes = http.get(`${BASE_URL}/recommendations/${userId}`);
  check(recRes, {
    'recommendations status is 200': (r) => r.status === 200,
  });

  // Test 3: Process event
  const eventPayload = JSON.stringify({
    event_id: `event-${Date.now()}-${__VU}`,
    event_type: getRandomEventType(),
    user_id: userId,
    product_id: PRODUCT_IDS[Math.floor(Math.random() * PRODUCT_IDS.length)],
    product_category: getRandomCategory(),
    product_brand: getRandomBrand(),
    metadata: {
      device: getRandomDevice(),
      country: getRandomCountry(),
    },
    occurred_at: new Date().toISOString(),
  });

  const eventRes = http.post(`${BASE_URL}/events`, eventPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(eventRes, {
    'event processing status is 202': (r) => r.status === 202,
  });

  sleep(0.5); // Reduced sleep for higher throughput
}

// Helper functions
function getRandomEventType() {
  const types = ['product_viewed', 'search_performed', 'product_added_cart', 'product_purchased'];
  return types[Math.floor(Math.random() * types.length)];
}

function getRandomCategory() {
  const categories = ['electronics', 'clothing', 'home', 'books', 'sports'];
  return categories[Math.floor(Math.random() * categories.length)];
}

function getRandomBrand() {
  const brands = ['Apple', 'Samsung', 'Nike', 'Adidas', 'Sony', 'LG'];
  return brands[Math.floor(Math.random() * brands.length)];
}

function getRandomDevice() {
  const devices = ['mobile', 'desktop', 'tablet'];
  return devices[Math.floor(Math.random() * devices.length)];
}

function getRandomCountry() {
  const countries = ['US', 'ES', 'UK', 'DE', 'FR', 'IT'];
  return countries[Math.floor(Math.random() * countries.length)];
}
