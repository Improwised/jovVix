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
      htmlAttrs: {
        lang: "en",
      },
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
        // Preconnect to Google Fonts for faster loading
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "",
        },
        // Load Google Fonts with font-display: swap - critical for performance
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap",
          media: "print",
          onload: "this.media='all'",
        },
              // Add your favicon here
      { rel: "icon", type: "image/x-icon", href: "/favicon.ico" },
      ],
    },
  },

  css: [
    "@/assets/scss/theme.scss",
    "@fortawesome/fontawesome-svg-core/styles.css",
  ],
  modules: [
    "@nuxt/test-utils/module",
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", "acceptHMRUpdate"],
      },
    ],
    (_options, nuxt) => {
      nuxt.hooks.hook("vite:extendConfig", (config) => {
        // @ts-expect-error error 'config.plugins' is possibly 'undefined'
        config.plugins.push(vuetify({ autoImport: true }));
      });
    },
  ],

  vite: {
    // Temporary solution to silence Bootstrap SCSS deprecation warnings
    // Reference: https://github.com/twbs/bootstrap/issues/40962
    css: {
      preprocessorOptions: {
        scss: {
          silenceDeprecations: [
            "mixed-decls",
            "color-functions",
            "global-builtin",
            "import",
          ],
        },
      },
    },
    define: {
      "process.env.DEBUG": false,
    },
    vue: {
      template: {
        transformAssetUrls,
      },
    },
    build: {
      // Code splitting optimizations - only include actual JS modules
      rollupOptions: {
        output: {
          manualChunks: (id) => {
            // Chunk large libraries separately
            if (id.includes("node_modules")) {
              if (id.includes("vuetify")) return "vendor-vuetify";
              if (id.includes("chart.js")) return "vendor-charts";
              if (id.includes("bootstrap")) return "vendor-bootstrap";
              if (id.includes("@fortawesome")) return "vendor-icons";
              if (id.includes("highlight.js")) return "vendor-highlight";
            }
          },
        },
      },
      // Reduce chunk size warnings
      chunkSizeWarningLimit: 1000,
    },
    // Optimize dependencies - only JS modules
    optimizeDeps: {
      include: ["vuetify", "chart.js", "bootstrap"],
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

  // Performance optimizations
  ssr: true, // Enable SSR for better performance
  experimental: {
    payloadExtraction: false, // Improve initial load
  },

  // Critical performance optimizations
  nitro: {
    compressPublicAssets: true, // Enable compression
    minify: true, // Minify output
  },

  plugins: ["@/plugins/chart.js"],
});
