// user_operations.js

import QuizHandler from "./quiz_operation";
import constants from "~~/config/constants";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler, errorHandler, successHandler) {
    const url = useRuntimeConfig().public;
    super(
      url.apiSocketUrl + "/join/" + code + "?username=" + username,
      code,
      handler
    );
    this.errorHandler = errorHandler;
    this.successHandler = successHandler;
    this.apiUrl = url.apiUrl;
    this.pingInterval = null;
    this.isWaiting = false;

    this.pingIntervalTime = 45000;
    this.sendErrorMessage = false;
    this.reconnect = false;

    this.connectUser();
    this.startPing();
  }

  connectUser() {
    this.connect(this);
  }

  // Method to start pinging through WebSocket
  startPing() {
    if (!this.pingInterval) {
      this.pingInterval = setInterval(() => {
        this.pingServer();
      }, this.pingIntervalTime);
    }
  }

  pingServer() {
    if (this.socket.readyState === WebSocket.OPEN) {
      console.log("pinging server");
      this.socket.send(
        JSON.stringify({ event: constants.EventPing, user: this.username })
      );
    }
  }

  stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  onClose() {
    if (this.isWaiting) {
      this.stopPing();
      this.handleConnectionProblem();
      this.startPing();
    }
  }

  handleConnectionProblem() {
    this.errorHandler();
    this.reconnect = true;
    this.connectUser();
  }

  onOpen(event) {
    if (this.reconnect) {
      this.successHandler();
    }
    super.onOpen(event);
  }

  async handleSendAnswer(answers, user_played_quiz, session_id) {
    const responseTime = this.getAnswerResponseTime();
    try {
      const response = await fetch(
        `${this.apiUrl}/quiz/answer?user_played_quiz=${user_played_quiz}&session_id=${session_id}`,
        {
          method: "POST",
          body: JSON.stringify({
            id: this.currentQuestion,
            keys: answers || [],
            response_time: responseTime,
          }),
          credentials: "include",
          mode: "cors",
        }
      );

      if (response.status !== 202) {
        const errorMessage = await response.text();
        return { error: `Failed to submit answer: ${errorMessage}` };
      }

      return { error: null };
    } catch (err) {
      return { error: err.message || "An unknown error occurred." };
    }
  }

  endQuiz() {
    this.close(1000);
    this.sendMessage("", "websocket_close", "");
  }

  onMessage(event) {
    const message = this.destructureMessage(event);

    if (message.component === constants.Waiting) {
      this.isWaiting = true;
    } else if (this.isWaiting && message.component === constants.Question) {
      console.log("stopping ping");
      this.stopPing();
      this.isWaiting = false;
    }

    super.onMessage(event);
  }
}
