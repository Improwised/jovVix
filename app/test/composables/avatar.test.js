import { describe, expect, it } from "vitest";
import {
  getAvatarUrlByName,
  getRandomAvatarName,
} from "../../composables/avatar";

const knownAvatars = [
  "Sophia",
  "Jude",
  "Jade",
  "Ryan",
  "Adrian",
  "Chase",
  "Nolan",
  "Sadie",
  "Brian",
  "Aidan",
  "Destiny",
  "Maria",
  "Kingston",
  "Andrea",
  "Vivian",
  "Eden",
  "Wyatt",
  "Sawyer",
  "Jocelyn",
];

describe("getAvatarUrlByName", () => {
  it("returns a dicebear URL for a known name", () => {
    const url = getAvatarUrlByName("Sophia");
    expect(url).toBe("https://api.dicebear.com/10.x/clay/svg?seed=Sophia");
  });

  it("falls back to Eden for unknown names", () => {
    const url = getAvatarUrlByName("UnknownPerson");
    expect(url).toBe("https://api.dicebear.com/10.x/clay/svg?seed=Eden");
  });
});

describe("getRandomAvatarName", () => {
  it("returns one of the known avatar names", () => {
    const name = getRandomAvatarName();
    expect(knownAvatars).toContain(name);
  });
});
