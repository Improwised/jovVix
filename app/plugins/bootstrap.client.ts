import { Modal, Dropdown } from "bootstrap";

export default defineNuxtPlugin(() => ({
  provide: {
    bootstrap: {
      Modal,
      Dropdown,
    },
  },
}));
