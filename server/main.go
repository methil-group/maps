package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Démarrage du Route Optimiser Server...")

	start := time.Now()
	// Initialiser le graphe
	graph := NewGraph()

	// Charger les données OSM
	fmt.Println("⏳ Chargement des données OSM...")
	
	files, err := os.ReadDir("data")
	if err != nil {
		log.Fatalf("❌ Erreur lors de la lecture du dossier data: %v", err)
	}

	loadedCount := 0
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".pbf" || strings.HasSuffix(file.Name(), ".osm.pbf")) {
			fmt.Printf("📦 Chargement de %s...\n", file.Name())
			err := graph.LoadFromPBF(filepath.Join("data", file.Name()))
			if err != nil {
				fmt.Printf("⚠️ Erreur lors du chargement de %s: %v\n", file.Name(), err)
				continue
			}
			loadedCount++
		}
	}

	if loadedCount == 0 {
		fmt.Println("⚠️ Aucun fichier .osm.pbf trouvé dans le dossier data/")
		fmt.Println("💡 Placez vos fichiers de carte dans le dossier server/data/")
	} else {
		fmt.Printf("✅ %d fichier(s) chargé(s) en %s. Nœuds: %d, Arêtes: %d\n", loadedCount, time.Since(start), len(graph.Nodes), graph.EdgeCount)
	}

	// Configurer l'API Gin
	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/api/route", func(c *gin.Context) {
		startLng, _ := strconv.ParseFloat(c.Query("startLng"), 64)
		startLat, _ := strconv.ParseFloat(c.Query("startLat"), 64)
		endLng, _ := strconv.ParseFloat(c.Query("endLng"), 64)
		endLat, _ := strconv.ParseFloat(c.Query("endLat"), 64)

		criterion := c.DefaultQuery("criterion", "time")
		fuelPrice, _ := strconv.ParseFloat(c.DefaultQuery("fuelPrice", "1.80"), 64)
		consumption, _ := strconv.ParseFloat(c.DefaultQuery("consumption", "6.5"), 64)
		tollPrice, _ := strconv.ParseFloat(c.DefaultQuery("tollPrice", "0.13"), 64)
		timeValue, _ := strconv.ParseFloat(c.DefaultQuery("timeValue", "20.0"), 64)

		if startLng == 0 || startLat == 0 || endLng == 0 || endLat == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètres startLng, startLat, endLng, endLat requis"})
			return
		}

		params := RoutingParams{
			Criterion:   criterion,
			FuelPrice:   fuelPrice,
			Consumption: consumption,
			TollPrice:   tollPrice,
			TimeValue:   timeValue,
		}

		fmt.Printf("🔍 Recherche d'itinéraire [%s] de [%f, %f] à [%f, %f]\n", criterion, startLat, startLng, endLat, endLng)
		routeStart := time.Now()

		result, err := graph.FindRoute(startLat, startLng, endLat, endLng, params)
		if err != nil {
			if errors.Is(err, ErrOutOfBounds) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "out_of_bounds", "message": "Les coordonnées sont en dehors de la carte chargée."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("✅ Itinéraire trouvé en %s (Distance: %.2fm)\n", time.Since(routeStart), result.Distance)

		c.JSON(http.StatusOK, result)
	})

	fmt.Println("🌍 Serveur API en écoute sur http://localhost:8080")
	r.Run(":8080")
}
