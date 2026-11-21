import { describe, expect, it } from "vitest";
import {
  applyTemplate,
  buildCsv,
  identityMapping,
  loadTemplates,
  saveTemplate,
} from "../../composables/csv_import";

const memoryStorage = () => {
  const data = {};
  return {
    getItem: (k) => data[k] ?? null,
    setItem: (k, v) => {
      data[k] = v;
    },
  };
};

describe("csv_import", () => {
  it("builds a CSV that quotes commas, quotes and newlines", () => {
    const csv = buildCsv(
      ["Question Text", "Option 1"],
      [{ values: { "Question Text": 'Say "hi", then\nwave', "Option 1": "A" } }]
    );
    expect(csv).toBe('Question Text,Option 1\n"Say ""hi"", then\nwave",A\n');
  });

  it("writes blank cells for missing values", () => {
    expect(buildCsv(["A", "B"], [{ values: { A: "1" } }])).toBe("A,B\n1,\n");
  });

  it("maps each field to itself", () => {
    expect(identityMapping(["A", "B"])).toEqual({ A: "A", B: "B" });
  });

  it("applies only template headers that the file has", () => {
    expect(
      applyTemplate({ "Question Text": "Prompt", Points: "Score" }, ["Prompt"])
    ).toEqual({ "Question Text": "Prompt" });
  });

  it("saves and loads templates, dropping unmapped fields", () => {
    const storage = memoryStorage();
    saveTemplate("Vendor", { "Question Text": "Prompt", Points: "" }, storage);
    expect(loadTemplates(storage)).toEqual({
      Vendor: { "Question Text": "Prompt" },
    });
  });

  it("survives corrupt or missing storage", () => {
    const storage = memoryStorage();
    storage.setItem("jovvix.csvMappingTemplates", "not json");
    expect(loadTemplates(storage)).toEqual({});
    const broken = {
      getItem: () => {
        throw new Error("blocked");
      },
      setItem: () => {
        throw new Error("blocked");
      },
    };
    expect(loadTemplates(broken)).toEqual({});
    expect(saveTemplate("X", { A: "B" }, broken)).toEqual({ X: { A: "B" } });
  });
});
