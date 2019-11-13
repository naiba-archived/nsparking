package controller

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/miekg/dns"

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
var resolver dns.Client

func init() {
	var err error
	serverRegexp = regexp.MustCompile(`(.{10})\.ns\d\.` + model.Domain)
	captcha = recaptcha.NewReCaptcha(model.GServer)
	db, err = gorm.Open("sqlite3", "nsparking.db")
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.Redirect{}, model.Stat{})
}

func getNS(domain string) ([]string, error) {
	q := new(dns.Msg)
	q.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	q.RecursionDesired = true
	msg, _, err := resolver.Exchange(q, "223.5.5.5:53")
	var ns []string
	for i := 0; i < len(msg.Answer); i++ {
		ns = append(ns, msg.Answer[i].(*dns.NS).Ns)
	}
	return ns, err
}

func getRedirectByDomain(domain string) (*model.Redirect, error) {
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}
	ns, err := getNS(domainutil.Domain(domain))
	if err != nil || len(ns) == 0 {
		return nil, fmt.Errorf("NS设置错误：%s", err)
	}
	var server string
	for i := 0; i < len(ns); i++ {
		matches := serverRegexp.FindStringSubmatch(ns[i])
		if len(matches) == 2 {
			server = matches[1]
			break
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
