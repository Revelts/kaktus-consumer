package main

import (
	"kaktus-task/Connections"
	"kaktus-task/Controllers"
	"log"
)

func main() {
	var err error
	err = Connections.InitiateConnection()
	if err != nil {
		log.Panicf("[ERROR] : %s", err)
		return
	}

	err = Connections.RabbitChannel.ExchangeDeclare("kaktus", "topic",
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR] : %s", err)
		return
	}

	Controllers.ViewAllThreads()
	Controllers.ViewThreadDetails()
	Controllers.CommentThread()
	Controllers.CreateThread()
	Controllers.LikeThread()
	Controllers.ReplyComment()

	var forever chan struct{}
	<-forever
}
