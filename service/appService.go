package service

import (
	"github.com/s1okouji/price_notify_api/dto"
	"github.com/s1okouji/price_notify_api/mapper"
)

// ChannelとGameを結びつける
// Gameが存在しないとき、新たに追加する
func AddGame(dto dto.CreateAppDTO) {
	mapper.CreateGame(dto.AppId, dto.ChannelId, DB)
}

func DeleteGame(dto dto.DeleteAppDTO) {
	mapper.DeleteGame(dto.AppId, dto.ChannelId, DB)
}

func GetGames() []dto.GetAppDTO {
	var apps []dto.GetAppDTO
	gamesEntity := mapper.GetGames(DB)
	for _, v := range gamesEntity {
		apps = append(apps, dto.GetAppDTO{
			AppId:           v.AppId,
			AppName:         v.AppName,
			BasePrice:       v.BasePrice,
			FinalPrice:      v.FinalPrice,
			DiscountPercent: v.DiscountPercent,
		})
	}
	return apps
}
