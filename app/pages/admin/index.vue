<template>
  <main
    class="flex flex-col gap-8 bg-jv-canvas sm:gap-10 px-4 sm:px-6 md:px-8 py-5 md:py-6"
  >
    <div class="mx-auto flex w-full flex-col gap-8 sm:gap-10">
      <div class="min-w-0">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Profile Settings
        </h1>
        <div
          class="ml-12 mt-1 h-3 w-28 rounded-full border-b-[3px] border-jv-yellow sm:ml-14 sm:w-32"
          aria-hidden="true"
        ></div>
      </div>

      <div v-if="userError?.data?.code == 401" class="sr-only">
        {{ navigateTo("/account/login") }}
      </div>

      <div
        v-if="userError"
        class="jv-border-rough bg-jv-white p-5 text-[18px] font-semibold text-jv-coral shadow-brutal-sm"
        role="alert"
      >
        {{ userError.data }}
      </div>

      <div
        v-else-if="userPending"
        class="jv-border-rough bg-jv-white p-5 text-center text-[18px] font-semibold text-jv-muted shadow-brutal-sm"
      >
        Pending...
      </div>

      <template v-else>
        <section
          class="rotate-[-0.2deg] jv-border-rough bg-jv-white shadow-brutal-lg"
        >
          <div class="p-5 sm:p-7 md:p-8">
            <div class="flex flex-col gap-5 md:flex-row md:items-center">
              <img
                :src="avatar"
                class="size-24 shrink-0 rounded-full border-[3px] border-jv-ink bg-jv-slate object-cover shadow-brutal-sm sm:size-28"
                :alt="userData.full_name"
              />

              <div class="min-w-0 flex-1">
                <h2
                  class="font-headings text-[32px] leading-tight text-jv-ink sm:text-[38px]"
                >
                  {{ userData.full_name }}
                </h2>
                <p
                  class="mt-1 break-all text-[16px] font-semibold text-jv-muted sm:text-[18px]"
                >
                  {{ userData.email }}
                </p>

                <div class="mt-5 flex flex-wrap gap-3">
                  <NavigationLink
                    :url-name="
                      userData.email_verify
                        ? 'Email Verified'
                        : 'Verify Yourself'
                    "
                    :disabled="userData.email_verify"
                    :title="
                      userData.email_verify
                        ? 'Email already verified'
                        : 'Click to verify your email'
                    "
                    :class="
                      userData.email_verify
                        ? 'rounded-[999px] bg-jv-slate'
                        : 'rounded-[999px] bg-jv-mint'
                    "
                    @click="handleEmailVerification"
                  />
                  <NavigationLink
                    url-name="Change Password"
                    class="rounded-[999px]"
                    @click="handleChangePasswordClick"
                  />

                  <NavigationLink
                    url-name="Delete Account"
                    class="rounded-[999px] bg-jv-coral text-white"
                    @click="deleteAccount"
                  />
                </div>
              </div>
            </div>
          </div>

          <form
            class="border-t-2 border-dashed border-jv-ink/15 p-5 sm:p-7 md:p-8"
            @submit.prevent="changeUserMetaData"
          >
            <div class="grid gap-5 md:grid-cols-2 md:gap-6">
              <label class="flex min-w-0 flex-col gap-2">
                <span
                  class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-muted"
                >
                  First Name
                </span>
                <input
                  id="first-name"
                  v-model="userData.first_name"
                  type="text"
                  class="h-14 w-full border-2 border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink outline-none transition-colors placeholder:text-jv-ink/35 focus:bg-jv-white"
                  placeholder="Pending..."
                  required
                  @focus="showCancleButton"
                />
              </label>

              <label class="flex min-w-0 flex-col gap-2">
                <span
                  class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-muted"
                >
                  Last Name
                </span>
                <input
                  id="last-name"
                  v-model="userData.last_name"
                  type="text"
                  class="h-14 w-full border-2 border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink outline-none transition-colors placeholder:text-jv-ink/35 focus:bg-jv-white"
                  placeholder="Pending..."
                  required
                  @focus="showCancleButton"
                />
              </label>
            </div>

            <label class="mt-6 flex min-w-0 flex-col gap-2">
              <span
                class="text-[13px] font-black uppercase tracking-[0.16em] text-jv-muted"
              >
                Email
              </span>
              <input
                id="email"
                v-model="userData.email"
                type="email"
                class="h-14 w-full border-2 border-jv-ink bg-jv-canvas px-4 text-[17px] font-semibold text-jv-ink outline-none disabled:text-jv-ink/80"
                placeholder="Pending..."
                disabled
              />
            </label>

            <div
              v-if="updateuserError"
              class="mt-5 jv-border-rough bg-jv-white p-3 text-[15px] font-semibold text-jv-coral shadow-brutal-sm"
              role="alert"
            >
              {{ updateuserError.data }}
            </div>

            <div class="mt-7 flex flex-wrap gap-3">
              <NavigationLink
                :url-name="saveButtonText"
                type="submit"
                :disabled="updateuserPending"
                class="bg-jv-coral text-white"
              >
                <Save class="size-5" :stroke-width="2.4" />
              </NavigationLink>

              <NavigationLink
                v-if="cancleButton"
                url-name="Cancel"
                class="bg-jv-white text-jv-ink"
                @click="hideCancleButton"
              />
            </div>
          </form>
        </section>

        <section
          class="rotate-[0.2deg] jv-border-rough bg-jv-white shadow-brutal-lg"
        >
          <div
            class="flex flex-col gap-4 border-b-2 border-jv-ink p-5 sm:flex-row sm:items-center sm:justify-between sm:p-6"
          >
            <div class="min-w-0">
              <h2
                class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[32px]"
              >
                Played Quizzes
              </h2>
              <p
                class="mt-1 text-[14px] font-semibold text-jv-muted sm:text-[15px]"
              >
                Review every quiz you've played so far.
              </p>
            </div>
            <NavigationLink
              url="/admin/quiz/list-quiz"
              url-name="View Played Quizzes"
              class="w-full bg-jv-yellow py-2 text-jv-ink sm:w-fit"
            />
          </div>
        </section>
      </template>
    </div>
  </main>
</template>

<script setup>
import { useUsersStore } from "~~/store/users";
import { getAvatarUrlByName } from "~~/composables/avatar";
import { usePush } from "notivue";
import { Save } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Admin Dashboard - jovVix",
  description:
    "Manage your quizzes, sessions, and reports from the jovVix admin dashboard.",
  robots: "noindex, nofollow",
});

const config = useRuntimeConfig();
const toast = usePush();
const userStore = useUsersStore();
const { getUserData, setUserData } = userStore;
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const updateuserError = ref(false);
const updateuserPending = ref(false);
const cancleButton = ref(false);
const avatar = computed(() => {
  const user = getUserData();
  return getAvatarUrlByName(user?.avatar);
});

const saveButtonText = computed(() =>
  updateuserPending.value ? "Pending..." : "Save Changes"
);

const {
  data: user,
  pending: userPending,
  error: userError,
} = useFetch(url.apiUrl + "/kratos/whoami", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
  server: false,
});

const userData = reactive({
  full_name: "",
  first_name: "",
  last_name: "",
  email: "",
  email_verify: false,
});

watch(
  [user, userError],
  () => {
    if (user.value) {
      userData.full_name =
        user.value.data.identity.traits.name.first +
        " " +
        user.value.data.identity.traits.name.last;
      userData.first_name = user.value.data.identity.traits.name.first;
      userData.last_name = user.value.data.identity.traits.name.last;
      userData.email = user.value.data.identity.traits.email;
      userData.email_verify =
        user.value.data.identity.verifiable_addresses[0].verified;
    }
    if (userError.value) {
      console.log("error");
    }
  },
  { immediate: true, deep: true }
);

const changeUserMetaData = async () => {
  updateuserPending.value = true;
  const { data, error } = await useFetch(url.apiUrl + "/kratos/user", {
    method: "PUT",
    headers: headers,
    mode: "cors",
    credentials: "include",
    body: {
      first_name: userData.first_name,
      last_name: userData.last_name,
    },
  });

  updateuserPending.value = false;
  if (error.value) {
    updateuserError.value = error.value.data;
    setTimeout(() => {
      updateuserError.value = false;
    }, 2000);
  } else {
    userData.full_name =
      data.value.data.first_name + " " + data.value.data.last_name;
    cancleButton.value = false;
  }
};

const showCancleButton = () => {
  cancleButton.value = true;
};

const hideCancleButton = () => {
  userData.first_name = user.value.data.identity.traits.name.first;
  userData.last_name = user.value.data.identity.traits.name.last;
  cancleButton.value = false;
};

const deleteAccount = async () => {
  const isconfirm = confirm("are you sure?");
  if (isconfirm) {
    try {
      await $fetch(`${url.apiUrl}/kratos/user`, {
        method: "DELETE",
        headers: headers,
        credentials: "include",
      });
      toast.success("user deleted succesfully!");
      setUserData(null);
      navigateTo("/");
    } catch (error) {
      console.error("Failed to delete the user", error);
      toast.error("Failed to delete the user.");
    }
  }
};

const isPrivilegeSessionValid = () => {
  if (!user.value?.data?.authenticated_at) return false;
  const authAt = new Date(user.value.data.authenticated_at).getTime();
  const now = Date.now();
  const diffInMinutes = (now - authAt) / (1000 * 60);
  return diffInMinutes <= config.public.privilegedSessionMaxAge;
};

const handleChangePasswordClick = async () => {
  if (!isPrivilegeSessionValid()) {
    toast.warning("Please re-enter your password to continue.");
    navigateTo("/account/login?returnTo=/account/change-password");
    return;
  }
  await nextTick();
  navigateTo("/account/change-password");
};

const handleEmailVerification = async () => {
  try {
    const verificationResponse = await fetch(
      `${url.kratosUrl}/self-service/verification/browser?return_to=${window.location.origin}/admin`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
      }
    );

    if (!verificationResponse.ok) {
      throw new Error("Failed to initiate verification flow");
    }

    const verification = await verificationResponse.json();
    const verificationPage = verification?.ui?.action;
    const flowId = verification?.id;
    const csrfToken = verification.ui.nodes.find(
      (node) => node.attributes.name === "csrf_token"
    ).attributes.value;

    const sendEmailResponse = await fetch(verificationPage, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({
        email: userData.email,
        csrf_token: csrfToken,
        method: "code",
      }),
    });

    if (!sendEmailResponse.ok) {
      throw new Error("Failed to send verification email");
    }

    setTimeout(() => {
      navigateTo(`/verification?flow=${flowId}`);
    }, 300);
  } catch (error) {
    console.error(error);
    toast.error("Failed to start verification");
  }
};
</script>
