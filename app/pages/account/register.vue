<script setup>
import { ref } from "vue";
import {
  User,
  Mail,
  Lock,
  Eye,
  EyeOff,
  Loader2,
  AlertCircle,
} from "lucide-vue-next";
import { useToast } from "vue-toastification";
import { useUsersStore } from "~~/store/users";
import { useUserPasswordRules } from "@/composables/user_password_rules";
import NavigationLink from "@/components/common/NavigationLink.vue";

definePageMeta({
  layout: "auth",
});

const userData = useUsersStore();
const { getUserData } = userData;
const route = useRoute();
const router = useRouter();
const toast = useToast();
const firstname = ref();
const lastname = ref();
const email = ref();
const password = ref();
const csrfToken = ref();
const showPassword = ref(false);
const submitted = ref(false);
const { kratosUrl } = useRuntimeConfig().public;
const errors = ref({
  email: "",
  password: "",
  firstname: "",
  lastname: "",
});
const registerURLWithFlowQuery = ref("");
const { passwordErrors } = useUserPasswordRules(password, firstname, lastname);

console.log(); // required so async IIFE below doesn't trigger Nuxt 5xx
function onSubmit(event) {
  submitted.value = true;
  if (passwordErrors.value.length > 0) {
    return;
  }
  event.target.submit();
}

(async () => {
  if (process.client) {
    const user = getUserData();
    if (user && user?.role == "admin-user") {
      navigateTo("/");
      return;
    }
    if (route.query.flow) {
      try {
        await $fetch(kratosUrl + "/self-service/registration/flows", {
          method: "GET",
          credentials: "include",
          headers: {
            Accept: "application/json",
          },
          query: {
            id: route.query.flow,
          },
          onResponseError({ response }) {
            if (response._data?.error?.code === 410) {
              navigateTo("/account/register");
            }

            if (response.status === 400) {
              toast.warning(
                "Please fill out the form correctly, password or email was incorrect"
              );
            }
          },
          onResponse({ response }) {
            response?._data?.ui?.messages?.forEach((message) => {
              if (message.type === "error") {
                if (message.text?.toLowerCase().includes("expired")) {
                  toast.error(
                    "Your session has expired. Please start the registration process again."
                  );
                  return;
                }
                if (message.id === 4000007) {
                  toast.error("An account with the same email exists already!");
                } else {
                  toast.error(message.text);
                }
              }
            });
            response?._data?.ui?.nodes?.forEach((node) => {
              if (node.attributes.name === "traits.email") {
                errors.value.email = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "password") {
                errors.value.password = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "traits.name.first") {
                errors.value.firstname = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
              if (node.attributes.name === "traits.name.last") {
                errors.value.lastname = node.messages
                  .map((message) => message.text)
                  .join(", ");
              }
            });
          },
        });
        await setFlowIDAndCSRFToken();
      } catch (error) {
        console.error(error);
        await setFlowIDAndCSRFToken();
      }
    } else {
      await setFlowIDAndCSRFToken();
    }
  } else {
    await setFlowIDAndCSRFToken();
  }
})();

async function setFlowIDAndCSRFToken() {
  try {
    const kratosResponse = await $fetch(
      kratosUrl + "/self-service/registration/browser",
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        credentials: "include",
        onResponseError({ response }) {
          console.error(
            "error while getting the flow id from the server",
            response
          );
        },
      }
    );

    router.push("?flow=" + kratosResponse?.id);
    csrfToken.value = kratosResponse?.ui?.nodes[0]?.attributes?.value;
    registerURLWithFlowQuery.value = kratosResponse?.ui?.action;
  } catch (error) {
    console.error(error);
  }
}
</script>

<template>
  <div class="relative w-full max-w-[420px] mx-auto z-10">
    <!-- Card with subtle tilt -->
    <div class="relative rotate-1">
      <!-- Decorative "clip" tab on top of the card -->
      <span
        class="absolute left-1/2 -top-[8px] z-20 h-3 w-12 -translate-x-1/2 jv-card border-2 border-jv-ink bg-jv-slate shadow-brutal-sm"
        aria-hidden="true"
      ></span>

      <div
        class="bg-jv-white border-2 border-jv-ink shadow-brutal-lg jv-card px-6 sm:px-8 py-7 sm:py-9"
      >
        <!-- Heading -->
        <header class="flex flex-col items-center gap-1.5 mb-5">
          <div class="relative inline-block">
            <h1
              class="font-headings text-jv-ink text-[28px] sm:text-[34px] leading-none m-0 whitespace-nowrap"
            >
              Create account
            </h1>
            <!-- Hand-drawn underline doodle -->
            <svg
              class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-mint"
              width="160"
              height="14"
              viewBox="0 0 160 14"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                d="M3 9 Q 25 1, 50 7 T 100 6 T 157 4"
                stroke="currentColor"
                stroke-width="2.5"
                stroke-linecap="round"
                fill="none"
              />
            </svg>
          </div>
          <p
            class="font-body text-jv-ink/70 text-sm sm:text-base m-0 mt-1 text-center"
          >
            Join the platform and start your journey in seconds.
          </p>
        </header>

        <!-- Form -->
        <form
          method="POST"
          :action="registerURLWithFlowQuery"
          enctype="application/json"
          class="flex flex-col gap-3"
          @submit.prevent="onSubmit"
        >
          <!-- First name -->
          <div class="flex flex-col gap-1.5">
            <label
              for="firstname"
              class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
            >
              First Name
            </label>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <User
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="firstname"
                v-model="firstname"
                type="text"
                name="traits.name.first"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                placeholder="Virat"
                required
              />
            </div>
            <p
              v-if="errors.firstname"
              class="font-body text-jv-accent-red text-xs px-0.5 m-0 flex items-center gap-1"
            >
              <AlertCircle class="size-3.5 shrink-0" :stroke-width="2.2" />
              {{ errors.firstname }}
            </p>
          </div>

          <!-- Last name -->
          <div class="flex flex-col gap-1.5">
            <label
              for="lastname"
              class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
            >
              Last Name
            </label>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <User
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="lastname"
                v-model="lastname"
                type="text"
                name="traits.name.last"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                placeholder="Kohli"
                required
              />
            </div>
            <p
              v-if="errors.lastname"
              class="font-body text-jv-accent-red text-xs px-0.5 m-0 flex items-center gap-1"
            >
              <AlertCircle class="size-3.5 shrink-0" :stroke-width="2.2" />
              {{ errors.lastname }}
            </p>
          </div>

          <!-- Email -->
          <div class="flex flex-col gap-1.5">
            <label
              for="email"
              class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
            >
              Email Address
            </label>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <Mail
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="email"
                v-model="email"
                type="email"
                name="traits.email"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                placeholder="virat.kohli@gmail.com"
                required
              />
            </div>
            <p
              v-if="errors.email"
              class="font-body text-jv-accent-red text-xs px-0.5 m-0 flex items-center gap-1"
            >
              <AlertCircle class="size-3.5 shrink-0" :stroke-width="2.2" />
              {{ errors.email }}
            </p>
          </div>

          <!-- Password -->
          <div class="flex flex-col gap-1.5">
            <label
              for="password"
              class="font-body text-jv-ink text-xs sm:text-[13px] font-bold tracking-wide uppercase px-0.5"
            >
              Password
            </label>
            <div
              class="flex items-center gap-2.5 bg-jv-white border-2 border-jv-ink jv-card shadow-brutal-sm focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none transition-all px-3 py-2.5"
            >
              <Lock
                class="size-[18px] text-jv-ink/70 shrink-0"
                :stroke-width="2.2"
              />
              <input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                name="password"
                class="flex-1 min-w-0 bg-transparent border-0 outline-none font-body text-jv-ink placeholder:text-jv-ink/40 text-sm sm:text-base"
                placeholder="••••••••"
                required
              />
              <button
                type="button"
                class="grid place-items-center size-6 text-jv-ink/70 hover:text-jv-ink transition-colors"
                :aria-label="showPassword ? 'Hide password' : 'Show password'"
                @click="showPassword = !showPassword"
              >
                <EyeOff
                  v-if="showPassword"
                  class="size-[18px]"
                  :stroke-width="2.2"
                />
                <Eye v-else class="size-[18px]" :stroke-width="2.2" />
              </button>
            </div>
            <ul
              v-if="submitted && passwordErrors.length"
              class="list-none p-0 m-0 mt-0.5 flex flex-col gap-0.5"
            >
              <li
                v-for="(err, i) in passwordErrors"
                :key="i"
                class="font-body text-jv-accent-red text-xs px-0.5 m-0 flex items-center gap-1"
              >
                <AlertCircle class="size-3.5 shrink-0" :stroke-width="2.2" />
                {{ err }}
              </li>
            </ul>
            <p
              v-else-if="errors.password"
              class="font-body text-jv-accent-red text-xs px-0.5 m-0 flex items-center gap-1"
            >
              <AlertCircle class="size-3.5 shrink-0" :stroke-width="2.2" />
              {{ errors.password }}
            </p>
          </div>

          <input type="hidden" name="csrf_token" :value="csrfToken" />
          <input type="hidden" name="method" value="password" />

          <!-- Submit button -->
          <NavigationLink
            type="submit"
            class="mt-1 w-full bg-jv-coral text-white py-2.5 sm:py-3 text-sm sm:text-base"
            :disabled="!csrfToken || !registerURLWithFlowQuery"
          >
            <template v-if="!csrfToken || !registerURLWithFlowQuery">
              <Loader2 class="size-[18px] animate-spin" :stroke-width="2.4" />
              <span>Loading…</span>
            </template>
            <template v-else>
              <span>Sign Up</span>
            </template>
          </NavigationLink>
        </form>

        <!-- Footer link -->
        <p class="text-center mt-5 font-body text-jv-ink/70 text-sm m-0">
          Already have an account?
          <NuxtLink
            to="/account/login"
            class="text-jv-coral underline underline-offset-4 decoration-2 hover:decoration-jv-coral font-semibold ml-1"
          >
            Sign In
          </NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
