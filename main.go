package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	log.Println("🚀 Starting venues-searchpos service...")
	
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI environment variable is not set")
	}
	log.Printf("✅ MONGO_URI found: %s", mongoURI[:20]+"...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("🔌 Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	log.Println("✅ MongoDB connected successfully")

	collection = client.Database("supporters").Collection("venues")
	log.Println("✅ Collection initialized")
}

func main() {
	log.Println("🌐 Setting up HTTP server...")
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🔌 Using port: %s", port)
	
	r := gin.Default()

	r.GET("/nearby", getNearbyStadiums)
	r.GET("/inside", isPointInsideStadium)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	log.Printf("🚀 Starting server on 0.0.0.0:%s", port)
	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

// GET /nearby?long=2.15&lat=41.38&maxDistance=5000
func getNearbyStadiums(c *gin.Context) {
	longStr, latStr, distStr := c.Query("long"), c.Query("lat"), c.Query("maxDistance")
	longitude, _ := strconv.ParseFloat(longStr, 64)
	latitude, _ := strconv.ParseFloat(latStr, 64)
	maxDistance, _ := strconv.ParseInt(distStr, 10, 64)

	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
				"$maxDistance": maxDistance,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GET /inside?long=2.1228&lat=41.3809
func isPointInsideStadium(c *gin.Context) {
	longStr, latStr := c.Query("long"), c.Query("lat")
	longitude, _ := strconv.ParseFloat(longStr, 64)
	latitude, _ := strconv.ParseFloat(latStr, 64)

	filter := bson.M{
		"location": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{longitude, latitude},
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, gin.H{"inside": false})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inside": true, "venue": result["name"]})
}

// EntryPoint is the function entry for HTTP request
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Running venue-searchpos-service!")
}
