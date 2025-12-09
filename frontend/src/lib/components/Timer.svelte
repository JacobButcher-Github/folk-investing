<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type { Settings } from "../api/fetchSettings";

  let { settings }: { settings: Settings} = $props();

  let remaining = $state();
  let intervalId: number;

  function updateTimer() {
    const now = Date.now();
    const lockoutTime = new Date(settings.lockout_time).getTime();

    const diff = lockoutTime - now;

    if (diff <= 0) {
      remaining = "Lockout active";
      return;
    }

    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff / (1000 * 60)) % 60);
    const seconds = Math.floor((diff / 1000) % 60);

    remaining = `${hours}h ${minutes}m ${seconds}s`;
  }

  onMount(() => {
    updateTimer();                 
    intervalId = setInterval(updateTimer, 1000); 
  });

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

<div class="timer">
  Next lockout in: <strong>{remaining}</strong>
</div>

<style>
  .timer {
    font-size: 1.3rem;
    color: white;
  }
</style>
