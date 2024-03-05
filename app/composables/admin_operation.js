import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, handler) {
    // get nuxt hooks
    const app = useNuxtApp();
    const url = useState("urls");

    // Initialize object
    super(url.value.socket_url, session_id, handler);

    // Initialize custom attribute
    this.app = app;
  }

  getAddress(currentObj = this) {
    return currentObj.socket_url + "/admin/arrange/" + currentObj.identifier;
  }

  quizStartRequest() {
    this.sendMessage(this.currentComponent, this.app.$StartQuiz);
  }
}
