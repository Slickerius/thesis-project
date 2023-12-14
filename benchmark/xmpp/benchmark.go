package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"mellium.im/communique/internal/client"
	"mellium.im/communique/internal/client/event"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/ping"
	"mellium.im/xmpp/stanza"
)

var (
	messageIds map[string]chan bool = make(map[string]chan bool)
)

func c2cMessageTest(ctx context.Context, c *client.Client, to jid.JID, id string) (float64, error) {
	message := event.ChatMessage{
		Message: stanza.Message{
			To:   to,
			Type: stanza.ChatMessage,
			ID:   id,
		},
		OriginID: stanza.OriginID{
			ID: id,
		},
		Body: "This is a test message",
		Sent: true,
	}

	start := time.Now()
	_, err := c.SendMessage(ctx, message)
	if err != nil {
		return 0, err
	}
	<-messageIds[id]

	elapsed := time.Since(start)
	fmt.Printf("Message %s took %s\n", id, elapsed)

	return elapsed.Seconds(), nil
}

func c2sPingTest(ctx context.Context, c *client.Client) (float64, error) {
	start := time.Now()
	err := ping.Send(ctx, c.Session, c.LocalAddr().Domain())
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)
	fmt.Printf("Ping took %s\n", elapsed)

	return elapsed.Seconds(), nil
}

func c2cBatchTest() float64 {
	var wg sync.WaitGroup
	var idList []string
	var (
		totalTime float64
		totalTest int
		totalMu   sync.Mutex
	)
	updateTotal := func(elapsedTime float64) {
		totalMu.Lock()
		defer totalMu.Unlock()
		totalTime += elapsedTime
		totalTest++
	}
	for i := 0; i < int(clientCount); i++ {
		id := randomID()
		idList = append(idList, id)
		messageIds[id] = make(chan bool)
	}
	for i := 0; i < int(clientCount); i++ {
		idx1 := i
		idx2 := (i + 1) % int(clientCount)
		// fmt.Printf("Starting Message test from client %d to client %d\n", idx1, idx2)
		go func() {
			wg.Add(1)
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			elapsedTime, err := c2cMessageTest(ctx, clientList[idx1], clientList[idx2].LocalAddr(), idList[idx1])
			if err != nil {
				fmt.Printf("Message from client %d failed\n", idx1)
				return
			}
			updateTotal(elapsedTime)
		}()
	}
	wg.Wait()
	for _, id := range idList {
		delete(messageIds, id)
	}

	return totalTime / float64(totalTest)
}

func c2sBatchTest() float64 {
	var wg sync.WaitGroup
	var (
		totalTime float64
		totalTest int
		totalMu   sync.Mutex
	)
	updateTotal := func(elapsedTime float64) {
		totalMu.Lock()
		defer totalMu.Unlock()
		totalTime += elapsedTime
		totalTest++
	}
	for i := 0; i < int(clientCount); i++ {
		idx := 1
		go func() {
			wg.Add(1)
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			elapsedTime, err := c2sPingTest(ctx, clientList[idx])
			if err != nil {
				fmt.Printf("Ping from client %d failed\n", idx)
				return
			}
			updateTotal(elapsedTime)
		}()
	}
	wg.Wait()

	return totalTime / float64(totalTest)
}
