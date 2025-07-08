// https://nuxt.com/docs/api/configuration/nuxt-config
import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

export default defineNuxtConfig({
  // please refer https://nuxt.com/docs/guide/going-further/runtime-config#environment-variables for setting up environment variables.
  runtimeConfig: {
    public: {
      baseUrl: process.env.NUXT_PUBLIC_BASE_URL || "http://127.0.0.1:3001",
      apiUrl: process.env.NUXT_PUBLIC_API_URL || "http://127.0.0.1:3000/api/v1",
      apiSocketUrl:
        process.env.NUXT_PUBLIC_API_SOCKET_URL ||
        "ws://127.0.0.1:3000/api/v1/socket",
      kratosUrl: process.env.NUXT_PUBLIC_KRATOS_URL || "http://127.0.0.1:4433",
    },
  },
  app: {
    head: {
      title: "Jovvix",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
      ],
      script: [
        //This is just for example how to add js
        //you can  include js  by this method direact include or via import individual method as per below link
        //https://github.com/Debonex/samples/blob/master/nuxt3-bootstrap5/app.vue
      ],
      link: [
        //This for just example how to add css
        { rel: "stylesheet", href: "" },
      ],
    },
  },

  css: ["@/assets/scss/theme.scss"], // add
  modules: [
    "@nuxt/test-utils/module",
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", "acceptHMRUpdate"],
      },
    ],
    "@pinia-plugin-persistedstate/nuxt",
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        // @ts-expect-error error 'config.plugins' is possibly 'undefined'
        config.plugins.push(vuetify({ autoImport: true }));
      });
    },
  ],

  vite: {
    define: {
      "process.env.DEBUG": false,
    },
    vue: {
      template: {
        transformAssetUrls,
      },
    },
  },

  build: {
    transpile: [
      "vue-toastification",
      "vuetify",
      "@fortawesome/vue-fontawesome",
      "@fortawesome/fontawesome-svg-core",
      "@fortawesome/pro-solid-svg-icons",
      "@fortawesome/pro-regular-svg-icons",
      "@fortawesome/free-brands-svg-icons",
    ],
  },

  plugins: ["@/plugins/chart.js"],
});
