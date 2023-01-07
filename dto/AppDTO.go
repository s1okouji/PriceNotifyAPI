package dto

import (
	"fmt"
)

type CreateAppDTO struct {
	AppId     int
	ChannelId string
}

type DeleteAppDTO struct {
	AppId     int
	ChannelId string
}

type GetAppDTO struct {
	AppId           int    `json:"id"`
	AppName         string `json:"name"`
	BasePrice       int    `json:"base_price"`
	FinalPrice      int    `json:"final_price"`
	DiscountPercent int    `json:"discount_percent"`
}

type NotifyDTO struct {
	ChannelId string
	Games     []GetAppDTO
}

func (dto *GetAppDTO) String() string {
	content := fmt.Sprintf("ゲーム名: %v\n定価: %v 円\n現在の価格: %v 円\n割引率: %v%%\n", dto.AppName, dto.BasePrice/100, dto.FinalPrice/100, dto.DiscountPercent)
	return content
}
