// Import only specific Vuetify components and styles used in the app
import { createVuetify } from "vuetify";
import {
  VBtn,
  VCard,
  VCardText,
  VCardTitle,
  VFileInput,
  VListItemTitle,
  VOtpInput,
  VProgressCircular,
  VProgressLinear,
  VSkeletonLoader,
} from "vuetify/components";
import { aliases, mdi } from "vuetify/iconsets/mdi-svg";
// Import only essential Vuetify styles - use the correct paths
import "vuetify/styles";
import "~/assets/scss/theme.scss";

export default defineNuxtPlugin((app) => {
  const vuetify = createVuetify({
    ssr: true,
    components: {
      VBtn,
      VCard,
      VCardText,
      VCardTitle,
      VFileInput,
      VSkeletonLoader,
      VProgressCircular,
      VProgressLinear,
      VListItemTitle,
      VOtpInput,
    },
    icons: {
      defaultSet: "mdi",
      aliases,
      sets: {
        mdi,
      },
    },
    theme: {
      defaultTheme: "light",
    },
  });
  app.vueApp.use(vuetify);
});
