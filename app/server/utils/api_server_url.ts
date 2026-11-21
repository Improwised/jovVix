// The browser reaches the API at the public URL (NUXT_PUBLIC_API_URL). The Nuxt server may not:
// in the compose stack that URL is 127.0.0.1:3000, which inside the web container is the web
// container itself. NUXT_API_SERVER_URL names the API as the server sees it (http://api:3000/api/v1).

// Returns the URL to fetch on the server: a URL under publicUrl moves to serverUrl, others stay.
export function toServerApiUrl(
  url: string,
  publicUrl: string,
  serverUrl: string
): string {
  if (!serverUrl || !publicUrl) return url;
  const from = publicUrl.replace(/\/+$/, "");
  if (
    url !== from &&
    !url.startsWith(from + "/") &&
    !url.startsWith(from + "?")
  )
    return url;
  return serverUrl.replace(/\/+$/, "") + url.slice(from.length);
}

// Wraps a fetch so every request for the public API URL goes to the server URL instead.
export function withServerApiUrl(
  fetchFn: typeof fetch,
  publicUrl: string,
  serverUrl: string
): typeof fetch {
  return (input, init) => {
    if (typeof input === "string" || input instanceof URL) {
      return fetchFn(toServerApiUrl(String(input), publicUrl, serverUrl), init);
    }
    const moved = toServerApiUrl(input.url, publicUrl, serverUrl);
    return fetchFn(
      moved === input.url ? input : new Request(moved, input),
      init
    );
  };
}
