package main

import (
	"WebTokenAuthorization/domain/model"
	"WebTokenAuthorization/infrastructure/datastore"
	"WebTokenAuthorization/infrastructure/interactor"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func main() {
	var err error

	datastore.DB, err = datastore.Connect("configure.yaml", &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Panicf("(error) failed to connect the database cause %v", err)
	}

	err = datastore.DB.AutoMigrate(&model.Collections{})
	if err != nil {
		log.Panicf("(error) failed to auto migrate cause %v", err)
	}

	err = NewRouter(":3000")
	if err != nil {
		log.Panicf("(error) failed to start the server cause %v", err)
	}
}

func NewRouter(address string) (err error) {
	r := gin.Default()

	if err = r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return err
	}

	api := r.Group("/api", func(c *gin.Context) { c.Next() })
	v1 := api.Group("/v1", func(c *gin.Context) { c.Next() })

	v1.GET("/receiver", interactor.Receiver)

	v1.POST("/create/collection", interactor.CreateCollection)
	v1.POST("/access/collection", interactor.AccessCollection)

	v1.GET("/get/collection/:id", interactor.GetCollectionById)
	v1.GET("/get/collections", interactor.GetAllCollection)

	v1.DELETE("/delete/soft/collection/:id", interactor.SoftDeleteCollectionById)
	v1.DELETE("/delete/hard/collection/:id", interactor.HardDeleteCollectionById)

	if err = r.Run(address); err != nil {
		return err
	}

	return nil
}
