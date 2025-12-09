export interface StockDatum {
  id: number;
  stock_id: number;
  event_label: string;
  value_dollars: number;
  value_cents: number;
}

interface CachedStockData {
  timestamp: number;
  stock_data: StockDatum[];
}

const MAX_AGE = 120 * 1000; //2 minute lockout

export async function getStockData(): Promise<StockDatum[]> {
  const key = `stockData`;
  const cachedRaw = localStorage.getItem(key);

  if (cachedRaw) {
    const cached: CachedStockData = JSON.parse(cachedRaw);
    if (Date.now() - cached.timestamp < MAX_AGE) {
      return cached.stock_data;
    }
  }

  const res: Response = await fetch("stocks/get_stocks_data");
  if (!res.ok) throw new Error("Failed to fetch stock data");
  const data: StockDatum[] = await res.json();

  const cacheEntry: CachedStockData = {
    timestamp: Date.now(),
    stock_data: data,
  };
  localStorage.setItem(key, JSON.stringify(cacheEntry));
  return data;
}
