package models

import "matching-service/bin/pkg/validator"

type PickupPassanger struct {
	PassangerID string `json:"passangerId" bson:"passangerId"`
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

type DriverMatch struct {
	DriverID  string  `json:"Name"`
	Longitude float64 `json:"Longitude"`
	Latitude  float64 `json:"Latitude"`
	Dist      float64 `json:"Dist"`
	GeoHash   int32   `json:"GeoHash"`
}

func (r *PickupPassanger) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
