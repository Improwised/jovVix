import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler) {
    const url = useState("urls");
    super(url.value.socket_url, code, handler, { username });
  }

  getAddress(currentObj) {
    return (
      currentObj.socket_url +
      "/join/" +
      currentObj.identifier +
      "?username=" +
      currentObj.others.username
    );
  }

  getDefaultState(currentObj = this) {
    throw Error("not implement yet", currentObj);
  }
}
