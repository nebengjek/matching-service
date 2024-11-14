package user

import (
	"context"

	"matching-service/bin/modules/driver/models"
	"matching-service/bin/pkg/utils"
)

type UsecaseQuery interface {
}

type UsecaseCommand interface {
	DriverAvailable(ctx context.Context, payload models.DriverAvailable) error
}

type MongodbRepositoryQuery interface {
}

type MongodbRepositoryCommand interface {
	NewObjectID(ctx context.Context) string
	UpsertDriver(ctx context.Context, data models.DriverAvailable) <-chan utils.Result
}
