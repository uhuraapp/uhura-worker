package models

import (
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/dukex/login2"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id            int64
	Name          string
	Link          string
	Picture       string
	Gender        string
	Locale        string
	GoogleId      string
	Email         string `sql:"not null;unique"`
	Password      string `sql:"type:varchar(100);"`
	WelcomeMail   bool
	CreatedAt     time.Time
	Provider      string `sql:"type:varchar(100);"`
	ProviderId    string `sql:"type:varchar(50);"`
	RememberToken string `sql:"type:varchar(100);unique"`
	ApiToken      string `sql:"type:varchar(100);unique"`
	LastVisitedAt time.Time
}

func (self User) TableName() string {
	return "users"
}

type UserHelpers struct {
	DB gorm.DB
}

func (h UserHelpers) SetupFromOAuth(provider string, u *login2.User, rawResponde *http.Response) (int64, error) {
	user, err := h.findByEmail(u.Email)
	if err != nil {
		if err == gorm.RecordNotFound {
			return h.createFromOAuth(provider, u)
		}
		return 0, err
	} else {
		return user.Id, nil
	}
}

func (h UserHelpers) ByToken(token string) (int64, bool) {
	var user User

	if token == "" {
		return 0, false
	}

	if CACHED_TOKEN[token] != 0 {
		return CACHED_TOKEN[token], true
	}

	err := h.DB.Table(User{}.TableName()).Where("api_token = ?", token).First(&user).Error
	if err != nil {
		return 0, false
	}

	CACHED_TOKEN[token] = user.Id
	return user.Id, true
}

//

func (h UserHelpers) findByEmail(email string) (user User, err error) {
	err = h.DB.Table(User{}.TableName()).Where("email = ?", email).First(&user).Error
	return
}

func (h UserHelpers) createFromOAuth(provider string, temp *login2.User) (int64, error) {
	user := User{
		Email:      temp.Email,
		Password:   uniuri.NewLen(6),
		Provider:   provider,
		ProviderId: temp.Id,
		Link:       temp.Link,
		Picture:    temp.Picture,
		Locale:     temp.Locale,
		Name:       temp.Name,
		ApiToken:   temp.Token,
	}
	err := h.DB.Table(User{}.TableName()).Save(&user).Error

	return user.Id, err
}

var CACHED_TOKEN map[string]int64

func init() {
	CACHED_TOKEN = make(map[string]int64, 0)
}
