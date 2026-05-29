# 🚀 Methil Maps

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Nuxt Version](https://img.shields.io/badge/Nuxt-4.0+-00DC82?style=flat&logo=nuxt.js)](https://nuxt.com/)
[![TailwindCSS](https://img.shields.io/badge/TailwindCSS-3.4+-06B6D4?style=flat&logo=tailwindcss)](https://tailwindcss.com/)

Planning a trip with multiple stops can be a headache, especially when you're trying to balance time, fuel costs, and toll roads. **Methil Maps** is a high-performance tool designed to make that process effortless, calculating the optimal path based on real-world OpenStreetMap data and your personal preferences.

---

## ✨ Key Features

*   **🎯 Multi-Criteria Optimization**: Choose between different routing strategies:
    *   **Fastest**: The quickest route regardless of cost.
    *   **Shortest**: The minimum distance path.
    *   **Economic**: Minimizes total cost (Fuel + Tolls).
    *   **Smart (Balanced)**: The perfect trade-off between time and money.
    *   **Toll-free**: Avoids toll roads to plan the best free route.
*   **⛽ Personalized Costs**: Input your car's actual consumption (L/100km) and select fuel types (Diesel, E10, SP95/98, LPG) with live-calculated average prices to get hyper-accurate cost estimates.
*   **🌐 Hybrid Routing Engine**: Seamlessly falls back to the public OSRM (Open Source Routing Machine) API if any coordinates are outside the downloaded map area (out-of-bounds), combining global routing capability with local high-performance custom-cost pathfinding.
*   **🗺️ Raw OSM Processing**: Unlike generic wrappers, our Go engine parses high-volume `.pbf` data directly, understanding road types, speed limits, and toll sections at a granular level.
*   **📍 Interactive Map Selection**: Long-press on the map to add destinations instantly, drag-and-drop stops list items to re-order, or use manual up/down controls.
*   **🗺️ Live Geolocation & Tracking**: Opt-in to track your position in real-time, rendered as a custom pulsing blue dot on the map. The map automatically centers on your location on startup or when clearing stops, and you can instantly add a static snapshot of your location to the route.
*   **📈 Rich Diagnostic & Leg Analysis**: An in-depth analysis modal offering:
    *   Estimated total trip budget and CO₂ footprint based on engine type.
    *   Fuel vs. Tolls budget share visualization progress bars.
    *   Leg-by-leg breakdowns showing distance, duration, fuel cost, toll cost, and total direct cost.
    *   Automatic highlighting of the **Most Expensive** and **Longest** legs.
*   **🚗 External Navigators Sync**:
    *   **Google Maps Redirection**: Export the exact computed route sequence with up to 75 coordinate waypoints to lock the path in Google Maps directions.
    *   **GPX Route Exporter**: Download standard GPX files containing stops as waypoints and the precise track route shape for outdoor navigators (Garmin, Strava, OsmAnd, etc.).
*   **🎨 Premium UI/UX Design**: A sleek, modern dashboard built with Nuxt 4, Vue 3, and Tailwind CSS, featuring CartoDB responsive light/dark maps, collapsible widgets, and glassmorphic glowing overlay cards.
*   **🌍 Multi-Language & Settings Persistence**: Fully localized in French and English with settings automatically saved in `localStorage` and URL query parameter sync for bookmarking or sharing.

---

## 🏗️ Architecture

The project is architected for speed and scalability:

1.  **Backend (Go)**: A custom-built routing engine using the A* algorithm. It handles PBF parsing, graph construction, fast bounding box bounds checks, and high-speed pathfinding.
2.  **Frontend (Nuxt)**: A responsive single-page application built with Nuxt 4, Vue 3, and Tailwind CSS, featuring collapsible option widgets and glassmorphic overlay cards. It computes segments sequentially to match the user's defined order.

---

## 🧠 How the Smart Criterion (Balanced) Works

The **Smart** optimization criterion balances travel time against financial cost (fuel and tolls) based on your custom **Value of Time** (expressed in €/h or equivalent currency per hour).

It minimizes the following generalized cost formula:
`Generalized Cost = Fuel Cost + Toll Cost + (Duration in hours * Value of Time)`

### Mathematical Example

Here is a comparison of three possible routes calculated by the server for a specific trip, with a **Value of Time** set to **18 €/h**:

1. **Option A: Economy (No tolls, slowest)**
   - Duration: 133 minutes (2.22 hours)
   - Fuel + Tolls: 20.73 € + 0.00 € = **20.73 €**
   - **Generalized Cost**: 20.73 € + (2.22 h * 18 €/h) = 20.73 € + 39.96 € = **60.69 €**

2. **Option B: Smart (Balanced compromise)**
   - Duration: 86 minutes (1.44 hours)
   - Fuel + Tolls: 20.48 € + 12.35 € = **32.83 €**
   - **Generalized Cost**: 32.83 € + (1.44 h * 18 €/h) = 32.83 € + 25.92 € = **58.75 €**

3. **Option C: Fastest (Time-focused, most expensive)**
   - Duration: 77 minutes (1.28 hours)
   - Fuel + Tolls: 20.43 € + 15.79 € = **36.22 €**
   - **Generalized Cost**: 36.22 € + (1.28 h * 18 €/h) = 36.22 € + 23.04 € = **59.26 €**

### Why Option B is Chosen

* **Option B vs Option A (Economy)**: By paying an extra **12.10 €** (32.83 € - 20.73 €), you save **47 minutes** (0.78 hours). The cost of saving this time is `12.10 € / 0.78 h = 15.51 €/h`. Since this is **less** than your time value of 18 €/h, the algorithm decides it is worth paying the toll to save the time.
* **Option C (Fastest) vs Option B**: To save an extra **9 minutes** (0.15 hours), you would have to pay **3.39 €** more (36.22 € - 32.83 €). The cost of saving this extra time is `3.39 € / 0.15 h = 22.60 €/h`. Since this is **more** than your time value of 18 €/h, the algorithm decides it is not worth paying the extra toll.

Therefore, **Option B (Smart)** is the mathematically optimal choice under a 18 €/h preference.

> [!NOTE]
> Dividing the total trip cost by the total duration (e.g. `32.83 € / 1.43 h = 22.91 €/h`) does not represent the criterion threshold, as fuel is consumed on any route and the parameter acts on the *marginal* rate of saving time between alternative options.

---

## 🛠️ Installation & Setup

### 1. Prerequisites
- [Go](https://go.dev/doc/install) (1.21 or later)
- [Node.js](https://nodejs.org/) (v18 or later)
- [wget](https://www.gnu.org/software/wget/) (for map downloads)

### 2. Backend Setup (Server)
```bash
cd server

# 1. Configure your maps
# Edit MAPS.txt to include the Geofabrik URLs you need
nano MAPS.txt

# 2. Download map data
make maps
# Or alternatively: chmod +x scripts/download_maps.sh && ./scripts/download_maps.sh

# 3. Start the engine
go run .
# Or use: make run
```
The server will start on `http://localhost:8080`.

### 3. Frontend Setup (Web)
```bash
cd web

# 1. Install dependencies
npm install

# 2. Launch dev server
npm run dev
```
The application will be available at `http://localhost:3000` (default Nuxt port).

---

## 🐳 Docker Deployment

### Backend (Go Engine)
You can run the backend engine using Docker:
```bash
docker build -t route-optimizer-server ./server
docker run -p 8080:8080 -v $(pwd)/server/data:/app/data route-optimizer-server
```

### Frontend (Nuxt Web UI)
You can run the frontend client using Docker (ideal for Coolify, CapRover, etc.):
```bash
docker build -t route-optimizer-web ./web
docker run -p 3000:3000 --env VITE_BACKEND_URL=http://localhost:8080 route-optimizer-web
```

---

## ⚙️ Configuration

### Server
- **`MAPS.txt`**: List of `.osm.pbf` URLs to download.
- **`data/`**: Directory where `.pbf` files are stored and indexed.

### Web
- **Environment variables** (or Coolify build variables):
    - `VITE_BACKEND_URL`: The base URL of the Go backend engine (defaults to `http://localhost:8080`).

---

## 🤝 Contributing

Contributions are welcome! Whether it's optimizing the A* implementation, adding new UI features, or improving documentation, feel free to open a Pull Request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/amazing-feature`)
3. Commit your Changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the Branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

*Made with ❤️ for road trippers everywhere.*
