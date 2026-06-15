<script setup>
/**
 * WinnerConfetti
 * --------------
 * A self-contained, dependency-free award-ceremony confetti effect.
 *
 * Two cannons sit just outside the left and right edges of the viewport and
 * fire occasional, short bursts of paper confetti inward. Pieces launch with
 * modest velocity, spread, then float gently down under gravity — the look of
 * a trophy presentation or concert finale rather than a high-pressure spray.
 *
 * The effect is intentionally minimal: low particle density, larger and
 * clearly distinguishable pieces, and a calm centre so the podium stays
 * readable. One full-screen <canvas> driven by a single requestAnimationFrame
 * loop; physics are time-scaled to 60fps and the pool is capped for smooth
 * performance on mobile.
 */

const props = defineProps({
  // Lets the parent pause/resume without unmounting. The component also stops
  // itself automatically when the tab is hidden.
  active: {
    type: Boolean,
    default: true,
  },
});

const canvas = ref(null);

// Theme colours (jv-*) mixed with a few classic celebratory tones.
const COLORS = [
  "#fde047", // jv-yellow
  "#ff6b6b", // jv-coral
  "#d1fae5", // jv-mint
  "#e4dff4", // jv-lavender
  "#d4b524", // gold
  "#ff4fa3", // pink
  "#a855f7", // purple
  "#ffffff", // white
];

let ctx = null;
let rafId = 0;
let running = false;
let lastTime = 0;
// Countdown (ms) until the next cannon pop; reset to a random delay each time.
let burstTimer = 0;

// Logical (CSS pixel) viewport size; kept in sync on resize.
let viewW = 0;
let viewH = 0;
let dpr = 1;

const particles = [];

// Tuning ---------------------------------------------------------------------
// Deliberately low cap — this is a minimal, premium effect, not a fountain.
const MAX_PARTICLES = 90;
const isSmallScreen = () => viewW < 640;

const rand = (min, max) => min + Math.random() * (max - min);
const pick = (arr) => arr[(Math.random() * arr.length) | 0];

/**
 * Fire one short burst from a single cannon.
 * @param {"left"|"right"} side
 */
function fireCannon(side) {
  if (particles.length >= MAX_PARTICLES) return;

  const small = isSmallScreen();
  // Few pieces per pop — quality over quantity.
  const count = small ? Math.round(rand(4, 7)) : Math.round(rand(6, 10));

  // Cannons just outside the edge, around mid-height, aimed up-and-inward so
  // confetti arcs over the podium and rains down onto the winners in the
  // centre of the stage.
  const originX = side === "left" ? -12 : viewW + 12;
  const originY = viewH * rand(0.5, 0.72);
  const inward = side === "left" ? 1 : -1;

  for (let i = 0; i < count; i++) {
    if (particles.length >= MAX_PARTICLES) break;

    // Flatter arc + higher speed so pieces actually travel toward the centre
    // before gravity takes over and drops them on the podium.
    const elevation = rand(0.45, 0.85); // ~26°–49° above horizontal
    const speed = rand(small ? 12 : 15, small ? 18 : 23);

    const roll = Math.random();
    const streamer = roll < 0.18; // occasional long ribbon
    const size = rand(10, 17);

    particles.push({
      x: originX,
      y: originY,
      vx: Math.cos(elevation) * speed * inward + rand(-0.8, 0.8),
      vy: -Math.sin(elevation) * speed + rand(-0.8, 0.8),
      streamer,
      // Paper: roughly square-ish. Streamer: long and thin.
      w: streamer ? rand(4, 6) : size,
      h: streamer ? rand(26, 40) : size * rand(0.55, 0.9),
      color: pick(COLORS),
      rotation: rand(0, Math.PI * 2),
      spin: rand(-0.14, 0.14),
      // Flutter drives both the edge-on flip and a gentle horizontal sway, so
      // pieces tumble and drift like real paper.
      flutter: rand(0, Math.PI * 2),
      flutterSpeed: rand(0.04, 0.09),
      sway: rand(0.4, 1.1),
      // Lighter air resistance so horizontal momentum carries pieces inward
      // toward the winners before they flutter down.
      drag: rand(0.982, 0.992),
      life: 0,
      maxLife: rand(240, 360),
    });
  }
}

function step(dt) {
  // Gentle gravity with a low terminal speed so confetti drifts down slowly.
  const gravity = 0.17 * dt;
  const terminal = 4;

  for (let i = particles.length - 1; i >= 0; i--) {
    const p = particles[i];

    p.vx *= Math.pow(p.drag, dt);
    p.vy = p.vy * Math.pow(p.drag, dt) + gravity;
    if (p.vy > terminal) p.vy = terminal;

    p.flutter += p.flutterSpeed * dt;
    p.x += (p.vx + Math.sin(p.flutter) * p.sway) * dt;
    p.y += p.vy * dt;
    p.rotation += p.spin * dt;
    p.life += dt;

    const gone =
      p.y > viewH + 50 || p.x < -80 || p.x > viewW + 80 || p.life > p.maxLife;
    if (gone) particles.splice(i, 1);
  }
}

function draw() {
  ctx.clearRect(0, 0, viewW, viewH);

  for (let i = 0; i < particles.length; i++) {
    const p = particles[i];

    // Soft fade-in at birth and fade-out at the end of life.
    const fadeOut =
      p.life > p.maxLife - 36 ? Math.max(0, (p.maxLife - p.life) / 36) : 1;
    const fadeIn = p.life < 8 ? p.life / 8 : 1;

    // cos(flutter) turns the piece edge-on; negative = we see its back, which
    // we dim slightly to read as a 3D flip.
    const flip = Math.cos(p.flutter);

    ctx.save();
    ctx.globalAlpha = Math.min(fadeIn, fadeOut) * (flip < 0 ? 0.72 : 1);
    ctx.translate(p.x, p.y);
    ctx.rotate(p.rotation);
    ctx.scale(Math.max(0.15, Math.abs(flip)), 1);
    ctx.fillStyle = p.color;
    ctx.fillRect(-p.w / 2, -p.h / 2, p.w, p.h);
    ctx.restore();
  }
}

function loop(now) {
  if (!running) return;

  // Normalise to 60fps units and clamp so a backgrounded tab returning to
  // focus doesn't teleport every particle across the screen.
  const deltaMs = lastTime ? now - lastTime : 16.67;
  lastTime = now;
  const dt = Math.min(deltaMs / 16.67, 2.5);

  // Both cannons fire together on every pop, like a paired stage launch. The
  // random gap between pops keeps the rhythm natural rather than mechanical.
  burstTimer -= deltaMs;
  if (burstTimer <= 0) {
    fireCannon("left");
    fireCannon("right");
    burstTimer = rand(1700, 3000);
  }

  step(dt);
  draw();

  rafId = requestAnimationFrame(loop);
}

function resize() {
  if (!canvas.value) return;
  // Cap DPR at 2 — beyond that the extra pixels cost fps for no visible gain.
  dpr = Math.min(window.devicePixelRatio || 1, 2);
  viewW = window.innerWidth;
  viewH = window.innerHeight;

  canvas.value.width = Math.floor(viewW * dpr);
  canvas.value.height = Math.floor(viewH * dpr);
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
}

function start() {
  if (running || !ctx) return;
  running = true;
  lastTime = 0;
  // Fire the first pop almost immediately so the screen isn't empty on entry.
  burstTimer = 150;
  rafId = requestAnimationFrame(loop);
}

function stop() {
  running = false;
  if (rafId) cancelAnimationFrame(rafId);
  rafId = 0;
}

function onVisibility() {
  if (document.hidden) {
    stop();
  } else if (props.active) {
    start();
  }
}

watch(
  () => props.active,
  (on) => {
    if (on && !document.hidden) start();
    else stop();
  }
);

onMounted(() => {
  if (!canvas.value) return;
  ctx = canvas.value.getContext("2d");
  resize();
  window.addEventListener("resize", resize, { passive: true });
  document.addEventListener("visibilitychange", onVisibility);
  if (props.active && !document.hidden) start();
});

onBeforeUnmount(() => {
  stop();
  window.removeEventListener("resize", resize);
  document.removeEventListener("visibilitychange", onVisibility);
  particles.length = 0;
  ctx = null;
});
</script>

<template>
  <!-- Sits above the page background but below the winner content (z-10). -->
  <canvas
    ref="canvas"
    class="pointer-events-none fixed inset-0 z-[5] h-full w-full"
    aria-hidden="true"
  ></canvas>
</template>
