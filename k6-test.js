import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Counter } from 'k6/metrics';
import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';

export const options = {
  stages: [
    { duration: '10s', target: 20 },  // Ramp-up
    { duration: '60s', target: 50 },  // Peak 50 virtual users
    { duration: '10s', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],    // Errors under 1%
    http_req_duration: ['p(95)<200'],  // 95% of requests under 200ms
  },
  summaryTrendStats: [
    'min',
    'max',
    'avg',
    'med',
    'p(90)',
    'p(95)',
    'p(99)',
  ],
};

// Custom metrics counters for each endpoint
const healthCounter = new Counter('health_requests');
const userProfileCounter = new Counter('user_profile_requests');
const bannerPublicCounter = new Counter('banner_public_requests');
const bannerUserCounter = new Counter('banner_user_requests');
const accountGetCounter = new Counter('account_get_requests');
const accountSetMainCounter = new Counter('account_set_main_requests');
const debitCardGetCounter = new Counter('debit_card_get_requests');
const debitCardStatusCounter = new Counter('debit_card_status_requests');
const transactionCounter = new Counter('transaction_requests');

const BASE_URL = __ENV.BASE_URL || 'http://127.0.0.1:8080';
const USER_ID = __ENV.USER_ID || '000018b0e1a211ef95a30242ac180002';

const jsonHeaders = {
  headers: { 'Content-Type': 'application/json' },
};

export default function () {
  // 1. Health API
  group('Health API', () => {
    const res = http.get(`${BASE_URL}/api/health`);
    check(res, {
      'health status is 200': (r) => r.status === 200,
    });
    healthCounter.add(1);
  });

  // 2. User Profile API
  group('User Profile API', () => {
    const res = http.get(`${BASE_URL}/api/user/${USER_ID}`,
      { tags: { endpoint: 'GET_UserProfile'},
    });
    check(res, {
      'user status is 200': (r) => r.status === 200,
    });
    userProfileCounter.add(1);
  });

  // 3. Banner APIs (Public & User-targeted)
  group('Banner API', () => {
    const publicRes = http.get(`${BASE_URL}/api/banner/`,
      { tags: { endpoint: 'GET_Banner'},
    });
    check(publicRes, {
      'public banner status is 200': (r) => r.status === 200,
    });
    bannerPublicCounter.add(1);

    const userBannerRes = http.get(`${BASE_URL}/api/banner/${USER_ID}`,
      { tags: { endpoint: 'GET_UserBanner'} }
    );
    check(userBannerRes, {
      'user banner status is 200': (r) => r.status === 200,
    });
    bannerUserCounter.add(1);
  });

  // 4. Account APIs (GET accounts & PUT set main account)
  group('Account API', () => {
    let accountId = '000020ece1a211ef95a30242ac180002';
    const getRes = http.get(`${BASE_URL}/api/account/${USER_ID}`,
      { tags: { endpoint: 'GET_Account'} }
    );
    check(getRes, {
      'account get status is 200': (r) => r.status === 200,
    });
    accountGetCounter.add(1);

    try {
      const data = getRes.json().data;
      if (Array.isArray(data) && data.length > 0 && data[0].account_id) {
        accountId = data[0].account_id;
      }
    } catch (_) {}

    const putPayload = JSON.stringify({
      user_id: USER_ID,
      account_id: accountId,
    });

    const putRes = http.put(`${BASE_URL}/api/account/main`, putPayload, {
      ...jsonHeaders,
      tags: { endpoint: 'PUT_UpdateMain' },
    });
    check(putRes, {
      'account set main status is 200': (r) => r.status === 200,
    });
    accountSetMainCounter.add(1);
  });

  // 5. Debit Card APIs (GET cards & PUT update status)
  group('Debit Card API', () => {
    let cardId = '000018cfe1a211ef95a30242ac180002';
    const getRes = http.get(`${BASE_URL}/api/debit-card/${USER_ID}`,
      { tags: { endpoint: 'GET_DebitCard'} }
    );
    check(getRes, {
      'debit card get status is 200': (r) => r.status === 200,
    });
    debitCardGetCounter.add(1);

    // Extract first card dynamically if available
    try {
      const data = getRes.json().data;
      if (Array.isArray(data) && data.length > 0 && data[0].card_id) {
        cardId = data[0].card_id;
      }
    } catch (_) {}

    const putPayload = JSON.stringify({
      user_id: USER_ID,
      card_id: cardId,
      status: 'Active',
    });
    const putRes = http.put(`${BASE_URL}/api/debit-card/status`, putPayload, {
      ...jsonHeaders,
      tags: { endpoint: 'PUT_UpdateStatusDebitCard' },
    });
    check(putRes, {
      'debit card update status is 200': (r) => r.status === 200,
    });
    debitCardStatusCounter.add(1);
  });

  // 6. Transaction API (GET paginated)
  group('Transaction API', () => {
    const res = http.get(`${BASE_URL}/api/transaction/${USER_ID}?page=1&limit=10`,
      { tags: { endpoint: 'GET_Transaction' },
    });
    check(res, {
      'transaction status is 200': (r) => r.status === 200,
    });
    transactionCounter.add(1);
  });

  sleep(1);
}

// Generate HTML and text summary reports
export function handleSummary(data) {
  return {
    'reports/summary.html': htmlReport(data),
    'reports/summary.json': JSON.stringify(data, null, 2),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}

// export function handleSummary(data) {
//   const metrics = data.metrics;
//   const endpoints = [
//     'GET_UserProfile',
//     'GET_Banner',
//     'GET_UserBanner',
//     'GET_Account',
//     'PUT_UpdateMain',
//     'GET_DebitCard',
//     'PUT_UpdateStatusDebitCard',
//     'GET_Transaction',
//   ];

//   let report = '\n===== API RESPONSE TIME REPORT =====\n';

//   for (const endpoint of endpoints) {
//     const metric = metrics.http_req_duration;
//     const result = metric.values;

//     report += `\n${endpoint}\n`;
//     report += `Min : ${result.min} ms\n`;
//     report += `Max : ${result.max} ms\n`;
//     report += `Avg : ${result.avg} ms\n`;
//     report += `P90 : ${result['p(90)']} ms\n`;
//     report += `P95 : ${result['p(95)']} ms\n`;
//   }

//   return {
//     stdout: report,
//   };
// }
