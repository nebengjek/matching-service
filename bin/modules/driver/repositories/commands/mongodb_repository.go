package commands

import (
	"context"
	"strconv"

	user "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	order "matching-service/bin/modules/passanger/models"
	"matching-service/bin/pkg/databases/mongodb"
	"matching-service/bin/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commandMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewCommandMongodbRepository(mongodb mongodb.MongoDBLogger) user.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
	}
}

func (c commandMongodbRepository) NewObjectID(ctx context.Context) string {
	return primitive.NewObjectID().Hex()
}

func (c commandMongodbRepository) UpsertDriver(ctx context.Context, data models.DriverAvailable) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		err := c.mongoDb.UpsertOne(mongodb.UpsertOne{
			CollectionName: "driver-available",
			Filter: bson.M{
				"driverId": data.MetaData.DriverID,
			},
			Document: bson.M{
				"driverId": data.MetaData.DriverID,
				"socketId": data.MetaData.SenderID,
				"status":   data.Available,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: nil,
		}

	}()

	return output
}

func (c commandMongodbRepository) UpdateOneTripOrder(ctx context.Context, data order.TripOrder) <-chan utils.Result {
	output := make(chan utils.Result)
	go func() {
		defer close(output)
		err := c.mongoDb.UpdateOne(mongodb.UpdateOne{
			CollectionName: "trip-orders",
			Filter: bson.M{
				"orderId": data.OrderID,
			},
			Document: bson.M{
				"passengerId":   data.PassengerID,
				"driverId":      data.DriverID,
				"origin":        data.Origin,
				"destination":   data.Destination,
				"status":        data.Status,
				"createdAt":     data.CreatedAt,
				"updatedAt":     data.UpdatedAt,
				"estimatedFare": data.EstimatedFare,
				"distanceKm":    data.DistanceKm,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: nil,
		}

	}()

	return output
}

func (c commandMongodbRepository) CompletedTripOrder(ctx context.Context, data order.TripOrder, trip models.TripTracker) <-chan utils.Result {
	output := make(chan utils.Result)
	realDistance, _ := strconv.ParseFloat(trip.Data.Distance, 64)
	go func() {
		defer close(output)
		err := c.mongoDb.UpdateOne(mongodb.UpdateOne{
			CollectionName: "trip-orders",
			Filter: bson.M{
				"orderId": data.OrderID,
			},
			Document: bson.M{
				"passengerId":   data.PassengerID,
				"driverId":      data.DriverID,
				"origin":        data.Origin,
				"destination":   data.Destination,
				"status":        data.Status,
				"createdAt":     data.CreatedAt,
				"updatedAt":     data.UpdatedAt,
				"estimatedFare": data.EstimatedFare,
				"distanceKm":    data.DistanceKm,
				"realDistance":  realDistance,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: nil,
		}

	}()

	return output
}
