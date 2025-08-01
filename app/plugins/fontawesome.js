import { library, config } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { fab } from "@fortawesome/free-brands-svg-icons";

// Configure FontAwesome to let Nuxt handle CSS injection for better performance
config.autoAddCss = true;

// Add selected icon packs (solid and brands) to the FontAwesome library
library.add(fas, fab);

// Register the FontAwesomeIcon component globally in the Nuxt app
export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.component("font-awesome-icon", FontAwesomeIcon);
});
