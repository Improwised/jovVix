// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      base_url: process.env.BASE_URL,
    },
  },
  app: {
    head: {
      title: "Nuxt3",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
      ],
      script: [
        //This is just for example how to add js
        //you can  include js  by this method direact include or via import individual method as per below link
        //https://github.com/Debonex/samples/blob/master/nuxt3-bootstrap5/app.vue
        { src: "https://awesome-lib.js" },
      ],
      link: [
        //This for just example how to add css
        { rel: "stylesheet", href: "https://awesome-lib.css" },
      ],
    },
  },

  css: ["bootstrap/dist/css/bootstrap.min.css"], // add
  modules: [
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", "acceptHMRUpdate"],
      },
    ],
  ],

  vite: {
    define: {
      "process.env.DEBUG": false,
    },
  },
});
