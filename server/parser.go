package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"github.com/qedus/osmpbf"
)

// Constantes pour estimer la vitesse (km/h) par type de route
var speedLimits = map[string]float64{
	"motorway":       130,
	"motorway_link":  90,
	"trunk":          110,
	"trunk_link":     80,
	"primary":        80,
	"primary_link":   60,
	"secondary":      80,
	"secondary_link": 60,
	"tertiary":       80,
	"tertiary_link":  50,
	"unclassified":   50,
	"residential":    30,
	"living_street":  20,
}

func (g *Graph) LoadFromPBF(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	d.SetBufferSize(osmpbf.MaxBlobSize)

	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return err
	}

	// Première passe : On cherche tous les ways (routes) pour savoir quels nodes on doit garder
	fmt.Println("  [1/4] Première passe : Identification des routes...")
	
	validNodes := make(map[int64]bool)
	var ways []osmpbf.Way

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				highway, isRoad := v.Tags["highway"]
				if isRoad && speedLimits[highway] > 0 {
					for _, nodeID := range v.NodeIDs {
						validNodes[nodeID] = true
					}
					ways = append(ways, *v)
				}
			}
		}
	}

	// Deuxième passe : Extraction des coordonnées des nodes valides
	f.Seek(0, 0)
	d = osmpbf.NewDecoder(f)
	d.SetBufferSize(osmpbf.MaxBlobSize)
	d.Start(runtime.GOMAXPROCS(-1))

	fmt.Printf("  [2/4] Deuxième passe : Extraction des nœuds routiers (validNodes: %d)...\n", len(validNodes))

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				if validNodes[v.ID] {
					if _, exists := g.OSMToInternal[v.ID]; !exists {
						internalID := int32(len(g.Nodes))
						g.Nodes = append(g.Nodes, Node{
							ID:  v.ID,
							Lat: v.Lat,
							Lng: v.Lon,
						})
						g.OSMToInternal[v.ID] = internalID

						if v.Lat < g.MinLat {
							g.MinLat = v.Lat
						}
						if v.Lat > g.MaxLat {
							g.MaxLat = v.Lat
						}
						if v.Lon < g.MinLng {
							g.MinLng = v.Lon
						}
						if v.Lon > g.MaxLng {
							g.MaxLng = v.Lon
						}
					}
				}
			}
		}
	}

	fmt.Println("  [3/4] Construction du graphe (Adjacence)...")
	// Ajuster la taille de g.Edges pour correspondre à g.Nodes
	if len(g.Edges) < len(g.Nodes) {
		newEdges := make([][]Edge, len(g.Nodes))
		copy(newEdges, g.Edges)
		g.Edges = newEdges
	}

	// Maintenant on construit les edges
	for _, w := range ways {
		highway := w.Tags["highway"]
		speedKmh := speedLimits[highway]
		
		// Gérer la vitesse max si précisée dans les tags
		if maxSpeedStr, ok := w.Tags["maxspeed"]; ok {
			if speed, err := strconv.ParseFloat(maxSpeedStr, 64); err == nil {
				speedKmh = speed
			}
		}
		
		oneway := w.Tags["oneway"] == "yes"
		speedMs := speedKmh / 3.6

		for i := 0; i < len(w.NodeIDs)-1; i++ {
			nodeA := w.NodeIDs[i]
			nodeB := w.NodeIDs[i+1]

			internalA, okA := g.OSMToInternal[nodeA]
			internalB, okB := g.OSMToInternal[nodeB]
			if !okA || !okB {
				continue
			}

			nA := g.Nodes[internalA]
			nB := g.Nodes[internalB]

			distance := Haversine(nA.Lat, nA.Lng, nB.Lat, nB.Lng)
			duration := distance / speedMs
			isToll := w.Tags["toll"] == "yes" || ((highway == "motorway" || highway == "motorway_link") && w.Tags["toll"] != "no")

			// Ajouter l'arête A -> B
			g.Edges[internalA] = append(g.Edges[internalA], Edge{
				To:       internalB,
				Distance: distance,
				Duration: duration,
				Highway:  highway,
				IsToll:   isToll,
			})
			g.EdgeCount++

			// Ajouter l'arête B -> A si ce n'est pas un sens unique
			if !oneway {
				g.Edges[internalB] = append(g.Edges[internalB], Edge{
					To:       internalA,
					Distance: distance,
					Duration: duration,
					Highway:  highway,
					IsToll:   isToll,
				})
				g.EdgeCount++
			}
		}
	}

	// Construire ou reconstruire le KDTree
	fmt.Println("  [4/4] Construction de l'index spatial KD-Tree...")
	g.KDTree = BuildKDTree(g.Nodes)

	return nil
}
