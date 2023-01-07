package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/s1okouji/pnabot_client/dto"
	"github.com/s1okouji/pnabot_client/service"
)

type InteractionEvent map[string]interface{}
type InteractionResponse struct {
	Type int                      `json:"type"`
	Data *InteractionCallBackData `json:"data"`
}

type InteractionCallBackData struct {
	Tts     bool   `json:"tts"`
	Content string `json:"content"`
}

const (
	PING                             = 1
	APPLICATION_COMMAND              = 2
	MESSAGE_COMPONENT                = 3
	APPLICATION_COMMAND_AUTOCOMPLETE = 4
	MODAL_SUBMIT                     = 5
)

func (event *InteractionEvent) Handle() error {
	data := *event
	switch int(data["type"].(float64)) {
	case PING:
	case APPLICATION_COMMAND, APPLICATION_COMMAND_AUTOCOMPLETE:
		mp := data["data"].(map[string]interface{})
		switch mp["name"].(string) {
		case "add":
			var dto dto.CreateAppDTO
			dto.AppId = int(mp["options"].([]interface{})[0].(map[string]interface{})["value"].(float64))
			dto.ChannelId = data["channel_id"].(string)
			err := createApp(&dto)
			if err != nil {
				return err
			}
			respondWithMessage(data["id"].(string), data["token"].(string), 4, "追加されました!")
		case "list":
			// sendGamesList(data["channel_id"].(string))
			games := *service.GetGamesWithChannel(data["channel_id"].(string))
			var content strings.Builder
			for _, v := range games {
				content.WriteString(v.String())
			}
			message := content.String()
			if message == "" {
				message = "empty .."
			}
			respondWithMessage(data["id"].(string), data["token"].(string), 4, message)
		}
	}
	return nil
}

func createApp(dto *dto.CreateAppDTO) error {
	err := service.AddGame(dto)
	if err != nil {
		return fmt.Errorf("addGame error")
	}
	return nil
}

func respondWithMessage(interaction_id string, interaction_token string, t int, message string) error {
	url := fmt.Sprintf("https://discord.com/api/v%v/interactions/%v/%v/callback", 10, interaction_id, interaction_token)
	var res InteractionResponse
	var data InteractionCallBackData
	data.Content = message
	data.Tts = false
	res.Type = t
	res.Data = &data
	req_body, _ := json.Marshal(res)
	request, err := http.NewRequest("POST", url, bytes.NewReader(req_body))
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
