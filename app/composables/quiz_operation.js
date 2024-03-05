export default class QuizHandler {
  JOIN_EVENT = "JOIN_REQUEST";

  constructor(socket_url, identifier, componentHandler, others) {
    if (!(socket_url && identifier && componentHandler)) {
      console.log(socket_url, identifier, componentHandler);
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

    // custom attributes
    this.socket = new WebSocket(this.getAddress(this));
    this.socket.onopen = (event) => this.onOpen(this, event);
    this.socket.onerror = (event) => this.onError(this, event);
    this.socket.onclose = (event) => this.onClose(this, event);
    this.socket.onmessage = (event) => this.onMessage(this, event);
  }

  getAddress(currentObj = this) {
    return currentObj.socket_url;
  }

  async handler(_, message) {
    if (!this.componentHandler) {
      throw Error(
        `Handler for component "${message.component}" is not registered`
      );
    }
    this.componentHandler(message);
  }

  onOpen(currentObj = this, event) {
    currentObj.log.push({ state: "Init", message: event });
    // console.log(event);
  }

  onError(currentObj = this, event) {
    console.log(currentObj, event);
    currentObj.log.push({ state: "err", message: event });
    // console.log(event);
  }

  onClose(currentObj = this, event) {
    currentObj.log.push({ state: "Init", message: event });
    // console.log(event);
  }

  onMessage(currentObj = this, event) {
    const message = currentObj.destructureMessage(event);
    this.currentComponent = message.component;
    this.currentEvent = message.event;
    this.log.push({ state: "Receive", message });
    currentObj.handler(currentObj, message, event);
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
