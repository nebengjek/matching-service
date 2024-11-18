package queries

import (
	"context"
	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	order "matching-service/bin/modules/passanger/models"
	"matching-service/bin/pkg/databases/mongodb"
	"matching-service/bin/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewQueryMongodbRepository(mongodb mongodb.MongoDBLogger) driver.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
	}
}

func (q queryMongodbRepository) FindDriver(ctx context.Context, userId string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var driver models.Driver
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &driver,
			CollectionName: "user",
			Filter: bson.M{
				"userId": userId,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: driver,
		}

	}()

	return output
}

func (q queryMongodbRepository) FindOrderPassanger(ctx context.Context, passangerId string) <-chan utils.Result {
	output := make(chan utils.Result)
	go func() {
		defer close(output)
		var trip order.TripOrder
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &trip,
			CollectionName: "trip-orders",
			Filter: bson.M{
				"passengerId": passangerId,
				"$or": []bson.M{
					{
						"status": bson.M{"$ne": "completed"},
					},
					{
						"status": bson.M{"$ne": "ontheway"},
					},
					{
						"status": "request-pickup",
					},
				},
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}
		output <- utils.Result{
			Data: trip,
		}

	}()

	return output
}

func (q queryMongodbRepository) FindActiveOrderPassanger(ctx context.Context, orderId string) <-chan utils.Result {
	output := make(chan utils.Result)
	go func() {
		defer close(output)
		var trip order.TripOrder
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &trip,
			CollectionName: "trip-orders",
			Filter: bson.M{
				"orderId": orderId,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}
		output <- utils.Result{
			Data: trip,
		}

	}()

	return output
}
