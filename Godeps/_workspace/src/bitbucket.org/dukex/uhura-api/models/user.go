package models

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/uniuri"
	authenticator "github.com/dukex/go-auth"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id            int64
	Name          string
	Email         string `sql:"not null;unique"`
	Password      string `sql:"type:varchar(100);"`
	Locale        string
	CreatedAt     time.Time
	LastVisitedAt time.Time
	Provider      string `sql:"type:varchar(100);"`
	ProviderId    string `sql:"type:varchar(50);"`
	RememberToken string `sql:"type:varchar(100);unique"`
	ApiToken      string `sql:"type:varchar(100);unique"`
	WelcomeMail   bool
}

type UserEntity struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	Email      string `json:"email"`
	ProviderId string `json:"provider_id"`
	ApiToken   string `json:"token"`
}

func (self User) TableName() string {
	return "users"
}

type UserHelper struct {
	DB gorm.DB
}

func (h *UserHelper) PasswordByEmail(email string) (string, bool) {
	var u struct {
		Password string
	}

	err := h.DB.Table(User{}.TableName()).
		Select("password").
		Where("email = ? ", email).Scan(&u).Error

	if err != nil {
		return "", false
	}

	return u.Password, true
}
func (h *UserHelper) FindUserDataByEmail(email string) (string, bool) {
	var user UserEntity
	err := h.DB.Table(User{}.TableName()).
		Where("email = ? ", email).First(&user).Error

	if err != nil {
		return "", false
	}

	userJSON, err := json.Marshal(&user)

	if err != nil {
		return "", false
	}

	return string(userJSON), true
}
func (h *UserHelper) FindUserByToken(token string) (string, bool) {
	var u struct {
		Id string
	}

	err := h.DB.Table(User{}.TableName()).
		Select("id").
		Where("api_token = ? ", token).Scan(&u).Error

	if err != nil {
		return "", false
	}

	return u.Id, true
}
func (h *UserHelper) FindUserFromOAuth(provider string, user *authenticator.User, rawResponse *http.Response) (userID string, err error) {
	var u struct {
		Id string
	}
	err = h.DB.Table(User{}.TableName()).
		Where("email = ?", user.Email).
		First(&u).
		Error
	if err != nil {
		if err == gorm.RecordNotFound {
			userData := User{
				Email:      user.Email,
				Password:   uniuri.NewLen(6),
				Provider:   provider,
				ProviderId: user.Id,
				Locale:     user.Locale,
				Name:       user.Name,
				ApiToken:   user.Token,
			}
			err = h.DB.Table(User{}.TableName()).Save(&userData).Error
			userID = strconv.Itoa(int(userData.Id))
		}
	} else {
		userID = u.Id
	}

	return
}

func NewUserHelper(db gorm.DB) *UserHelper {
	return &UserHelper{db}
}
