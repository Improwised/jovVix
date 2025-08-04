import http from "k6/http";
import ws from "k6/ws";
import { check, sleep } from "k6";

export const options = {
  vus: 400,
  duration: "300s",
};

export default function () {
  const vuID = __VU;
  const uniqueUsername = `uou_${vuID}`;
  const quizCode = "923852";

  function retryHttpPost(url, payload = null, params = {}, retries = 3) {
    for (let i = 0; i < retries; i++) {
      const res = http.post(url, payload, params);
      if (res.status === 200) return res;
      sleep(0.5 + Math.random()); // slight backoff
    }
    return null;
  }

  const authRes = retryHttpPost(
    `https://quiz.i8d.in/api/v1/user/${uniqueUsername}?avatar_name=Andrea`
  );
  const cookie = authRes?.cookies["user"]?.[0]?.value;

  if (!cookie) {
    console.log(`VU ${vuID}: Failed to retrieve session cookie`);
    return;
  }

  const registerRes = retryHttpPost(
    `https://quiz.i8d.in/api/v1/user_played_quizes/${quizCode}`,
    null,
    {
      headers: {
        "Cookie": `user=${cookie}`,
        "Content-Type": "application/json",
      },
    }
  );

  if (!registerRes) {
    console.log(`VU ${vuID}: Failed to register to quiz`);
    return;
  }

  let sessionID, currentQuiz;
  try {
    const body = JSON.parse(registerRes.body);
    sessionID = body?.data?.session_id;
    currentQuiz = body?.data?.user_played_quiz;
  } catch (e) {
    console.log(`VU ${vuID}: Failed to parse register response`);
    return;
  }

  if (!sessionID || sessionID.length !== 36 || !currentQuiz) {
    console.log(`VU ${vuID}: Invalid session_id or quiz`);
    return;
  }

  sleep(Math.random() * 2); // slight jitter

  const url = `wss://quiz.i8d.in/api/v1/socket/join/${quizCode}?quiz_id=${currentQuiz}&session_id=${sessionID}`;
  const params = {
    tags: { name: "QuizWebSocketTest" },
    headers: {
      "Cookie": `user=${cookie}`,
    },
  };

  let reconnectAttempts = 0;
  const maxReconnects = 2;

  function connectWebSocket() {
    return ws.connect(url, params, function (socket) {
      socket.on("open", function () {
        console.log(`VU ${vuID}: Connected`);

        socket.setInterval(function () {
          socket.send(JSON.stringify({ event: "ping" }));
        }, 10000);
      });

      socket.on("message", function (data) {
        let message;
        try {
          message = JSON.parse(data);
        } catch (e) {
          console.log(`VU ${vuID}: Invalid message`);
          return;
        }

        const questionID = message?.data?.data?.data?.id;
        if (!questionID || questionID.length !== 36) return;

        sleep(Math.random() * 2);

        const answerPayload = {
          id: questionID,
          keys: [2],
          response_time: 2000,
        };

        const answerRes = http.post(
          `https://quiz.i8d.in/api/v1/quiz/answer?user_played_quiz=${currentQuiz}&session_id=${sessionID}`,
          JSON.stringify(answerPayload),
          {
            headers: {
              "Cookie": `user=${cookie}`,
              "Content-Type": "application/json",
            },
          }
        );

        if (answerRes.status !== 200) {
          console.log(`VU ${vuID}: Answer failed`);
        }
      });

      socket.on("error", function (e) {
        console.log(`VU ${vuID}: WS Error - ${e.error()}`);
      });

      socket.on("close", function () {
        reconnectAttempts++;
        if (reconnectAttempts > maxReconnects) {
          console.log(`VU ${vuID}: Max reconnects reached`);
          return;
        }

        const delay = 1 + Math.random() * 2;
        console.log(`VU ${vuID}: Reconnecting in ${delay.toFixed(1)}s`);
        sleep(delay);
        connectWebSocket();
      });
    });
  }

  const wsRes = connectWebSocket();
  check(wsRes, { "WebSocket connected": (r) => r && r.status === 101 });

  sleep(999); // keep alive
}
