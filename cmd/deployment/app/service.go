package app

// Sort packages based on stlibs, remote, local.
import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"awesomeProject/internal/database"
	"awesomeProject/internal/helpers"
	"awesomeProject/internal/model"
	proto "awesomeProject/pb"
)

func startService() {
	var url = fmt.Sprintf("amqp://%s:%s@%s:%d/",
		DeploymentServiceConfig.Queue.Username,
		DeploymentServiceConfig.Queue.Password,
		DeploymentServiceConfig.Queue.Address,
		DeploymentServiceConfig.Queue.Port,
	)

	// Connect to RabbitMQ.
	conn, err := amqp.Dial(url)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Connect to channel.
	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare queue deployment.
	q, err := ch.QueueDeclare("deployment", false, false, false, false, nil)
	helpers.FailOnError(err, "Failed to declare a queue")

	// Register consumer to queue.
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	helpers.FailOnError(err, "Failed to register a consumer")

	// Create a go channel to forever till interrupt.
	foreverChannel := make(chan bool)

	// Create a go channel to receive log events.
	logChannel := make(chan *zerolog.Event)

	// Start goroutine.
	go func() {
		for d := range msgs {
			var order model.Order

			if err := json.Unmarshal(d.Body, &order); err != nil {
				log.Err(err).Msg("Cannot unmarshal json data")
			}

			var deployment = model.Deployment{
				State:         "deploying",
				ApplicationID: order.ApplicationID,
			}

			database.DB.Create(&deployment)

			// Here I simulate a long-running deployment, off course this should be a real deployment,
			database.DB.Joins("Application").Find(&deployment)
			logChannel <- log.Info().Str("application", deployment.Application.Name).Str("state", deployment.State).Uint("order", order.ID)
			time.Sleep(time.Second * 5)

			// Here I simulate a long-running configuration, off course this should be a real configuration,
			database.DB.Model(&deployment).Update("state", "configuring")
			logChannel <- log.Info().Str("application", deployment.Application.Name).Str("state", deployment.State).Uint("order", order.ID)
			time.Sleep(time.Second * 5)

			// Here the deployment and configuration is done. So change the order state to done over GRPC.
			// In a real situation i would also use RabbitMQ there is no need here to use RabbitMQ and GRPC for the communication back.
			// I would consider using GRPC in situations where latency and bidirectional communcation is needed.
			database.DB.Model(&deployment).Update("state", "done")
			logChannel <- log.Info().Str("application", deployment.Application.Name).Str("state", deployment.State).Uint("order", order.ID)

			target := fmt.Sprintf("%s:%d", DeploymentServiceConfig.GRPC.Address, DeploymentServiceConfig.GRPC.Port)
			conn, _ := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			client := proto.NewOrderServiceClient(conn)

			// Note how we are calling the UpdateOrder method on the server.
			// This is available to us through the auto-generated code.
			_, err := client.UpdateOrderToDone(context.Background(), &proto.OrderRequest{
				Id:    strconv.Itoa(int(order.ID)),
				State: "deployed",
			})
			helpers.FailOnError(err, "Failed to update order")
		}
	}()

	// Log events from within the goroutine.
	// Without this channel the log will not be visible,
	for logEvent := range logChannel {
		logEvent.Send()
	}

	// Block forever unit interrupt.
	// For production use you want gracefully shutdown.
	<-foreverChannel
}
