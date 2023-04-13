package client

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"time"

	clientmqtt "github.com/eclipse/paho.mqtt.golang"
)

var f clientmqtt.MessageHandler = func(client clientmqtt.Client, msg clientmqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func CLientStart() {
	clientmqtt.DEBUG = log.New(os.Stdout, "", 0)
	clientmqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := clientmqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("emqx_test_client")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetProtocolVersion(3)
	opts.SetPingTimeout(1 * time.Second)

	c := clientmqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	if token := c.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":80"))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
