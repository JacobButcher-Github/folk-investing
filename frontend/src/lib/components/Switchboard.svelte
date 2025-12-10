<script lang="ts">
  import type { Stock } from "../api/fetchStocks";
  import type { idToData } from "../api/fetchStockData";
  import type { Settings } from "../api/fetchSettings";

  let { stocks, data, settings }: { stocks: Stock[]; data: idToData; settings: Settings} = $props();
  
  let quantity = 1;

  function purchase(stock: Stock) {
    console.log("Buying", quantity, "shares of", stock.name);
  }
</script>

<div id="board">
  {#each stocks as stock (stock.id)}
  {/each}
</div>
