package sqs

import (
	"time"
	"log"
	"github.com/crowdmob/goamz/aws"
	goamz_sqs "github.com/crowdmob/goamz/sqs"
)

const (
	DefaultMessagesToDequeue = 1
	DefaultBatchSize = 10
)

// Takes an array and turns it in to chunks of size size. i.e. if size is 10
// and arr has 20 arrays, 2 arrays of 10 are returned.
func chunk(arr []string, size int) [][]string {
	if len(arr) == 0 {
		return [][]string{}
	}

	// That was easy...
	if len(arr) < size {
		return [][]string{arr}
	}

	results := make([][]string, 0)
	offset := 0

	var current []string

	for i := range arr {
		if i % size == 0 {
			current = make([]string, size)
			results = append(results, current)

			if i > 0 {
				offset++
			}
		}

		current[i % size] = arr[i]
	}

	return results
}

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

func (self *Queue) PublishMessage(message string) error {
	_, err := self.queue.SendMessage(message)
	return err
}

func (self *Queue) PublishMessages(messages []string) error {
	arr := chunk(messages, DefaultBatchSize)

	for _, a := range arr {
		sent := false
		var err error

		// Loop is for retries.
		for i := 0; i < 10; i++ {
			_, err = self.queue.SendMessageBatchString(a)

			if err == nil {
				sent = true
				break
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}

		// Tell someone it failed entirely.
		if !sent {
			return err
		}
	}

	return nil
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
