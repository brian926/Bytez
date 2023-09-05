package controllers

import (
	"net/http"
	"net/url"

	"github.com/brian926/Bytez/server/forms"
	"github.com/brian926/Bytez/server/shortener"
	"github.com/brian926/Bytez/server/store"
	"github.com/gin-gonic/gin"
)

type UrlController struct{}

var storeModel = new(store.UrlModel)

func (crtl UrlController) CreateShortUrl(c *gin.Context) {
	var creationRequest forms.UrlCreationRequest

	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := url.ParseRequestURI(creationRequest.LongUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Please use links starting with https:// or http://"})
		return
	}

	creationRequest.LongUrl = u.String()

	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	creationRequest.ShortUrl = shortUrl

	save, err := storeModel.SaveUrlMapping(creationRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	host := "https://bytez.us/api/"

	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + save.ShortUrl,
	})
}

func (crtl UrlController) HandleShortUrlRedirect(c *gin.Context) {
	var creationRequest forms.UrlCreationRequest
	creationRequest.ShortUrl = c.Param("shortUrl")
	initialUrl, err := storeModel.RetrieveInitialUrl(creationRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(302, initialUrl.LongUrl)
}

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Pong",
	})
}
