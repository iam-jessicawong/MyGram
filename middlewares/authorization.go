package middlewares

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/iam-jessicawong/mygram/repositories"
)

var ErrUnauthorized error //for custom error

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "PUT" {
			ErrUnauthorized = errors.New("you are not allowed to edit this product")
		} else if c.Request.Method == "DELETE" {
			ErrUnauthorized = errors.New("you are not allowed to delete this product")
		}

		id, errConvert := strconv.Atoi(c.Param("id"))
		if errConvert != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid id",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		var userIDData uint //to hold the user_id obtained from database
		var err error
		// check which endpoint is being hit to get user id from that specific database
		if strings.Contains(c.Request.URL.Path, "photo") {
			photo, dberr := repositories.FindPhotoUser(uint(id))
			err = dberr
			userIDData = photo.UserID
		} else if strings.Contains(c.Request.URL.Path, "comment") {
			comment, dberr := repositories.FindCommentUser(uint(id))
			err = dberr
			userIDData = comment.UserID
		} else if strings.Contains(c.Request.URL.Path, "social-media") {
			sm, dberr := repositories.FindSocialMediaUser(uint(id))
			err = dberr
			userIDData = sm.UserID
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": err.Error(),
			})
			return
		}

		if userIDData != userID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": ErrUnauthorized.Error(),
			})
			return
		}

		c.Next()
	}
}
