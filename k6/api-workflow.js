import http from "k6/http";
import { check, sleep, fail } from "k6";
import { Counter, Trend } from "k6/metrics";

// Base URL and Port for API endpoints
//const BASE_URL = "http://synopsy.ddns.net";
const BASE_URL = "http://localhost";
const BASE_PORT = 8000;

let successfulRequests = new Counter("successful_requests");
let failedAttempts = new Counter("failed_attempts");

let createTrend = new Trend("create_trend");
let loginTrend = new Trend("login_trend");
let getTrend = new Trend("get_trend");
let deleteTrend = new Trend("delete_trend");

export const options = {
  thresholds: {
    //http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true }], // http errors should be less than 1%; abort after first failure
    http_req_duration: ["p(99)<1000"], // 99% of requests should be below 1s
  },
  stages: [
    { duration: "10s", target: 10 },
    { duration: "20s", target: 20 },
    { duration: "50s", target: 40 },
    { duration: "50s", target: 60 },
    { duration: "50s", target: 0 },
  ],
};

// Helper function to handle HTTP requests and validate response
function makeRequestAndCheck(
  method,
  path,
  body = null,
  headers = {},
  expectedStatus = 200,
  successMessage = "",
) {
  let url = `${BASE_URL}:${BASE_PORT}${path}`;
  let res;

  if (method === "GET") {
    res = http.get(url, { headers });
  } else if (method === "POST") {
    res = http.post(url, body, { headers });
  } else if (method === "PUT") {
    res = http.put(url, body, { headers });
  } else if (method === "DELETE") {
    res = http.del(url, null, { headers });
  }

  const isSuccess = check(res, {
    [successMessage]: (r) => r.status === expectedStatus,
  });

  if (isSuccess) {
    successfulRequests.add(1);
  } else {
    failedAttempts.add(1);
    throw new Error(`Request to ${url} failed with status ${res.status}`);
  }

  return res;
}

// Stage 1: Create a user
function stageCreateUser(username, email, password) {
  let res = makeRequestAndCheck(
    "POST",
    "/api/user",
    JSON.stringify({
      username: username,
      email: email,
      password: password,
    }),
    { "Content-Type": "application/json" },
    201,
    "User created successfully",
  );

  createTrend.add(res.timings.duration);
}

// Stage 2: Login and retrieve JWT
function stageLoginUser(username, password) {
  let res = makeRequestAndCheck(
    "POST",
    "/api/user/login",
    JSON.stringify({
      username: username,
      password: password,
    }),
    { "Content-Type": "application/json" },
    200,
    "Login successful",
  );
  loginTrend.add(res.timings.duration);
  return res.json("access_token");
}

// Stage 3: Get user information
function stageGetUser(jwtToken) {
  let res = makeRequestAndCheck(
    "GET",
    "/api/user",
    null,
    { Authorization: `Bearer ${jwtToken}` },
    200,
    "Get user successful",
  );
  getTrend.add(res.timings.duration);
}

// Stage 4: Modify user information
function stageModifyUser(jwtToken) {
  makeRequestAndCheck(
    "PUT",
    "/modify-user",
    JSON.stringify({ attribute: "newValue" }),
    {
      "Content-Type": "application/json",
      Authorization: `Bearer ${jwtToken}`,
    },
    200,
    "Modify user successful",
  );
}

// Stage 5: Delete user
function stageDeleteUser(jwtToken) {
  let res = makeRequestAndCheck(
    "DELETE",
    "/api/user",
    null,
    { Authorization: `Bearer ${jwtToken}` },
    200,
    "Delete user successful",
  );

  deleteTrend.add(res.timings.duration);
}

// Stage 6: Verify user deletion
function stageVerifyUserDeletion(jwtToken) {
  let res = http.get(`${BASE_URL}:${BASE_PORT}/api/user`, {
    headers: { Authorization: `Bearer ${jwtToken}` },
  });
  let userCheckFailed = check(res, {
    "Get user failed as expected": (r) => r.status === 500,
  });

  if (!userCheckFailed) {
    failedAttempts.add(1);
    retryDeleteUser(jwtToken);
  }
}

// Retry logic for deleting user
function retryDeleteUser(jwtToken) {
  for (let i = 0; i < 2; i++) {
    try {
      stageDeleteUser(jwtToken);
      break; // Exit loop if retry is successful
    } catch (e) {
      failedAttempts.add(1);
      if (i === 1) {
        fail("Stop test after 2 failed retry attempts");
      }
    }
  }
}

// Main function coordinating the stages
export default function () {
  let uuid = Math.random().toString(36).substring(7);
  let username = `John Doe ${uuid}`;
  let email = `john.doe${uuid}@icloud.com`;
  let password = "test";
  try {
    stageCreateUser(username, email, password);
    let jwtToken = stageLoginUser(username, password);
    stageGetUser(jwtToken);
    //stageModifyUser(jwtToken);
    stageDeleteUser(jwtToken);
    //stageVerifyUserDeletion(jwtToken);
  } catch (e) {
    console.error(e.message);
  }

  //sleep(1); // Delay before next iteration for realistic load simulation
}
