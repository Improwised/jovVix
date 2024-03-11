import QuizHandler from "./quiz_operation";

export default class UserOperation extends QuizHandler {
  constructor(code, username, handler, errorHandler) {
    const url = useState("urls");
    super(url.value.socket_url +
      "/join/" +
      code +
      "?username=" +
      username, code, handler, { username, url: 
      url.value.socket_url +
      "/join/" +
      code +
      "?username=" +
      username
    });
    this.errorHandler = errorHandler;
    this.api_url = url.value.api_url;
    console.log(url.value.api_url);
  }

  handleConnectionProblem(self) {
    self.errorHandler("problem in connecting with server");
  }

  async handleSendAnswer(answers) {
    let error;
    const responseTime = this.getAnswerResponseTime();
    try {
      const response = await fetch(this.api_url + "/quiz/answer", {
        method: "POST",
        body: JSON.stringify({
          id: this.currentQuestion,
          keys: answers || [],
          response_time: responseTime,
        }),
        credentials: "include",
        mode: "cors",
      });
      const responseJson = await response.json();
      error = responseJson.data;
    } catch (err) {
      error = err;
    }

    return { error };
  }
}
