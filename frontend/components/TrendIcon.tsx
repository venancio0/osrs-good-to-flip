import { TrendType } from "@/lib/api";

interface TrendIconProps {
  trend: TrendType;
  className?: string;
}

export default function TrendIcon({ trend, className = "" }: TrendIconProps) {
  const getIcon = () => {
    switch (trend) {
      case "UP":
        return "↑";
      case "DOWN":
        return "↓";
      case "FLAT":
        return "→";
      default:
        return "→";
    }
  };

  const getColor = () => {
    switch (trend) {
      case "UP":
        return "text-green-600";
      case "DOWN":
        return "text-red-600";
      case "FLAT":
        return "text-gray-500";
      default:
        return "text-gray-500";
    }
  };

  return (
    <span className={`font-bold ${getColor()} ${className}`}>
      {getIcon()}
    </span>
  );
}

