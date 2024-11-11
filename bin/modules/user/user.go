package user

import (
	"context"

	"matching-service/bin/modules/user/models"
	"matching-service/bin/pkg/utils"
	//"go.mongodb.org/mongo-driver/bson"
)

type UsecaseQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	GetUser(userId string, ctx context.Context) utils.Result
	FindDriver(userId string, ctx context.Context) utils.Result
}

type UsecaseCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	PostLocation(userId string, payload models.LocationSuggestionRequest, ctx context.Context) utils.Result
}

type MongodbRepositoryQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	FindOne(userId string, ctx context.Context) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	NewObjectID(ctx context.Context) string
}
