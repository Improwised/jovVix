// https://nuxt.com/docs/api/configuration/nuxt-config
import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      base_url: process.env.BASE_URL,
      api_url: process.env.API_URL,
      socket_url: process.env.API_SOCKET_URL,
      kratos_url: process.env.KRATOS_URL,
    },
  },
  app: {
    head: {
      title: "Quiz App",
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

  css: [
    "@/assets/scss/theme.scss",
  ], // add
  modules: [
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
    transpile: ["vue-toastification", "vuetify"],
  },
});
