export interface Settings {
  number_of_events: number;
  value_symbol: string;
  event_label: string;
  title: string;
  give_each_day: number;
  lockout: number;
  lockout_time: Date;
}

interface cachedSettings {
  timestamp: number;
  settings: Settings;
}
const MAX_AGE = 120 * 1000;

export async function getSiteSettings(): Promise<Settings> {
  const key = `settings`;
  const cachedRaw = localStorage.getItem(key);

  if (cachedRaw) {
    const cached: cachedSettings = JSON.parse(cachedRaw);
    if (Date.now() - cached.timestamp < MAX_AGE) {
      return cached.settings;
    }
  }

  const res: Response = await fetch("settings/");
  if (!res.ok) throw new Error("Failed to settings");
  const data: Settings = await res.json();

  const cacheEntry: cachedSettings = {
    timestamp: Date.now(),
    settings: data,
  };
  localStorage.setItem(key, JSON.stringify(cacheEntry));
  return data;
}
