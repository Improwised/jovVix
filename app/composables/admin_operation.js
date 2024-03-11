import constants from "~~/config/constants";
import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler, errorHandler) {
    // get nuxt hooks
    const url = useState("urls");

    // Initialize object
    super(
      url.value.socket_url + "/admin/arrange/" + session_id,
      session_id,
      handler
    );

    // Initialize custom attribute
    this.api_url = url.value.api_url;
    this.errorHandler = errorHandler;
  }

  quizStartRequest() {
    this.sendMessage(this.currentComponent, constants.StartQuiz);
  }

  handleConnectionProblem() {
    this.errorHandler("problem in connecting with server");
  }

  async handler(message) {
    if (this.currentEvent == constants.NextQuestionAsk) {
      this.sendMessage(this.currentComponent, this.currentEvent);
    }

    super.handler(message, preventAssignment);
  }
}
