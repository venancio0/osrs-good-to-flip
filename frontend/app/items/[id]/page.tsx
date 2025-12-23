"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { getItemById, getPriceHistory, ItemPrice, PriceHistoryEntry } from "@/lib/api";
import PriceBadge from "@/components/PriceBadge";
import TrendIcon from "@/components/TrendIcon";
import PriceChart from "@/components/PriceChart";

export default function ItemDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [item, setItem] = useState<ItemPrice | null>(null);
  const [history, setHistory] = useState<PriceHistoryEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [historyLoading, setHistoryLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchItem = async () => {
      if (!params.id || typeof params.id !== "string") {
        setError("Invalid item ID");
        setLoading(false);
        return;
      }

      try {
        const data = await getItemById(params.id);
        setItem(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch item");
      } finally {
        setLoading(false);
      }
    };

    fetchItem();
  }, [params.id]);

  useEffect(() => {
    const fetchHistory = async () => {
      if (!params.id || typeof params.id !== "string") {
        setHistoryLoading(false);
        return;
      }

      try {
        const data = await getPriceHistory(params.id, 7);
        setHistory(data);
      } catch (err) {
        console.error("Failed to fetch price history:", err);
      } finally {
        setHistoryLoading(false);
      }
    };

    if (item) {
      fetchHistory();
    }
  }, [params.id, item]);

  if (loading) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-600 dark:text-gray-400">Loading item details...</p>
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
        <p className="text-red-800 dark:text-red-300">Error: {error || "Item not found"}</p>
        <button
          onClick={() => router.push("/")}
          className="mt-4 px-4 py-2 bg-blue-600 dark:bg-blue-700 text-white rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors"
        >
          Back to List
        </button>
      </div>
    );
  }

  return (
    <div>
      <button
        onClick={() => router.push("/")}
        className="mb-6 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      >
        ← Back to List
      </button>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg dark:shadow-gray-900/50 p-8 border border-gray-200 dark:border-gray-700">
        <div className="flex items-start justify-between mb-6">
          <h1 className="text-4xl font-bold text-gray-900 dark:text-white">{item.name}</h1>
          <TrendIcon trend={item.trend} className="text-4xl" />
        </div>

        <div className="space-y-4">
          <div className="border-b border-gray-200 dark:border-gray-700 pb-4">
            <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">
              Current Price
            </h3>
            <PriceBadge price={item.price} className="text-2xl" />
          </div>

          <div className="border-b border-gray-200 dark:border-gray-700 pb-4">
            <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">
              24 Hour Average
            </h3>
            <PriceBadge price={item.avg_24h} className="text-xl" />
          </div>

          <div className="border-b border-gray-200 dark:border-gray-700 pb-4">
            <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">
              7 Day Average
            </h3>
            <PriceBadge price={item.avg_7d} className="text-xl" />
          </div>

          <div>
            <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">Trend</h3>
            <div className="flex items-center gap-2">
              <TrendIcon trend={item.trend} className="text-2xl" />
              <span className="text-lg font-semibold text-gray-700 dark:text-gray-300">
                {item.trend}
              </span>
            </div>
          </div>

          <div className="pt-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Last updated: {new Date(item.updated_at).toLocaleString()}
            </p>
          </div>
        </div>

        {/* Price History Chart */}
        <div className="mt-8 pt-8 border-t border-gray-200 dark:border-gray-700">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
            Price History (Last 7 Days)
          </h2>
          {historyLoading ? (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
              Loading price history...
            </div>
          ) : (
            <PriceChart data={history} />
          )}
        </div>
      </div>
    </div>
  );
}

