package controllers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/iam-jessicawong/mygram/helpers"
	"github.com/iam-jessicawong/mygram/models"
	"github.com/iam-jessicawong/mygram/repositories"
	"gorm.io/gorm"
)

// GetAllPhoto godoc
// @Summary Get details
// @Description Get details of all photo
// @Tags photo
// @Accept json
// @Produce json
// @Success 200 {object} models.Photo
// @Param Authorization header string true "Type Bearer your_token"
// @Router /photo [get]
func GetAllPhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	photo, err := repositories.FindAllPhoto(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting photo data",
			"err":     err.Error(),
		})
		return
	}

	for _, photo := range photo {
		photo.User.Password = ""
	}
	c.JSON(http.StatusOK, photo)
}

// GetOnePhoto godoc
// @Summary Get details for a given id
// @Description Get details of photo corresponding to the input id
// @Tags photo
// @Accept json
// @Produce json
// @Param id path int true "ID of the photo"
// @Success 200 {object} models.Photo
// @Param Authorization header string true "Type Bearer your_token"
// @Router /photo/{id} [get]
func GetOnePhoto(c *gin.Context) {
	photoID, _ := strconv.Atoi(c.Param("id"))
	photo, err := repositories.FindByIdPhoto(uint(photoID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "photo not found",
				"err":     "not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting photo",
			"err":     err.Error(),
		})
		return
	}

	photo.User.Password = ""
	c.JSON(http.StatusOK, &photo)
}

// this struct is for swagger custom body
type InputPhoto struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
}

// CreatePhoto godoc
// @Summary Post new photo
// @Description Post details of new photo corresponding to the input
// @Tags photo
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param models.Photo body InputPhoto true "create photo"
// @Success 201 {object} models.Photo
// @Router /photo [post]
func CreatePhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = userID
	Photo.Title = strings.TrimSpace(Photo.Title)
	Photo.PhotoURL = strings.TrimSpace(Photo.PhotoURL)
	Photo.Caption = strings.TrimSpace(Photo.Caption)

	err := repositories.CreatePhoto(&Photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &Photo)
}

// UpdatePhoto godoc
// @Summary Update photo for a given id
// @Description Update the photo corresponding to the input photo id
// @Tags photo
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the photo to be updated"
// @Param models.Photo body InputPhoto true "update photo"
// @Success 201 {object} models.Photo
// @Router /photo/{id} [put]
func UpdatePhoto(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	PhotoInput := models.Photo{}

	photoID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&PhotoInput)
	} else {
		c.ShouldBind(&PhotoInput)
	}

	PhotoInput.UserID = userID
	PhotoInput.ID = uint(photoID)
	PhotoInput.Title = strings.TrimSpace(PhotoInput.Title)
	PhotoInput.PhotoURL = strings.TrimSpace(PhotoInput.PhotoURL)
	PhotoInput.Caption = strings.TrimSpace(PhotoInput.Caption)

	// validate url, if user update the photo url
	_, err := url.ParseRequestURI(PhotoInput.PhotoURL)
	if err != nil && PhotoInput.PhotoURL != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": "invalid url",
		})
		return
	}

	updatedPhoto, err := repositories.UpdatePhoto(&PhotoInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &updatedPhoto)
}

// DeletePhoto godoc
// @Summary Delete photo for a given id
// @Description Update the photo corresponding to the input photo id
// @Tags photo
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the photo to be delete"
// @Success 200 "deleted"
// @Router /photo/{id} [delete]
func DeletePhoto(c *gin.Context) {
	photoID, _ := strconv.Atoi(c.Param("id"))

	err := repositories.DeletePhoto(uint(photoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Can't delete photo",
		})
		return
	}

	c.JSON(http.StatusOK, "Photo successfully deleted")
}
