package dto

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
