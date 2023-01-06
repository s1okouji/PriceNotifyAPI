package main

type Payload struct {
	Op int         `json:"op"`
	D  interface{} `json:"d"`
	S  int         `json:"s"`
	T  string      `json:"t"`
}

type Identify struct {
	Token      string      `json:"token"`
	Properties interface{} `json:"properties"`
	Intents    int         `json:"intents"`
}

type EventHandler interface {
	Handle() error
}
