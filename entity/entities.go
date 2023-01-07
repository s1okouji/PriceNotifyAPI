package entity

type Game struct {
	AppId           int
	AppName         string
	BasePrice       int
	FinalPrice      int
	DiscountPercent int
	History         []History
}

type History struct {
	AppId           int
	Day             int
	FinalPrice      int
	DiscountPercent int
}

type Channel struct {
	ChannelId string
	Type      int
	Users     []string
	AppIds    []int
}

type Admin struct {
	Id string
}
