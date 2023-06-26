package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for user login
type LoginResponse struct {
	Token string `json:"Login is Successfull."`
}

// SignupResponse represent the response body for user signup
type SignUpResponse struct {
	Token string `json:""Signup is Successfull""`
}

// MongoDB configuration
const (
	uri        = "mongodb+srv://chamith:123@cluster0.ujlq82i.mongodb.net/?retryWrites=true&w=majority"
	dbName     = "sample_db2"
	collection = "Users2"
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
func CreateUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = coll.InsertOne(context.TODO(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	token := generateToken(user.Username)
	signUpResponse := SignUpResponse{
		Token: token,
	}

	return c.JSON(http.StatusCreated, signUpResponse)
}

// GetUser retrieves a user record by username
func GetUser(c echo.Context) error {
	username := c.Param("username")

	var user User
	err := coll.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user record by username
func UpdateUser(c echo.Context) error {
	username := c.Param("username")

	var updatedUser User
	err := c.Bind(&updatedUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = coll.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{
		"$set": bson.M{
			"email":    updatedUser.Email,
			"password": updatedUser.Password,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser deletes a user record by username
func DeleteUser(c echo.Context) error {
	username := c.Param("username")

	_, err := coll.DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "User deleted successfully")
}

// Login handles user login
func Login(c echo.Context) error {
	var loginRequest LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Check if the user exists and password is correct
	var user User
	err = coll.FindOne(context.TODO(), bson.M{
		"username": loginRequest.Username,
		"password": loginRequest.Password,
	}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	// Generate and return a token as the response

	token := generateToken(user.Username)
	loginResponse := LoginResponse{
		Token: token,
	}

	return c.JSON(http.StatusOK, loginResponse)
}
func getAll(c echo.Context) error {
	var users []bson.D
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		c.Logger().Print(err)
		return err
	}
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.Logger().Print(err)
		return err
	}

	var response []User
	for _, result := range users {
		var user User
		for _, elem := range result {
			switch elem.Key {
			case "username":
				user.Username = elem.Value.(string)
			case "email":
				user.Email = elem.Value.(string)
			case "password":
				user.Password = elem.Value.(string)
			}
		}
		response = append(response, user)
	}

	return c.JSON(http.StatusOK, response)
}

// generateToken generates a token for the given username
func generateToken(username string) string {

	return username
}

func main() {
	// Initialize MongoDB connection
	err := Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Create Echo instance
	e := echo.New()

	// Define routes
	e.POST("/users", CreateUser)
	e.GET("/users/:username", GetUser)
	e.PUT("/users/:username", UpdateUser)
	e.DELETE("/users/:username", DeleteUser)
	e.POST("/login", Login)
	e.GET("/getAll", getAll)

	// Start the server
	log.Println("Server started on port 8080")
	log.Fatal(e.Start(":8080"))
}
