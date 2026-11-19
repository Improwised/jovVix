import { test, expect } from "@playwright/test";
import { execFileSync } from "node:child_process";
import { randomUUID } from "node:crypto";
const api = process.env.ROOM_TEST_API_URL || "http://127.0.0.1:3000/api/v1";
const sql = (query) =>
  execFileSync(
    "docker",
    [
      "exec",
      process.env.ROOM_TEST_DB_CONTAINER || "jovvix-db",
      "psql",
      "-v",
      "ON_ERROR_STOP=1",
      "-U",
      "jovvix",
      "-d",
      "jovvix",
      "-At",
      "-c",
      query,
    ],
    { encoding: "utf8" },
  ).trim();
const quizA = randomUUID(),
  quizB = randomUUID(),
  questionA = randomUUID(),
  questionB = randomUUID();
const suffix = Date.now().toString().slice(-6);
const titleA = `Room A ${suffix}`,
  titleB = `Room B ${suffix}`;
async function data(response) {
  const body = await response.json();
  expect(response.ok(), JSON.stringify(body)).toBeTruthy();
  return body.data;
}
async function guest(context, name) {
  return data(
    await context.request.post(`${api}/user/${name}?avatar_name=Bear`),
  );
}
async function join(page, code, name) {
  await page.goto("/join");
  await page.getByLabel("Invitation Code").fill(String(code));
  await page.getByLabel("Your quiz name").fill(name);
  await page.getByRole("button", { name: "Enter Lobby" }).click();
  await expect(page).toHaveURL(/\/join\/play\//);
}
test.beforeAll(() => {
  sql(`INSERT INTO quizzes(id,title,is_public) VALUES('${quizA}','${titleA}',true),('${quizB}','${titleB}',true);
 INSERT INTO questions(id,question,options,answers,points,duration_in_seconds,type) VALUES
 ('${questionA}','Question A ${suffix}','{"1":"Correct","2":"Wrong"}','[1]',10,30,1),
 ('${questionB}','Question B ${suffix}','{"1":"Correct","2":"Wrong"}','[1]',10,30,1);
 INSERT INTO quiz_questions(id,quiz_id,question_id) VALUES('${randomUUID()}','${quizA}','${questionA}'),('${randomUUID()}','${quizB}','${questionB}');`);
});
test.afterAll(() => {
  sql(`BEGIN;
 DELETE FROM user_quiz_responses WHERE user_played_quiz_id IN(SELECT p.id FROM user_played_quizzes p JOIN active_quizzes a ON a.id=p.active_quiz_id WHERE a.quiz_id IN('${quizA}','${quizB}'));
 DELETE FROM user_played_quizzes WHERE active_quiz_id IN(SELECT id FROM active_quizzes WHERE quiz_id IN('${quizA}','${quizB}'));
 DELETE FROM active_quiz_questions WHERE active_quiz_id IN(SELECT id FROM active_quizzes WHERE quiz_id IN('${quizA}','${quizB}'));
 UPDATE active_quizzes SET room_id=NULL WHERE quiz_id IN('${quizA}','${quizB}');
 DELETE FROM quiz_room_members WHERE room_id IN(SELECT r.id FROM quiz_rooms r JOIN active_quizzes a ON a.id=r.current_session_id WHERE a.quiz_id IN('${quizA}','${quizB}'));
 DELETE FROM quiz_rooms WHERE current_session_id IN(SELECT id FROM active_quizzes WHERE quiz_id IN('${quizA}','${quizB}'));
 DELETE FROM active_quizzes WHERE quiz_id IN('${quizA}','${quizB}');
 DELETE FROM quiz_questions WHERE quiz_id IN('${quizA}','${quizB}');DELETE FROM quizzes WHERE id IN('${quizA}','${quizB}');DELETE FROM questions WHERE id IN('${questionA}','${questionB}');COMMIT;`);
});
test("retains players and code across quizzes, recovers refreshes, admits new players, and closes safely", async ({
  browser,
}) => {
  const contexts = await Promise.all(
    Array.from({ length: 4 }, () => browser.newContext()),
  );
  const [hostCtx, oneCtx, twoCtx, newCtx] = contexts;
  try {
    await guest(hostCtx, `Host${suffix}`);
    const sessionA = await data(
      await hostCtx.request.post(`${api}/quizzes/${quizA}/public_session`),
    );
    const host = await hostCtx.newPage();
    const hostErrors = [];
    host.on("pageerror", (error) => hostErrors.push(error.message));
    await host.goto(`/admin/arrange/${sessionA}`);
    await expect(
      host.getByRole("button", { name: "Start Quiz", exact: true }),
    ).toBeVisible();
    const roomA = await data(
      await hostCtx.request.get(`${api}/rooms/lookup?session_id=${sessionA}`),
    );
    const code = roomA.code,
      roomID = roomA.id;
    const one = await oneCtx.newPage(),
      two = await twoCtx.newPage();
    await join(one, code, "Alpha");
    await join(two, code, "Beta");
    await expect(host.getByText("Alpha", { exact: true })).toBeVisible();
    await expect(host.getByText("Beta", { exact: true })).toBeVisible();
    const participantA = await data(
      await oneCtx.request.post(`${api}/user_played_quizes/${code}`),
    );
    // A player cannot operate the room or read another host's choices.
    expect(
      (await oneCtx.request.get(`${api}/rooms/${roomID}/quizzes`)).status(),
    ).toBe(403);
    expect(
      (await oneCtx.request.post(`${api}/rooms/${roomID}/close`)).status(),
    ).toBe(403);
    await two.reload();
    await expect(
      two.getByRole("heading", { name: "Ready, Steady, Go!" }),
    ).toBeVisible();
    await host.getByRole("button", { name: "Start Quiz", exact: true }).click();
    await expect(
      one.getByText(`Question A ${suffix}`, { exact: true }),
    ).toBeVisible();
    await one.getByText("Correct", { exact: true }).click();
    // Stop quiz through the same API used by End Quiz.
    await host
      .getByRole("button", { name: "End quiz for all players" })
      .click();
    await host
      .getByRole("dialog")
      .getByRole("button", { name: "End Quiz", exact: true })
      .click();
    await expect(host).toHaveURL(/\/admin\/scoreboard/);
    await expect(one).toHaveURL(/\/scoreboard/);
    await expect(two).toHaveURL(/\/scoreboard/);
    await host.getByRole("link", { name: "Choose next quiz" }).click();
    await expect(host.getByLabel("Choose the next quiz")).toBeVisible();
    await host.reload();
    await expect(host.getByLabel("Choose the next quiz")).toBeVisible();
    await one.reload();
    await expect(
      one.getByText("Waiting for the host to choose the next quiz"),
    ).toBeVisible();
    // Explicit leave removes membership; stale play links cannot silently rejoin.
    await two.getByRole("button", { name: "Leave room", exact: true }).click();
    await expect(two).toHaveURL(/\/join$/);
    expect((await twoCtx.request.get(`${api}/rooms/${roomID}`)).status()).toBe(
      403,
    );
    // New players join while prior results remain on screen.
    const newcomer = await newCtx.newPage();
    await newcomer.goto("/join");
    await newcomer.getByLabel("Invitation Code").fill(String(code));
    await newcomer.getByLabel("Your quiz name").fill("Gamma");
    await newcomer.getByRole("button", { name: "Enter Lobby" }).click();
    await expect(newcomer).toHaveURL(new RegExp(`/room/${roomID}`));
    await expect(
      newcomer.getByText("Waiting for the host", { exact: false }),
    ).toBeVisible();
    // A returning player stays on results while the host is choosing.
    await one.goto(`/join/play/${code}`);
    await expect(one).toHaveURL(/\/scoreboard/);
    await expect(
      one.getByText("Waiting for the host to choose the next quiz"),
    ).toBeVisible();
    await expect(one.getByRole("link", { name: "Waiting room" })).toHaveCount(
      0,
    );
    await host.screenshot({
      path: "/tmp/jovvix-room-host.png",
      fullPage: true,
    });
    await host.getByRole("combobox", { name: "Choose the next quiz" }).click();
    await expect(host.getByRole("button", { name: "Private", exact: true })).toHaveAttribute("aria-pressed", "true");
    await host.getByRole("button", { name: "Public", exact: true }).click();
    await host.getByRole("option", { name: titleB }).click();
    await host.getByRole("button", { name: "Continue to lobby" }).click();
    await expect(host).toHaveURL(/\/admin\/arrange\//);
    await expect(one).toHaveURL(new RegExp(`/join/play/${code}`));
    await expect(newcomer).toHaveURL(new RegExp(`/join/play/${code}`));
    const roomB = await data(
      await hostCtx.request.get(`${api}/rooms/${roomID}`),
    );
    expect(roomB.code).toBe(code);
    expect(roomB.session_id).not.toBe(sessionA);
    expect(roomB.members.map((m) => m.name).sort()).toEqual(["Alpha", "Gamma"]);
    const participantB = await data(
      await oneCtx.request.post(`${api}/user_played_quizes/${code}`),
    );
    expect(participantB.user_played_quiz).not.toBe(
      participantA.user_played_quiz,
    );
    expect(participantB.player_name).toBe("Alpha");
    expect(
      sql(
        `SELECT sum(calculated_score) FROM user_quiz_responses WHERE user_played_quiz_id='${participantB.user_played_quiz}'`,
      ),
    ).toBe("0");
    // Duplicate/stale selection cannot create an extra round.
    expect(
      (
        await hostCtx.request.post(`${api}/rooms/${roomID}/next`, {
          data: { quiz_id: quizB, expected_session_id: sessionA },
        })
      ).status(),
    ).toBe(409);
    expect(
      (await twoCtx.request.post(`${api}/user_played_quizes/${code}`)).status(),
    ).toBe(403);
    await one.reload();
    await expect(
      one.getByRole("heading", { name: "Ready, Steady, Go!" }),
    ).toBeVisible();
    await host.reload();
    await expect(
      host.getByRole("button", { name: "Start Quiz", exact: true }),
    ).toBeVisible();
    await host.getByRole("button", { name: "Start Quiz", exact: true }).click();
    await expect(
      one.getByText(`Question B ${suffix}`, { exact: true }),
    ).toBeVisible();
    await expect(
      newcomer.getByText(`Question B ${suffix}`, { exact: true }),
    ).toBeVisible();
    const stale = await oneCtx.request.post(
      `${api}/quiz/answer?user_played_quiz=${participantA.user_played_quiz}&session_id=${sessionA}`,
      { data: { id: questionA, keys: [1], response_time: 100 } },
    );
    expect(stale.ok()).toBeFalsy();
    await host
      .getByRole("button", { name: "End quiz for all players" })
      .click();
    await host
      .getByRole("dialog")
      .getByRole("button", { name: "End Quiz", exact: true })
      .click();
    await host.getByRole("link", { name: "Choose next quiz" }).click();
    await expect(host.getByLabel("Choose the next quiz")).toBeVisible();
    expect(
      (
        await oneCtx.request.get(
          `${api}/final_score/user?user_played_quiz=${participantA.user_played_quiz}`,
        )
      ).ok(),
    ).toBeTruthy();
    await data(await hostCtx.request.post(`${api}/rooms/${roomID}/close`));
    await expect(
      one.getByText("The host ended this room.", { exact: true }),
    ).toBeVisible();
    expect(hostErrors).toEqual([]);
    expect(
      (
        await newCtx.request.post(`${api}/rooms/join/${code}`, {
          data: { player_name: "Gamma" },
        })
      ).status(),
    ).toBe(410);
  } finally {
    await Promise.all(contexts.map((c) => c.close()));
  }
});

test("recovers a missed quiz end and next-round selection after a socket outage", async ({
  browser,
}) => {
  const hostCtx = await browser.newContext(),
    playerCtx = await browser.newContext();
  try {
    await guest(hostCtx, `Rec${suffix}`);
    const sessionA = await data(
      await hostCtx.request.post(`${api}/quizzes/${quizA}/public_session`),
    );
    const host = await hostCtx.newPage();
    const hostErrors = [];
    host.on("pageerror", (error) => hostErrors.push(error.message));
    await host.goto(`/admin/arrange/${sessionA}`);
    await expect(
      host.getByRole("button", { name: "Start Quiz", exact: true }),
    ).toBeVisible();
    const room = await data(
      await hostCtx.request.get(`${api}/rooms/lookup?session_id=${sessionA}`),
    );
    const player = await playerCtx.newPage();
    let refused = 0;
    // Refuse the first two handshakes, then allow the single retry scheduler through.
    await player.routeWebSocket("**/api/v1/socket/join/**", (socket) => {
      if (refused++ < 2)
        socket.close({ code: 1013, reason: "Test connection outage" });
      else socket.connectToServer();
    });
    await join(player, room.code, "Recover");
    await expect(
      player.getByRole("heading", { name: "Ready, Steady, Go!" }),
    ).toBeVisible();
    await expect.poll(() => refused).toBeGreaterThanOrEqual(3);
    await expect(
      player.getByText("Connecting to lobby...", { exact: true }),
    ).toHaveCount(0);
    await player.reload();
    await expect(
      player.getByRole("heading", { name: "Ready, Steady, Go!" }),
    ).toBeVisible();
    await host.getByRole("button", { name: "Start Quiz", exact: true }).click();
    await expect(
      player.getByText(`Question A ${suffix}`, { exact: true }),
    ).toBeVisible();
    await playerCtx.setOffline(true);
    await host
      .getByRole("button", { name: "End quiz for all players" })
      .click();
    await host
      .getByRole("dialog")
      .getByRole("button", { name: "End Quiz", exact: true })
      .click();
    await host.getByRole("link", { name: "Choose next quiz" }).click();
    await expect(host.getByLabel("Choose the next quiz")).toBeVisible();
    await host.getByRole("combobox", { name: "Choose the next quiz" }).click();
    await host.getByRole("button", { name: "Public", exact: true }).click();
    await host.getByRole("option", { name: titleB }).click();
    await host.getByRole("button", { name: "Continue to lobby" }).click();
    await expect(
      host.getByRole("button", { name: "Start Quiz", exact: true }),
    ).toBeVisible();
    const next = await data(
      await hostCtx.request.get(`${api}/rooms/${room.id}`),
    );
    await playerCtx.setOffline(false);
    await expect(player).toHaveURL(new RegExp(`/join/play/${room.code}`));
    await expect(
      player.getByRole("heading", { name: "Ready, Steady, Go!" }),
    ).toBeVisible();
    const saved = await player.evaluate(
      () =>
        JSON.parse(sessionStorage.getItem("session-store-play")).playSession,
    );
    expect(saved.sessionId).toBe(next.session_id);
    // An old host URL resolves the current round instead of reopening A.
    await host.goto(`/admin/arrange/${sessionA}`);
    await expect(host).toHaveURL(
      new RegExp(`/admin/arrange/${next.session_id}`),
    );
    await expect(
      host.getByRole("button", { name: "Start Quiz", exact: true }),
    ).toBeVisible();
    await host.getByRole("button", { name: "Start Quiz", exact: true }).click();
    await expect(
      player.getByText(`Question B ${suffix}`, { exact: true }),
    ).toBeVisible();
    await player.getByText("Correct", { exact: true }).click();
    await expect
      .poll(() =>
        sql(
          `SELECT count(*) FROM user_quiz_responses WHERE user_played_quiz_id='${saved.userPlayedQuiz}' AND is_attend=true`,
        ),
      )
      .toBe("1");
    // All players answered, so the quiz advances automatically to scores.
    await host.getByRole("button", { name: "Finish", exact: true }).click();
    await host.getByRole("link", { name: "Choose next quiz" }).click();
    await expect(host.getByLabel("Choose the next quiz")).toBeVisible();
    await data(await hostCtx.request.post(`${api}/rooms/${room.id}/close`));
    expect(hostErrors).toEqual([]);
  } finally {
    await Promise.all([hostCtx.close(), playerCtx.close()]);
  }
});

test("a new guest host can enter their name before any room heartbeat starts", async ({
  browser,
}) => {
  const context = await browser.newContext();
  try {
    const page = await context.newPage();
    const heartbeat = [];
    page.on("request", (request) => {
      if (request.url().includes("/rooms/lookup"))
        heartbeat.push(request.url());
    });
    await page.goto(`/admin/arrange/new?public=1&quiz_id=${quizA}`);
    await expect(page.getByRole("dialog")).toBeVisible();
    expect(heartbeat).toHaveLength(0);
    await expect(page).toHaveURL(/arrange\/new/);
  } finally {
    await context.close();
  }
});
