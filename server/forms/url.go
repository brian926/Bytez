package forms

type UrlCreationRequest struct {
	ShortUrl string
	LongUrl  string `json:"long_url" binding:"required"`
	UserId   string `json:"user_id" binding:"required"`
}
