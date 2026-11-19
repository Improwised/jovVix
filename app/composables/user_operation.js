// user_operations.js

import QuizHandler from "./quiz_operation";
import constants from "~~/config/constants";

export default class UserOperation extends QuizHandler {
  constructor(
    code,
    username,
    handler,
    errorHandler,
    successHandler,
    sessionUnavailableHandler,
    reconnectHandler
  ) {
    const url = useRuntimeConfig().public;
    super(
      url.apiSocketUrl + "/join/" + code + "?username=" + username,
      code,
      handler
    );
    this.errorHandler = errorHandler;
    this.successHandler = successHandler;
    // Fires when the socket opens but is closed by the server before we ever
    // reach the Waiting room (i.e. the session no longer exists or has been
    // terminated). The page uses this to bounce the player back to /join.
    // shortcut: kept for API compat but no longer invoked — room recovery handles it.
    this.reconnectHandler = reconnectHandler;
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
    this.stopPing();
    if (!this.disposed) this.handleConnectionProblem();
  }

  handleConnectionProblem() {
    if (this.disposed || this.reconnectTimer) return;
    this.reconnectAttempts = (this.reconnectAttempts || 0) + 1;
    this.errorHandler();
    this.reconnect = true;
    const delay = Math.min(
      2000 * Math.pow(1.5, this.reconnectAttempts - 1),
      30000
    );
    this.reconnectTimer = setTimeout(async () => {
      this.reconnectTimer = null;
      try {
        if (this.reconnectHandler && !(await this.reconnectHandler())) return;
        if (this.disposed) return;
        this.connectUser();
        this.startPing();
      } catch {
        this.handleConnectionProblem();
      }
    }, delay);
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
        return { error: `Failed to submit answer` };
      }

      return { error: null };
    } catch (err) {
      return { error: err.message || "An unknown error occurred." };
    }
  }

  // shortcut: no longer called — the dispose() path covers clean-up.
  // endQuiz() { this.close(1000); this.sendMessage("", "websocket_close", ""); }

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
