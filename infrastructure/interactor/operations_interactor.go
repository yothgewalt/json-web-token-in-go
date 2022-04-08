package interactor

import (
	"WebTokenAuthorization/domain/model"
	"WebTokenAuthorization/infrastructure/datastore"
	"WebTokenAuthorization/infrastructure/presenter"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetCollectionById(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, but the header does not contain cookies for the token."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSignature, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the signature is invalid."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the token is invalid."})
		return
	}

	var collection model.Collections

	id := c.Param("id")
	collectionID, _ := strconv.Atoi(id)

	tx := datastore.DB.Table("collections").Model(&model.Collections{}).Select("id", "username", "firstname", "lastname", "email_address", "phone_number", "facebook_link").Where("id = ?", collectionID).Take(&collection)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the collection could not be found."})
			return
		}
	}

	c.JSON(http.StatusFound, presenter.ResponseCollection{
		ID:           collection.ID,
		Username:     collection.Username,
		Firstname:    collection.Firstname,
		Lastname:     collection.Lastname,
		PhoneNumber:  collection.PhoneNumber,
		EmailAddress: collection.EmailAddress,
		FacebookLink: collection.FacebookLink,
	})
}

func GetAllCollection(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, but the header does not contain cookies for the token."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSignature, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the signature is invalid."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the token is invalid."})
		return
	}

	var collections []model.Collections

	tx := datastore.DB.Table("collections").Model(&model.Collections{}).Select("id", "username", "firstname", "lastname", "email_address", "phone_number", "facebook_link").Find(&collections)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the collection could not be found."})
			return
		}
	}

	var responses []presenter.ResponseCollection
	for _, field := range collections {
		response := presenter.ResponseCollection{
			ID:           field.ID,
			Username:     field.Username,
			Firstname:    field.Firstname,
			Lastname:     field.Lastname,
			PhoneNumber:  field.PhoneNumber,
			EmailAddress: field.EmailAddress,
			FacebookLink: field.FacebookLink,
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusFound, responses)
}

func SoftDeleteCollectionById(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, but the header does not contain cookies for the token."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSignature, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the signature is invalid."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the token is invalid."})
		return
	}

	var collection model.Collections

	id := c.Param("id")
	collectionID, _ := strconv.Atoi(id)

	tx := datastore.DB.Table("collections").Where("id = ?", collectionID).Delete(&collection).Error
	if tx != nil {
		if errors.Is(tx, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, presenter.ResponseOnlyMessage{Message: "Sorry, the collection could not be found."})
			return
		}
	}

	c.JSON(http.StatusOK, presenter.ResponseOnlyMessage{Message: "successfully, the collection has been deleted (soft)."})
}

func HardDeleteCollectionById(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, but the header does not contain cookies for the token."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSignature, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the signature is invalid."})
			return
		}

		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, the token is invalid."})
		return
	}
	
	var collection model.Collections

	id := c.Param("id")
	collectionID, _ := strconv.Atoi(id)

	tx := datastore.DB.Table("collections").Unscoped().Where("id = ?", collectionID).Delete(&collection).Error
	if tx != nil {
		if errors.Is(tx, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, presenter.ResponseOnlyMessage{Message: "Sorry, the collection could not be found."})
			return
		}
	}

	c.JSON(http.StatusOK, presenter.ResponseOnlyMessage{Message: "successfully, the collection has been deleted (hard)."})
}
