# Go SQS Wrapper

This is a wrapper for quickly/easily reading from and publishing to SQS. It's
meant to be pulled in really quickly in places where you may need it.

## Delivery Guarantees

SQS has it's own delivery guarantees that I can't control. That said, the
implementation of this wrapper may or may not extend that--not really sure yet.

## Example

```go
package main

import (
  "fmt"
  "github.com/bradhe/sqs"
)

func main() {
  queue, err := sqs.NewQueue("my-queue", "<Your AWS Access Key>", "<Your AWS Secrety Key>", "us-east-1")

  if err != nil {
    panic(err)
  }

  // Send a message to all subscribers.
  queue.PublishMessage("Hello, World!")

  // Start reading messages from the queue.
  messages := queue.ReadMessages()

  for message := range messages {
    fmt.Printf("Found message: %v", message)
  }
}
```
