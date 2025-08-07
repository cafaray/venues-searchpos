package main

import (
	"context"
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

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("supporters").Collection("venues")

	r := gin.Default()

	r.GET("/nearby", getNearbyStadiums)
	r.GET("/inside", isPointInsideStadium)

	r.Run(":8080")
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
