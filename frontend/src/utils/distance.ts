/**
 * ハバーシン公式を使って2点間の距離（メートル）を計算する
 * @param lat1 地点1の緯度
 * @param lng1 地点1の経度
 * @param lat2 地点2の緯度
 * @param lng2 地点2の経度
 * @returns 距離（メートル）
 */
export const calculateDistance = (
  lat1: number,
  lng1: number,
  lat2: number,
  lng2: number,
): number => {
  const R = 6371000; // 地球の半径（メートル）
  const toRad = (deg: number) => (deg * Math.PI) / 180;

  const dLat = toRad(lat2 - lat1);
  const dLng = toRad(lng2 - lng1);

  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(lat1)) *
      Math.cos(toRad(lat2)) *
      Math.sin(dLng / 2) *
      Math.sin(dLng / 2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
};

/**
 * 距離（メートル）を表示用文字列にフォーマットする
 * 1000m未満は「XXXm」、1000m以上は「X.Xkm」
 * @param meters 距離（メートル）
 * @returns フォーマット済み文字列
 */
export const formatDistance = (meters: number): string => {
  if (meters < 1000) {
    return `ここから${Math.round(meters)}m`;
  }
  return `${(meters / 1000).toFixed(1)}km`;
};
