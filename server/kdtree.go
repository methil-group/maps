package main

import (
	"math"
	"sort"
)

// KDNode représente un nœud dans le KD-Tree
type KDNode struct {
	Index    int32
	Lat, Lng float64
	Left     *KDNode
	Right    *KDNode
	Axis     int // 0 pour Latitude, 1 pour Longitude
}

// KDTree est la structure de l'index spatial KD-Tree
type KDTree struct {
	Root *KDNode
}

// BuildKDTree construit récursivement un KD-Tree à partir des nœuds du graphe
func BuildKDTree(nodes []Node) *KDTree {
	if len(nodes) == 0 {
		return &KDTree{Root: nil}
	}
	indices := make([]int32, len(nodes))
	for i := range indices {
		indices[i] = int32(i)
	}
	return &KDTree{
		Root: buildKDTreeRec(indices, nodes, 0),
	}
}

func buildKDTreeRec(indices []int32, nodes []Node, depth int) *KDNode {
	if len(indices) == 0 {
		return nil
	}
	axis := depth % 2

	// Trier les indices selon l'axe sélectionné
	if axis == 0 {
		sort.Slice(indices, func(i, j int) bool {
			return nodes[indices[i]].Lat < nodes[indices[j]].Lat
		})
	} else {
		sort.Slice(indices, func(i, j int) bool {
			return nodes[indices[i]].Lng < nodes[indices[j]].Lng
		})
	}

	medianIdx := len(indices) / 2
	medianNodeIdx := indices[medianIdx]

	return &KDNode{
		Index:  medianNodeIdx,
		Lat:    nodes[medianNodeIdx].Lat,
		Lng:    nodes[medianNodeIdx].Lng,
		Axis:   axis,
		Left:   buildKDTreeRec(indices[:medianIdx], nodes, depth+1),
		Right:  buildKDTreeRec(indices[medianIdx+1:], nodes, depth+1),
	}
}

// FindNearest recherche l'indice du nœud le plus proche des coordonnées spécifiées
func (t *KDTree) FindNearest(lat, lng float64) int32 {
	if t.Root == nil {
		return -1
	}
	bestIndex := t.Root.Index
	bestDistSq := planeDistanceSq(lat, lng, t.Root.Lat, t.Root.Lng)

	t.searchNearest(t.Root, lat, lng, &bestIndex, &bestDistSq)
	return bestIndex
}

func (t *KDTree) searchNearest(curr *KDNode, lat, lng float64, bestIndex *int32, bestDistSq *float64) {
	if curr == nil {
		return
	}

	// Calcul de la distance plane locale au carré
	distSq := planeDistanceSq(lat, lng, curr.Lat, curr.Lng)
	if distSq < *bestDistSq {
		*bestDistSq = distSq
		*bestIndex = curr.Index
	}

	var diff float64
	if curr.Axis == 0 {
		diff = lat - curr.Lat
	} else {
		diff = lng - curr.Lng
	}

	var nearChild, farChild *KDNode
	if diff < 0 {
		nearChild = curr.Left
		farChild = curr.Right
	} else {
		nearChild = curr.Right
		farChild = curr.Left
	}

	// Recherche dans le sous-arbre le plus proche
	t.searchNearest(nearChild, lat, lng, bestIndex, bestDistSq)

	// Déterminer si une partie de la boîte de découpage opposée peut être plus proche que la meilleure distance actuelle
	var axisPlaneDistSq float64
	if curr.Axis == 0 {
		axisPlaneDistSq = diff * diff
	} else {
		// Ajustement de la longitude par le cosinus de la latitude de recherche
		cosLat := math.Cos(lat * math.Pi / 180.0)
		axisPlaneDistSq = (diff * cosLat) * (diff * cosLat)
	}

	if axisPlaneDistSq < *bestDistSq {
		t.searchNearest(farChild, lat, lng, bestIndex, bestDistSq)
	}
}

// planeDistanceSq est une approximation plane rapide de la distance au carré en degrés décimaux
func planeDistanceSq(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := lat2 - lat1
	latMoyRad := (lat1 + lat2) * 0.5 * (math.Pi / 180.0)
	dLng := (lng2 - lng1) * math.Cos(latMoyRad)
	return dLat*dLat + dLng*dLng
}
