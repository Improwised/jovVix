import { library, config } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { fab } from "@fortawesome/free-brands-svg-icons";

// This is important, we are going to let Nuxt worry about the CSS
config.autoAddCss = true;

// You can add your icons directly in this plugin. See other examples for how you
// can add other styles or just individual icons.
library.add(fas, fab);

export default defineNuxtPlugin((nuxtApp) => {
  css: ["@fortawesome/fontawesome-svg-core/styles.css"],
    nuxtApp.vueApp.component("font-awesome-icon", FontAwesomeIcon, {});
});
