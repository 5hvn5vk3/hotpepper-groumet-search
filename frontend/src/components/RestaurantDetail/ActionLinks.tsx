import React from "react";

interface ActionLinksProps {
  googleMapsUrl: string;
  hotpepperUrl: string;
}

export const ActionLinks: React.FC<ActionLinksProps> = ({
  googleMapsUrl,
  hotpepperUrl,
}) => {
  return (
    <div className="flex flex-wrap gap-4 mt-4">
      <div className="pt-4">
        <a
          href={googleMapsUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block bg-green-600 text-white px-6 py-3 rounded-md hover:bg-yellow-700 transition-colors"
        >
          Googleマップで場所を見る
        </a>
      </div>

      <div className="pt-4">
        <a
          href={hotpepperUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block bg-red-600 text-white px-6 py-3 rounded-md hover:bg-red-700 transition-colors"
        >
          ホットペッパーで詳細を見る
        </a>
      </div>
    </div>
  );
};
