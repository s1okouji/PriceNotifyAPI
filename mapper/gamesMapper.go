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

func GetGamesWithChannel(db *gorm.DB, channel_id string) []entity.Game {
	var gamesEntity []entity.Game
	var gamesModel []model.Game

	// FIX: gormを使用すると、どこにwhere句をおいても最後に処理されるためサブクエリを用いて処理順を指定する
	result := db.Table("games").Joins("left join channel_games on channel_games.app_id = games.app_id").Where("channel_games.channel_id = ?", channel_id).Find(&gamesModel)
	fmt.Printf("%+v", gamesModel)
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

func GetDiscountedGames(db *gorm.DB) *[]entity.Game {
	var gamesEntity []entity.Game
	var gamesModel []model.Game
	result := db.Find(&gamesModel)
	if result.Error != nil {
		panic("Database error on GetGames")
	}

	for _, v := range gamesModel {
		if v.DiscountPercent == 0 {
			continue
		}
		gamesEntity = append(gamesEntity, entity.Game{
			AppId:           v.AppId,
			AppName:         v.AppName,
			BasePrice:       v.BasePrice,
			FinalPrice:      v.FinalPrice,
			DiscountPercent: v.DiscountPercent,
			History:         nil,
		})
	}

	return &gamesEntity
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

func GetChannels(db *gorm.DB) *[]entity.Channel {
	// TypeとUsersは現在使用する予定がないので中身は入ってない
	var channel_gameModels []model.Channel_Game
	var channel_entities []entity.Channel
	mp := map[string]*[]int{}
	result := db.Find(&channel_gameModels)
	if result.Error != nil {
		panic("Database error on Channel_Games")
	}
	for _, v := range channel_gameModels {
		var appids []int
		if val, ok := mp[v.ChannelId]; ok {
			appids = *val
		}
		appids = append(appids, v.AppId)
		mp[v.ChannelId] = &appids
	}

	for k, v := range mp {
		channel_entities = append(channel_entities, entity.Channel{
			ChannelId: k,
			AppIds:    *v,
		})
	}

	return &channel_entities
}

// TODO: 一回でchannelとgameのmodelを取り出せるようにする。
// 多分新しくmodelを作ればできそう
func GetChannelsHaveDiscountedGames(db *gorm.DB) *[]entity.Channel {
	var channel_entities []entity.Channel
	var channelGameModels []model.Channel_Game
	mp := map[string]*[]int{}
	// FIX: gormを使用すると、どこにwhere句をおいても最後に処理されるためサブクエリを用いて処理順を指定する
	result := db.Table("games").Joins("left join channel_games on channel_games.app_id = games.app_id").Where("games.discount_percent > 0").Find(&channelGameModels)
	if result.Error != nil {
		panic("Database error on GetGames")
	}

	for _, v := range channelGameModels {
		var appids []int
		if val, ok := mp[v.ChannelId]; ok {
			appids = *val
		}
		appids = append(appids, v.AppId)
		mp[v.ChannelId] = &appids
	}

	for k, v := range mp {
		channel_entities = append(channel_entities, entity.Channel{
			ChannelId: k,
			AppIds:    *v,
		})
	}

	return &channel_entities
}
