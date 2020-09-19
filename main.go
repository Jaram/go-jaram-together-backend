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

func deletePod(group string, podName string, uid string) bool {
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	_, e := client.Collection("groupInfo").Doc(group).Collection("pods").Doc(uid).Delete(ctx)
	if e == nil {
		return true
	}
	defer client.Close()
	return false
}

func msgRequest(c *gin.Context) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	podTime := c.PostFormArray("podTime")
	msgType := c.PostForm("msgType")
	podName := c.PostForm("podName")
	group := c.PostForm("group")
	uid := c.PostForm("uid")

	title := ""
	body := ""
	topic := ""

	if msgType == "newPod" {
		title = "새로운 팟이 생성되었습니다!"
		body = podName + "팟에 참가해 보시는건 어떠신가요?"
		topic = group
		go scheduler(podTime[0], podTime[1], podTime[2], podTime[3], podName, group, uid)
	} else if msgType == "completePod" {
		title = podName + "팟이 결성되었습니다!"
		body = podName + "팟에 참가하지 못하신분은 아쉽지만 다음에 노립시다"
		topic = group
	} else if msgType == "readyReminder" {
		title = podName + "팟이 10분뒤 시작됩니다!"
		body = podName + "팟 여러분들은 준비해 주시기 바랍니다!"
		topic = group + "/" + podName
	} else if msgType == "startPod" {
		title = podName + "팟이 시작합니다!"
		body = podName + "팟 여러분 다들 모이셨나요?"
		topic = group + "/" + podName
	} else if msgType == "endPod" {
		title = podName + "팟이 끝났어요!"
		body = podName + "팟 여러분들 수고하셨습니다!"
		topic = group + "/" + podName
	} else if msgType == "deletePod" {
		if deletePod(group, podName, uid) {
			c.String(200, "Successfully deleted")
		} else {
			c.String(500, "delete failed")
		}
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
		Topic: topic,
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
