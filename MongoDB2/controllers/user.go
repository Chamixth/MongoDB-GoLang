package controllers

import (
	"MongoDB2/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//extract the id from the url
	id := p.ByName("id")

	//convert the id into an object id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Access the users collection in the database
	collection := uc.client.Database("Test-Mongo").Collection("users")
	u := models.User{}

	//find the user with the specified id
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//converr the user object to json string
	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set the response header and write the json response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//create an empty user object
	u := models.User{}

	//decode the request body into the user object
	err := json.NewDecoder(r.Body).Decode(&u)
	//if there is an error, return a bad request status
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//access the users collection in the database
	collection := uc.client.Database("Test-Mongo").Collection("users")

	//insert the user object into the database
	_, err = collection.InsertOne(context.TODO(), u)
	//if there is an error, return a server error status
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//convert the user object to json string
	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set the response header and write the json response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//extract the id from the url
	id := p.ByName("id")
	//convert the id into an object id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//create an empty user object
	u := models.User{}
	//decode the request body into the user object
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//access the users collection in the database
	collection := uc.client.Database("Test-Mongo").Collection("users")
	//define the filter to find the user with the specified id
	filter := bson.M{"_id": objID}
	//define the update operation to set the new user object
	update := bson.M{"$set": u}

	//update the user with the specified id
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Return a status ok response
	w.WriteHeader(http.StatusOK)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//extract the id from the url
	id := p.ByName("id")
	//convert the id into an object id
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//access the users collection in the database
	collection := uc.client.Database("Test-Mongo").Collection("users")
	//define the filter to find the user with the specified id
	filter := bson.M{"_id": objID}
	//delete the user with the specified id
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Return a status ok response
	w.WriteHeader(http.StatusOK)
}

//https://www.youtube.com/watch?v=zICaTPBkupY
