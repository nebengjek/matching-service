package user

import (
	"context"

	"matching-service/bin/modules/passanger/models"
	"matching-service/bin/pkg/utils"
	//"go.mongodb.org/mongo-driver/bson"
)

type UsecaseQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	DetailTrip(ctx context.Context, psgId string, driverId string) utils.Result
}

type UsecaseCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	BroadcastPickupPassanger(ctx context.Context, payload models.RequestRide) error
}

type MongodbRepositoryQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	FindDriver(ctx context.Context, driverId string) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	NewObjectID(ctx context.Context) string
}
