<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type { Stock } from "../api/fetchStocks";
  import type { idToData } from "../api/fetchStockData";
  import type { Settings } from "../api/fetchSettings";
  import {breakpoint, initBreakpointWatcher } from "../stores/screenBreakpoints";

  let { stocks, data, settings }: { stocks: Stock[]; data: idToData; settings: Settings} = $props();

  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D | null = null;

  let parent: HTMLDivElement;
  let width = 1200;
  let height = 600;


  let unsubscribeBreakpoint: () => void;

  onMount(() => {

    const cleanup = initBreakpointWatcher();
    unsubscribeBreakpoint = cleanup;

    const unsub = breakpoint.subscribe(() => {
      resizeCanvas();
      drawCanvas();
    });

    onDestroy(() => {
      unsub();
      cleanup();
    });


    if (stocks && data && settings){
      resizeCanvas()
      drawCanvas()
    }
    ctx = canvas.getContext("2d");
  });



  function resizeCanvas() {
    canvas.width = width;
    canvas.height = height;
  }

  function drawCanvas() {
    if (!ctx) return;

    ctx.clearRect(0, 0, width, height);

    drawGrid(ctx);
    drawStocks(ctx);
  }

  function drawGrid(ctx: CanvasRenderingContext2D) {
    const cols = settings?.number_of_events ?? 10;
    const colWidth = width / cols;

    ctx.strokeStyle = "white";
    ctx.lineWidth = 0.5;

    // vertical grid
    for (let x = 0; x <= width; x += colWidth) {
      ctx.beginPath();
      ctx.moveTo(x, 0);
      ctx.lineTo(x, height);
      ctx.stroke();
    }

    // horizontal grid
    const rows = 10;
    for (let y = 0; y <= height; y += height / rows) {
      ctx.beginPath();
      ctx.moveTo(0, y);
      ctx.lineTo(width, y);
      ctx.stroke();
    }
  }

  async function drawStocks(ctx: CanvasRenderingContext2D) {
    const cols = settings?.number_of_events ?? 10;
    const priceScale = height / 120; 

    for (const stock of stocks) {
      const entries = data[stock.id];
      if (!entries || entries.length === 0) continue;


      for (let i = 0; i < entries.length; i++) {
        let entry = entries[i];
        const price = entry.value_dollars * 100 + entry.value_cents; 
        let nextEntry = entries[i + 1];
        const nextPrice = nextEntry.value_dollars * 100 + nextEntry.value_cents; 
        ctx.beginPath();
        ctx.lineWidth = 5;
        if (nextPrice < price) {
          ctx.strokeStyle = "red";
        } else if (nextPrice > price) {
          ctx.strokeStyle = "green";
        } else if (nextPrice == price) {
          ctx.strokeStyle = "gray";
        }
        ctx.moveTo(
          i * cols,
          canvas.height - price / 100 * priceScale - 20
        );
        ctx.lineTo(
          (i + 1) * cols,
          canvas.height - nextPrice / 100 * priceScale - 20 
        );
        ctx.stroke();
      }

      const last = entries.at(-1)!;
      const lastPrice = last.value_dollars * 100 + last.value_cents;
      if (lastPrice <= 0) continue;

      await drawStockImage(ctx, stock.name, stock.image_path, lastPrice);
    }
  }

  //Draw image at the last stock location.
  function drawStockImage(ctx: CanvasRenderingContext2D, name: string, src: string, price: number) {
    return new Promise<void>((resolve) => {
      const img = new Image();
      img.id = name;
      img.src = src;

      img.onload = () => {
        const priceScale = height / 120; 
        const y = height - price / 100 * priceScale - 20;

        ctx.drawImage(img, 20, y, 32, 32);
        resolve();
      };
    });
  }
</script>

<div class="graph-container" bind:this={parent}>
  <canvas bind:this={canvas}></canvas>
</div>
