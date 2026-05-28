const OSRM_API_URL = "https://router.project-osrm.org/route/v1/driving";

export interface OSRMRouteResult {
  distance: number; // in meters
  duration: number; // in seconds
  path: [number, number][]; // [lat, lng] coordinates
}

export const fetchOSRMRoute = async (
  startLat: number,
  startLng: number,
  endLat: number,
  endLng: number
): Promise<OSRMRouteResult> => {
  const url = `${OSRM_API_URL}/${startLng},${startLat};${endLng},${endLat}?overview=full&geometries=geojson`;
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`OSRM API error: ${response.statusText}`);
  }
  const data = await response.json();
  if (data.routes && data.routes.length > 0) {
    const route = data.routes[0];
    const path = route.geometry.coordinates.map((coord: [number, number]) => [coord[1], coord[0]] as [number, number]);
    return {
      distance: route.distance,
      duration: route.duration,
      path,
    };
  }
  throw new Error("No route found in OSRM response");
};
