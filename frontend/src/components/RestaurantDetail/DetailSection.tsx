import React from "react";

interface DetailSectionProps {
  title: string;
  value: React.ReactNode;
}

export const DetailSection: React.FC<DetailSectionProps> = ({
  title,
  value,
}) => {
  return (
    <div>
      <h4 className="font-semibold text-gray-700 mb-1">{title}</h4>
      <p className="text-gray-600">{value}</p>
    </div>
  );
};
