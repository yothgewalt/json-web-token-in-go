package datastore

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

type DatabaseConfig struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Port     string `yaml:"port"`
		Sslmode  string `yaml:"ssl_mode"`
		Timezone string `yaml:"timezone"`
	}
}

func Connect(pathYaml string, gormConfig *gorm.Config) (db *gorm.DB, err error) {
	yamlBytes, err := os.ReadFile(pathYaml)
	if err != nil {
		return nil, err
	}

	configure := &DatabaseConfig{}
	err = yaml.Unmarshal(yamlBytes, configure)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", configure.Database.Host, configure.Database.User, configure.Database.Password, configure.Database.Dbname, configure.Database.Port, configure.Database.Sslmode, configure.Database.Timezone)
	db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
