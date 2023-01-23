package util

import (
	"time"

	"github.com/s1okouji/pnabot_client/dto"
	"github.com/s1okouji/pnabot_client/entity"
)

func GetData(str string) string {
	state := 0
	start := 0
	end := len(str) - 1
	for i := 0; i < len(str); i++ {
		if str[i] == '{' {
			if state == 1 {
				start = i
				break
			}
			state++
		} else if str[i] == '}' {
			break
		}
	}
	return str[start:end]
}

func RegularRequest(f func()) {
	for {
		y, m, d := time.Now().AddDate(0, 0, 1).Date()
		date := time.Date(y, m, d, 10, 0, 0, 0, time.UTC)
		time.Sleep(time.Until(date))
		f()
	}
}

func Convert(games *[]entity.Game) *map[int]entity.Game {
	mp := map[int]entity.Game{}
	for _, v := range *games {
		mp[v.AppId] = v
	}
	return &mp
}

func Mapping(game *entity.Game) *dto.GetAppDTO {
	// 型スイッチで対応させる型を増やすかもしれない
	apps := dto.GetAppDTO{}
	apps.AppId = game.AppId
	apps.AppName = game.AppName
	apps.BasePrice = game.BasePrice
	apps.DiscountPercent = game.DiscountPercent
	apps.FinalPrice = game.FinalPrice
	return &apps
}
