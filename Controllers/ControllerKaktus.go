package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"kaktus-task/Connections"
	"kaktus-task/Controllers/Dto/request"
	"kaktus-task/Controllers/Dto/response"
	"kaktus-task/Repository"
	"log"
	"strconv"
	"time"
)

func ViewAllThreads() {

	q, err := Connections.RabbitChannel.QueueDeclare("consumer_view_all_thread",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "view_all_thread", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	go func() {
		for d := range data {
			var queueError error
			repository := Repository.AllRepository.Threads
			resp, queueError := repository.ViewThreads()
			if queueError != nil {
				log.Panicf("[ERROR]: %s", queueError)
				return
			}

			var dataBytes []byte
			dataBytes, queueError = json.Marshal(resp)
			if queueError != nil {
				log.Panicf("[ERROR]: %s", queueError)
				return
			}
			queueError = Connections.RabbitChannel.Publish("", d.ReplyTo,
				false,
				false,
				amqp091.Publishing{
					ContentType:   "application/json",
					Body:          dataBytes,
					CorrelationId: d.CorrelationId,
				})
			log.Printf(" [x] Sent %s\n", string(dataBytes))
			if queueError != nil {
				log.Panicf("[ERROR]: %s", queueError)
				return
			}
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}

func ViewThreadDetails() {
	q, err := Connections.RabbitChannel.QueueDeclare("consumer_view_thread_details",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "view_thread_details", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	go func() {
		for d := range data {
			var tId string
			err = json.Unmarshal(d.Body, &tId)

			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			threadId, err := strconv.Atoi(tId)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			keys := fmt.Sprintf("thread::%d", threadId)
			public := Repository.AllRepository.Public
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
			var respMap map[string]string
			respMap, err = public.RedisHGetAllCache(keys)

			if len(respMap) != 0 {
				err = Connections.RabbitChannel.Publish("", d.ReplyTo,
					false,
					false,
					amqp091.Publishing{
						ContentType:   "application/json",
						Body:          []byte(respMap[keys]),
						CorrelationId: d.CorrelationId,
					})
				if err != nil {
					log.Panicf("[ERROR]: %s", err)
					return
				}
				log.Printf(" [x] Sent %s\n", respMap[keys])
				return
			}
			repository := Repository.AllRepository.Threads
			var resp response.ViewThreadDetail
			resp, err = repository.ViewThreadDetail(threadId)

			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			mapTask := make(map[string]interface{})
			dataBytes, _ := json.Marshal(resp)
			mapTask[keys] = dataBytes
			err = public.RedisSetHMCache(keys, mapTask, time.Minute*10)

			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			err = Connections.RabbitChannel.Publish("", d.ReplyTo,
				false,
				false,
				amqp091.Publishing{
					ContentType:   "application/json",
					Body:          dataBytes,
					CorrelationId: d.CorrelationId,
				})
			log.Printf(" [x] Sent %s\n", string(dataBytes))

			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}

func CommentThread() {
	q, err := Connections.RabbitChannel.QueueDeclare("consumer_comment_thread",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "comment_thread", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	var requestBody request.CommentThread
	go func() {
		for d := range data {
			err = json.Unmarshal(d.Body, &requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			repository := Repository.AllRepository.Threads
			resp, err := repository.CommentThread(requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			err = UpdateRedisViewThreadDetailIfExists(resp.ThreadId)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}

func ReplyComment() {
	q, err := Connections.RabbitChannel.QueueDeclare("consumer_reply_comment",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "reply_comment", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	var requestBody request.ReplyComment
	go func() {
		for d := range data {
			err = json.Unmarshal(d.Body, &requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			repository := Repository.AllRepository.Threads
			resp, err := repository.ReplyComment(requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			err = UpdateRedisViewThreadDetailIfExists(resp.ThreadId)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}

func CreateThread() {
	q, err := Connections.RabbitChannel.QueueDeclare("consumer_create_thread",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "create_thread", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	var requestBody request.CreateThread
	go func() {
		for d := range data {
			err = json.Unmarshal(d.Body, &requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			repository := Repository.AllRepository.Threads
			_, err := repository.CreateThread(requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}

func LikeThread() {
	q, err := Connections.RabbitChannel.QueueDeclare("consumer_like_thread",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	err = Connections.RabbitChannel.QueueBind("", "like_thread", "kaktus", false, nil)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	data, err := Connections.RabbitChannel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	var requestBody request.LikeThread
	go func() {
		for d := range data {
			err = json.Unmarshal(d.Body, &requestBody)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			repository := Repository.AllRepository.Threads
			_, err := repository.LikeThread(requestBody)

			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			err = UpdateRedisViewThreadDetailIfExists(requestBody.ThreadId)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}
			log.Printf("[Queue : %s] Received a message: %s", q.Name, d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages from [%s]. To exit press CTRL+C", q.Name)
}
func UpdateRedisViewThreadDetailIfExists(threadId int) (err error) {

	// Check If Redis Key is Set
	keys := fmt.Sprintf("thread::%d", threadId)
	public := Repository.AllRepository.Public
	var respMap map[string]string
	respMap, err = public.RedisHGetAllCache(keys)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	if len(respMap) != 0 {
		// If Set then update the viewthreaddetail with new value
		repository := Repository.AllRepository.Threads
		var resp response.ViewThreadDetail
		resp, err = repository.ViewThreadDetail(threadId)

		var dataBytes []byte
		mapTask := make(map[string]interface{})
		dataBytes, err = json.Marshal(resp)
		mapTask[keys] = dataBytes
		err = public.RedisSetHMCache(keys, mapTask, time.Minute*10)
		if err != nil {
			log.Panicf("[ERROR]: %s", err)
			return
		}
		return
	}
	return
}
