package main

import (
	"math"
)

// Node représente une intersection routière
type Node struct {
	ID  int64   // ID OpenStreetMap original
	Lat float64
	Lng float64
}

// Edge représente un segment de route entre deux Node
type Edge struct {
	To       int32   // ID interne (Internal_ID) de la destination
	Distance float64 // en mètres
	Duration float64 // en secondes (basé sur la vitesse max)
	Highway  string  // type de route (motorway, primary, etc.)
	IsToll   bool    // si la route est à péage
}

// Graph représente le réseau routier complet optimisé
type Graph struct {
	Nodes         []Node
	Edges         [][]Edge
	OSMToInternal map[int64]int32
	EdgeCount     int
	MinLat        float64
	MaxLat        float64
	MinLng        float64
	MaxLng        float64
	KDTree        *KDTree
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:         make([]Node, 0),
		Edges:         make([][]Edge, 0),
		OSMToInternal: make(map[int64]int32),
		MinLat:        90.0,
		MaxLat:        -90.0,
		MinLng:        180.0,
		MaxLng:        -180.0,
	}
}

// Haversine calcule la distance en mètres entre deux points GPS (utilisé pour les frontières et calcul final précis)
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Rayon de la terre en mètres
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
