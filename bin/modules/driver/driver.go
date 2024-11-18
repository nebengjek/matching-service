package driver

import (
	"context"

	"matching-service/bin/modules/driver/models"
	orderModel "matching-service/bin/modules/passanger/models"
	"matching-service/bin/pkg/utils"
)

type UsecaseQuery interface {
}

type UsecaseCommand interface {
	DriverAvailable(ctx context.Context, payload models.DriverAvailable) error
	PickupPassanger(ctx context.Context, userId string, payload models.PickupPassanger) utils.Result
	CompletedTrip(ctx context.Context, userId string, payload models.Trip) utils.Result
}

type MongodbRepositoryQuery interface {
	FindDriver(ctx context.Context, userId string) <-chan utils.Result
	FindOrderPassanger(ctx context.Context, psgId string) <-chan utils.Result
	FindActiveOrderPassanger(ctx context.Context, psgId string) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	NewObjectID(ctx context.Context) string
	UpsertDriver(ctx context.Context, data models.DriverAvailable) <-chan utils.Result
	UpdateOneTripOrder(ctx context.Context, data orderModel.TripOrder) <-chan utils.Result
	CompletedTripOrder(ctx context.Context, data orderModel.TripOrder, trip models.TripTracker) <-chan utils.Result
}
