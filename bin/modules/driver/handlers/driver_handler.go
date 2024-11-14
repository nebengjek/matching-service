package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	driver "matching-service/bin/modules/driver"
	"matching-service/bin/modules/driver/models"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
	"matching-service/bin/pkg/log"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type passangerHandler struct {
	driverUsecaseCommand driver.UsecaseCommand
}

func NewDriverConsumer(dc driver.UsecaseCommand) kafkaPkgConfluent.ConsumerHandler {
	return &passangerHandler{
		driverUsecaseCommand: dc,
	}
}

func (i passangerHandler) HandleMessage(message *k.Message) {
	log.GetLogger().Info("consumer", fmt.Sprintf("Partition: %v - Offset: %v", message.TopicPartition.Partition, message.TopicPartition.Offset.String()), *message.TopicPartition.Topic, string(message.Value))

	var msg models.DriverAvailable
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.GetLogger().Error("consumer", "unmarshal-data", err.Error(), string(message.Value))
		return
	}

	if err := i.driverUsecaseCommand.DriverAvailable(context.Background(), models.DriverAvailable{
		Longitude: msg.Longitude,
		Latitude:  msg.Latitude,
		MetaData:  msg.MetaData,
		Available: msg.Available,
	}); err != nil {
		log.GetLogger().Error("consumer", "BroadcastPickupPassanger", err.Error(), string(message.Value))
		return
	}

	return
}
