package service

import (
	"errors"
	"time"

	back "project-2x"
	"project-2x/pkg/database"

	"github.com/dgrijalva/jwt-go"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const (
	salt = "asdasdasdasdasd"
    accessTokenTTL  = time.Hour * 1    // Access token живет 1 час
    refreshTokenTTL = time.Hour * 24 * 7  // Refresh token живет 7 дней
    signingKey      = "asdasdasdasdasd" // Ваш секретный ключ для подписи JWT
)

type tokenClaims struct{
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

type AuthService struct {
	db database.Authorization
}

func NewAuthService(db database.Authorization) *AuthService{
	return &AuthService{db: db}
}

func (s *AuthService) CreateUser(user initdata.InitData) (back.User, error){
	return s.db.GetOrCreateUser(user)
}


func (s *AuthService) GenerateAccessToken(userId int64) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
        UserId: userId,                                                         // telegram user id 
    })
    return token.SignedString([]byte(signingKey))
}


func (s *AuthService) ParseToken(accessToken string) (int64, error){
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("invalid signin method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	clames, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return clames.UserId, nil 
}
