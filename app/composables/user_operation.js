import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler) {
    const app = useNuxtApp();
    const url = useState("urls");
    const cookie = useCookie(app.$UserIdentifier);
    super(url.value.socket_url, username, code, handler, cookie);
  }

  getAddress(currentObj) {
    return (
      currentObj.api_url +
      "/join/" +
      currentObj.identifier +
      "?username=" +
      currentObj.username
    );
  }

  getDefaultState(currentObj = this) {
    throw Error("not implement yet", currentObj);
  }
}
