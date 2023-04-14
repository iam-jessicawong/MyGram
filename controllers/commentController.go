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

// GetAllComment godoc
// @Summary Get details
// @Description Get details of all comment
// @Tags comment
// @Accept json
// @Produce json
// @Success 200 {object} models.Comment
// @Param Authorization header string true "Type Bearer your_token"
// @Router /comment [get]
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

	for _, comment := range comment {
		comment.User.Password = ""
	}
	c.JSON(http.StatusOK, comment)
}

// GetOneComment godoc
// @Summary Get details for a given id
// @Description Get details of comment corresponding to the input id
// @Tags comment
// @Accept json
// @Produce json
// @Param id path int true "ID of the comment"
// @Success 200 {object} models.Comment
// @Param Authorization header string true "Type Bearer your_token"
// @Router /comment/{id} [get]
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

type CommentInput struct {
	Message string `json:"message" form:"message"`
}

// CreateComment godoc
// @Summary Post new comment
// @Description Post details of new comment corresponding to the input
// @Tags comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param photoId path int true "id of the photo to comment"
// @Param models.Comment body CommentInput true "create comment"
// @Success 201 {object} models.Comment
// @Router /comment/{photoId} [post]
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

// UpdateComment godoc
// @Summary Update comment for a given id
// @Description Update the comment corresponding to the input comment id
// @Tags comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the comment to be updated"
// @Param models.Comment body CommentInput true "update comment"
// @Success 201 {object} models.Comment
// @Router /comment/{id} [put]
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

// DeleteComment godoc
// @Summary Delete comment for a given id
// @Description Update the comment corresponding to the input comment id
// @Tags comment
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the comment to be delete"
// @Success 200 "deleted"
// @Router /comment/{id} [delete]
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
