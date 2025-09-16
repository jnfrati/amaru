package stunning

// Package to handle STUN requests using AF_XDP? Maybe...
// BIG TODO
// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/u-utils/stun/stun"
// )

// func main() {
// 	fmt.Println("Starting stun test")
// 	client := stun.NewClient()

// 	message, err := stun.NewMessage(
// 		stun.MethodBinding,
// 		stun.ClassRequest,
// 	)

// 	if err != nil {
// 		log.Fatalf("Error creating message: %v", err)
// 	}

// 	stunServerAddr := "stun.actionvoip.com:3478"

// 	err = client.Send(
// 		context.Background(),
// 		stunServerAddr,
// 		message,
// 	)
// 	if err != nil {
// 		log.Fatalf("Error sending message to stun addr %s: %v", stunServerAddr, err)
// 	}
// }
