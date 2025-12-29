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
  const url = `${API_URL}${endpoint}`;
  
  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Add credentials for CORS if needed
      credentials: 'omit',
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(
        `API error (${response.status}): ${response.statusText}. ${errorText}`
      );
    }

    return response.json();
  } catch (error) {
    if (error instanceof TypeError && error.message.includes('fetch')) {
      throw new Error(
        `Network error: Unable to reach API at ${url}. Check CORS configuration and API URL.`
      );
    }
    throw error;
  }
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

