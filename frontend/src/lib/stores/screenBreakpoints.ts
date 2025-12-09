import { writable } from "svelte/store";

export type Breakpoint = "small" | "medium" | "large";

export const breakpoint = writable<Breakpoint>("large");

const bpDefs = {
  small: 320,
  medium: 768,
  large: 1024,
};

function getCurrent(): Breakpoint {
  const w = window.innerWidth;
  if (w < bpDefs.medium) return "small";
  if (w < bpDefs.large) return "medium";
  return "large";
}

export function initBreakpointWatcher() {
  let current = getCurrent();
  breakpoint.set(current);

  const handler = () => {
    const now = getCurrent();
    if (now !== current) {
      current = now;
      breakpoint.set(now);
    }
  };

  window.addEventListener("resize", handler);
  return () => window.removeEventListener("resize", handler);
}
