"use client";

import { useState, useEffect } from "react";
import { getItems, ItemPrice } from "@/lib/api";
import ItemsTable from "@/components/ItemsTable";
import SearchBar from "@/components/SearchBar";

export default function Home() {
  const [items, setItems] = useState<ItemPrice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (query: string) => {
    setLoading(true);
    setError(null);
    try {
      const data = await getItems(query);
      setItems(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch items");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    handleSearch("");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
          Grand Exchange Items
        </h2>
        <SearchBar onSearch={handleSearch} />
      </div>

      {loading && (
        <div className="text-center py-12">
          <p className="text-gray-600 dark:text-gray-400">Loading items...</p>
        </div>
      )}

      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-6">
          <p className="text-red-800 dark:text-red-300">Error: {error}</p>
        </div>
      )}

      {!loading && !error && items.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-600 dark:text-gray-400">No items found.</p>
        </div>
      )}

      {!loading && !error && items.length > 0 && (
        <ItemsTable items={items} />
      )}
    </div>
  );
}

