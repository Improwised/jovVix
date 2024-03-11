import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler, errorHandler) {
    // get nuxt hooks
    const app = useNuxtApp();
    const url = useState("urls");

    // Initialize object
    super(url.value.socket_url + "/admin/arrange/" + session_id, session_id, handler, {url: url.value.socket_url + "/admin/arrange/" + session_id});

    // Initialize custom attribute
    this.app = app;
    this.errorHandler = errorHandler;
  }
  
  quizStartRequest() {
    this.sendMessage(this, this.currentComponent, this.app.$StartQuiz);
  }

  handleConnectionProblem(self) {
    self.errorHandler("problem in connecting with server");
  }
}
