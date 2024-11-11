package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id           string `json:"_id" bson:"_id"`
	FullName     string `json:"fullName" bson:"fullName" validate:"required,min=3,max=100"`
	MobileNumber string `json:"mobileNumber" bson:"mobileNumber" validate:"required"`
	Completed    bool   `'json:"completed" bson:"completed"`
}

type BeaconRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Status    string  `json:"status" validate:"required"`
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

func (r *BeaconRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
