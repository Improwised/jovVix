import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, username, handler) {
    // get nuxt hooks
    const cookie = useCookie("user");
    const url = useState("urls");
    const app = useNuxtApp();

    // Initialize object
    super(url.value.socket_url, username, session_id, handler, cookie);

    // Initialize custom attribute
    this.app = app;
  }

  getAddress(currentObj = this) {
    return currentObj.api_url + "/admin/arrange/" + currentObj.identifier;
  }

  quizStartRequest() {
    this.sendMessage(this.currentComponent, this.app.$StartQuiz);
  }
}
