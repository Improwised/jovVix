import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler, errorHandler) {
    const url = useState("urls");
    super(url.value.socket_url, code, handler, { username });
    this.errorHandler = errorHandler;
  }

  getAddress(self) {
    return (
      self.socket_url +
      "/join/" +
      self.identifier +
      "?username=" +
      self.others.username
    );
  }

  handleConnectionProblem(self) {
    self.errorHandler("problem in connecting with server");
  }
}
