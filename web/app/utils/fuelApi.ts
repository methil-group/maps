export type FuelType = "Gazole" | "E10" | "SP98" | "SP95" | "E85" | "GPLc";

export interface FuelPriceData {
  type: FuelType;
  price: number;
}

const FUEL_API_URL = "https://data.economie.gouv.fr/api/explore/v2.1/catalog/datasets/prix-des-carburants-en-france-flux-instantane-v2/records?limit=20";

export const fetchAverageFuelPrice = async (fuelType: FuelType): Promise<number | null> => {
  try {
    const response = await fetch(FUEL_API_URL);
    if (!response.ok) {
      throw new Error("Failed to fetch fuel prices");
    }

    const data = await response.json();
    let totalPrice = 0;
    let count = 0;

    const priceKey = `${fuelType.toLowerCase()}_prix`;

    for (const record of data.results) {
      if (record[priceKey] !== null && record[priceKey] !== undefined) {
        totalPrice += record[priceKey];
        count++;
      }
    }

    if (count === 0) return null;
    return totalPrice / count;
  } catch (error) {
    console.error("Error fetching fuel price:", error);
    return null;
  }
};
