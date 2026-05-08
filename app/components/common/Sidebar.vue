<template>
  <!-- Desktop sidebar (lg and up) -->
  <aside
    class="hidden lg:flex bg-jv-ivory min-w-56 h-screen flex-col justify-between border-r-[4px] border-jv-ink px-8 py-5 overflow-y-hidden"
  >
    <div class="flex flex-col gap-4">
      <NuxtLink
        to="/"
        class="flex items-center gap-[6px] text-[30px] font-black tracking-[-0.75px] text-jv-ink no-underline"
      >
        <span
          class="grid size-9 rotate-[-6deg] place-items-center rounded-[135px_120px_120px_135px] border-2 border-jv-ink bg-jv-yellow text-jv-ink shadow-[2px_2px_0_#2D2D2D]"
        >
          <Zap class="size-5 rotate-[-6deg]" :stroke-width="2.4" />
        </span>
        <span>JovVix</span>
      </NuxtLink>

      <NavigationLink url="/admin/quiz/create-quiz" url-name="Create Quiz">
        <Plus class="size-[18px]" />
      </NavigationLink>

      <nav class="mt-5 flex flex-col gap-2">
        <NuxtLink
          v-for="item in navItems"
          :key="item.url"
          :to="item.url"
          :class="navItemClass(item)"
          class="flex items-center gap-3 px-4 py-2.5 text-[18px] font-semibold text-jv-ink/65 no-underline transition-transform hover:rotate-[1deg] hover:text-jv-ink"
        >
          <component
            :is="item.icon"
            class="size-5"
            :class="item.active ? 'text-jv-coral' : 'text-jv-ink/45'"
            :stroke-width="2.3"
          />
          <span>{{ item.label }}</span>
        </NuxtLink>
      </nav>
    </div>

    <div class="flex flex-col gap-4">
      <button
        v-if="showAdminNav"
        type="button"
        class="w-fit px-0 text-left text-[18px] font-semibold text-jv-ink/60 underline decoration-wavy underline-offset-2 transition-colors hover:text-jv-ink"
        @click="handleLogout"
      >
        Log Out
      </button>
      <template v-else>
        <NavigationLink url="/account/login" url-name="Sign In" />
        <NavigationLink url="/account/register" url-name="Sign Up" />
      </template>
    </div>
  </aside>

  <!-- Mobile/Tablet header (md and below) -->
  <header
    class="lg:hidden sticky top-0 z-40 border-b-[3px] border-jv-ink bg-jv-ivory/95 px-4 sm:px-6 py-3 backdrop-blur"
  >
    <div class="flex items-center justify-between gap-3">
      <NuxtLink
        to="/"
        class="flex items-center gap-[6px] text-[22px] sm:text-[26px] font-black tracking-[-0.75px] text-jv-ink no-underline"
      >
        <span
          class="grid size-8 sm:size-9 rotate-[-6deg] place-items-center rounded-[135px_120px_120px_135px] border-2 border-jv-ink bg-jv-yellow text-jv-ink shadow-[2px_2px_0_#2D2D2D]"
        >
          <Zap class="size-4 sm:size-5 rotate-[-6deg]" :stroke-width="2.4" />
        </span>
        <span>JovVix</span>
      </NuxtLink>

      <button
        type="button"
        class="grid size-10 place-items-center rounded-[10px] border-[3px] border-jv-ink bg-jv-yellow text-jv-ink shadow-brutal-sm active:shadow-none active:translate-x-[2px] active:translate-y-[2px] transition-transform"
        :aria-expanded="open"
        aria-label="Toggle navigation menu"
        aria-controls="mobile-nav"
        @click="open = !open"
      >
        <X v-if="open" class="size-5" :stroke-width="2.4" />
        <Menu v-else class="size-5" :stroke-width="2.4" />
      </button>
    </div>

    <nav
      v-show="open"
      id="mobile-nav"
      class="mt-4 grid grid-cols-1 min-[420px]:grid-cols-2 sm:grid-cols-3 gap-3"
    >
      <NuxtLink
        v-for="item in mobileNavItems"
        :key="item.url"
        :to="item.url"
        :class="[
          'inline-flex items-center justify-center gap-1.5 jv-border-uneven px-4 py-2 text-sm font-headings text-jv-ink no-underline shadow-brutal-sm transition-transform hover:rotate-[2deg]',
          item.active ? 'bg-jv-yellow/70' : 'bg-jv-white',
          item.highlight ? 'bg-jv-coral !text-white' : '',
        ]"
        @click="open = false"
      >
        <component
          :is="item.icon"
          v-if="item.icon"
          class="size-4"
          :stroke-width="2.4"
        />
        <span>{{ item.label }}</span>
      </NuxtLink>
      <button
        v-if="showAdminNav"
        type="button"
        class="inline-flex items-center justify-center gap-1.5 jv-border-uneven bg-jv-white px-4 py-2 text-sm font-headings text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg]"
        @click="handleMobileLogout"
      >
        <LogOut class="size-4" :stroke-width="2.4" />
        <span>Log Out</span>
      </button>
    </nav>
  </header>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import {
  BarChart3,
  HelpCircle,
  Home,
  LogOut,
  Menu,
  Plus,
  UserRound,
  X,
  Zap,
} from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";
import { setUserDataStore } from "@/composables/auth";
import { useUsersStore } from "~~/store/users";

const route = useRoute();
const open = ref(false);
const userDataStore = useUsersStore();

const isQuizListPage = computed(() =>
  route.path.startsWith("/admin/quiz/list-quiz")
);

const isLoggedInAdmin = computed(
  () => userDataStore.getUserData()?.role === "admin-user"
);

const showAdminNav = computed(
  () => isLoggedInAdmin.value || isQuizListPage.value
);

const isActiveRoute = (url) => {
  if (url === "/") return route.path === "/";
  if (url === "/admin/quiz/list-quiz")
    return route.path.startsWith("/admin/quiz");
  if (url === "/admin/reports") return route.path.startsWith("/admin/reports");
  if (url === "/admin") return route.path === "/admin";

  return route.path === url;
};

const navItems = computed(() => {
  if (showAdminNav.value) {
    return [
      { label: "Home", url: "/", icon: Home, active: isActiveRoute("/") },
      {
        label: "Quizzes",
        url: "/admin/quiz/list-quiz",
        icon: HelpCircle,
        active: isActiveRoute("/admin/quiz/list-quiz"),
      },
      {
        label: "Reports",
        url: "/admin/reports",
        icon: BarChart3,
        active: isActiveRoute("/admin/reports"),
      },
      {
        label: "Profile",
        url: "/admin",
        icon: UserRound,
        active: isActiveRoute("/admin"),
      },
    ];
  }

  return [{ label: "Home", url: "/", icon: Home, active: true }];
});

const mobileNavItems = computed(() => {
  if (showAdminNav.value) {
    return [
      { label: "Home", url: "/", icon: Home, active: isActiveRoute("/") },
      {
        label: "Quizzes",
        url: "/admin/quiz/list-quiz",
        icon: HelpCircle,
        active: isActiveRoute("/admin/quiz/list-quiz"),
      },
      {
        label: "Reports",
        url: "/admin/reports",
        icon: BarChart3,
        active: isActiveRoute("/admin/reports"),
      },
      {
        label: "Profile",
        url: "/admin",
        icon: UserRound,
        active: isActiveRoute("/admin"),
      },
      {
        label: "Create Quiz",
        url: "/admin/quiz/create-quiz",
        icon: Plus,
        highlight: true,
      },
    ];
  }

  return [
    { label: "Home", url: "/", icon: Home, active: true },
    { label: "Enter Code", url: "/join", highlight: true },
    { label: "Create Quiz", url: "/admin/quiz/create-quiz", icon: Plus },
    { label: "Sign In", url: "/account/login" },
    { label: "Sign Up", url: "/account/register" },
  ];
});

const navItemClass = (item) =>
  item.active
    ? "rounded-full border-2 border-dashed border-jv-ink/40 bg-jv-yellow/20 text-jv-ink"
    : "";

const handleMobileLogout = async () => {
  open.value = false;
  await handleLogout();
};

onMounted(() => {
  setUserDataStore();
});
</script>
