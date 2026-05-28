import { fetchOSRMRoute } from "./osrm";

export interface Location {
  id: string;
  name: string;
  lat: number;
  lng: number;
  order: number;
}

export interface SegmentDetail {
  fromLat: number;
  fromLng: number;
  toLat: number;
  toLng: number;
  distance: number;
  duration: number;
  highway: string;
  cost: number;
}

export type OptimizationCriterion = "distance" | "time" | "economy" | "smart" | "toll" | "fuel";

export class RouteOptimizer {
  public distanceMatrix: Map<string, Map<string, number>>;
  public durationMatrix: Map<string, Map<string, number>>;
  private pathMatrix: Map<string, Map<string, [number, number][]>>;
  public fuelCostMatrix: Map<string, Map<string, number>>;
  public tollCostMatrix: Map<string, Map<string, number>>;
  public segmentsMatrix: Map<string, Map<string, SegmentDetail[]>>;
  private locations: Location[];
  public isUsingOSRM: boolean = false;
  
  public fuelConsumptionLitersPer100Km: number = 7;
  public fuelPricePerLiter: number = 1.8;
  public tollPricePerKm: number = 0.13;
  public timeValuePerHour: number = 20.0;
  public optimizationCriterion: OptimizationCriterion = "distance";

  constructor() {
    this.distanceMatrix = new Map();
    this.durationMatrix = new Map();
    this.pathMatrix = new Map();
    this.fuelCostMatrix = new Map();
    this.tollCostMatrix = new Map();
    this.segmentsMatrix = new Map();
    this.locations = [];
  }

  async buildGraph(locations: Location[]): Promise<void> {
    this.locations = [...locations];
    this.distanceMatrix.clear();
    this.durationMatrix.clear();
    this.pathMatrix.clear();
    this.fuelCostMatrix.clear();
    this.tollCostMatrix.clear();
    this.segmentsMatrix.clear();
    this.isUsingOSRM = false;

    for (let i = 0; i < locations.length - 1; i++) {
      const from = locations[i];
      const to = locations[i + 1];
      const { distance, duration, path, fuelCost, tollCost, segments } = await this.calculateRoute(from, to);

      let fromDistances = this.distanceMatrix.get(from.id);
      if (!fromDistances) {
        fromDistances = new Map<string, number>();
        this.distanceMatrix.set(from.id, fromDistances);
      }
      fromDistances.set(to.id, distance);

      let fromDurations = this.durationMatrix.get(from.id);
      if (!fromDurations) {
        fromDurations = new Map<string, number>();
        this.durationMatrix.set(from.id, fromDurations);
      }
      fromDurations.set(to.id, duration);

      let fromPaths = this.pathMatrix.get(from.id);
      if (!fromPaths) {
        fromPaths = new Map<string, [number, number][]>();
        this.pathMatrix.set(from.id, fromPaths);
      }
      fromPaths.set(to.id, path);

      let fromFuelCosts = this.fuelCostMatrix.get(from.id);
      if (!fromFuelCosts) {
        fromFuelCosts = new Map<string, number>();
        this.fuelCostMatrix.set(from.id, fromFuelCosts);
      }
      fromFuelCosts.set(to.id, fuelCost);

      let fromTollCosts = this.tollCostMatrix.get(from.id);
      if (!fromTollCosts) {
        fromTollCosts = new Map<string, number>();
        this.tollCostMatrix.set(from.id, fromTollCosts);
      }
      fromTollCosts.set(to.id, tollCost);

      let fromSegments = this.segmentsMatrix.get(from.id);
      if (!fromSegments) {
        fromSegments = new Map<string, SegmentDetail[]>();
        this.segmentsMatrix.set(from.id, fromSegments);
      }
      fromSegments.set(to.id, segments);
    }
  }

  private async calculateRoute(
    from: Location,
    to: Location
  ): Promise<{ distance: number; duration: number; path: [number, number][]; fuelCost: number; tollCost: number; segments: SegmentDetail[] }> {
    try {
      const backendUrl = import.meta.env?.VITE_BACKEND_URL || "http://localhost:8080";
      const response = await fetch(
        `${backendUrl}/api/route?startLat=${from.lat}&startLng=${from.lng}&endLat=${to.lat}&endLng=${to.lng}&criterion=${this.optimizationCriterion}&fuelPrice=${this.fuelPricePerLiter}&consumption=${this.fuelConsumptionLitersPer100Km}&tollPrice=${this.tollPricePerKm}&timeValue=${this.timeValuePerHour}`
      );
      
      if (!response.ok) {
        if (response.status === 400) {
          try {
            const errData = await response.json();
            if (errData.error === "out_of_bounds") {
              this.isUsingOSRM = true;
              const osrmRes = await fetchOSRMRoute(from.lat, from.lng, to.lat, to.lng);
              const fuelCost = (osrmRes.distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
              const speedKmh = osrmRes.duration > 0 ? (osrmRes.distance / osrmRes.duration) * 3.6 : 0;
              const tollCost = speedKmh >= 110 ? (osrmRes.distance / 1000) * this.tollPricePerKm : 0;
              return { ...osrmRes, fuelCost, tollCost, segments: [] };
            }
          } catch (_) {
            // Ignore error parsing, fallback to standard error
          }
        }
        throw new Error("Failed to fetch route from local server");
      }
      
      const data = await response.json();

      if (data.path && data.path.length > 0) {
        return {
          distance: data.distance,
          duration: data.duration,
          path: data.path.map((coord: [number, number]) => [coord[1], coord[0]]),
          fuelCost: data.fuelCost,
          tollCost: data.tollCost,
          segments: data.segments || [],
        };
      }
      throw new Error("No route found in response");
    } catch (error) {
      console.error("Error calculating route with local server:", error);
      try {
        console.log("Attempting OSRM fallback...");
        this.isUsingOSRM = true;
        const osrmRes = await fetchOSRMRoute(from.lat, from.lng, to.lat, to.lng);
        const fuelCost = (osrmRes.distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
        const speedKmh = osrmRes.duration > 0 ? (osrmRes.distance / osrmRes.duration) * 3.6 : 0;
        const tollCost = speedKmh >= 110 ? (osrmRes.distance / 1000) * this.tollPricePerKm : 0;
        return { ...osrmRes, fuelCost, tollCost, segments: [] };
      } catch (osrmError) {
        console.error("OSRM fallback also failed:", osrmError);
        const distance = this.haversineDistance(from, to);
        const duration = distance / 13.88; // ~50km/h fallback
        const fuelCost = (distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
        const tollCost = 0;
        return {
          distance,
          duration,
          path: [
            [from.lat, from.lng],
            [to.lat, to.lng],
          ],
          fuelCost,
          tollCost,
          segments: [],
        };
      }
    }
  }

  private haversineDistance(from: Location, to: Location): number {
    const R = 6371000;
    const dLat = this.toRad(to.lat - from.lat);
    const dLon = this.toRad(to.lng - from.lng);
    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos(this.toRad(from.lat)) *
        Math.cos(this.toRad(to.lat)) *
        Math.sin(dLon / 2) *
        Math.sin(dLon / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
  }

  private toRad(degrees: number): number {
    return degrees * (Math.PI / 180);
  }

  findOptimalRoute(startId: string): {
    path: [number, number][];
    totalDistance: number;
    totalDuration: number;
    totalFuelCost: number;
    totalTollCost: number;
    orderedLocations: Location[];
    segments: SegmentDetail[];
  } {
    if (this.locations.length <= 1) {
      return { path: [], totalDistance: 0, totalDuration: 0, totalFuelCost: 0, totalTollCost: 0, orderedLocations: [], segments: [] };
    }

    const route = this.locations.map((loc) => loc.id);
    return this.constructRouteResult(route);
  }

  public getSegmentCost(fromId: string, toId: string): number {
    if (this.optimizationCriterion === "fuel") {
      const exactFuel = this.fuelCostMatrix.get(fromId)?.get(toId);
      if (exactFuel !== undefined) return exactFuel;
    }
    if (this.optimizationCriterion === "toll") {
      const exactToll = this.tollCostMatrix.get(fromId)?.get(toId);
      if (exactToll !== undefined) return exactToll;
    }

    const distance = this.distanceMatrix.get(fromId)?.get(toId) || 0;
    const duration = this.durationMatrix.get(fromId)?.get(toId) || 0;
    
    switch (this.optimizationCriterion) {
      case "distance":
        return distance;
      case "time":
        return duration;
      case "fuel":
        return (distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
      case "toll": {
        const speedKmh = duration > 0 ? (distance / duration) * 3.6 : 0;
        if (speedKmh >= 110) {
          return (distance / 1000) * this.tollPricePerKm;
        }
        return 0;
      }
      case "economy": {
        const exactFuel = this.fuelCostMatrix.get(fromId)?.get(toId);
        const exactToll = this.tollCostMatrix.get(fromId)?.get(toId);
        const f = exactFuel !== undefined ? exactFuel : (distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
        const t = exactToll !== undefined ? exactToll : (duration > 0 && (distance / duration) * 3.6 >= 110 ? (distance / 1000) * this.tollPricePerKm : 0);
        return f + t;
      }
      case "smart": {
        const exactFuel = this.fuelCostMatrix.get(fromId)?.get(toId);
        const exactToll = this.tollCostMatrix.get(fromId)?.get(toId);
        const sf = exactFuel !== undefined ? exactFuel : (distance / 1000 / 100) * this.fuelConsumptionLitersPer100Km * this.fuelPricePerLiter;
        const st = exactToll !== undefined ? exactToll : (duration > 0 && (distance / duration) * 3.6 >= 110 ? (distance / 1000) * this.tollPricePerKm : 0);
        const stime = (duration / 3600) * this.timeValuePerHour;
        return sf + st + stime;
      }
      default:
        return distance;
    }
  }

  private constructRouteResult(route: string[]): {
    path: [number, number][];
    totalDistance: number;
    totalDuration: number;
    totalFuelCost: number;
    totalTollCost: number;
    orderedLocations: Location[];
    segments: SegmentDetail[];
  } {
    const path: [number, number][] = [];
    let totalDistance = 0;
    let totalDuration = 0;
    let totalFuelCost = 0;
    let totalTollCost = 0;
    const orderedLocations: Location[] = [];
    const segments: SegmentDetail[] = [];

    for (let i = 0; i < route.length - 1; i++) {
      const fromId = route[i];
      const toId = route[i + 1];

      const fromLocation = this.locations.find((loc) => loc.id === fromId);
      if (fromLocation) {
        orderedLocations.push({
          ...fromLocation,
          order: orderedLocations.length + 1,
        });
      }

      const segmentPath = this.pathMatrix.get(fromId)?.get(toId) || [];
      path.push(...segmentPath);

      const d = this.distanceMatrix.get(fromId)?.get(toId) || 0;
      const t = this.durationMatrix.get(fromId)?.get(toId) || 0;
      totalDistance += d;
      totalDuration += t;
      
      totalFuelCost += this.fuelCostMatrix.get(fromId)?.get(toId) || 0;
      totalTollCost += this.tollCostMatrix.get(fromId)?.get(toId) || 0;

      const segmentDetail = this.segmentsMatrix.get(fromId)?.get(toId) || [];
      segments.push(...segmentDetail);
    }

    const lastId = route[route.length - 1];
    const lastLocation = this.locations.find((loc) => loc.id === lastId);
    if (lastLocation) {
      orderedLocations.push({
        ...lastLocation,
        order: orderedLocations.length + 1,
      });
    }

    return { path, totalDistance, totalDuration, totalFuelCost, totalTollCost, orderedLocations, segments };
  }
}
