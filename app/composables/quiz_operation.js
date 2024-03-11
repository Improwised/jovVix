import constants from "~~/config/constants";

export default class QuizHandler {
  JOIN_EVENT = "JOIN_REQUEST";

  constructor(socket_url, identifier, componentHandler, others) {
    if (!(socket_url && identifier && componentHandler)) {
      throw Error("all demanded parameters are necessary");
    }

    if (typeof componentHandler !== "function") {
      throw Error("component handler must be a function");
    }

    // general attributes
    this.socket_url = socket_url;
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
    this.connect(this);
  }

  // get address
  getAddress(self = this) {
    return self.socket_url;
  }

  // connect and set websocket
  connect(self = this) {
    self.socket = new WebSocket(self.getAddress(self));
    self.socket.onopen = (event) => self.onOpen(self, event);
    self.socket.onerror = (event) => self.onError(self, event);
    self.socket.onclose = (event) => self.onClose(self, event);
    self.socket.onmessage = (event) => self.onMessage(self, event);
  }

  // handle open event
  onOpen(self = this, event) {
    self.isOpen = true;
    self.retrying = 0;
    self.log.push({
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
  onMessage(self = this, event) {
    const message = self.destructureMessage(event);
    self.currentComponent = message.component;
    self.currentEvent = message.event;
    self.log.push({
      state: "Receive",
      message,
      time: new Date().toLocaleString(),
    });
    self.handler(self, message, event);
  }

  // pre-process message if needed and assign to outside function
  async handler(self, message) {
    if (
      self.currentComponent == constants.Question &&
      self.currentEvent == constants.GetQuestion
    ) {
      self.currentQuestion = message.data.id;
      self.currentQuestionGetTime = new Date();
    } else {
      self.currentQuestionGetTime = null;
    }

    if (message.event == constants.TerminateQuiz) {
      await self.handleTerminate();
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
    self = this,
    component = this.currentComponent,
    event = this.currentEvent,
    data = ""
  ) {
    if (self.socket.readyState == WebSocket.CLOSED) {
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
      console.err(err);
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
        console.log(message);
      } else {
        console.error(message);
      }
    });
  }

  // handle error and re-connecting
  onError(self = this, event) {
    // websocket can not connect to the server
    if (self.socket.readyState == WebSocket.CLOSED && !self.isOpen) {
      // sent alert
      self.handleConnectionProblem(self);

      // check if retrying process is undergoing
      if (self.retrying == 0) {
        // check for every 2 second
        const id = setInterval(() => {
          if (self.retrying < 3 && !self.isOpen) {
            self.retrying += 1;
            self.connect(self);
          } else {
            clearInterval(id);
            self.retrying = -1;
          }
        }, 2000);
      }
    }
    self.log.push({
      state: "err",
      message: event,
      time: new Date().toLocaleString(),
    });
  }

  // handle reconnecting problem - this function needs to override by the child
  handleConnectionProblem(self = this) {
    console.log("connection problem, retrying", self.retrying);
  }

  // handle close event
  onClose(self = this, event) {
    self.log.push({
      state: "Init",
      message: event,
      time: new Date().toLocaleString(),
    });
    console.clear();
    console.log("start printing log...");
    self.printLog();
    console.log("end log...");
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
