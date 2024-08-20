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

  async handleSendAnswer(answers, user_played_quiz) {
    let error;
    const responseTime = this.getAnswerResponseTime();
    try {
      const response = await fetch(
        `${this.api_url}/quiz/answer?user_played_quiz=${user_played_quiz}`,
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
      const responseJson = await response.json();
      error = responseJson.data;
    } catch (err) {
      error = err;
    }

    return { error };
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
      this.isWaiting = false;
    }

    super.onMessage(event);
  }
}
