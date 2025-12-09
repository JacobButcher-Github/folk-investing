<script lang="ts">
  import { onMount } from "svelte";
  import type { Stock } from "../api/fetchStocks";
  import { getStocks } from "../api/fetchStocks";
  import type { StockDatum, idToData } from "../api/fetchStockData";
  import { getStockData} from "../api/fetchStockData";
  import type { Settings } from "../api/fetchSettings";
  import { getSiteSettings } from "../api/fetchSettings";

  let settings: Settings;
  let canvas: HTMLCanvasElement;
  let data: idToData;
  let stocks: Stock[];
  let oneTick: number;

  onMount(() => {
    getInformation()
    const ctx = canvas.getContext("2d")!;
    const width = canvas.width;
    const height = canvas.height;
    oneTick = canvas.width / 10; //settings.number_of_events;
    ctx.clearRect(0, 0, width, height);
    drawGrid(ctx);
  })

  async function getInformation() {
    data = await getStockData()
    stocks = await getStocks()
    settings = await getSiteSettings()
  }

  function drawGrid(ctx: CanvasRenderingContext2D) {
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
      let stock = stocks[i];
      let imagePath = stock.image_path;
      let price = 0
      if (data[stock.id].at(-1) != undefined) {
        price = data[stock.id].at(-1)!.value_dollars * 100 + data[stock.id].at(-1)!.value_cents;
      }
      if (price === 0){
        continue;
      }
      drawStockImage(ctx, stock.name, imagePath, price)
    }
  }

  function drawStockImage(ctx: CanvasRenderingContext2D, name: string, imagePath: string, price: number) {
    return new Promise((resolve) => {
      const img = new Image();
      img.id = name;
      img.src = imagePath;
      img.onload = () => {
        ctx.drawImage(
          img,
          canvas.height - oneTick - 20,
          canvas.height - ((price / 100) * (canvas.height / 120) + 20)
        );
        resolve("resolved");
      };
    });
  }

</script>

<canvas bind:this={canvas} width="1220" height="660"></canvas>
