// Types matching backend domain
export type TrendType = "UP" | "DOWN" | "FLAT";

export interface ItemPrice {
  item_id: number;
  name: string;
  price: number;
  high: number;
  low: number;
  volume: number;
  avg_24h: number;
  avg_7d: number;
  trend: TrendType;
  updated_at: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchAPI<T>(endpoint: string): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`);

  if (!response.ok) {
    throw new Error(`API error: ${response.statusText}`);
  }

  return response.json();
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export async function getItems(
  query?: string,
  page: number = 1,
  limit: number = 20
): Promise<PaginatedResponse<ItemPrice>> {
  const params = new URLSearchParams();
  if (query) {
    params.append("q", query);
  }
  params.append("page", page.toString());
  params.append("limit", limit.toString());
  
  const endpoint = `/items?${params.toString()}`;
  return fetchAPI<PaginatedResponse<ItemPrice>>(endpoint);
}

export async function getItemById(id: string): Promise<ItemPrice> {
  return fetchAPI<ItemPrice>(`/items/${id}`);
}

export interface PriceHistoryEntry {
  date: string;
  price: number;
}

export async function getPriceHistory(
  id: string,
  days?: number
): Promise<PriceHistoryEntry[]> {
  const endpoint =
    days && days > 0
      ? `/items/${id}/history?days=${days}`
      : `/items/${id}/history`;
  return fetchAPI<PriceHistoryEntry[]>(endpoint);
}

