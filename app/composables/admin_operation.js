import constants from "~~/config/constants";
import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(
    session_id,
    handler,
    errorHandler,
    skipConfirmHandler,
    reconnectHandler
  ) {
    const url = useRuntimeConfig().public;

    // Initialize object
    super(
      url.apiSocketUrl + "/admin/arrange/" + session_id,
      session_id,
      handler
    );

    // Initialize custom attribute
    this.apiUrl = url.apiUrl;
    this.errorHandler = errorHandler;
    this.skipHandler = skipConfirmHandler;
    this.reconnectHandler = reconnectHandler;
    this.pingIntervalTime = 45000;
    this.pingInterval = null;
    this.isWaiting = true;

    this.startPing();
  }

  continueAdmin() {
    this.continue(this);
  }

  connectAdmin() {
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

  quizStartRequest() {
    this.sendMessage(this.currentComponent, constants.StartQuiz);
  }

  requestSkip(force = false) {
    this.sendMessage(
      this.currentComponent,
      force ? constants.AskForceSkip : constants.AskSkip
    );
  }

  requestSkipTimer() {
    this.sendMessage(this.currentComponent, constants.SkipTimer);
  }

  handleConnectionProblem() {
    this.errorHandler("problem in connecting with server");
    if (this.disposed || this.reconnectTimer) return;
    this.reconnectAttempts = (this.reconnectAttempts || 0) + 1;
    const delay = Math.min(
      2000 * Math.pow(1.5, this.reconnectAttempts - 1),
      30000
    );
    this.reconnectTimer = setTimeout(async () => {
      this.reconnectTimer = null;
      try {
        if (this.reconnectHandler && !(await this.reconnectHandler())) return;
        if (this.disposed) return;
        this.connect();
        this.startPing();
      } catch {
        this.handleConnectionProblem();
      }
    }, delay);
  }

  async handler(message) {
    if (this.currentEvent == constants.NextQuestionAsk) {
      this.sendMessage(this.currentComponent, this.currentEvent);
    } else if (this.currentEvent == constants.AskSkip) {
      this.skipHandler(message);
    } else if (this.currentEvent === constants.Counter && this.isWaiting) {
      this.isWaiting = false;
    }
    super.handler(message);
  }

  requestTerminateQuiz() {
    this.close(1000);
  }

  onClose(event) {
    console.log("stoping ping of admin");
    this.stopPing();
    super.onClose(event);
    setSocketObject(null);
    if (!this.disposed) this.handleConnectionProblem();
  }

  requestPauseQuiz(isPause) {
    this.sendMessage(this.currentComponent, constants.PauseQuiz, isPause);
  }
}
