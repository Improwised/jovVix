import { mount, flushPromises } from "@vue/test-utils";
import { describe, it, expect, beforeEach } from "vitest";
import { mockNuxtImport, registerEndpoint } from "@nuxt/test-utils/runtime";
import { readMultipartFormData } from "h3";
import CsvImportModal from "~/components/Quiz/CsvImportModal.vue";

const FIELDS = [
  { name: "Question Text", required: true },
  { name: "Question Type", required: true },
  { name: "Option 1", required: true },
  { name: "Option 2", required: true },
  { name: "Correct Answer", required: true },
  { name: "Resource", required: false },
];
const HEADERS = ["Prompt", "Kind", "A", "B", "Right"];
const MAPPED = {
  "Question Text": "Prompt",
  "Question Type": "Kind",
  "Option 1": "A",
  "Option 2": "B",
  "Correct Answer": "Right",
};
const row = (text, errors) => ({
  row: 2,
  values: {
    "Question Text": text,
    "Question Type": "single answer",
    "Option 1": "x",
    "Option 2": "y",
    "Correct Answer": "1",
  },
  errors,
});

// $fetch is a module import in the Nuxt test app, so requests go to registered endpoints.
// A relative API URL lets registerEndpoint intercept them.
mockNuxtImport("useRuntimeConfig", (original) => () => {
  const config = original();
  return { ...config, public: { ...config.public, apiUrl: "/api/v1" } };
});

// calls holds each request as { path, mapping, attachment }; replies are queued per test.
const calls = [];
const replies = [];
const respond = async (event) => {
  const parts = (await readMultipartFormData(event)) || [];
  const part = (name) => parts.find((p) => p.name === name);
  calls.push({
    path: event.path.replace(/^\/_/, ""), // test-utils serves endpoints under /_
    mapping: part("mapping")
      ? JSON.parse(part("mapping").data.toString())
      : null,
    attachment: part("attachment")?.data.toString(),
  });
  return replies.shift();
};
const BASE = "/api/v1/quizzes/quiz-1/questions";
registerEndpoint(`${BASE}/upload/preview`, {
  method: "POST",
  handler: respond,
});
registerEndpoint(`${BASE}/upload`, { method: "POST", handler: respond });
beforeEach(() => {
  calls.length = 0;
  replies.length = 0;
});

const ModalStub = {
  props: ["modelValue", "title", "description"],
  template:
    '<div><h2>{{ title }}</h2><p class="desc">{{ description }}</p><slot /><slot name="footer" /></div>',
};

const mountModal = () =>
  mount(CsvImportModal, {
    props: { modelValue: true, quizId: "quiz-1" },
    global: { stubs: { Modal: ModalStub, NavigationLink: true } },
  });

// A request round-trips through the mock server, which takes a few macrotasks.
const settle = async () => {
  for (let i = 0; i < 10; i++) {
    await flushPromises();
    await new Promise((resolve) => setTimeout(resolve, 0));
  }
};

const button = (wrapper, text) =>
  wrapper.findAll("button").find((b) => b.text().includes(text));

describe("CsvImportModal", () => {
  it("maps columns, previews rows, re-checks edits and imports", async () => {
    const wrapper = mountModal();

    // 1. Choosing a file asks for a preview with no mapping.
    replies.push({
      data: {
        fields: FIELDS,
        headers: HEADERS,
        mapping: {},
        unmapped_required: FIELDS.filter((f) => f.required).map((f) => f.name),
        rows: [],
      },
    });
    const input = wrapper.find('input[type="file"]');
    Object.defineProperty(input.element, "files", {
      value: [new File(["Prompt,Kind\n"], "vendor.csv", { type: "text/csv" })],
    });
    await input.trigger("change");
    await settle();
    expect(calls[0].path).toBe(`${BASE}/upload/preview`);
    expect(calls[0].mapping).toBeNull();
    expect(wrapper.text()).toContain("Match your columns");

    // 2. Previewing with required fields unmapped is refused before any request.
    await button(wrapper, "Preview rows").trigger("click");
    await settle();
    expect(calls).toHaveLength(1);

    // 3. Map every required field, then preview: a row with errors blocks the import.
    const selects = wrapper.findAll("label select");
    for (const select of selects) {
      const label = select.element.closest("label").textContent;
      const field = Object.keys(MAPPED).find((f) => label.includes(f));
      if (field) await select.setValue(MAPPED[field]);
    }
    replies.push({
      data: { rows: [row("", ["question text is required"])] },
    });
    await button(wrapper, "Preview rows").trigger("click");
    await settle();
    expect(calls[1].mapping).toEqual(MAPPED);
    expect(wrapper.text()).toContain("question text is required");
    expect(button(wrapper, "Import").attributes("disabled")).toBeDefined();

    // 4. Fix the cell and check again: the edited rows go up as a CSV with field headers.
    await wrapper
      .find('input[aria-label="Row 2 Question Text"]')
      .setValue("Q1");
    expect(button(wrapper, "Import").attributes("disabled")).toBeDefined();
    replies.push({ data: { rows: [row("Q1", [])] } });
    await button(wrapper, "Check edits").trigger("click");
    await settle();
    const rebuilt = calls[2].attachment;
    expect(rebuilt.split("\n")[0]).toBe(
      "Question Text,Question Type,Option 1,Option 2,Correct Answer"
    );
    expect(rebuilt).toContain("Q1,single answer,x,y,1");
    expect(calls[2].mapping["Question Text"]).toBe("Question Text");

    // 5. No errors left: import sends the same file and mapping, then reports the count.
    replies.push({ data: ["id-1"] });
    await button(wrapper, "Import 1 questions").trigger("click");
    await settle();
    expect(calls[3].path).toBe(`${BASE}/upload`);
    expect(calls[3].attachment).toBe(rebuilt);
    expect(wrapper.emitted("imported")).toEqual([[1]]);
    expect(wrapper.emitted("update:modelValue")).toEqual([[false]]);
  });

  // Import uploads the checked file, not the edited cells, so an unchecked edit must block it.
  it("blocks import after a cell edit until the edits are checked", async () => {
    const wrapper = mountModal();
    replies.push({
      data: {
        fields: FIELDS,
        headers: HEADERS,
        mapping: MAPPED,
        unmapped_required: [],
        rows: [],
      },
    });
    const input = wrapper.find('input[type="file"]');
    Object.defineProperty(input.element, "files", {
      value: [new File(["Prompt,Kind\n"], "vendor.csv", { type: "text/csv" })],
    });
    await input.trigger("change");
    await settle();

    replies.push({ data: { rows: [row("Q1", [])] } });
    await button(wrapper, "Preview rows").trigger("click");
    await settle();
    const importButton = () => button(wrapper, "Import 1 questions");
    expect(importButton().attributes("disabled")).toBeUndefined();

    // A valid row edited into another valid row: only the edit blocks the import.
    await wrapper
      .find('input[aria-label="Row 2 Question Text"]')
      .setValue("Q1 edited");
    expect(importButton().attributes("disabled")).toBeDefined();

    replies.push({ data: { rows: [row("Q1 edited", [])] } });
    await button(wrapper, "Check edits").trigger("click");
    await settle();
    expect(importButton().attributes("disabled")).toBeUndefined();
  });
});
