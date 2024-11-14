package handlers

import (
	driver "matching-service/bin/modules/driver"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
)

func InitPassangerEventHandler(driver driver.UsecaseCommand, kc kafkaPkgConfluent.Consumer) {

	kc.SetHandler(NewDriverConsumer(driver))
	kc.Subscribe("driver-available")

}
