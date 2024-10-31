package handlers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
)

var firebaseApp *firebase.App

func initFirebase() {
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	firebaseApp = app
}

func init() {
	initFirebase()
}
