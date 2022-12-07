package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/brian926/UrlShorterGo/server/packages/config"
	"github.com/brian926/UrlShorterGo/server/packages/db"
	"github.com/brian926/UrlShorterGo/server/packages/utils"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	db.User
	jwt.StantardClaims
}

func Pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hey Weclome to the URL Shortener API",
	})
}

func CreateUser(c *gin.Context) {
	user := &db.User{}

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errs := utils.ValidateUser(*user); len(errs) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		return
	}

	if user.UserExists() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	user.HashPassword()
	_, err := db.Exec(db.CreateUserQuery, user.Name, user.Password, user.Email)
	if err != nil {
		return err
	}

	c.JSON(200, gin.H{"message": "successfully created user"})
}

func Session(c *gin.Context) {
	tokenUser := c.BindJSON("user").(*jwt.Token)
	claims := tokenUser.Claims.(jwt.MapClaims)
	userID, ok := claims["id"].(string)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal error"})
		return
	}

	user := &db.User{}
	if err := db.Exec(db.GetUserByIDQuery, userID).
		Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
			return
		}
	}
	user.Password = ""
	c.JSON(200, gin.H{"message": "session created"})
}

func Login(c *gin.Context) {
	loginUser := &db.User{}

	if err := c.ShouldBindJSON(loginUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &db.User{}
	if err := db.Exec(db.GetUserByEmailQuery, loginUser.Email).
		Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	match := utils.ComparePassword(user.Password, loginUser.Password)
	if !match {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		return
	}

	//expiration time of the token ->30 mins
	expirationTime := time.Now().Add(30 * time.Minute)

	user.Password = ""
	claims := &Claims{
		*user,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtKey = []byte(config.Config[config.JWT_KEY])
	tokenValue, err := token.SignedString(jwtKey)

	if err != nil {
		return err
	}

	c.SetCookie(
		"token",
		tokenValue,
		expirationTime,
		"/session",
		config.Config[config.CLIENT_URL],
		true,
		true,
	)

	c.JSON(200, gin.H{"user": claims.User, "token": tokenValue})
	return
}

func Logout(c *gin.Context) {
	c.SetCookie("name", "", 1, "/", "localhost", false, true)
	c.String(200, "Deleted cookie")
}
