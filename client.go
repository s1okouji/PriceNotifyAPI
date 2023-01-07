package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/s1okouji/pnabot_client/dto"
	"github.com/s1okouji/pnabot_client/service"
	"github.com/s1okouji/pnabot_client/util"
	"golang.org/x/net/websocket"
)

var last_sequence int
var resume_url string
var ws *websocket.Conn
var interval int
var start_url string
var session_id string
var errch chan error
var bot_token string

func main() {
	service.SetUp()
	start_url = "gateway.discord.gg/?v=10&encoding=json"
	bot_token = os.Getenv("bot_token")
	ws, interval = Connect(start_url)
	defer ws.Close()
	props := &map[string]string{"os": "linux", "browser": "my_app", "device": "my_app"}
	// props := &map[string]string{"os": "windows", "browser": "my_app", "device": "my_app"}
	intent := 1<<15 + 1<<12 + 1<<9
	quit := make(chan os.Signal, 10)
	sync := make(chan string)
	errch = make(chan error, 10)
	go BeatHeart(sync, quit)
	sync <- "start"
	SendIdentifiy(bot_token, props, intent)
	go route(sync, quit)
	go loggingError()
	util.RegularRequest(func() {
		service.UpdateGames()
		Notify()
	})
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Bye!")
}

func loggingError() {
	for {
		err := <-errch
		log.Fatal(err)
		err = nil
	}
}

func Beat(ws *websocket.Conn, s int) error {
	fmt.Println("Start Beat!")
	var pl Payload
	pl.Op = 1
	if s != 0 {
		pl.D = s
	}
	err := SendJson(ws, &pl)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func BeatHeart(sync <-chan string, quit <-chan os.Signal) {
	for {
		select {
		case <-quit:
			fmt.Println("Stop Beat...")
			return
		default:
		}
		<-sync
		time.Sleep(time.Millisecond * time.Duration(interval))
		err := Beat(ws, last_sequence)
		if err != nil {
			log.Fatal(err)
			fmt.Println("Connection restart from resume_url")
			// TODO Resume
			fmt.Println("Connection restarted is completed")
		}
	}
}

func SendIdentifiy(bot_token string, conn_props *map[string]string, intents int) {
	identify := new(Identify)
	identify.Token = bot_token
	identify.Intents = intents
	identify.Properties = conn_props
	var pl Payload
	pl.Op = 2
	pl.D = identify
	printFormattedJson(&pl)
	SendJson(ws, &pl)
	var data map[string]interface{}
	err := ReceiveJson(ws, &data)
	if err != nil {
		errch <- fmt.Errorf("failed to receive ready event")
	}
	if data["t"] != "READY" {
		errch <- fmt.Errorf("failed to receive ready event")
	}

	last_sequence = int(data["s"].(float64))
	d := data["d"].(map[string]interface{})

	session_id = d["session_id"].(string)
	resume_url = d["resume_gateway_url"].(string)
	fmt.Printf("session_id: %v resume_url:%v\n", session_id, resume_url)
}

func printFormattedJson(pl *Payload) {
	bytes, err := json.Marshal(pl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func route(sync chan<- string, quit <-chan os.Signal) {
	for {
		select {
		case <-quit:
			fmt.Println("route function is done")
		default:
		}
		var data map[string]interface{}
		err := ReceiveJson(ws, &data)
		if err != nil {
			errch <- err
			err = nil
			continue
		}

		if len(data) == 0 {
			errch <- err
			err = nil
			continue
		}
		op := int(data["op"].(float64))

		switch op {
		case 0:
			event_name := data["t"]
			d := data["d"].(map[string]interface{})
			switch event_name {
			case "MESSAGE_CREATE":
				var event MessageCreateEvent
				event = d
				err = EventHandler.Handle(&event)
			case "INTERACTION_CREATE":
				var event InteractionEvent
				event = d
				err = EventHandler.Handle(&event)
			}
			last_sequence = int(data["s"].(float64))
		case 11:
			sync <- "OK"
			fmt.Println("HeartBeat ACK Event")
		}
		if err != nil {
			errch <- err
			err = nil
		}
	}
}

func Notify() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("failed to send messages")
		}
	}()
	d := *service.GetChannelsWithDiscountedGames()
	for _, v := range d {
		SendGamesList(v.ChannelId, &v.Games)
	}
}

func SendGamesList(channel_id string, games *[]dto.GetAppDTO) error {
	var content strings.Builder
	for _, v := range *games {
		content.WriteString(v.String())
	}

	url := fmt.Sprintf("https://discord.com/api/v%v/channels/%v/messages", 10, channel_id)
	fmt.Printf(`{"content": "%v","tts": false}`, content.String())
	request, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(`{"content": "%v","tts": false}`, content.String())))
	if err != nil {
		return fmt.Errorf("http new request error")
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bot %v", bot_token))
	request.Header.Add("Content-Type", "application/json")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("send messages error")
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return fmt.Errorf("send messages error")
	}

	return nil
}
