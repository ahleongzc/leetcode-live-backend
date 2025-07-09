package messagequeue

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(
	config *config.MessageQueueConfig,
) MessageQueue {
	producerClient := New(
		config.Host,
		config.Queues,
		config.ReconnectionDelay,
		config.ReinitializationDelay,
	)
	consumerClient := New(
		config.Host,
		config.Queues,
		config.ReconnectionDelay,
		config.ReinitializationDelay,
	)

	return &RabbitMQ{
		producerClient: producerClient,
		consumerClient: consumerClient,
		resendDelay:    config.ResendDelay,
	}
}

type RabbitMQ struct {
	producerClient *RabbitMQClient
	consumerClient *RabbitMQClient
	resendDelay    time.Duration
}

func (r *RabbitMQ) Close() error {
	if err := r.producerClient.Close(); err != nil {
		return err
	}
	if err := r.consumerClient.Close(); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) StartConsuming(ctx context.Context, queue string) (<-chan *Delivery, error) {
	deliveryChan := make(chan *Delivery)

	r.consumerClient.m.Lock()
	ready := r.consumerClient.isReady
	r.consumerClient.m.Unlock()

	if !ready {
		err := waitForConnection(config.MESSAGE_QUEUE_CONNECTION_TIMEOUT, r.consumerClient)
		if err != nil {
			close(deliveryChan)
			return nil, err
		}
	}

	go func() {
		defer close(deliveryChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				deliveries, err := r.consumerClient.Consume(ctx, queue)
				if err != nil {
					r.consumerClient.logger.Printf("failed to consume from queue %s, %s", queue, err)
					return
				}
				for delivery := range deliveries {
					message := &Delivery{
						Body: delivery.Body,
						Acknowledger: &amqpAcknowledger{
							delivery: delivery,
						},
					}
					select {
					case deliveryChan <- message:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return deliveryChan, nil
}

func (r *RabbitMQ) Push(ctx context.Context, data []byte, queue string) error {
	r.producerClient.m.Lock()
	ready := r.consumerClient.isReady
	r.producerClient.m.Unlock()

	if !ready {
		err := waitForConnection(config.MESSAGE_QUEUE_CONNECTION_TIMEOUT, r.producerClient)
		if err != nil {
			return err
		}
	}

	for {
		err := r.producerClient.UnsafePush(ctx, data, queue)
		if err != nil {
			select {
			case <-r.producerClient.done:
				return common.ErrShutdown
			case <-time.After(r.resendDelay):
			}
			continue
		}
		confirm := <-r.producerClient.notifyConfirm
		if confirm.Ack {
			return nil
		}
	}
}

func waitForConnection(timeout time.Duration, client *RabbitMQClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for RabbitMQ connection: %w", common.ErrInternalServerError)
		case <-ticker.C:
			client.m.Lock()
			ready := client.isReady
			client.m.Unlock()
			if ready {
				return nil
			}
		}
	}
}

type amqpAcknowledger struct {
	delivery amqp.Delivery
}

func (a *amqpAcknowledger) Ack() error {
	return a.delivery.Ack(false)
}

func (a *amqpAcknowledger) Nack(requeue bool) error {
	return a.delivery.Nack(false, requeue)
}

func (a *amqpAcknowledger) Reject(requeue bool) error {
	return a.delivery.Reject(requeue)
}

type RabbitMQClient struct {
	m                     *sync.Mutex
	queues                []string
	logger                *log.Logger
	connection            *amqp.Connection
	channel               *amqp.Channel
	done                  chan bool
	notifyConnClose       chan *amqp.Error
	notifyChanClose       chan *amqp.Error
	notifyConfirm         chan amqp.Confirmation
	isReady               bool
	reconnectionDelay     time.Duration
	reinitializationDelay time.Duration
}

func New(addr string, queues []string, reconnectionDelay, reinitializationDelay time.Duration) *RabbitMQClient {
	client := RabbitMQClient{
		m:                     &sync.Mutex{},
		logger:                log.New(os.Stdout, "", log.LstdFlags),
		queues:                queues,
		done:                  make(chan bool),
		reconnectionDelay:     reconnectionDelay,
		reinitializationDelay: reinitializationDelay,
	}
	go client.handleReconnect(addr)
	return &client
}

func (c *RabbitMQClient) Consume(ctx context.Context, queue string) (<-chan amqp.Delivery, error) {
	c.m.Lock()
	if !c.isReady {
		c.m.Unlock()
		return nil, common.ErrNotConnected
	}
	c.m.Unlock()

	if !slices.Contains(c.queues, queue) {
		return nil, fmt.Errorf("queue %s not initialized: %w", queue, common.ErrInternalServerError)
	}

	if err := c.channel.Qos(
		1,
		0,
		false,
	); err != nil {
		return nil, err
	}

	return c.channel.ConsumeWithContext(
		ctx,
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *RabbitMQClient) handleReconnect(addr string) {
	for {
		c.m.Lock()
		c.isReady = false
		c.m.Unlock()

		conn, err := c.connect(addr)

		if err != nil {
			c.logger.Println("Failed to connect. Retrying...")

			select {
			case <-c.done:
				return
			case <-time.After(c.reconnectionDelay):
			}
			continue
		}

		if done := c.handleReInit(conn); done {
			break
		}
	}
}

func (c *RabbitMQClient) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	for _, queue := range c.queues {
		_, err = ch.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("unable to create queue %s, %s: %w", queue, err, common.ErrInternalServerError)
		}
	}

	c.changeChannel(ch)
	c.m.Lock()
	c.isReady = true
	c.m.Unlock()

	return nil
}

func (c *RabbitMQClient) changeChannel(channel *amqp.Channel) {
	c.channel = channel
	c.notifyChanClose = make(chan *amqp.Error, 1)
	c.notifyConfirm = make(chan amqp.Confirmation, 1)
	c.channel.NotifyClose(c.notifyChanClose)
	c.channel.NotifyPublish(c.notifyConfirm)
}

func (c *RabbitMQClient) handleReInit(conn *amqp.Connection) bool {
	for {
		c.m.Lock()
		c.isReady = false
		c.m.Unlock()

		err := c.init(conn)
		if err != nil {
			c.logger.Println("Failed to initialize channel. Retrying...")

			select {
			case <-c.done:
				return true
			case <-c.notifyConnClose:
				c.logger.Println("Connection closed. Reconnecting...")
				return false
			case <-time.After(c.reinitializationDelay):
			}
			continue
		}

		select {
		case <-c.done:
			return true
		case <-c.notifyConnClose:
			c.logger.Println("Connection closed. Reconnecting...")
			return false
		case <-c.notifyChanClose:
			c.logger.Println("Channel closed. Re-running init...")
		}
	}
}

func (c *RabbitMQClient) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	c.changeConnection(conn)
	return conn, nil
}

func (c *RabbitMQClient) changeConnection(connection *amqp.Connection) {
	c.connection = connection
	c.notifyConnClose = make(chan *amqp.Error, 1)
	c.connection.NotifyClose(c.notifyConnClose)
}

func (c *RabbitMQClient) UnsafePush(ctx context.Context, data []byte, queue string) error {
	c.m.Lock()
	if !c.isReady {
		c.m.Unlock()
		return fmt.Errorf("rabbit mq client not connected %w", common.ErrNotConnected)
	}
	c.m.Unlock()

	queueExists := slices.Contains(c.queues, queue)
	if !queueExists {
		return fmt.Errorf("queue %s not initialized: %w", queue, common.ErrInternalServerError)
	}

	ctx, cancel := context.WithTimeout(ctx, config.PUBLISHER_TIMEOUT)
	defer cancel()

	return c.channel.PublishWithContext(
		ctx,
		"",
		queue,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

// Close will cleanly shut down the channel and connection.
func (c *RabbitMQClient) Close() error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.isReady {
		return common.ErrAlreadyClosed
	}

	close(c.done)
	err := c.channel.Close()
	if err != nil {
		return err
	}

	err = c.connection.Close()
	if err != nil {
		return err
	}

	c.isReady = false
	return nil
}
