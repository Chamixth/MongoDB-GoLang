package main

import (
	"MongoDB2/controllers"
	"context"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func main() {
	//Initialize the router
	r := httprouter.New()
	//Crete a new instance of the user controller
	uc := controllers.NewUserController(getClient())
	//Define the routed and their respective handlers
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user/", uc.CreateUser)
	r.PATCH("/user/:id", uc.UpdateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	//Start the server
	http.ListenAndServe("Localhost:8000", r)
}

func getClient() *mongo.Client {
	//MongoDB connection string uri
	uri := "mongodb+srv://chamith:123@cluster0.ujlq82i.mongodb.net/?retryWrites=true&w=majority"

	//Set client options
	clientOptions := options.Client().ApplyURI(uri)
	//Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	//Ping the MongoDB server to check if its running
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return client
}
