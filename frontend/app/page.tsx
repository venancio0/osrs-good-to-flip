"use client";

import { useState, useEffect, useCallback } from "react";
import { getItems, ItemPrice, PaginatedResponse } from "@/lib/api";
import ItemsTable from "@/components/ItemsTable";
import SearchBar from "@/components/SearchBar";
import Pagination from "@/components/Pagination";

export default function Home() {
  const [items, setItems] = useState<ItemPrice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");
  const limit = 20;

  const fetchItems = useCallback(
    async (query: string, page: number) => {
      setLoading(true);
      setError(null);
      try {
        const data: PaginatedResponse<ItemPrice> = await getItems(
          query,
          page,
          limit
        );
        setItems(data.data);
        setTotalPages(data.total_pages);
        setTotal(data.total);
        setCurrentPage(data.page);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch items");
      } finally {
        setLoading(false);
      }
    },
    [limit]
  );

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    setCurrentPage(1);
    fetchItems(query, 1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    fetchItems(searchQuery, page);
  };

  // Initial load
  useEffect(() => {
    fetchItems("", 1);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Auto-refresh every 5 minutes
  useEffect(() => {
    const interval = setInterval(() => {
      fetchItems(searchQuery, currentPage);
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(interval);
  }, [searchQuery, currentPage, fetchItems]);

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
        <>
          <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
            Showing {items.length} of {total} items
          </div>
          <ItemsTable items={items} />
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </>
      )}
    </div>
  );
}

