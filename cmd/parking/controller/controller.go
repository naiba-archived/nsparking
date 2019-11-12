package controller

import (
	"fmt"
	"net"
	"regexp"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/jinzhu/gorm"

	// ..
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/naiba/nsparking/model"
	"github.com/naiba/nsparking/pkg/recaptcha"
)

var db *gorm.DB
var captcha *recaptcha.ReCaptcha
var serverRegexp *regexp.Regexp

func init() {
	var err error
	serverRegexp = regexp.MustCompile(``)
	captcha = recaptcha.NewReCaptcha(model.GServer)
	db, err = gorm.Open("sqlite3", "nsparking.db")
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.Redirect{}, model.Stat{})
}

func getRedirectByDomain(domain string) (*model.Redirect, error) {
	ns, err := net.LookupNS(domainutil.Domain(domain))
	if err != nil || len(ns) == 0 {
		return nil, fmt.Errorf("NS设置错误：%s", err)
	}
	var server string
	for i := 0; i < len(ns); i++ {
		if ns[i] == nil {
			continue
		}
	}
	if len(server) == 0 {
		return nil, fmt.Errorf("NS设置错误：%v", ns)
	}
	var rd model.Redirect
	err = db.Model(&model.Redirect{}).Where("server = ?", server).First(&rd).Error
	return &rd, err
}

func parseDomain(domain string) (string, string, string) {
	return domainutil.Domain(domain), domainutil.DomainPrefix(domain), domainutil.DomainSuffix(domain)
}
