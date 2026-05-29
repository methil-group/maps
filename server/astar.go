package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
)

// --- Priority Queue implementation for A* (Zero Allocations on heap operations) ---
type Item struct {
	NodeID   int32
	Priority float64 // f-score (g + h)
	Index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// --- End Priority Queue ---

// FindNearestNode recherche le nœud le plus proche en O(log N) à l'aide du KD-Tree
func (g *Graph) FindNearestNode(lat, lng float64) (int32, error) {
	if g.KDTree == nil || g.KDTree.Root == nil {
		return 0, errors.New("graphe vide")
	}
	bestID := g.KDTree.FindNearest(lat, lng)
	if bestID == -1 {
		return 0, errors.New("aucun nœud trouvé")
	}
	return bestID, nil
}

type RoutingParams struct {
	Criterion   string  // "distance", "time", "fuel", "toll", "balanced"
	FuelPrice   float64 // prix au litre
	Consumption float64 // L/100km
	TollPrice   float64 // prix au km (ex: 0.13)
	TimeValue   float64 // prix de l'heure (ex: 20.0)
	TollWeight  float64 // facteur pour les péages (ex: 1.0)
}

type SegmentDetail struct {
	FromLat  float64 `json:"fromLat"`
	FromLng  float64 `json:"fromLng"`
	ToLat    float64 `json:"toLat"`
	ToLng    float64 `json:"toLng"`
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Highway  string  `json:"highway"`
	Cost     float64 `json:"cost"`
}

type RouteResult struct {
	Path     [][2]float64    `json:"path"`
	Distance float64         `json:"distance"`
	Duration float64         `json:"duration"`
	FuelCost float64         `json:"fuelCost"`
	TollCost float64         `json:"tollCost"`
	Segments []SegmentDetail `json:"segments"`
}

func (g *Graph) getEdgeCost(edge Edge, params RoutingParams) float64 {
	switch params.Criterion {
	case "distance":
		return edge.Distance
	case "economy":
		fuelLitres := (edge.Distance / 1000.0) * (params.Consumption / 100.0)
		fuelCost := fuelLitres * params.FuelPrice
		tollCost := 0.0
		if edge.IsToll {
			tollCost = (edge.Distance / 1000.0) * params.TollPrice
		}
		return fuelCost + tollCost
	case "toll":
		cost := 0.0
		if edge.IsToll {
			cost = (edge.Distance / 1000.0) * params.TollPrice
		}
		return cost + (edge.Distance / 1000000.0)
	case "time":
		return edge.Duration
	case "smart":
		fuelLitres := (edge.Distance / 1000.0) * (params.Consumption / 100.0)
		fuelCost := fuelLitres * params.FuelPrice
		tollCost := 0.0
		if edge.IsToll {
			tollCost = (edge.Distance / 1000.0) * params.TollPrice
		}
		timeCost := (edge.Duration / 3600.0) * params.TimeValue
		return fuelCost + tollCost + timeCost
	default:
		return edge.Duration
	}
}

var ErrOutOfBounds = errors.New("out_of_bounds")

func (g *Graph) IsCoordinatesInMapBounds(lat, lng float64) bool {
	if len(g.Nodes) == 0 {
		return false
	}
	clampLat := math.Max(g.MinLat, math.Min(lat, g.MaxLat))
	clampLng := math.Max(g.MinLng, math.Min(lng, g.MaxLng))
	return Haversine(lat, lng, clampLat, clampLng) <= 10000
}

// PlaneDistance calcule la distance plane locale approximée (très rapide, pas de trigonométrie lourde)
func PlaneDistance(lat1, lng1, lat2, lng2, factorLng float64) float64 {
	const R = 6371000.0
	const degToRad = math.Pi / 180.0
	dLat := (lat2 - lat1) * degToRad
	dLng := (lng2 - lng1) * degToRad * factorLng
	return R * math.Sqrt(dLat*dLat + dLng*dLng)
}

func (g *Graph) FindRoute(startLat, startLng, endLat, endLng float64, params RoutingParams) (*RouteResult, error) {
	if !g.IsCoordinatesInMapBounds(startLat, startLng) || !g.IsCoordinatesInMapBounds(endLat, endLng) {
		return nil, ErrOutOfBounds
	}

	startID, err := g.FindNearestNode(startLat, startLng)
	if err != nil {
		return nil, err
	}

	// Vérifie si les coordonnées de départ sont sur la carte chargée (distance < 10 km)
	startDist := Haversine(startLat, startLng, g.Nodes[startID].Lat, g.Nodes[startID].Lng)
	if startDist > 10000 {
		return nil, ErrOutOfBounds
	}

	endID, err := g.FindNearestNode(endLat, endLng)
	if err != nil {
		return nil, err
	}

	// Vérifie si les coordonnées d'arrivée sont sur la carte chargée (distance < 10 km)
	endDist := Haversine(endLat, endLng, g.Nodes[endID].Lat, g.Nodes[endID].Lng)
	if endDist > 10000 {
		return nil, ErrOutOfBounds
	}

	if startID == endID {
		return &RouteResult{
			Path: [][2]float64{{startLng, startLat}},
		}, nil
	}

	endNode := g.Nodes[endID]

	// 1. Précalcul du facteur de projection pour la longitude (trigonométrie unique)
	factorLng := math.Cos(((startLat + endLat) / 2.0) * math.Pi / 180.0)

	// 2. Remplacement des maps par des slices pré-alloués (Évite les overheads du Garbage Collector)
	numNodes := len(g.Nodes)
	gScore := make([]float64, numNodes)
	for i := range gScore {
		gScore[i] = math.MaxFloat64
	}
	gScore[startID] = 0

	cameFrom := make([]int32, numNodes)
	for i := range cameFrom {
		cameFrom[i] = -1
	}

	edgeUsed := make([]Edge, numNodes)
	visited := make([]bool, numNodes)

	// 4. File d'attente à priorité réutilisable (zéro allocations pendant l'exploration)
	items := make([]Item, numNodes)
	for i := range items {
		items[i].NodeID = int32(i)
		items[i].Priority = math.MaxFloat64
		items[i].Index = -1
	}

	pq := make(PriorityQueue, 0, 1024)
	heap.Init(&pq)

	items[startID].Priority = 0
	heap.Push(&pq, &items[startID])

	found := false

	steps := 0
	for pq.Len() > 0 {
		steps++
		item := heap.Pop(&pq).(*Item)
		current := item.NodeID

		if current == endID {
			found = true
			break
		}

		if visited[current] {
			continue
		}
		visited[current] = true

		for _, edge := range g.Edges[current] {
			if visited[edge.To] {
				continue
			}

			cost := g.getEdgeCost(edge, params)
			tentativeGScore := gScore[current] + cost

			if tentativeGScore < gScore[edge.To] {
				cameFrom[edge.To] = current
				edgeUsed[edge.To] = edge
				gScore[edge.To] = tentativeGScore

				// 3. Heuristique Plane ultra-rapide (pas de cos/sin dans la boucle)
				hDist := PlaneDistance(g.Nodes[edge.To].Lat, g.Nodes[edge.To].Lng, endNode.Lat, endNode.Lng, factorLng)

				var h float64
				switch params.Criterion {
				case "distance":
					h = hDist
				case "time":
					h = hDist / 36.1 // 130 km/h (temps minimum théorique)
				case "economy":
					h = (hDist / 1000.0) * (params.Consumption / 100.0) * params.FuelPrice
				case "toll":
					h = 0
				case "smart":
					hTime := (hDist / 36.1) / 3600.0 * params.TimeValue
					hFuel := (hDist / 1000.0) * (params.Consumption / 100.0) * params.FuelPrice
					h = hTime + hFuel
				default:
					h = hDist / 36.1
				}

				priority := tentativeGScore + h
				neighborItem := &items[edge.To]

				if neighborItem.Index == -1 {
					neighborItem.Priority = priority
					heap.Push(&pq, neighborItem)
				} else if priority < neighborItem.Priority {
					neighborItem.Priority = priority
					heap.Fix(&pq, neighborItem.Index)
				}
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("aucun itinéraire trouvé (startID=%d, endID=%d, steps=%d, pqLen=%d, criterion=%s)", startID, endID, steps, pq.Len(), params.Criterion)
	}

	// Reconstruire le chemin et agréger les données
	var path [][2]float64
	var segments []SegmentDetail
	var totalDist, totalDuration, totalFuel, totalToll float64

	curr := endID
	for {
		node := g.Nodes[curr]
		path = append([][2]float64{{node.Lng, node.Lat}}, path...)

		parent := cameFrom[curr]
		if parent == -1 {
			break
		}

		edge := edgeUsed[curr]
		parentNode := g.Nodes[parent]

		totalDist += edge.Distance
		totalDuration += edge.Duration

		fuelLitres := (edge.Distance / 1000.0) * (params.Consumption / 100.0)
		totalFuel += fuelLitres * params.FuelPrice

		toll := 0.0
		if edge.IsToll {
			toll = (edge.Distance / 1000.0) * params.TollPrice
			totalToll += toll
		}

		segments = append([]SegmentDetail{{
			FromLat:  parentNode.Lat,
			FromLng:  parentNode.Lng,
			ToLat:    node.Lat,
			ToLng:    node.Lng,
			Distance: edge.Distance,
			Duration: edge.Duration,
			Highway:  edge.Highway,
			Cost:     g.getEdgeCost(edge, params),
		}}, segments...)

		curr = parent
	}

	return &RouteResult{
		Path:     path,
		Distance: totalDist,
		Duration: totalDuration,
		FuelCost: totalFuel,
		TollCost: totalToll,
		Segments: segments,
	}, nil
}
