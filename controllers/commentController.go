package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/iam-jessicawong/mygram/helpers"
	"github.com/iam-jessicawong/mygram/models"
	"github.com/iam-jessicawong/mygram/repositories"
	"gorm.io/gorm"
)

func GetAllComment(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	comment, err := repositories.FindAllComment(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting comment data",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func GetOneComment(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("id"))
	comment, err := repositories.FindByIdComment(uint(commentID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
				"err":     "not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting comment",
			"err":     err.Error(),
		})
		return
	}

	comment.User.Password = ""
	c.JSON(http.StatusOK, &comment)
}

func CreateComment(c *gin.Context) {
	photoID, errConvert := strconv.Atoi(c.Param("photoId"))
	if errConvert != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	_, err := repositories.FindByIdPhoto(uint(photoID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error":   "data not found",
			"message": "photo is not exist",
		})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.PhotoID = uint(photoID)
	Comment.Message = strings.TrimSpace(Comment.Message)

	err = repositories.CreateComment(&Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &Comment)
}

func UpdateComment(c *gin.Context) {
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}

	commentID, _ := strconv.Atoi(c.Param("id"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.ID = uint(commentID)
	Comment.Message = strings.TrimSpace(Comment.Message)

	err := repositories.UpdateComment(&Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your comment successfully updated",
		"data":    Comment,
	})
}

func DeleteComment(c *gin.Context) {
	commentID, _ := strconv.Atoi(c.Param("id"))

	err := repositories.DeleteComment(uint(commentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Can't delete comment",
		})
		return
	}

	c.JSON(http.StatusOK, "Comment successfully deleted")
}
