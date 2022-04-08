package interactor

import (
	"WebTokenAuthorization/domain/model"
	"WebTokenAuthorization/infrastructure/datastore"
	"WebTokenAuthorization/infrastructure/presenter"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/matthewhartstonge/argon2"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	GeneralEmailRegexRFC5322 = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
	GeneralUsernameRegex     = "([A-Za-z0-9_])\\w+"
	GeneralRealNameRegex     = "([A-Za-z])\\w+"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSignature = []byte("signed_signature")

func CreateCollection(c *gin.Context) {
	var collection model.Collections
	argon := argon2.DefaultConfig()

	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ResponseOnlyError{Error: err})
		return
	}

	username, _ := regexp.MatchString(GeneralUsernameRegex, collection.Username)
	if (username != false && len(collection.Username) > 0) && (len(collection.Username) >= 6 && len(collection.Username) <= 32) {
		collection.Username = strings.ToLower(collection.Username)
	} else {
		c.JSON(http.StatusUnprocessableEntity, presenter.ResponseOnlyMessage{Message: "Sorry, but the username Only [a-z, A-Z, 0-9, _] can be used and must not be less than 6 characters and more than 32 characters."})
		return
	}

	firstname, _ := regexp.MatchString(GeneralRealNameRegex, collection.Firstname)
	lastname, _ := regexp.MatchString(GeneralRealNameRegex, collection.Lastname)
	if (firstname != false && lastname != false) && (len(collection.Firstname) > 0 && len(collection.Lastname) > 0) {
		collection.Firstname = strings.ToLower(collection.Firstname)
		collection.Lastname = strings.ToLower(collection.Lastname)
	} else {
		c.JSON(http.StatusUnprocessableEntity, presenter.ResponseOnlyMessage{Message: "Sorry, but the firstname and lastname Only [a-z, A-Z] can be used and the fields of firstname, lastname cannot empty text."})
		return
	}

	if len(collection.Password) >= 6 {
		encoded, err := argon.HashEncoded([]byte(collection.Password))
		if err != nil {
			c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
			return
		}

		collection.Password = string(encoded)
	} else {
		c.JSON(http.StatusUnprocessableEntity, presenter.ResponseOnlyMessage{Message: "Sorry, but the password must not be less than 6 characters."})
		return
	}

	email, _ := regexp.MatchString(GeneralEmailRegexRFC5322, collection.EmailAddress)
	if email != false && len(collection.EmailAddress) > 0 {
		collection.EmailAddress = strings.ToLower(collection.EmailAddress)
	} else {
		c.JSON(http.StatusUnprocessableEntity, presenter.ResponseOnlyMessage{Message: "Sorry, but you must use a valid email address. (example@email.com)"})
		return
	}

	tx := datastore.DB.Table("collections").Select("username").Where("username = ?", collection.Username).Take(&collection.Username).Error
	if tx != nil {
		if errors.Is(tx, gorm.ErrRecordNotFound) {
			if tx := datastore.DB.Table("collections").Model(&model.Collections{}).Create(&collection).Error; tx != nil {
				if tx != nil {
					c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyMessage{Message: "Sorry, but you cannot request to create a collection."})
					return
				}
			}
		} else if !errors.Is(tx, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyMessage{Message: "Sorry, but you cannot create a collection with the json data because the json data has created in database."})
			return
		}
	}

	c.JSON(http.StatusCreated, presenter.ResponseOnlyMessage{Message: "Congratulation, the collection has been created."})
}

func AccessCollection(c *gin.Context) {
	var credential Credentials

	if err := c.ShouldBindJSON(&credential); err != nil {
		c.JSON(http.StatusBadRequest, presenter.ResponseOnlyError{Error: err})
		return
	}

	payload := &Credentials{
		Username: credential.Username,
		Password: credential.Password,
	}

	tx := datastore.DB.Table("collections").Select("username", "password").Where("username = ?", credential.Username).Take(&credential).Error
	if tx != nil {
		if errors.Is(tx, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, presenter.ResponseOnlyMessage{Message: "Sorry, but we not found that username."})
			return
		}
	}

	ok, err := argon2.VerifyEncoded([]byte(payload.Password), []byte(credential.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, presenter.ResponseOnlyMessage{Message: "Sorry, but your password is invalid."})
		return
	}

	expirationTime := jwt.NewNumericDate(time.Now().Add(5 * time.Minute))
	claims := &Claims{
		Username: credential.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString(jwtSignature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.ResponseOnlyError{Error: err})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenSigned, 300, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenSigned,
	})
}
