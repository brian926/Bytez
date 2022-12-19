package controllers

import (
	"fmt"
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
		fmt.Println("error from save")
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	host := "http://localhost:9808/"

	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + save.ShortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shortUrl)
	c.Redirect(302, initialUrl)
}
