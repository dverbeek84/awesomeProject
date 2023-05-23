package app

// Sort packages based on stlibs, remote, local.
import (
	"awesomeProject/internal/database"
	ctx "context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"

	"awesomeProject/internal/helpers"
	"awesomeProject/internal/model"
)

var secureConfig = secure.DefaultConfig()

// Normally I would create the handler package in the internal folder.
// Also, I should move the whole RabbitMQ to its own package.
func CreateOrder(context *gin.Context) {
	var order model.Order

	//  Bind the JSON data to the order struct.
	if err := context.BindJSON(&order); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Always set the state to ordered when creating to order.
	order.State = "ordered"

	// Here starts the RabbiqMQ stuff. In a real situation I should move this to its onw package.
	// The connection should be moved to the startup op the application.
	var url = fmt.Sprintf("amqp://%s:%s@%s:%d/",
		OrderServiceConfig.Queue.Username,
		OrderServiceConfig.Queue.Password,
		OrderServiceConfig.Queue.Address,
		OrderServiceConfig.Queue.Port,
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

	queueCTX, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(order)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Publish message.
	err = ch.PublishWithContext(queueCTX, "", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save object to database.
	if err := database.CreateOrder(&order); err != nil {
		context.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, order)
}

func startRESTServer() {
	gin.SetMode(gin.ReleaseMode)

	// Setup Gin engine
	r := gin.New()

	// Add logger middleware to us zerolog for no memory allocations.
	r.Use(logger.SetLogger())

	// Add recovery middleware to recover from panic.
	r.Use(gin.Recovery())

	// Add secure middleware for the strict secure settings.
	secureConfig.IsDevelopment = OrderServiceConfig.Debug
	r.Use(secure.New(secureConfig))

	// Add CORS Middleware
	r.Use(cors.Default())

	// OrderService API V1. I use the prefix to easily add multiple version if needed in the future.
	APIv1 := r.Group("/api/v1")

	// In real production env you would have a fully working CRUD api.
	order := APIv1.Group("/order")
	{
		order.POST("", CreateOrder)
	}

	// Start the gin server
	var address = fmt.Sprintf("%s:%d", OrderServiceConfig.Application.Address, OrderServiceConfig.Application.Port)
	log.Info().Msg("REST server listening on " + address)

	err := r.Run(address)
	helpers.FailOnError(err, "Cannot start REST server")
}
