import globals from "globals";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";

export default tseslint.config(
  {
    ignores: [
      "node_modules",
      ".nuxt",
      ".nitro",
      ".cache",
      ".output",
      ".data",
      "dist",
      "coverage",
    ],
  },

  ...tseslint.configs.recommended,
  ...pluginVue.configs["flat/recommended"],
  eslintPluginPrettierRecommended,

  {
    languageOptions: {
      ecmaVersion: 2020,
      sourceType: "module",
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
  },

  {
    files: ["**/*.vue"],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
      },
    },
  },

  {
    rules: {
      "vue/multi-word-component-names": 0,
      // Added to "recommended" in eslint-plugin-vue 10; existing components
      // intentionally mark props as required while still providing defaults.
      "vue/no-required-prop-with-default": "off",
    },
  }
);
