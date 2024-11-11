import { check } from "k6";
import http from "k6/http";

// 1. init code
export function setup() {
  // 2. Setup code
}
export const options = {
  thresholds: {
    http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true }], // http errors should be less than 1%; abort after first failure
    http_req_duration: ["p(99)<1000"], // 99% of requests should be below 1s
  },
  scenarios: {
    average_load: {
      executor: "ramping-vus",
      //startRate: 50,
      //timeUnit: "1s",
      //preAllocatedVUs: 50,
      //maxVUs: 100,
      stages: [
        //// ramp up to average load of 20 virtual users
        //{ duration: "10s", target: 20 },
        //// maintain load
        //{ duration: "50s", target: 20 },
        //// ramp down to zero
        //{ duration: "5s", target: 0 },
        { duration: "10s", target: 20 },
        { duration: "50s", target: 20 },
        { duration: "50s", target: 40 },
        { duration: "50s", target: 60 },
        { duration: "50s", target: 80 },
        { duration: "50s", target: 100 },
        { duration: "50s", target: 120 },
        { duration: "50s", target: 140 },
      ],
    },
  },
};

// 3. default function; run the tests
export default function () {
  // define URL and payload
  const url = "http://localhost:8000/api/user/all";
  const payload = JSON.stringify({
    username: "Johne Doe",
    email: "",
    password: "test",
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const res = http.get(url, payload, params);
  //console.log(res.body);
  check(res, {
    "response code was 200": (res) => res.status == 200,
  });
}

//export default function () {
//  // define URL and payload
//  const url = "http://localhost:8000/api/user/";
//  const payload = JSON.stringify({
//    username: "Johne Doe",
//    email: "test@test.com",
//        password: "test",
//  });
//
//  const params = {
//    headers: {
//      "Content-Type": "application/json",
//    },
//  };
//
//  const res = http.post(url, payload, params);
//  console.log(res.body);
//}
//

export function teardown() {
  // 4. teardown code
}
