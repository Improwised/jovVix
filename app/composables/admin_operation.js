import QuizHandler from "./quiz_operation";

export default class AdminOperations extends QuizHandler {
  constructor(url, session_id, username, handler) {
    // get nuxt hooks
    const cookie = useCookie("user");

    // Initialize object
    super(url, username, session_id, handler, cookie);
  }

  getAddress(currentObj = this) {
    return currentObj.api_url + "/admin/arrange/" + currentObj.identifier;
  }

  quizStartRequest() {
    this.sendMessage();
  }
}
