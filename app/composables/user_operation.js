import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(url, code, username, handler) {
    const cookie = useCookie("user");
    console.log(url, username, code, handler, cookie);
    super(url, username, code, handler, cookie);
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
