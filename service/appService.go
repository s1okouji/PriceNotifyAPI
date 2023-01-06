package service

import (
	"github.com/s1okouji/pnabot_client/dto"
	"github.com/s1okouji/pnabot_client/entity"
	"github.com/s1okouji/pnabot_client/mapper"
	"github.com/s1okouji/pnabot_client/util"
)

// ChannelとGameを結びつける
// Gameが存在しないとき、新たに追加する
func AddGame(dto *dto.CreateAppDTO) error {
	err := mapper.CreateGame(dto.AppId, dto.ChannelId, DB)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGame(dto dto.DeleteAppDTO) {
	mapper.DeleteGame(dto.AppId, dto.ChannelId, DB)
}

func GetGames() *[]dto.GetAppDTO {
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
	return &apps
}

func UpdateGames() error {
	gamesEntity := mapper.GetGames(DB)
	var newGamesEntity []entity.Game
	for _, v := range gamesEntity {
		newGamesEntity = append(newGamesEntity, *util.GetGameEntity(v.AppId))
	}

	mapper.UpdateGames(&newGamesEntity, DB)

	return nil
}
