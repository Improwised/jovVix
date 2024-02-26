import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(url, code, username) {
    const cookie = useCookie("user");
    super(url, code, username, cookie);
    console.log(code, username);
  }

  getDefaultState(currentObj = this) {
    throw Error("not implement yet", currentObj);
  }
}
