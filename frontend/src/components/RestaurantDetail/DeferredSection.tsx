import React from "react";
import { DetailSection } from "./DetailSection";

interface DeferredSectionProps {
  title: string;
  value?: string;
  renderDeferredValue: (value?: string) => string;
}

export const DeferredSection: React.FC<DeferredSectionProps> = ({
  title,
  value,
  renderDeferredValue,
}) => {
  return <DetailSection title={title} value={renderDeferredValue(value)} />;
};
