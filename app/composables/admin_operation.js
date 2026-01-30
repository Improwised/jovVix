import constants from "~~/config/constants";
import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler, errorHandler, skipConfirmHandler) {
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
    this.connect();
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
  }

  requestPauseQuiz(isPause) {
    this.sendMessage(this.currentComponent, constants.PauseQuiz, isPause);
  }
}
