package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/julienschmidt/httprouter"
)

var app *firebase.App
var err error

//Index Function for each address
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome, it's home page")
}

//SendMsg Function for send message
func SendMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	topic := "LOL"
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
		},
		Notification: &messaging.Notification{
			"Title",
			"test",
			"",
		},
		Topic: topic,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Fprintf(w, "Successfully sent message:", response)
}

// Main Function
func main() {

	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Declare router variable
	router := httprouter.New()

	// Connect each function to each address
	router.GET("/", Index)
	router.GET("/msg", SendMsg)

	log.Fatal(http.ListenAndServe(":8080", router))
}
