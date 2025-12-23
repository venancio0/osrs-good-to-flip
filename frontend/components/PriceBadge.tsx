interface PriceBadgeProps {
  price: number;
  className?: string;
}

export default function PriceBadge({ price, className = "" }: PriceBadgeProps) {
  const formatPrice = (price: number): string => {
    if (price >= 1000000) {
      return `${(price / 1000000).toFixed(2)}M GP`;
    } else if (price >= 1000) {
      return `${(price / 1000).toFixed(1)}k GP`;
    }
    return `${price} GP`;
  };

  const getColorClass = () => {
    if (price >= 1000000) {
      return "bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200";
    } else if (price >= 100000) {
      return "bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200";
    } else if (price >= 10000) {
      return "bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200";
    }
    return "bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200";
  };

  return (
    <span
      className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold ${getColorClass()} ${className}`}
    >
      {formatPrice(price)}
    </span>
  );
}

