package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var client *mongo.Client

type Trade struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date 	  time.Time           `json:"date,omitempty" bson:"date,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Stock     Stock				  `json:"stock,omitempty" bson:"stock,omitempty"'`
}

type Stock struct {
	Name       string   `json:"name,omitempty" bson:"name,omitempty"`
	Value      int   	`json:"value,omitempty" bson:"value,omitempty"`
	Volume 	   int   	`json:"volume,omitempty" bson:"volume,omitempty"`
	Total	   int	  	`json:"total,omitempty" bson:"total,omitempty"`
}

func createTrade(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var trade Trade
	trade.Date = time.Now()
	_ = json.NewDecoder(request.Body).Decode(&trade)
	collection := client.Database("fdm_trades").Collection("trades")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, trade)
	json.NewEncoder(response).Encode(result)
}

func getTradeById(response http.ResponseWriter, request *http.Request) { }
func getTrades(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var trades []Trade
	collection := client.Database("fdm_trades").Collection("trades")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var trade Trade
		cursor.Decode(&trade)
		trades = append(trades, trade)
	}
	if err := cursor.Err(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(writer).Encode(trades)
}
func getTradesByUser(response http.ResponseWriter, request *http.Request) { }



func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://jWhitfield:qwertyPotato99@fdm05-06-t1f4g.mongodb.net/fdm_trades?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println("connected to database...")
	router := mux.NewRouter()
	router.HandleFunc("/api/trades", getTrades).Methods("GET") //TODO admin only
	router.HandleFunc("/api/trades/{username}", getTradesByUser).Methods("GET")
	router.HandleFunc("/api/trades/{id}", getTradeById).Methods("GET")
	router.HandleFunc("/api/trades/create", createTrade).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))
}

