export interface Stock {
  id: number;
  name: string;
  image_path: string;
}

interface CachedStock {
  timestamp: number;
  stocks: Stock[];
}

const MAX_AGE = 120 * 1000; //2 minute lockout

export async function getStocks(): Promise<Stock[]> {
  const key = `stocks`;
  const cachedRaw = localStorage.getItem(key);

  if (cachedRaw) {
    const cached: CachedStock = JSON.parse(cachedRaw);
    if (Date.now() - cached.timestamp < MAX_AGE) {
      return cached.stocks;
    }
  }

  const res: Response = await fetch("stocks/list_stocks");
  if (!res.ok) throw new Error("Failed to fetch stocks");
  const data: Stock[] = await res.json();

  const cacheEntry: CachedStock = {
    timestamp: Date.now(),
    stocks: data,
  };

  localStorage.setItem(key, JSON.stringify(cacheEntry));
  return data;
}
