// user_operations.js

import QuizHandler from "./quiz_operation";
import constants from "~~/config/constants";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler, errorHandler, successHandler) {
    const url = useRuntimeConfig().public;

    super(
      url.socket_url + "/join/" + code + "?username=" + username,
      code,
      handler
    );
    this.errorHandler = errorHandler;
    this.successHandler = successHandler;
    this.api_url = url.api_url;
    this.pingInterval = null;
    this.isWaiting = false;

    this.pingIntervalTime = 5000;
    this.sendErrorMessage = false;
    this.reconnect = false;

    this.startPing();
  }

  // Method to start pinging through WebSocket
  startPing() {
    if (!this.pingInterval) {
      this.pingInterval = setInterval(() => {
        this.pingServer();
      }, this.pingIntervalTime);
    }
  }

  // Method to send a ping through WebSocket
  pingServer() {
    if (this.socket.readyState === WebSocket.OPEN) {
      console.log("pinging server");
      this.socket.send(JSON.stringify({ event: "ping", user: this.username }));
    } else {
      this.handleConnectionProblem();
    }
  }

  handleConnectionProblem() {
    this.errorHandler();
    this.reconnect = true;
    this.pingIntervalTime = 10000;
    this.connect();
  }

  onOpen(event) {
    if (this.reconnect) {
      this.successHandler();
    }
    super.onOpen(event);
  }

  async handleSendAnswer(answers) {
    let error;
    const responseTime = this.getAnswerResponseTime();
    try {
      const response = await fetch(this.api_url + "/quiz/answer", {
        method: "POST",
        body: JSON.stringify({
          id: this.currentQuestion,
          keys: answers || [],
          response_time: responseTime,
        }),
        credentials: "include",
        mode: "cors",
      });
      const responseJson = await response.json();
      error = responseJson.data;
    } catch (err) {
      error = err;
    }

    return { error };
  }

  endQuiz() {
    this.sendMessage("", "websocket_close", "");
  }

  // Method to stop pinging through WebSocket
  stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  onMessage(event) {
    const message = this.destructureMessage(event);

    if (message.component === constants.Waiting) {
      this.isWaiting = true;
      this.startPing();
    } else if (this.isWaiting && message.component === constants.Question) {
      console.log("stopping ping");
      this.stopPing();
      this.isWaiting = false;
    }

    super.onMessage(event);
  }
}
