import http from 'k6/http';
import { check, sleep } from 'k6';

// Stress test configuration
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '5s', target: 10 },   // Ramp up to 10 users
    { duration: '10s', target: 50 },  // Ramp up to 50 users
    { duration: '20s', target: 100 }, // Ramp up to 100 users
    { duration: '30s', target: 100 }, // Stay at 100 users
    { duration: '10s', target: 50 },  // Ramp down to 50 users
    { duration: '5s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% of requests must complete below 1s
    http_req_failed: ['rate<0.05'],   // Error rate must be less than 5%
  },
};

const USER_IDS = Array.from({ length: 100 }, (_, i) => `user-${String(i + 1).padStart(3, '0')}`);

export default function () {
  const userId = USER_IDS[Math.floor(Math.random() * USER_IDS.length)];
  
  // Test: Get recommendations under stress
  const recRes = http.get(`${BASE_URL}/recommendations/${userId}`);
  
  check(recRes, {
    'recommendations status is 200': (r) => r.status === 200,
    'recommendations response time < 1s': (r) => r.timings.duration < 1000,
    'recommendations is valid JSON': (r) => {
      try {
        JSON.parse(r.body);
        return true;
      } catch (e) {
        return false;
      }
    },
  });

  sleep(0.5); // Faster iterations for stress test
}
