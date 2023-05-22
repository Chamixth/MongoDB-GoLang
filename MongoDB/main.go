package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user entity
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// MongoDB configuration
const (
	uri        = "mongodb+srv://chamith:123@cluster0.ujlq82i.mongodb.net/?retryWrites=true&w=majority"
	dbName     = "sample_db"
	collection = "users"
	timeout    = 5 * time.Second
)

// MongoDB client and collection instances
var (
	client *mongo.Client
	coll   *mongo.Collection
)

// Initialize initializes the MongoDB client and collection
func Initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	coll = client.Database(dbName).Collection(collection)
	return nil
}

// CreateUser creates a new user record
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = coll.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser retrieves a user record by username
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var user User
	err := coll.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user record by username
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = coll.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{
		"$set": bson.M{
			"email":    updatedUser.Email,
			"password": updatedUser.Password,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser deletes a user record by username
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	_, err := coll.DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}

func main() {
	// Initialize MongoDB connection
	err := Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{username}", GetUser).Methods("GET")
	router.HandleFunc("/users/{username}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{username}", DeleteUser).Methods("DELETE")

	// Start the server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/
//https://github.com/Dankuzo1/convencionAPI/blob/master/convencionAPI/01EjemploMongodb/main.go
//https://github.com/AthithyaJayadevan/mongo_go/blob/master/main.go
//https://www.leonvillamayor.org/2018/12/api-restful-con-golang/
//https://www.youtube.com/watch?v=VzBGi_n65iU
//https://www.youtube.com/watch?v=zICaTPBkupY
