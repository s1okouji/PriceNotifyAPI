package main

import "fmt"

type MessageCreateEvent map[string]interface{}

func (event *MessageCreateEvent) Handle() error {
	data := *event
	content := data["content"]
	channel_id := data["channel_id"]
	fmt.Printf("content: %v\n", content)
	fmt.Printf("channel_id: %v\n", channel_id)
	return nil
}
