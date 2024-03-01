import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(session_id, username, handler) {
    // get nuxt hooks
    const app = useNuxtApp();
    const cookie = useCookie(app.$UserIdentifier);
    const url = useState("urls");

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
