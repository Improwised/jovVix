import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler, errorHandler) {
    const url = useState("urls");
    super(url.value.socket_url, code, handler, { username });
    this.errorHandler = errorHandler;
    this.api_url = url.value.api_url;
    console.log(url.value.api_url);
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

  async handleSendAnswer(answers) {
    let error;
    try {
      const response = await useFetch(this.api_url + "/quiz/answer", {
        method: "POST",
        body: JSON.stringify({
          id: this.currentQuestion,
          keys: answers || [],
        }),
        credentials: "include",
        mode: "cors",
      });
      error = response.error.value;
    } catch (err) {
      console.log(err, "---------------------------------");
      error = err;
    }

    return { error };
  }
}
