package handlers

import (
	passanger "matching-service/bin/modules/passanger"
	kafkaPkgConfluent "matching-service/bin/pkg/kafka/confluent"
)

func InitPassangerEventHandler(passanger passanger.UsecaseCommand, kc kafkaPkgConfluent.Consumer) {

	kc.SetHandler(NewPassangerConsumer(passanger))
	kc.Subscribe("request-ride")

}
