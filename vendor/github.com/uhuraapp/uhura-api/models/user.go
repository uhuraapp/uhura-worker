package models

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/uhuraapp/uhura-api/entities"

	"github.com/dchest/uniuri"
	authenticator "github.com/dukex/go-auth"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id                           int64
	Name                         string
	Email                        string `sql:"not null;unique"`
	Password                     string `sql:"type:varchar(100);"`
	Locale                       string
	CreatedAt                    time.Time
	LastVisitedAt                time.Time
	Provider                     string `sql:"type:varchar(100);"`
	ProviderId                   string `sql:"type:varchar(50);"`
	RememberToken                string `sql:"type:varchar(100);unique"`
	ApiToken                     string `sql:"type:varchar(100);unique"`
	WelcomeMail                  bool
	AgreeWithTheTermsAndPolicyAt time.Time
	AgreeWithTheTermsAndPolicyIn string
	OptInAt                      time.Time
	DeletedAt                    time.Time
	ExportKey                    string `sql:"DEFAULT:uuid_generate_v4()"`
}

func (self User) TableName() string {
	return "users"
}

type UserHelper struct {
	DB *gorm.DB
}

func (h *UserHelper) PasswordByEmail(email string) (string, bool) {
	var u User

	err := h.DB.Table(User{}.TableName()).
		Select("password").
		Where("email = ? ", email).First(&u).Error

	if err != nil {
		return "", false
	}

	return u.Password, true
}

func (h *UserHelper) FindUserDataByEmail(email string) (string, string, bool) {
	var user entities.User

	err := h.DB.Table(User{}.TableName()).
		Where("email = ? ", email).First(&user).Error

	if err != nil {
		return "", "", false
	}

	userId := strconv.Itoa(int(user.Id))

	user.OptIn = !user.OptInAt.IsZero()

	userJSON, err := json.Marshal(&user)

	if err != nil {
		log.Println("ERROR", err)
		return "", "", false
	}

	return string(userJSON), userId, true
}

func (h *UserHelper) FindUserByToken(token string) (string, bool) {
	var u User

	err := h.DB.Table(User{}.TableName()).
		Select("id").
		Where("api_token = ? ", token).First(&u).Error

	if err != nil {

		return "", false
	}
	id := strconv.Itoa(int(u.Id))

	return id, true
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
		if err == gorm.ErrRecordNotFound {
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

func NewUserHelper(db *gorm.DB) *UserHelper {
	return &UserHelper{db}
}
