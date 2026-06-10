// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  // please refer https://nuxt.com/docs/guide/going-further/runtime-config#environment-variables for setting up environment variables.
  runtimeConfig: {
    public: {
      baseUrl: process.env.NUXT_PUBLIC_BASE_URL || "http://127.0.0.1:3001",
      apiUrl: process.env.NUXT_PUBLIC_API_URL || "http://127.0.0.1:3000/api/v1",
      maxImageFileSize: parseInt(
        process.env.NUXT_PUBLIC_MAX_IMAGE_FILE_SIZE || "512000"
      ), // 500 KB default (bytes)
      apiSocketUrl:
        process.env.NUXT_PUBLIC_API_SOCKET_URL ||
        "ws://127.0.0.1:3000/api/v1/socket",
      kratosUrl: process.env.NUXT_PUBLIC_KRATOS_URL || "http://127.0.0.1:4433",
      privilegedSessionMaxAge: parseInt(
        process.env.NUXT_PUBLIC_PRIVILEGED_SESSION_MAX_AGE || "15"
      ),
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
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "",
        },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap",
          media: "print",
          onload: "this.media='all'",
        },
        { rel: "icon", type: "image/x-icon", href: "/favicon.ico" },
        {
          rel: "icon",
          type: "image/png",
          sizes: "16x16",
          href: "/favicon-16x16.png",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "32x32",
          href: "/favicon-32x32.png",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "48x48",
          href: "/favicon-48x48.png",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "180x180",
          href: "/favicon-180x180.png",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "256x256",
          href: "/favicon-256x256.png",
        },
      ],
    },
  },

  css: [
    "@/assets/css/main.css",
    "notivue/notification.css",
    "notivue/animations.css",
    "notivue/notification-progress.css",
  ],
  modules: [
    "@nuxt/test-utils/module",
    [
      "@pinia/nuxt",
      {
        autoImports: ["defineStore", "acceptHMRUpdate"],
      },
    ],
    "shadcn-nuxt",
    "notivue/nuxt",
  ],
  notivue: {
    position: "top-center",
    limit: 4,
    enqueue: true,
    avoidDuplicates: true,
    notifications: {
      global: {
        duration: 8000,
      },
    },
  },
  shadcn: {
    prefix: "",
    componentDir: "@/components/ui",
  },
  vite: {
    define: {
      "process.env.DEBUG": false,
    },
    vue: {
      template: {
        transformAssetUrls: {
          includeAbsolute: false,
        },
      },
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: (id) => {
            if (id.includes("node_modules")) {
              if (id.includes("chart.js")) return "vendor-charts";
              if (id.includes("highlight.js")) return "vendor-highlight";
            }
          },
        },
      },
      chunkSizeWarningLimit: 1000,
    },
    optimizeDeps: {
      include: ["chart.js"],
    },
    plugins: [tailwindcss()],
  },

  ssr: true,
  experimental: {
    payloadExtraction: false,
  },

  nitro: {
    compressPublicAssets: true,
    minify: true,
  },

  plugins: ["@/plugins/chart.js"],
});
