<script lang="ts">
  import Timer from "./Timer.svelte";
  import type { Settings } from "../api/fetchSettings";

  let { settings }: { settings: Settings} = $props();
</script>

<div class="info-container">
  <h1 class="disclaimer">WARNING THIS SITE DOES NOT DEAL WITH REAL MONEY THIS IS PURELY FOR FUN.</h1>
  <h1 class="title">Welcome to {settings.title}!</h1>
  <Timer settings={settings}/>
</div>
