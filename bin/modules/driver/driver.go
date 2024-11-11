package user

import (
	"context"

	"matching-service/bin/modules/driver/models"
	"matching-service/bin/pkg/utils"
	//"go.mongodb.org/mongo-driver/bson"
)

type UsecaseQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	DetailTrip(psgId string, driverId string, ctx context.Context) utils.Result
}

type UsecaseCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	BroadcastPickupPassanger(payload models.RequestRide, ctx context.Context) error
}

type MongodbRepositoryQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	FindWorkLog(driverId string, date string, ctx context.Context) <-chan utils.Result
	FindDriver(userId string, ctx context.Context) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	NewObjectID(ctx context.Context) string
	UpsertBeacon(data models.WorkLog, ctx context.Context) <-chan utils.Result
}
