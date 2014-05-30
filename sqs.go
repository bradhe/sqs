package sqs

import (
	"log"
	"github.com/crowdmob/goamz/aws"
	goamz_sqs "github.com/crowdmob/goamz/sqs"
)

const (
	DefaultMessagesToDequeue = 1
)

type Queue struct {
	service *goamz_sqs.SQS
	queue *goamz_sqs.Queue
	ch chan string
}

func (self *Queue) runReadMessages(ch chan string) {
	for {
		resp, err := self.queue.ReceiveMessage(DefaultMessagesToDequeue)

		if err != nil {
			// Report the error. What do?
			log.Printf("Error receiving message from queue: %v", err)
		} else {
			// We got a message so let's disbatch it!
			for _, msg := range resp.Messages {
				ch <- msg.Body

				// Blindly confirm that we got it. Fantastic, innit?
				self.queue.DeleteMessage(&msg)
			}
		}
	}
}

func (self *Queue) ReadMessages() chan string {
	return self.ch
}

func (self *Queue) PublishMessages(message string) error {
	_, err := self.queue.SendMessage(message)
	return err
}

func NewQueue(queueName string, auth aws.Auth, region string) (*Queue, error) {
	queue := new(Queue)
	queue.service = goamz_sqs.New(auth, aws.Regions[region])

	// TODO: Figure out how to make sure we don't duplicate this across things?
	// Might not be a big deal.
	q, err := queue.service.CreateQueue(queueName)

	if err != nil {
		return nil, err
	}

	queue.queue = q
	queue.ch = make(chan string)
	go queue.runReadMessages(queue.ch)

	return queue, nil
}
