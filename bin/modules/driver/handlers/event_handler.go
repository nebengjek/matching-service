package handlers

import (
	driver "matching-service/bin/modules/driver"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
)

func InitPassangerEventHandler(passanger driver.UsecaseCommand, kc kafkaPkgConfluent.Consumer) {

	kc.SetHandler(NewPassangerConsumer(passanger))
	kc.Subscribe("request-ride")

}
