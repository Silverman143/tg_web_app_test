package back

type User struct {
	Telegram_id int64 `json:"telegram_id" db:"telegram_id"`
	Name string `json:"user_name" db:"user_name" binding:"required"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	Stars int `json:"stars" db:"stars_balance"`
	GlobalRank int `json:"rank" db:"rank"`
}