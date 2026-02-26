import React from "react";

interface FacilityItem {
  label: string;
  value?: string;
}

interface FacilitiesSectionProps {
  items: FacilityItem[];
  renderDeferredValue: (value?: string) => string;
}

export const FacilitiesSection: React.FC<FacilitiesSectionProps> = ({
  items,
  renderDeferredValue,
}) => {
  return (
    <div>
      <h4 className="font-semibold text-gray-700 mb-1">設備・条件</h4>
      {items.map((item) => (
        <p key={item.label} className="text-gray-600">
          {item.label}: {renderDeferredValue(item.value)}
        </p>
      ))}
    </div>
  );
};
