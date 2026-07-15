<script setup>
import { computed, onMounted, ref } from "vue";
import { Check, Pencil, Search, Tag, Trash2, X } from "lucide-vue-next";
import { usePush } from "notivue";
import NavigationLink from "@/components/common/NavigationLink.vue";
import { useUsersStore } from "~~/store/users";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Quiz Categories - jovVix",
  description: "Manage the categories used to group public quizzes on jovVix.",
  robots: "noindex, nofollow",
});

const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = usePush();
const usersStore = useUsersStore();

const canCreatePublicQuiz = computed(
  () => !!usersStore.userData?.canCreatePublicQuiz
);

// The API enforces the public-quiz admin allowlist on every write; this guard
// just redirects non-admins away instead of showing them a broken page.
onMounted(async () => {
  if (!usersStore.userData) {
    await setUserDataStore();
  }
  if (!canCreatePublicQuiz.value) {
    navigateTo("/admin/quiz/list-quiz");
  }
});

const {
  data: categoriesData,
  pending: categoriesPending,
  error: categoriesError,
  refresh,
} = useFetch(`${url.apiUrl}/categories`, {
  method: "GET",
  headers: headers,
  credentials: "include",
});

const categories = computed(() => categoriesData.value?.data || []);

const searchQuery = ref("");
const creating = ref(false);
const editingId = ref("");
const editingName = ref("");
const savingEdit = ref(false);
const deletingId = ref("");

const formatDate = (iso) => {
  if (!iso) return "";
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return "";
  const day = String(date.getDate()).padStart(2, "0");
  const month = date.toLocaleString("en", { month: "short" });
  return `${day}-${month}-${date.getFullYear()}`;
};

const filteredCategories = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return categories.value;

  return categories.value.filter((category) =>
    category.name.toLowerCase().includes(query)
  );
});

// The create option only shows when the searched name doesn't exist yet.
const hasExactMatch = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return true;
  return categories.value.some(
    (category) => category.name.toLowerCase() === query
  );
});

const handleCreateFromSearch = async () => {
  const name = searchQuery.value.trim();
  if (!name || creating.value) return;

  try {
    creating.value = true;
    await $fetch(`${url.apiUrl}/categories`, {
      method: "POST",
      body: { name },
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Category created successfully.");
    searchQuery.value = "";
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while creating category."
    );
  } finally {
    creating.value = false;
  }
};

const startEdit = (category) => {
  editingId.value = category.id;
  editingName.value = category.name;
};

const cancelEdit = () => {
  editingId.value = "";
  editingName.value = "";
};

const handleUpdateCategory = async (category) => {
  const name = editingName.value.trim();
  if (!name || savingEdit.value) return;
  if (name === category.name) {
    cancelEdit();
    return;
  }

  try {
    savingEdit.value = true;
    await $fetch(`${url.apiUrl}/categories/${category.id}`, {
      method: "PUT",
      body: { name },
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Category updated successfully.");
    cancelEdit();
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while updating category."
    );
  } finally {
    savingEdit.value = false;
  }
};

const handleDeleteCategory = async (category) => {
  if (deletingId.value) return;

  try {
    deletingId.value = category.id;
    await $fetch(`${url.apiUrl}/categories/${category.id}`, {
      method: "DELETE",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
    });
    toast.success("Category deleted successfully.");
    refresh();
  } catch (error) {
    toast.error(
      error?.data?.message || error?.message || "Error while deleting category."
    );
  } finally {
    deletingId.value = "";
  }
};
</script>

<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <div
      class="flex flex-col gap-5 md:flex-row md:items-start md:justify-between"
    >
      <div class="min-w-0">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Quiz Categories
        </h1>
        <div
          class="ml-1 mt-1 h-3 w-24 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-28"
          aria-hidden="true"
        ></div>
      </div>
    </div>

    <div
      class="flex flex-col gap-3 sm:flex-row sm:items-center sm:gap-4 justify-between"
    >
      <label
        class="flex h-11 w-full min-w-0 rotate-[1deg] items-center gap-2.5 jv-border-rough bg-jv-white px-3 shadow-brutal-sm sm:h-12 sm:max-w-[448px] sm:gap-3 sm:px-4"
      >
        <Search class="size-5 shrink-0 text-jv-ink/40" :stroke-width="2.4" />
        <input
          v-model="searchQuery"
          type="search"
          maxlength="50"
          class="h-full min-w-0 flex-1 bg-transparent text-[15px] text-jv-ink outline-none placeholder:text-jv-ink/45 sm:text-[17px]"
          placeholder="Search or add a category..."
          @keyup.enter="!hasExactMatch && handleCreateFromSearch()"
        />
      </label>
      <NavigationLink
        :url-name="creating ? 'Creating...' : 'Create Category'"
        class="w-full bg-jv-mint py-2 font-[500] sm:w-fit"
        :disabled="creating || hasExactMatch"
        @click="handleCreateFromSearch"
      />
    </div>

    <div
      v-if="categoriesPending"
      class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
    >
      Loading...
    </div>

    <div
      v-else-if="categoriesError"
      class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
    >
      {{ categoriesError.message }}
    </div>

    <section
      v-else-if="categories.length < 1 && !searchQuery.trim()"
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          No Categories Yet!
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          Type a name in the search box above to create your first category.
        </p>
      </div>
    </section>

    <section
      v-else-if="filteredCategories.length < 1"
      class="grid min-h-[320px] place-items-center px-2 text-center sm:min-h-[420px]"
    >
      <div class="max-w-md py-8">
        <h2
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[34px]"
        >
          No Category Found!
        </h2>
        <p class="mt-2 text-[17px] text-jv-muted">
          There is no category named "{{ searchQuery.trim() }}" yet. Use the
          Create Category button to add it.
        </p>
      </div>
    </section>

    <div v-else class="flex flex-col gap-3">
      <div
        class="hidden lg:grid grid-cols-[2fr_1fr_180px] gap-4 px-5 py-2 text-[12px] font-black uppercase tracking-[0.1em] text-jv-muted"
      >
        <div>Category Name</div>
        <div class="text-center">Created At</div>
        <div class="text-right">Actions</div>
      </div>

      <div
        v-for="category in filteredCategories"
        :key="category.id"
        class="grid grid-cols-1 lg:grid-cols-[2fr_1fr_180px] gap-4 items-center jv-border-rough bg-jv-white px-5 py-4 shadow-brutal-sm hover:translate-y-[-2px] hover:shadow-brutal-md transition-all"
      >
        <template v-if="editingId === category.id">
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="grid size-10 shrink-0 place-items-center border-2 border-jv-ink bg-jv-yellow/60 text-jv-ink"
            >
              <Tag class="size-5" :stroke-width="2.4" />
            </span>
            <input
              v-model.trim="editingName"
              type="text"
              required
              maxlength="50"
              class="h-10 min-w-0 flex-1 border-[3px] border-jv-ink bg-jv-canvas px-3 text-[16px] font-semibold text-jv-ink outline-none"
              @keyup.enter="handleUpdateCategory(category)"
              @keyup.esc="cancelEdit"
            />
          </div>
          <div class="hidden lg:block"></div>
          <div class="flex items-center gap-2 lg:justify-end">
            <button
              type="button"
              :disabled="savingEdit || !editingName"
              class="grid size-9 shrink-0 place-items-center border-2 border-jv-ink bg-jv-mint text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg] disabled:cursor-not-allowed disabled:opacity-60"
              aria-label="Save category name"
              @click="handleUpdateCategory(category)"
            >
              <Check class="size-4" :stroke-width="2.5" />
            </button>
            <button
              type="button"
              class="grid size-9 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg]"
              aria-label="Cancel rename"
              @click="cancelEdit"
            >
              <X class="size-4" :stroke-width="2.5" />
            </button>
          </div>
        </template>

        <template v-else>
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="grid size-10 shrink-0 place-items-center border-2 border-jv-ink bg-jv-yellow/60 text-jv-ink"
            >
              <Tag class="size-5" :stroke-width="2.4" />
            </span>
            <h3
              class="min-w-0 truncate font-headings text-[22px] font-bold leading-tight text-jv-ink sm:text-[24px]"
            >
              {{ category.name }}
            </h3>
          </div>

          <div class="lg:text-center">
            <span
              class="lg:hidden mr-2 text-[14px] font-bold uppercase tracking-wide text-jv-muted"
            >
              Created At:
            </span>
            <span class="text-[15px] font-bold text-jv-ink">
              {{ formatDate(category.created_at) }}
            </span>
          </div>

          <div class="flex items-center gap-2 lg:justify-end">
            <button
              type="button"
              class="grid size-9 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg]"
              aria-label="Rename category"
              @click="startEdit(category)"
            >
              <Pencil class="size-4" :stroke-width="2.4" />
            </button>
            <button
              type="button"
              :disabled="deletingId === category.id"
              class="grid size-9 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-all hover:rotate-[3deg] hover:text-jv-coral disabled:cursor-wait disabled:opacity-60"
              aria-label="Delete category"
              @click="handleDeleteCategory(category)"
            >
              <Trash2 class="size-4" :stroke-width="2.4" />
            </button>
          </div>
        </template>
      </div>
    </div>
  </main>
</template>
