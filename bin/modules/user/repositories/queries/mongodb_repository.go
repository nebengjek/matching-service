package queries

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"matching-service/bin/modules/user"
	"matching-service/bin/modules/user/models"
	"matching-service/bin/pkg/databases/mongodb"
	"matching-service/bin/pkg/utils"
)

type queryMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewQueryMongodbRepository(mongodb mongodb.MongoDBLogger) user.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
	}
}

func (q queryMongodbRepository) FindOne(userId string, ctx context.Context) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var user models.User
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &user,
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
			Data: user,
		}

	}()

	return output
}
