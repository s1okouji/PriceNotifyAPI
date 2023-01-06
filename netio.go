package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/websocket"
)

func Connect(host string) (*websocket.Conn, int) {
	origin := "https://" + host
	url := "wss://" + host
	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		log.Fatal(err)
	}
	var d net.Dialer
	d.Timeout = time.Minute
	d.LocalAddr, err = net.ResolveTCPAddr("tcp", ":50000")

	if err != nil {
		log.Fatal(err)
	}

	ws, err = websocket.DialConfig(config)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	err = ReceiveJson(ws, &data)
	if err != nil {
		log.Fatal(err)
	}
	if int(data["op"].(float64)) != 10 {
		log.Fatal("This is not Hello Event!")
	}
	ret := int(data["d"].(map[string]interface{})["heartbeat_interval"].(float64))
	return ws, ret
}

func SendJson(ws *websocket.Conn, pl *Payload) error {
	fmt.Printf("%+v", pl)
	err := websocket.JSON.Send(ws, pl)
	if err != nil {
		return err
	}
	return nil
}

func ReceiveJson(ws *websocket.Conn, mp *map[string]interface{}) error {
	var message string
	websocket.Message.Receive(ws, &message)

	fmt.Println(message)
	json.Unmarshal([]byte(message), mp)
	return nil
}
