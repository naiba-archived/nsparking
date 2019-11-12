package controller

import (
	"github.com/jinzhu/gorm"
	// ..
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/naiba/nsparking/model"
	"github.com/naiba/nsparking/pkg/recaptcha"
)

var db *gorm.DB
var captcha *recaptcha.ReCaptcha

func init() {
	var err error
	captcha = recaptcha.NewReCaptcha(model.GServer)
	db, err = gorm.Open("sqlite3", "nsparking.db")
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.Redirect{}, model.Stat{})
}
