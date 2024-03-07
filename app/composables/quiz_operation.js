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

    // custom attributes
    this.socket = null;
    this.connect(this);
  }

  getAddress(self = this) {
    return self.socket_url;
  }

  connect(self = this) {
    self.socket = new WebSocket(self.getAddress(self));
    self.socket.onopen = (event) => self.onOpen(self, event);
    self.socket.onerror = (event) => self.onError(self, event);
    self.socket.onclose = (event) => self.onClose(self, event);
    self.socket.onmessage = (event) => self.onMessage(self, event);
  }

  async handler(_, message) {
    if (!this.componentHandler) {
      throw Error(
        `Handler for component "${message.component}" is not registered`
      );
    }
    this.componentHandler(message);
  }

  onOpen(self = this, event) {
    self.isOpen = true;
    self.retrying = 0;
    self.log.push({ state: "Init", message: event });
    // console.log(event);
  }

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
    self.log.push({ state: "err", message: event });
    // console.log(event);
  }

  handleConnectionProblem(self = this) {
    console.log("connection problem, retrying", self.retrying);
  }

  onClose(self = this, event) {
    self.log.push({ state: "Init", message: event });
    // console.log(event);
  }

  onMessage(self = this, event) {
    const message = self.destructureMessage(event);
    this.currentComponent = message.component;
    this.currentEvent = message.event;
    this.log.push({ state: "Receive", message });
    self.handler(self, message, event);
  }

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

  sendMessage(
    component = this.currentComponent,
    event = this.currentEvent,
    data = ""
  ) {
    const message = {
      component: component,
      event: event,
      data: data,
    };
    this.log.push({ state: "Sent", message });
    this.socket.send(JSON.stringify(message));
  }

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
}
