// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  ssr: true,
  modules: [
    "@nuxtjs/tailwindcss",
    "shadcn-nuxt",
    "@nuxtjs/color-mode",
    "@nuxt/fonts",
    "@vueuse/nuxt",
  ],
  fonts: {
    defaults: {
      weights: [400],
      styles: ["normal", "italic"],
      subsets: [
        "cyrillic-ext",
        "cyrillic",
        "greek-ext",
        "greek",
        "vietnamese",
        "latin-ext",
        "latin",
      ],
    },
    families: [
      { name: "Dela Gothic One", provider: "google" },
      { name: "Inter", provider: "google" },
      { name: "Geist", provider: "google" },
    ],
  },
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: "",
    /**
     * Directory that the component lives in.
     * @default "./components/ui"
     */
    componentDir: "./components/ui",
  },
  colorMode: {
    preference: "system", // Default theme
    dataValue: "theme", // Adds `data-theme="dark"` to <html>
    classSuffix: "",
  },
  runtimeConfig: {
    public: {
      googleAuth: "",
      googleRedirectUrl: "",
      apiUrl: "",
      siteUrl: "",
      sessionKey: "syllabye.session",
    },
  },
  app: {
    head: {
      title: "Syllabye",
      htmlAttrs: {
        lang: "en",
      },
      link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
    },
  },
});
