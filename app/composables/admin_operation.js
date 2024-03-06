import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler, errorHandler) {
    // get nuxt hooks
    const app = useNuxtApp();
    const url = useState("urls");

    // Initialize object
    super(url.value.socket_url, session_id, handler);

    // Initialize custom attribute
    this.app = app;
    this.errorHandler = errorHandler;
  }

  getAddress(self = this) {
    return self.socket_url + "/admin/arrange/" + self.identifier;
  }

  quizStartRequest() {
    this.sendMessage(this.currentComponent, this.app.$StartQuiz);
  }

  handleConnectionProblem(self) {
    self.errorHandler("problem in connecting with server");
  }
}
