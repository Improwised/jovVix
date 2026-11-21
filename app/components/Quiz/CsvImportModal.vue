<script setup>
import { computed, ref, watch } from "vue";
import { Download, Upload } from "lucide-vue-next";
import { usePush } from "notivue";
import { Modal } from "@/components/ui/modal";
import NavigationLink from "@/components/common/NavigationLink.vue";
import { readApiError } from "@/composables/apiError";
import {
  applyTemplate,
  buildCsv,
  identityMapping,
  loadTemplates,
  saveTemplate,
} from "@/composables/csv_import";

// Import questions from any CSV: map its columns to question fields, check every row in a
// preview, fix cells there, then import. The API validates; this component only collects.

const props = defineProps({
  modelValue: { type: Boolean, required: true },
  quizId: { type: String, required: true },
});
const emit = defineEmits(["update:modelValue", "imported"]);

const url = useRuntimeConfig().public;
const toast = usePush();

const step = ref("file"); // file -> map -> preview
const pending = ref(false);
const file = ref(null);
const fileName = ref("");
const fields = ref([]);
const headers = ref([]);
const mapping = ref({});
const rows = ref([]);
const previewFields = ref([]);
const edited = ref(false);
const templates = ref({});
const selectedTemplate = ref("");
const templateName = ref("");

const reset = () => {
  step.value = "file";
  file.value = null;
  fileName.value = "";
  fields.value = [];
  headers.value = [];
  mapping.value = {};
  rows.value = [];
  previewFields.value = [];
  edited.value = false;
  selectedTemplate.value = "";
  templateName.value = "";
  templates.value = loadTemplates();
};

watch(
  () => props.modelValue,
  (open) => open && reset(),
  { immediate: true }
);

const close = () => emit("update:modelValue", false);

const errorRows = computed(
  () => rows.value.filter((row) => row.errors.length > 0).length
);
const missingRequired = computed(() =>
  fields.value
    .filter((field) => field.required && !mapping.value[field.name])
    .map((field) => field.name)
);
const canImport = computed(
  () =>
    step.value === "preview" &&
    rows.value.length > 0 &&
    errorRows.value === 0 &&
    !edited.value &&
    !pending.value
);

// Every field gets a key, "" when unmapped, so each select shows a choice.
const fullMapping = (map) =>
  Object.fromEntries(
    fields.value.map((field) => [field.name, map?.[field.name] || ""])
  );

const formFor = (blob, map) => {
  const form = new FormData();
  form.append("attachment", blob, fileName.value || "questions.csv");
  if (map) form.append("mapping", JSON.stringify(map));
  return form;
};

const requestPreview = async (blob, map) => {
  const response = await $fetch(
    `${url.apiUrl}/quizzes/${props.quizId}/questions/upload/preview`,
    {
      method: "POST",
      headers: { Accept: "application/json" },
      body: formFor(blob, map),
      credentials: "include",
    }
  );
  return response?.data || {};
};

const handleFile = async (event) => {
  const chosen = event.target.files?.[0];
  event.target.value = "";
  if (!chosen) return;
  pending.value = true;
  try {
    fileName.value = chosen.name;
    const data = await requestPreview(chosen, null);
    file.value = chosen;
    fields.value = data.fields || [];
    headers.value = data.headers || [];
    mapping.value = fullMapping(data.mapping);
    step.value = "map";
  } catch (error) {
    toast.error(readApiError(error, "Could not read this CSV file."));
  } finally {
    pending.value = false;
  }
};

const useTemplate = () => {
  const template = templates.value[selectedTemplate.value];
  if (template)
    mapping.value = fullMapping(applyTemplate(template, headers.value));
};

const storeTemplate = () => {
  const name = templateName.value.trim();
  if (!name) {
    toast.error("Name the template first.");
    return;
  }
  templates.value = saveTemplate(name, mapping.value);
  selectedTemplate.value = name;
  templateName.value = "";
  toast.success(`Template "${name}" saved in this browser.`);
};

const showPreview = async () => {
  if (missingRequired.value.length) {
    toast.error(`Map these columns: ${missingRequired.value.join(", ")}`);
    return;
  }
  pending.value = true;
  try {
    const cleanMapping = Object.fromEntries(
      Object.entries(mapping.value).filter(([, header]) => header)
    );
    const data = await requestPreview(file.value, cleanMapping);
    previewFields.value = fields.value
      .map((field) => field.name)
      .filter((name) => cleanMapping[name]);
    rows.value = data.rows || [];
    edited.value = false;
    step.value = "preview";
  } catch (error) {
    toast.error(readApiError(error, "Could not check this CSV file."));
  } finally {
    pending.value = false;
  }
};

// After cell edits, the preview rows become the file: a CSV with field names as headers.
const revalidate = async () => {
  pending.value = true;
  try {
    const blob = new Blob([buildCsv(previewFields.value, rows.value)], {
      type: "text/csv",
    });
    const map = identityMapping(previewFields.value);
    const data = await requestPreview(blob, map);
    file.value = blob;
    headers.value = [...previewFields.value];
    mapping.value = fullMapping(map);
    rows.value = data.rows || [];
    edited.value = false;
  } catch (error) {
    toast.error(readApiError(error, "Could not check the edited rows."));
  } finally {
    pending.value = false;
  }
};

const importRows = async () => {
  if (!canImport.value) return;
  pending.value = true;
  try {
    const cleanMapping = Object.fromEntries(
      Object.entries(mapping.value).filter(([, header]) => header)
    );
    await $fetch(`${url.apiUrl}/quizzes/${props.quizId}/questions/upload`, {
      method: "POST",
      headers: { Accept: "application/json" },
      body: formFor(file.value, cleanMapping),
      credentials: "include",
    });
    emit("imported", rows.value.length);
    close();
  } catch (error) {
    toast.error(readApiError(error, "Failed to import CSV."));
  } finally {
    pending.value = false;
  }
};

const stepTitle = computed(
  () =>
    ({
      file: "Import questions from CSV",
      map: "Match your columns",
      preview: "Check the rows",
    })[step.value]
);
const stepDescription = computed(
  () =>
    ({
      file: "Any column names work: you match them to question fields next.",
      map: `Pick the column in ${fileName.value} for each field. Fields marked * are required.`,
      preview: errorRows.value
        ? `${errorRows.value} of ${rows.value.length} rows need fixing. Edit them here, then check again.`
        : `${rows.value.length} rows are ready to import.`,
    })[step.value]
);

const buttonClass =
  "inline-flex h-12 items-center justify-center border-[3px] border-jv-ink px-6 text-[17px] font-black text-jv-ink shadow-brutal-sm transition-transform active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-60";
const fieldClass =
  "h-11 w-full min-w-0 border-[3px] border-jv-ink bg-jv-canvas px-3 text-[15px] font-semibold text-jv-ink";
</script>

<template>
  <Modal
    :model-value="modelValue"
    size="xl"
    :close-on-backdrop="false"
    :title="stepTitle"
    :description="stepDescription"
    class="md:max-w-[960px] lg:max-w-[1120px]"
    @update:model-value="(open) => !open && close()"
  >
    <div class="max-h-[62dvh] overflow-y-auto pb-3 pr-1">
      <div v-if="step === 'file'" class="grid gap-4">
        <label class="grid gap-2">
          <span
            class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-ink"
          >
            Choose CSV File <span class="text-jv-coral">*</span>
          </span>
          <span
            class="flex h-14 cursor-pointer border-[3px] border-jv-ink bg-jv-canvas"
          >
            <span
              class="inline-flex h-full items-center gap-2 bg-jv-ink px-4 text-[16px] font-black text-jv-white"
            >
              <Upload class="size-4" :stroke-width="2.4" />
              {{ pending ? "Reading..." : "Choose File" }}
            </span>
            <span
              class="flex min-w-0 flex-1 items-center px-4 text-[15px] font-semibold text-jv-muted"
            >
              {{ fileName || "No file chosen" }}
            </span>
            <input
              type="file"
              class="hidden"
              accept=".csv,text/csv"
              :disabled="pending"
              @change="handleFile"
            />
          </span>
        </label>
        <p class="text-[15px] leading-[1.6] text-jv-muted">
          New rows are appended to this quiz. The sample file shows every field.
        </p>
        <NavigationLink
          url="/files/demo.csv"
          url-name="Download Sample"
          external
          download="demo.csv"
          class="w-fit rounded-[999px] bg-jv-white"
        >
          <Download class="size-4" :stroke-width="2.4" />
        </NavigationLink>
      </div>

      <div v-else-if="step === 'map'" class="grid gap-5">
        <div class="grid gap-3 sm:grid-cols-2">
          <label class="grid gap-2">
            <span class="text-[13px] font-black uppercase tracking-[0.16em]">
              Saved template
            </span>
            <span class="flex gap-2">
              <select
                v-model="selectedTemplate"
                :class="fieldClass"
                aria-label="Saved mapping template"
              >
                <option value="">None</option>
                <option
                  v-for="(_, name) in templates"
                  :key="name"
                  :value="name"
                >
                  {{ name }}
                </option>
              </select>
              <button
                type="button"
                :class="[buttonClass, 'h-11 bg-jv-white px-4 text-[15px]']"
                :disabled="!selectedTemplate"
                @click="useTemplate"
              >
                Apply
              </button>
            </span>
          </label>
          <label class="grid gap-2">
            <span class="text-[13px] font-black uppercase tracking-[0.16em]">
              Save this mapping as
            </span>
            <span class="flex gap-2">
              <input
                v-model="templateName"
                type="text"
                maxlength="40"
                placeholder="Template name"
                :class="fieldClass"
              />
              <button
                type="button"
                :class="[buttonClass, 'h-11 bg-jv-white px-4 text-[15px]']"
                @click="storeTemplate"
              >
                Save
              </button>
            </span>
          </label>
        </div>

        <div class="grid gap-3 sm:grid-cols-2 md:grid-cols-3">
          <label v-for="field in fields" :key="field.name" class="grid gap-1">
            <span class="text-[14px] font-black text-jv-ink">
              {{ field.name }}
              <span v-if="field.required" class="text-jv-coral">*</span>
            </span>
            <select
              v-model="mapping[field.name]"
              :class="[
                fieldClass,
                field.required && !mapping[field.name] ? 'border-jv-coral' : '',
              ]"
            >
              <option value="">Not mapped</option>
              <option v-for="header in headers" :key="header" :value="header">
                {{ header }}
              </option>
            </select>
          </label>
        </div>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full border-collapse text-left text-[14px]">
          <thead>
            <tr class="bg-jv-ink text-jv-white">
              <th class="px-2 py-2">Row</th>
              <th
                v-for="name in previewFields"
                :key="name"
                class="min-w-[140px] px-2 py-2"
              >
                {{ name }}
              </th>
            </tr>
          </thead>
          <tbody>
            <template v-for="row in rows" :key="row.row">
              <tr
                :class="
                  row.errors.length ? 'bg-jv-coral/15' : 'odd:bg-jv-canvas'
                "
              >
                <td class="px-2 py-1 font-black">{{ row.row }}</td>
                <td v-for="name in previewFields" :key="name" class="px-1 py-1">
                  <input
                    v-model="row.values[name]"
                    type="text"
                    class="h-9 w-full min-w-0 border-2 border-jv-ink/40 bg-jv-white px-2"
                    :aria-label="`Row ${row.row} ${name}`"
                    @input="edited = true"
                  />
                </td>
              </tr>
              <tr v-if="row.errors.length" class="bg-jv-coral/15">
                <td></td>
                <td
                  :colspan="previewFields.length"
                  class="px-2 pb-2 text-[13px] font-semibold text-jv-coral"
                  role="alert"
                >
                  {{ row.errors.join("; ") }}
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <template #footer>
      <div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button
          v-if="step === 'file'"
          type="button"
          :class="[buttonClass, 'bg-jv-white']"
          @click="close"
        >
          Cancel
        </button>
        <template v-else>
          <button
            type="button"
            :class="[buttonClass, 'bg-jv-white']"
            :disabled="pending"
            @click="step = 'file'"
          >
            Upload another file
          </button>
          <button
            v-if="step === 'map'"
            type="button"
            :class="[buttonClass, 'bg-jv-accent-green']"
            :disabled="pending"
            @click="showPreview"
          >
            {{ pending ? "Checking..." : "Preview rows" }}
          </button>
          <template v-else>
            <button
              type="button"
              :class="[buttonClass, 'bg-jv-white']"
              :disabled="pending"
              @click="step = 'map'"
            >
              Change mapping
            </button>
            <button
              type="button"
              :class="[buttonClass, 'bg-jv-yellow']"
              :disabled="pending || !edited"
              @click="revalidate"
            >
              Check edits
            </button>
            <button
              type="button"
              :class="[buttonClass, 'bg-jv-accent-green']"
              :disabled="!canImport"
              @click="importRows"
            >
              {{ pending ? "Importing..." : `Import ${rows.length} questions` }}
            </button>
          </template>
        </template>
      </div>
    </template>
  </Modal>
</template>
