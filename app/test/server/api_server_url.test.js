import { describe, it, expect, vi } from "vitest";
import {
  toServerApiUrl,
  withServerApiUrl,
} from "~~/server/utils/api_server_url";

const PUBLIC = "http://127.0.0.1:3000/api/v1";
const SERVER = "http://api:3000/api/v1";

// In the compose stack the web container cannot reach the public API URL, so server-side
// rendering must fetch the API through the server URL.
describe("server API URL", () => {
  it("moves API URLs and leaves others alone", () => {
    expect(
      toServerApiUrl(`${PUBLIC}/quizzes/q1/questions`, PUBLIC, SERVER)
    ).toBe(`${SERVER}/quizzes/q1/questions`);
    expect(toServerApiUrl(`${PUBLIC}?x=1`, PUBLIC, SERVER)).toBe(
      `${SERVER}?x=1`
    );
    expect(toServerApiUrl(PUBLIC, PUBLIC + "/", SERVER + "/")).toBe(SERVER);
    // Same host, other path; another prefix sharing the text; other hosts.
    expect(
      toServerApiUrl("http://127.0.0.1:3000/metrics", PUBLIC, SERVER)
    ).toBe("http://127.0.0.1:3000/metrics");
    expect(toServerApiUrl(`${PUBLIC}2/x`, PUBLIC, SERVER)).toBe(`${PUBLIC}2/x`);
    expect(
      toServerApiUrl("http://127.0.0.1:4433/sessions/whoami", PUBLIC, SERVER)
    ).toBe("http://127.0.0.1:4433/sessions/whoami");
  });

  it("changes nothing when no server URL is set", () => {
    expect(toServerApiUrl(`${PUBLIC}/x`, PUBLIC, "")).toBe(`${PUBLIC}/x`);
  });

  it("sends string, URL and Request inputs to the server URL", async () => {
    const seen = [];
    const fetchFn = vi.fn(async (input, init) => {
      seen.push([typeof input === "string" ? input : input.url, init]);
      return new Response("{}");
    });
    const wrapped = withServerApiUrl(fetchFn, PUBLIC, SERVER);
    const init = { headers: { cookie: "ory_kratos_session=s" } };

    await wrapped(`${PUBLIC}/user/who`, init);
    await wrapped(new URL(`${PUBLIC}/categories`));
    await wrapped(
      new Request(`${PUBLIC}/quizzes`, { method: "POST", body: "{}" })
    );
    await wrapped("https://example.com/x");

    expect(seen.map(([url]) => url)).toEqual([
      `${SERVER}/user/who`,
      `${SERVER}/categories`,
      `${SERVER}/quizzes`,
      "https://example.com/x",
    ]);
    // Headers such as the forwarded cookie go along unchanged.
    expect(seen[0][1]).toBe(init);
  });
});
