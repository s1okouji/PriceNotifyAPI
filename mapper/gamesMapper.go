package mapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/s1okouji/price_notify_api/entity"
	"github.com/s1okouji/price_notify_api/model"
	"github.com/s1okouji/price_notify_api/util"
	"gorm.io/gorm"
)

func GetGames(db *gorm.DB) []entity.Game {
	var gamesEntity []entity.Game
	var gamesModel []model.Game
	result := db.Find(&gamesModel)
	if result.Error != nil {
		panic("Database error on GetGames")
	}

	for _, v := range gamesModel {
		gamesEntity = append(gamesEntity, entity.Game{
			AppId:           v.AppId,
			AppName:         v.AppName,
			BasePrice:       v.BasePrice,
			FinalPrice:      v.FinalPrice,
			DiscountPercent: v.DiscountPercent,
			History:         nil,
		})
	}

	return gamesEntity
}

// Entityからデータベースへアクセスする
func CreateGame(appId int, channelId string, db *gorm.DB) {
	var gameModel model.Game
	err := db.First(&gameModel, appId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// gameをテーブルに追加
			gameModel = util.GetGameModel(appId)
			db.Create(&gameModel)
			db.Create(&model.History{
				AppId:           appId,
				FinalPrice:      gameModel.FinalPrice,
				DiscountPercent: gameModel.DiscountPercent,
				Day:             time.Now().Format("20060102"),
			})
		} else {
			panic("Database error on CreateGame")
		}
	}

	var channel_gameModel model.Channel_Game
	err = db.Where("app_id=? AND channel_id=?", appId, channelId).First(&channel_gameModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Create(&model.Channel_Game{
				AppId:     appId,
				ChannelId: channelId,
			})
		}
	} else {
		fmt.Println("Already Exists!")
	}
}

func DeleteGame(appId int, channelId string, db *gorm.DB) {
	var channel_gameModel model.Channel_Game
	err := db.Where("app_id=? AND channel_id=?", appId, channelId).First(&channel_gameModel).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic("Database error on DeleteGame")
		} else {
			return
		}
	}

	db.Delete(&model.Channel_Game{
		Id:        channel_gameModel.Id,
		AppId:     appId,
		ChannelId: channelId,
	})
}
