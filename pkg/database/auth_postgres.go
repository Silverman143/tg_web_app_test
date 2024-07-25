package database

import (
	"database/sql"
	"fmt"
	"time"

	back "project-2x"

	// "github.com/go-playground/locales/currency"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) GetOrCreateUser(userInitData initdata.InitData) (back.User, error) {
	var user back.User

	user.Telegram_id = userInitData.User.ID
	user.Name = userInitData.User.Username
	user.AvatarUrl = userInitData.User.PhotoURL
	user.Stars = 0

	referrerKey := uuid.New().String()


	var id int
	currentTime := time.Now()
	
	// Вставка пользователя
	query := fmt.Sprintf(
		"INSERT INTO %s (telegram_id, user_name, lang, registration_date, avatar_url, stars_balance, referrer_key) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (telegram_id) DO NOTHING RETURNING id", 
		usersTable)

	err := r.db.QueryRow(query, 
		userInitData.User.ID, 
		userInitData.User.Username, 
		userInitData.User.LanguageCode, 
		currentTime, 
		userInitData.User.PhotoURL,
		0,
		referrerKey,
		).Scan(&id)
	
	if err != nil && err != sql.ErrNoRows{
		return back.User{}, err
	}

	// Если пользователь уже существует, получаем его id
	if err == sql.ErrNoRows {
		query = fmt.Sprintf("SELECT id, stars_balance FROM %s WHERE telegram_id = $1", usersTable)
		err = r.db.QueryRow(query, userInitData.User.ID).Scan(&id, &user.Stars)
		if err != nil {
			return back.User{}, err
		}
	}

	// Вставка пользователя в таблицу All_Time_Leaders, если он не существует
	query = fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING", allTimeLeadersTable)
	err = r.db.QueryRow(query, id).Scan()
	if err != nil && err != sql.ErrNoRows {
		return back.User{}, err
	}

	// Получение позиции пользователя по количеству звезд
	query = fmt.Sprintf("SELECT rank FROM (SELECT user_id, RANK() OVER (ORDER BY stars DESC) AS rank FROM %s) subquery WHERE user_id = $1", allTimeLeadersTable)
	err = r.db.QueryRow(query, id).Scan(&user.GlobalRank)
	if err != nil {
		return back.User{}, err
	}

	return user, nil
}