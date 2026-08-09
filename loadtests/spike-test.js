import http from 'k6/http';
import { check } from 'k6';

// Spike test configuration - sudden traffic spikes
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '10s', target: 5 },    // Normal traffic
    { duration: '1s', target: 100 },  // SPIKE! 5 -> 100 users in 1s
    { duration: '10s', target: 100 },  // Stay at spike level
    { duration: '5s', target: 5 },     // Return to normal
    { duration: '10s', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // More lenient during spike
    http_req_failed: ['rate<0.1'],     // Allow 10% errors during spike
  },
};

const USER_IDS = ['user-001', 'user-002', 'user-003', 'user-004', 'user-005'];

export default function () {
  const userId = USER_IDS[Math.floor(Math.random() * USER_IDS.length)];
  
  // Test: Health check during spike
  const healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
  });

  // Test: Get recommendations during spike
  const recRes = http.get(`${BASE_URL}/recommendations/${userId}`);
  check(recRes, {
    'recommendations status is 200': (r) => r.status === 200,
    'recommendations response time < 2s': (r) => r.timings.duration < 2000,
  });
}
