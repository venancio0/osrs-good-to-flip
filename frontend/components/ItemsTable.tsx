"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { ItemPrice, PriceHistoryEntry, getPriceHistory } from "@/lib/api";
import TrendIcon from "./TrendIcon";
import PriceBadge from "./PriceBadge";
import PriceChart from "./PriceChart";

interface ItemsTableProps {
  items: ItemPrice[];
}

// Helper functions for calculations
function calculateMargin(buyPrice: number, sellPrice: number): number {
  if (buyPrice === 0) return 0;
  return ((sellPrice - buyPrice) / buyPrice) * 100;
}

function calculateGETax(sellPrice: number): number {
  return Math.floor(sellPrice * 0.01);
}

function calculateExpectedProfit(buyPrice: number, sellPrice: number): number {
  const tax = calculateGETax(sellPrice);
  return sellPrice - buyPrice - tax;
}

function formatNumber(num: number): string {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(2)}M`;
  } else if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}k`;
  }
  return num.toString();
}

export default function ItemsTable({ items }: ItemsTableProps) {
  const [pinnedItems, setPinnedItems] = useState<Set<number>>(new Set());
  const [expandedCharts, setExpandedCharts] = useState<Set<number>>(new Set());
  const [chartData, setChartData] = useState<Record<number, PriceHistoryEntry[]>>({});
  const [loadingCharts, setLoadingCharts] = useState<Set<number>>(new Set());

  // Load pinned items from localStorage
  useEffect(() => {
    const saved = localStorage.getItem("pinnedItems");
    if (saved) {
      try {
        const pinned = JSON.parse(saved) as number[];
        setPinnedItems(new Set(pinned));
      } catch (e) {
        console.error("Failed to load pinned items", e);
      }
    }
  }, []);

  // Save pinned items to localStorage
  useEffect(() => {
    if (pinnedItems.size > 0) {
      localStorage.setItem("pinnedItems", JSON.stringify(Array.from(pinnedItems)));
    } else {
      localStorage.removeItem("pinnedItems");
    }
  }, [pinnedItems]);

  const togglePin = (itemId: number) => {
    setPinnedItems((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(itemId)) {
        newSet.delete(itemId);
      } else {
        newSet.add(itemId);
      }
      return newSet;
    });
  };

  const toggleChart = async (itemId: number) => {
    if (expandedCharts.has(itemId)) {
      setExpandedCharts((prev) => {
        const newSet = new Set(prev);
        newSet.delete(itemId);
        return newSet;
      });
    } else {
      setExpandedCharts((prev) => new Set(prev).add(itemId));
      
      // Load chart data if not already loaded
      if (!chartData[itemId]) {
        setLoadingCharts((prev) => new Set(prev).add(itemId));
        try {
          const data = await getPriceHistory(itemId.toString(), 7);
          setChartData((prev) => ({ ...prev, [itemId]: data }));
        } catch (err) {
          console.error("Failed to load chart data", err);
        } finally {
          setLoadingCharts((prev) => {
            const newSet = new Set(prev);
            newSet.delete(itemId);
            return newSet;
          });
        }
      }
    }
  };

  // Sort items: pinned first, then by item_id
  const sortedItems = [...items].sort((a, b) => {
    const aPinned = pinnedItems.has(a.item_id);
    const bPinned = pinnedItems.has(b.item_id);
    if (aPinned && !bPinned) return -1;
    if (!aPinned && bPinned) return 1;
    return a.item_id - b.item_id;
  });

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
        <thead className="bg-gray-50 dark:bg-gray-900">
          <tr>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Pin
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Chart
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Item Name
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Buy (Low)
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Sell (High)
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Margin
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              GE Tax
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Profit
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Volume
            </th>
            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
              Trend
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
          {sortedItems.map((item) => {
            const margin = calculateMargin(item.low, item.high);
            const geTax = calculateGETax(item.high);
            const profit = calculateExpectedProfit(item.low, item.high);
            const isPinned = pinnedItems.has(item.item_id);
            const isChartExpanded = expandedCharts.has(item.item_id);

            return (
              <>
                <tr
                  key={item.item_id}
                  className={`hover:bg-gray-50 dark:hover:bg-gray-700 ${
                    isPinned ? "bg-yellow-50 dark:bg-yellow-900/20" : ""
                  }`}
                >
                  <td className="px-4 py-4 whitespace-nowrap">
                    <button
                      onClick={() => togglePin(item.item_id)}
                      className="text-yellow-500 hover:text-yellow-600 dark:text-yellow-400 dark:hover:text-yellow-300"
                      aria-label={isPinned ? "Unpin item" : "Pin item"}
                    >
                      {isPinned ? (
                        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                          <path d="M5 4a2 2 0 012-2h6a2 2 0 012 2v14l-5-2.5L5 18V4z" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 5a2 2 0 012-2h6a2 2 0 012 2v14l-5-2.5L5 18V5z" />
                        </svg>
                      )}
                    </button>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">
                    <button
                      onClick={() => toggleChart(item.item_id)}
                      className="text-blue-500 hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300"
                      aria-label="Toggle chart"
                    >
                      {isChartExpanded ? (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                        </svg>
                      ) : (
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      )}
                    </button>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">
                    <Link
                      href={`/items/${item.item_id}`}
                      className="text-sm font-medium text-gray-900 dark:text-white hover:text-blue-600 dark:hover:text-blue-400"
                    >
                      {item.name}
                    </Link>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                    {formatNumber(item.low)} GP
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                    {formatNumber(item.high)} GP
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`font-medium ${
                        margin > 0
                          ? "text-green-600 dark:text-green-400"
                          : "text-red-600 dark:text-red-400"
                      }`}
                    >
                      {margin.toFixed(2)}%
                    </span>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                    {formatNumber(geTax)} GP
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`font-medium ${
                        profit > 0
                          ? "text-green-600 dark:text-green-400"
                          : "text-red-600 dark:text-red-400"
                      }`}
                    >
                      {formatNumber(profit)} GP
                    </span>
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                    {formatNumber(item.volume)}
                  </td>
                  <td className="px-4 py-4 whitespace-nowrap">
                    <TrendIcon trend={item.trend} className="text-xl" />
                  </td>
                </tr>
                {isChartExpanded && (
                  <tr>
                    <td colSpan={10} className="px-4 py-4 bg-gray-50 dark:bg-gray-900">
                      {loadingCharts.has(item.item_id) ? (
                        <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                          Loading chart...
                        </div>
                      ) : (
                        <div className="max-w-4xl mx-auto">
                          <PriceChart data={chartData[item.item_id] || []} />
                        </div>
                      )}
                    </td>
                  </tr>
                )}
              </>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

