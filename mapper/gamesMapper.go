package mapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/s1okouji/pnabot_client/entity"
	"github.com/s1okouji/pnabot_client/model"
	"github.com/s1okouji/pnabot_client/util"
	"gorm.io/gorm"
)

const (
	FORMAT = "20060102"
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
func CreateGame(appId int, channelId string, db *gorm.DB) error {
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
				Day:             time.Now().Format(FORMAT),
			})
		} else {
			return fmt.Errorf("database error on CreateGame")
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
	return nil
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

func UpdateGames(gamesEntity *[]entity.Game, db *gorm.DB) error {
	var gamesModel []model.Game
	var newGamesModel []model.Game
	var gamesHistory []model.History
	result := db.Find(&gamesModel)
	if result.Error != nil {
		return fmt.Errorf("database error on GetGames")
	}

	for _, v := range gamesModel {
		gamesHistory = append(gamesHistory, model.History{
			AppId:           v.AppId,
			Day:             v.CreatedAt.Format(FORMAT),
			FinalPrice:      v.FinalPrice,
			DiscountPercent: v.DiscountPercent,
		})
	}

	for _, v := range *gamesEntity {
		newGamesModel = append(newGamesModel, model.Game{
			AppId:           v.AppId,
			AppName:         v.AppName,
			BasePrice:       v.BasePrice,
			FinalPrice:      v.FinalPrice,
			DiscountPercent: v.DiscountPercent,
			CreatedAt:       time.Now(),
		})
	}

	err := db.Create(&gamesHistory)
	if err != nil {
		return fmt.Errorf("failed to create history columns")
	}
	err = db.Create(&newGamesModel)
	if err != nil {
		return fmt.Errorf("failed to create new games model columns")
	}
	return nil
}
