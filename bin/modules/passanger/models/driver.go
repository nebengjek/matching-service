package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           string `json:"_id" bson:"_id"`
	FullName     string `json:"fullName" bson:"fullName" validate:"required,min=3,max=100"`
	MobileNumber string `json:"mobileNumber" bson:"mobileNumber" validate:"required"`
	Completed    bool   `'json:"completed" bson:"completed"`
}

type DriverAvailable struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	DriverID string             `bson:"driverId"`
	SocketID string             `bson:"socketId"`
	Status   bool               `bson:"status"`
}

type WorkLog struct {
	DriverID string        `bson:"driverId" json:"driverId"`
	WorkDate string        `bson:"workdate" json:"workdate"`
	Log      []LogActivity `bson:"log" json:"log"`
}

type LogActivity struct {
	WorkTime time.Time `bson:"worktime" json:"worktime"`
	Active   bool      `bson:"active" json:"active"`
	Status   string    `bson:"status" json:"status"`
}

type LocationRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type Route struct {
	Origin      LocationRequest `json:"origin" `
	Destination LocationRequest `json:"destination"`
}

type RouteSummary struct {
	Route             Route   `json:"route"`
	MinPrice          float64 `json:"minPrice"`
	MaxPrice          float64 `json:"maxPrice"`
	BestRouteKm       float64 `json:"bestRouteKm"`
	BestRoutePrice    float64 `json:"bestRoutePrice"`
	BestRouteDuration string  `json:"bestRouteDuration"`
}

type RequestRide struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	UserId       string       `json:"userId" bson:"userId"`
}

type BroadcastPickupPassanger struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	DriverID     string       `json:"driverId" bson:"driverId"`
	SocketID     string       `json:"socketId" bson:"socketId"`
}

type DriverMatch struct {
	DriverID  string  `json:"Name"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
	Dist      float64 `json:"Dist"`
	GeoHash   int32   `json:"GeoHash"`
}
