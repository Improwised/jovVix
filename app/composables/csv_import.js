// Helpers for the CSV import modal: building a CSV from edited preview rows, and mapping
// templates kept in this browser.
// shortcut: templates live in localStorage, per browser; move them to the API if authors need
// them on several devices or shared with co-editors.

export const TEMPLATE_STORAGE_KEY = "jovvix.csvMappingTemplates";

const escapeCell = (value) => {
  const text = value == null ? "" : String(value);
  return /[",\r\n]/.test(text) ? `"${text.replaceAll('"', '""')}"` : text;
};

// buildCsv writes rows ({ values: { [field]: string } }) under the field names as headers, so the
// result imports with the identity mapping.
export const buildCsv = (fieldNames, rows) =>
  [fieldNames, ...rows.map((row) => fieldNames.map((f) => row.values?.[f]))]
    .map((cells) => cells.map(escapeCell).join(","))
    .join("\n") + "\n";

export const identityMapping = (fieldNames) =>
  Object.fromEntries(fieldNames.map((name) => [name, name]));

// applyTemplate keeps only the template entries whose header exists in this file.
export const applyTemplate = (template, headers) =>
  Object.fromEntries(
    Object.entries(template || {}).filter(([, header]) =>
      headers.includes(header)
    )
  );

const storageOrNull = (storage) => {
  try {
    return storage ?? globalThis.localStorage ?? null;
  } catch {
    return null;
  }
};

export const loadTemplates = (storage) => {
  try {
    const parsed = JSON.parse(
      storageOrNull(storage)?.getItem(TEMPLATE_STORAGE_KEY) || "{}"
    );
    return parsed && typeof parsed === "object" && !Array.isArray(parsed)
      ? parsed
      : {};
  } catch {
    return {};
  }
};

// saveTemplate stores the non-empty entries of mapping under name and returns all templates.
export const saveTemplate = (name, mapping, storage) => {
  const templates = loadTemplates(storage);
  templates[name] = Object.fromEntries(
    Object.entries(mapping).filter(([, header]) => header)
  );
  try {
    storageOrNull(storage)?.setItem(
      TEMPLATE_STORAGE_KEY,
      JSON.stringify(templates)
    );
  } catch {
    // Private mode or blocked storage: the template lasts for this session only.
  }
  return templates;
};
