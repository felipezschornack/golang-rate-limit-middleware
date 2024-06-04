import { check } from 'k6';
import http from 'k6/http';
import { SharedArray } from 'k6/data';

export const API_HOST = 'http://rate-limiter'
export const API_PORT = '8080'

const keys = new SharedArray('api_keys.json', function () {
    return JSON.parse(open('./api_keys_1.json')).api;
  });

export const options = {
    scenarios: {
        my_scenario1: {
          executor: 'constant-arrival-rate',
          duration: '30s',
          preAllocatedVUs: 50,    
          rate: 50,
          timeUnit: '1s',
        },
      },
}

export default function () {
    
    const endpoint = API_HOST + ':' + API_PORT;
    const response_1 = http.get(endpoint);
    check(response_1, {
        'IP: status 200': (r) => r.status === 200,
        'IP: status 429': (r) => r.status === 429,
    });

    const apiKey = keys[Math.floor(Math.random() * keys.length)].key;
    const response_2 = http.get(endpoint, { headers: { API_KEY: apiKey } });
    check(response_2, {
        'API_KEY: status 200': (r) => r.status === 200,
        'API_KEY: status 429': (r) => r.status === 429,
    });
}