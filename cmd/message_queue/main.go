package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue := NewMessageQueue(ctx, 30, 5000)
	queue.Start()
}

// MessageQueue is a message queue
type MessageQueue struct {
	ctx           context.Context
	queue         []string
	producerMutex *sync.Mutex
	consumerMutex *sync.Mutex
	BatchSize     int
	Intervals     time.Duration

	sizeTrigger chan bool
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(ctx context.Context, batchSize, intervalsMillisecond int) *MessageQueue {
	return &MessageQueue{
		ctx:           ctx,
		queue:         []string{},
		producerMutex: &sync.Mutex{},
		consumerMutex: &sync.Mutex{},
		BatchSize:     batchSize,
		Intervals:     time.Duration(intervalsMillisecond) * time.Millisecond,
		sizeTrigger:   make(chan bool),
	}
}

// Start starts the message queue
func (m *MessageQueue) Start() {
	// start the producer
	go func() {
		for {
			m.producer()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// start the consumer
	go func() {
		ticker := time.NewTicker(m.Intervals)
		defer ticker.Stop()

		go m.monitorQueue()

		for {
			select {
			case <-ticker.C:
				fmt.Println("[consumer] time trigger")
				m.consumer()
			case <-m.sizeTrigger:
				fmt.Println("[consumer] size trigger")
				m.consumer()
			}

		}
	}()

	quit := make(chan sync.Singal, 1)
}

func (m *MessageQueue) monitorQueue() {
	for {
		size := len(m.queue)
		if size >= m.BatchSize {
			m.sizeTrigger <- true      // 队列长度达到阈值时发送信号
			for size == len(m.queue) { // 等待队列长度发生变化
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (m *MessageQueue) producer() {
	m.producerMutex.Lock()
	defer m.producerMutex.Unlock()

	// add the message to the queue
	m.queue = append(m.queue, "message")
}

func (m *MessageQueue) consumer() {
	m.consumerMutex.Lock()
	defer m.consumerMutex.Unlock()

	size := m.BatchSize
	if len(m.queue) < size {
		size = len(m.queue)
	}
	// process the messages
	fmt.Printf("[consumer] batch processing, size:%d, queue:%d\n", size, len(m.queue))
	// clear the processed messages
	m.queue = m.queue[size:]
}
