import { withServerApiUrl } from "../utils/api_server_url";

// Server-side rendering fetches the API through globalThis.fetch (useFetch, $fetch and plain
// fetch all end there), so one wrapper covers every page. Unset NUXT_API_SERVER_URL: no change.
export default defineNitroPlugin(() => {
  const config = useRuntimeConfig();
  const serverUrl = config.apiServerUrl as string;
  if (!serverUrl || serverUrl === config.public.apiUrl) return;
  globalThis.fetch = withServerApiUrl(
    globalThis.fetch,
    config.public.apiUrl as string,
    serverUrl
  );
});
