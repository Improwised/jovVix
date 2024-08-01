import constants from "~~/config/constants";
import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler, errorHandler, skipConfirmHandler) {
    const url = useRuntimeConfig().public;

    // Initialize object
    super(url.socket_url + "/admin/arrange/" + session_id, session_id, handler);

    // Initialize custom attribute
    this.api_url = url.api_url;
    this.errorHandler = errorHandler;
    this.skipHandler = skipConfirmHandler;
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
    }

    super.handler(message, preventAssignment);
  }

  onClose(event) {
    // Check the close event code to determine if it was an error or proper closure
    if (event.code !== 1000 && this.currentEvent != constants.TerminateQuiz) {
      // 1000 indicates a normal closure
      console.log("Closed due to error, retrying...");
      this.connect();
    } else {
      console.log("Closed Properly");
    }

    super.onClose(event);
  }
}
