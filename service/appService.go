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
