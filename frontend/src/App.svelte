<!-- App.svelte -->
<script lang="ts">
  import "./stylesheets/main.scss"
  import Graph from "./lib/components/Graph.svelte"
  import { getStockData, type idToData } from "./lib/api/fetchStockData";
  import { getStocks, type Stock } from "./lib/api/fetchStocks";
  import { getSiteSettings, type Settings } from "./lib/api/fetchSettings";
    import { onMount } from "svelte";


  let { stocks, data, settings }: { stocks: Stock[]; data: idToData; settings?: Settings} = $props();

  onMount(async () => {
    loadEverything()
  })

  async function loadEverything() {
    [stocks, data, settings] = await Promise.all([
      getStocks(),
      getStockData(),
      getSiteSettings()
    ]);
  }
</script>

<main>
  <div class="main">
    <div class="top-main">
      <Graph stocks={stocks} data={data} settings={settings}/>
      <!-- <Info /> -->
    </div>
    <div class="bottom-main">
      <!-- <Switchboard /> -->
      <!-- <Controller /> -->
    </div>
  </div>
</main>
