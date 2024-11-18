package commands

import (
	"context"

	user "matching-service/bin/modules/passanger"
	"matching-service/bin/modules/passanger/models"
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

func (c commandMongodbRepository) CreateTripOrder(ctx context.Context, data models.TripOrder) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		err := c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "trip-orders",
			Document: bson.M{
				"orderId":       data.OrderID,
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
			return
		}

		output <- utils.Result{
			Data: data.OrderID,
		}
	}()

	return output
}

func (c commandMongodbRepository) UpdateOneTripOrder(ctx context.Context, orderId string, data models.TripOrder) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		err := c.mongoDb.UpsertOne(mongodb.UpsertOne{
			CollectionName: "trip-orders",
			Filter: bson.M{
				"orderId": orderId,
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
