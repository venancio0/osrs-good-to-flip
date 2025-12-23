import Link from "next/link";
import { ItemPrice } from "@/lib/api";
import PriceBadge from "./PriceBadge";
import TrendIcon from "./TrendIcon";

interface ItemCardProps {
  item: ItemPrice;
}

function formatPrice(price: number): string {
  if (price >= 1000000) {
    return `${(price / 1000000).toFixed(2)}M GP`;
  } else if (price >= 1000) {
    return `${(price / 1000).toFixed(1)}k GP`;
  }
  return `${price} GP`;
}

export default function ItemCard({ item }: ItemCardProps) {
  return (
    <Link href={`/items/${item.item_id}`}>
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md dark:shadow-gray-900/50 p-6 hover:shadow-lg dark:hover:shadow-gray-900 transition-shadow cursor-pointer border border-gray-200 dark:border-gray-700">
        <div className="flex items-start justify-between mb-3">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">{item.name}</h3>
          <TrendIcon trend={item.trend} className="text-2xl" />
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <PriceBadge price={item.price} />
          <span className="text-sm text-gray-500 dark:text-gray-400">
            Avg 24h: {formatPrice(item.avg_24h)}
          </span>
        </div>
      </div>
    </Link>
  );
}

