package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
)

var (
	app *firebase.App
	err error
	ctx context.Context
)

func firebaseApp() {
	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
}

func msgRequest(c *gin.Context) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	msgType := c.PostForm("msgType")
	podName := c.PostForm("podName")

	title := ""
	body := ""

	if msgType == "newPod" {
		title = "새로운 팟이 생성되었습니다!"
		body = podName + "팟에 참가해 보시는건 어떠신가요?"
	} else if msgType == "completePod" {
		title = podName + "팟이 결성되었습니다!"
		body = podName + "팟은 아쉽지만 다음에 노립시다"
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
		},
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: "",
		},
		Topic: podName,
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	c.String(200, "Successfully sent message:", response)
}

// Main Function
func main() {
	firebaseApp()
	r := gin.Default()
	r.Use(gin.Logger())
	r.POST("/msgRequest", msgRequest)

	log.Fatal(http.ListenAndServe(":8080", r))
}
