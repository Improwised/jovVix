import constants from "~~/config/constants";

export default class QuizHandler {
  constructor(socket_url, identifier, componentHandler, others) {
    if (!(socket_url && identifier && componentHandler)) {
      throw Error("all demanded parameters are necessary");
    }

    if (typeof componentHandler !== "function") {
      throw Error("component handler must be a function");
    }

    // general attributes
    this.others = others;
    this.identifier = identifier;
    this.componentHandler = componentHandler;

    // states and log
    this.currentComponent = null;
    this.currentEvent = null;
    this.log = [];
    this.isOpen = false;
    this.retrying = 0;

    // handle question
    this.currentQuestion = null;
    this.currentQuestionGetTime = null;

    // custom attributes
    this.socket = null;
    this.socket_url = socket_url;
    this.connect(this);
  }

  // connect and set websocket
  connect() {
    console.log("web socket connect() called");
    this.socket = new WebSocket(this.socket_url);
    this.socket.onopen = (event) => this.onOpen(event);
    this.socket.onerror = (event) => this.onError(event);
    this.socket.onclose = (event) => this.onClose(event);
    this.socket.onmessage = (event) => this.onMessage(event);
    this.close = (code) => this.socket.close(code);
  }

  // handle open event
  onOpen(event) {
    this.isOpen = true;
    this.retrying = 0;
    this.log.push({
      state: "Init",
      message: event,
      time: new Date().toLocaleString(),
    });
  }

  // destructure up-coming messages
  destructureMessage(e) {
    try {
      const obj = JSON.parse(e.data);
      obj["event"] = obj.data.event;
      delete obj.data["event"];
      obj["action"] = obj.data.data.action;
      delete obj.data.data["action"];
      obj["component"] = obj.data.data.component;
      delete obj.data.data["component"];
      obj["data"] = obj.data.data.data;
      return obj;
    } catch (err) {
      this.onError(err);
    }
  }

  // onmessage handler ans setup self object
  onMessage(event) {
    const message = this.destructureMessage(event);
    if (
      this.currentEvent == constants.GetQuestion &&
      message.event == constants.AdminDisconnected
    ) {
      this.componentHandler(message);
    } else {
      this.currentComponent = message.component;
      this.currentEvent = message.event;
      this.log.push({
        state: "Receive",
        message,
        time: new Date().toLocaleString(),
      });
      this.handler(message, event);
    }
  }

  // pre-process message if needed and assign to outside function
  async handler(message) {
    if (
      this.currentComponent == constants.Question &&
      this.currentEvent == constants.GetQuestion
    ) {
      this.currentQuestion = message.data.id;
      this.currentQuestionGetTime = new Date();
    } else if (
      this.currentComponent == constants.Score &&
      this.currentEvent == constants.ShowScore
    ) {
      const options = message.data.options;
      const answers = message.data.answers;

      const new_options = {};

      for (const key in options) {
        const element = options[key];
        let isAnswer = false;
        for (let answerIndex = 0; answerIndex < answers.length; answerIndex++) {
          if (answers[answerIndex] === parseInt(key)) {
            isAnswer = true;
            break;
          }
        }

        new_options[key] = { value: element, isAnswer };
      }
      delete message.data.answers;
      message.data.options = new_options;
    } else {
      this.currentQuestionGetTime = null;
    }
    if (message.event == constants.TerminateQuiz) {
      await this.handleTerminate();
    }
    this.componentHandler(message);
  }

  getAnswerResponseTime() {
    if (this.currentQuestionGetTime != null) {
      return new Date() - this.currentQuestionGetTime;
    }
  }

  // send message through socket
  sendMessage(
    component = this.currentComponent,
    event = this.currentEvent,
    data = ""
  ) {
    if (this.socket.readyState == WebSocket.CLOSED) {
      this.handleConnectionProblem(this);
    }
    const message = {
      component: component,
      event: event,
      data: data,
    };
    try {
      this.socket.send(JSON.stringify(message));
    } catch (err) {
      console.error(err);
    } finally {
      this.log.push({
        state: "Sent",
        message,
        time: new Date().toLocaleString(),
      });
    }
  }

  // print full communication if needed
  printLog() {
    this.log.forEach((message) => {
      if (message.state == "Sent") {
        console.warn(message);
      } else if (message.state == "Receive") {
        console.warn(message);
      } else if (message.state == "Init") {
        console.warn(message);
      } else {
        console.error(message);
      }
    });
  }

  // handle error and re-connecting
  onError(event) {
    // websocket can not connect to the server
    if (this.socket.readyState == WebSocket.CLOSED && !this.isOpen) {
      // sent alert
      this.handleConnectionProblem();

      console.log("in on error");
      console.log(this.socket.readyState);



      // check if retrying process is undergoing
      if (this.retrying == 0) {
        console.log("retrying: ", this.retrying);
        // check for every 2 second
        const id = setInterval(() => {
          if (this.retrying < 3 && !this.isOpen) {
            this.retrying += 1;
            console.log("connection called");
            this.connect();
          } else {
            clearInterval(id);
            this.retrying = -1;
          }
        }, 2000);
      }
    }
    this.log.push({
      state: "err",
      message: event,
      time: new Date().toLocaleString(),
    });
  }

  // handle reconnecting problem - this function needs to override by the child
  handleConnectionProblem() {
    console.log("connection problem, retrying", this.retrying);
  }

  // handle close event
  onClose(event) {
    this.log.push({
      state: "Init",
      message: event,
      time: new Date().toLocaleString(),
    });
    this.printLog()
  }

  async handleTerminate() {
    let error;
    try {
      const response = await useFetch(this.api_url + "/quiz/terminate", {
        method: "GET",
        credentials: "include",
        mode: "cors",
      });
      error = response.error.value;
    } catch (err) {
      error = err;
    }

    return { error };
  }
}
