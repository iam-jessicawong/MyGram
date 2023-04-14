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

// GetAllSocialMedia godoc
// @Summary Get details
// @Description Get details of all social media
// @Tags social-media
// @Accept json
// @Produce json
// @Success 200 {object} models.SocialMedia
// @Param Authorization header string true "Type Bearer your_token"
// @Router /social-media [get]
func GetAllSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	social_media, err := repositories.FindAllSocialMedia(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting social media data",
			"err":     err.Error(),
		})
		return
	}

	for _, social_media := range social_media {
		social_media.User.Password = ""
	}
	c.JSON(http.StatusOK, &social_media)
}

// GetOneSocialMedia godoc
// @Summary Get details for a given id
// @Description Get details of social media corresponding to the input id
// @Tags social-media
// @Accept json
// @Produce json
// @Param id path int true "ID of the social media"
// @Success 200 {object} models.SocialMedia
// @Param Authorization header string true "Type Bearer your_token"
// @Router /social-media/{id} [get]
func GetOneSocialMedia(c *gin.Context) {
	socialmediaID, _ := strconv.Atoi(c.Param("id"))
	social_media, err := repositories.FindByIdSocialMedia(uint(socialmediaID))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "social media not found",
				"err":     "not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error getting social media",
			"err":     err.Error(),
		})
		return
	}

	social_media.User.Password = ""
	c.JSON(http.StatusOK, &social_media)
}

type SocialMediaInput struct {
	Name           string `json:"name" form:"name"`
	SocialMediaURL string `json:"social_media_url" form:"social_media_url"`
}

// CreateSocialMedia godoc
// @Summary Post new social media
// @Description Post details of new social media corresponding to the input
// @Tags social-media
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param models.SocialMedia body SocialMediaInput true "create social media"
// @Success 201 {object} models.SocialMedia
// @Router /social-media [post]
func CreateSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID
	SocialMedia.Name = strings.TrimSpace(SocialMedia.Name)
	SocialMedia.SocialMediaURL = strings.TrimSpace(SocialMedia.SocialMediaURL)

	err := repositories.CreateSocialMedia(&SocialMedia)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "your social media successfully added",
		"data":    &SocialMedia,
	})
}

// UpdateSocialMedia godoc
// @Summary Update social media for a given id
// @Description Update the social media corresponding to the input social media id
// @Tags social-media
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the social media to be updated"
// @Param models.SocialMedia body SocialMediaInput true "update social media"
// @Success 201 {object} models.SocialMedia
// @Router /social-media/{id} [put]
func UpdateSocialMedia(c *gin.Context) {
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMediainput := models.SocialMedia{}

	socialMediaID, _ := strconv.Atoi(c.Param("id"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMediainput)
	} else {
		c.ShouldBind(&SocialMediainput)
	}

	SocialMediainput.UserID = userID
	SocialMediainput.ID = uint(socialMediaID)
	SocialMediainput.Name = strings.TrimSpace(SocialMediainput.Name)
	SocialMediainput.SocialMediaURL = strings.TrimSpace(SocialMediainput.SocialMediaURL)

	// validate url, if user update the socialMedia url
	_, err := url.ParseRequestURI(SocialMediainput.SocialMediaURL)
	if err != nil && SocialMediainput.SocialMediaURL != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": "invalid url",
		})
		return
	}

	updatedSocialMedia, err := repositories.UpdateSocialMedia(&SocialMediainput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your social media is successfully updated",
		"data":    &updatedSocialMedia,
	})
}

// DeleteSocialMedia godoc
// @Summary Delete social media for a given id
// @Description Update the social media corresponding to the input social media id
// @Tags social-media
// @Accept json
// @Produce json
// @Param Authorization header string true "Type Bearer your_token"
// @Param id path int true "ID of the social media to be delete"
// @Success 200 "deleted"
// @Router /social-media/{id} [delete]
func DeleteSocialMedia(c *gin.Context) {
	socialmediaID, _ := strconv.Atoi(c.Param("id"))

	err := repositories.DeleteSocialMedia(uint(socialmediaID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Can't delete social media",
		})
		return
	}

	c.JSON(http.StatusOK, "Social media successfully deleted")
}
