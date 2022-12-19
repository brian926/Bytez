package forms

type UrlCreationRequest struct {
	ShortUrl string
	LongUrl  string `form:"LongUrl" json:"long_url" binding:"required"`
	UserId   string `form:"UserId" json:"user_id" binding:"required"`
}
