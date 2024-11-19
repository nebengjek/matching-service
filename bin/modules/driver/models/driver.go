package models

import (
	"matching-service/bin/pkg/validator"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PickupPassanger struct {
	PassangerID string `json:"passangerId" bson:"passangerId"`
}

type TripTracker struct {
	Data DataTrip `json:"data"`
}

type DataTrip struct {
	DriverID string `json:"driverId" bson:"driverId"`
	Distance string `json:"distance" bson:"distance"`
}

type Driver struct {
	Id           string `json:"_id" bson:"_id"`
	FullName     string `json:"fullName" bson:"fullName" validate:"required,min=3,max=100"`
	MobileNumber string `json:"mobileNumber" bson:"mobileNumber" validate:"required"`
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

type TripOrderCompleted struct {
	OrderID        string  `json:"orderId" bson:"orderId"`
	RealDistance   float64 `json:"realDistance" bson:"realDistance"`
	FarePercentage float64 `json:"farePercentage" bson:"farePercentage"`
}

type MetaData struct {
	SenderID string `json:"senderId" bson:"senderId"`
	UserID   string `json:"userId" bson:"userId"`
	DriverID string `json:"driverId" bson:"driverId"`
}

type DriverAvailable struct {
	Longitude string   `json:"longitude" bson:"longitude"`
	Latitude  string   `json:"latitude" bson:"longitude"`
	MetaData  MetaData `json:"metadata" bson:"metadata"`
	Available bool     `json:"available" bson:"available"`
}

type StatusDriver struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	DriverID string             `json:"driverId" bson:"driverId"`
	SocketID string             `json:"socketId" bson:"socketId"`
	Status   bool               `json:"status" bson:"status"`
}

type DriverMatch struct {
	DriverID  string  `json:"Name"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
	Dist      float64 `json:"Dist"`
	GeoHash   int32   `json:"GeoHash"`
}

type Trip struct {
	OrderID        string  `json:"orderId" validate:"required"`
	FarePercentage float64 `json:"farePercentage" validate:"required"`
}

func (r *PickupPassanger) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
