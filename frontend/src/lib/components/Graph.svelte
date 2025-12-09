<script lang="ts">
  import { onMount } from "svelte";
  import type { Stock } from "../api/fetchStocks";
  import { getStocks } from "../api/fetchStocks";
  import type { StockDatum } from "../api/fetchStockData";
  import { getStockData} from "../api/fetchStockData";


  let canvas: HTMLCanvasElement;
  let data: StockDatum[];
  let stocks: Stock[];

  onMount(() => {
    getInformation()
    const ctx = canvas.getContext("2d")!;
    const width = canvas.width;
    const height = canvas.height;
    const oneTick = canvas.width / 10;
    ctx.clearRect(0, 0, width, height);
    drawGrid(ctx, oneTick);
  })

  async function getInformation() {
    data = await getStockData()
    stocks = await getStocks()
  }

  function drawGrid(ctx: CanvasRenderingContext2D, oneTick: number) {
    ctx.beginPath();
    ctx.strokeStyle = "white";
    ctx.lineWidth = 0.5;
    for (let i = 0; i < canvas.width + 1; i += oneTick) {
      ctx.moveTo(i, 0);
      ctx.lineTo(i, canvas.height);
      ctx.stroke();
    }
    for (let i = 0; i < canvas.height + 1; i += canvas.height / 10) {
      ctx.moveTo(0, i);
      ctx.lineTo(canvas.width, i);
      ctx.stroke();
    }
  }

  function prepareStockImages(ctx: CanvasRenderingContext2D) {
    for (let i = 0; i < stocks.length; i++) {
      item = stocks[i]
    }
  }
</script>

<canvas bind:this={canvas} width="1220" height="660"></canvas>
