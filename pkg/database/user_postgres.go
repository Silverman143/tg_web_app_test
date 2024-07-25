package database

import (
	"database/sql"
	"fmt"
	back "project-2x"
	tgbotapi "project-2x/pkg/telegramBot/all"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

type UserData struct {
	TelegramID string `db:"telegram_id"`
	Username   string `db:"username"`
	Avatar     string `db:"avatar"`
	Rank       int64    `db:"rank"`  
	Stars      int64    `db:"stars"`
}

func NewUserPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}


func (r *UsersPostgres) GetUserData(telegramID int64) (back.User, error) {
	var userData back.User
	query := fmt.Sprintf("SELECT telegram_id, user_name, avatar_url, stars_balance FROM %s WHERE telegram_id=$1", usersTable)

	err := r.db.Get(&userData, query, telegramID)
	if err != nil{
		return back.User{}, err
	}

	query = fmt.Sprintf("SELECT rank FROM ( SELECT u.telegram_id, ROW_NUMBER() OVER (ORDER BY atl.stars DESC) AS rank FROM %s atl JOIN %s u ON atl.user_id = u.id ) subquery WHERE subquery.telegram_id = $1;", allTimeLeadersTable, usersTable)
	err = r.db.QueryRow(query, telegramID).Scan(&userData.GlobalRank)
	if err != nil {
		return back.User{}, err
	}

	return userData, nil
}

type UserBonusData struct {
    DaysCounter  int       `db:"days_counter"`
    LastCollected *time.Time `db:"last_collected"`
}

func hasMoreThan24HoursPassed(dateTime *time.Time) bool {
    if dateTime == nil {
        return true
    }
    duration := time.Since(*dateTime)
    return duration.Hours() > 24
}

func hasMoreThan48HoursPassed(dateTime *time.Time) bool {
    if dateTime == nil {
        return true
    }
    duration := time.Since(*dateTime)
    return duration.Hours() > 48
}


func (r *UsersPostgres) GetUserDailyBonusData(telegramID int64) ([]BonusInfo, int, bool, error) {
	var userBonusData UserBonusData
	var available bool

	query := fmt.Sprintf("SELECT udb.days_counter, udb.last_collected FROM %s u JOIN %s udb ON u.id = udb.user_id WHERE u.telegram_id = $1;", usersTable, usersDailyBonusesTable)
	err := r.db.Get(&userBonusData, query, telegramID)

	if err != nil{
		if err == sql.ErrNoRows {
            insertQuery := fmt.Sprintf("INSERT INTO %s (user_id, days_counter) SELECT u.id, 0 FROM %s u WHERE u.telegram_id = $1 RETURNING 0;", usersDailyBonusesTable, usersTable)
            _, err = r.db.Exec(insertQuery, telegramID)
            if err != nil {
                return nil, 0, available, err
            }
            userBonusData.DaysCounter = 0 // Новая запись, поэтому days = 0
			available = true
        } else {
            return nil, 0, available, err
        }
	}else{
		available = hasMoreThan24HoursPassed(userBonusData.LastCollected)
	}

	if hasMoreThan48HoursPassed(userBonusData.LastCollected) && userBonusData.DaysCounter > 0{
		updateQuery := fmt.Sprintf("UPDATE %s udb SET days_counter = 0 FROM %s u WHERE u.id = udb.user_id AND u.telegram_id = $1;", usersDailyBonusesTable, usersTable)
        _, err = r.db.Exec(updateQuery, telegramID)
        if err != nil {
            return nil, 0, false, err
        }
	}

	dailyBonusInfo, err := r.GetDailyBonusesInfo()
	
	userBonusData.LastCollected = nil

	if err != nil{
		return nil, 0, available,  err
	}

	return dailyBonusInfo, userBonusData.DaysCounter, available, nil
}
func (r *UsersPostgres) ClaimDailyBonys(telegramID int64) (int, error) {
	dailyBonusInfo, daysCounter, available, err := r.GetUserDailyBonusData(telegramID)
    if err != nil {
        return 0, err
    }
    if !available {
        return 0, fmt.Errorf("daily bonus not available yet")
    }

    // Определяем бонус на основе daysCounter
    var bonusAmount int
    if daysCounter < len(dailyBonusInfo) {
        bonusAmount = dailyBonusInfo[daysCounter].Price
    } else {
        bonusAmount = dailyBonusInfo[0].Price
    }

	err = r.AddPayment(telegramID, bonusAmount, "daily_bonus")
	
	if err != nil{
		return 0, err
	}

	tx, err := r.db.Begin()
    if err != nil {
        return 0, err
    }
    // Обновляем Users_daily_bonuses
    updateDailyBonusesQuery := "UPDATE Users_daily_bonuses SET last_collected = CURRENT_DATE, days_counter = days_counter + 1 WHERE user_id = (SELECT id FROM Users WHERE telegram_id = $1)"
    _, err = tx.Exec(updateDailyBonusesQuery, telegramID)
    if err != nil {
		fmt.Println("error - 4")
        return 0, err
    }
	defer tx.Rollback()

	if err = tx.Commit(); err != nil {
		return  0, err
	}

    return bonusAmount, nil
}

//Referrals

func (r *UsersPostgres) GetOrCreateReferralCode(telegramID int64) (string, error){
	var refKey sql.NullString
	query := fmt.Sprintf("SELECT referrer_key FROM %s WHERE telegram_id=$1", usersTable)

	err := r.db.Get(&refKey, query, telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no referrer_key found, create a new one
			refKey.String = uuid.New().String()
			refKey.Valid = true

			updateQuery := fmt.Sprintf("UPDATE %s SET referrer_key = $1 WHERE telegram_id = $2", usersTable)
			_, err := r.db.Exec(updateQuery, refKey, telegramID)
			if err != nil {
				return "", err
			}
			return refKey.String, nil
		} else {
			return "", err
		}
	}

		// If referrer_key is NULL, generate a new one
		if !refKey.Valid {
			refKey.String = uuid.New().String()
			refKey.Valid = true
	
			updateQuery := fmt.Sprintf("UPDATE %s SET referrer_key = $1 WHERE telegram_id = $2", usersTable)
			_, err := r.db.Exec(updateQuery, refKey.String, telegramID)
			if err != nil {
				return "", err
			}
		}

	return refKey.String, nil
}

func (r *UsersPostgres) GetUserReferrals(telegramID int64, offset int, pageSize int ) ([]UserData, error){
	var referrals []UserData
	query := fmt.Sprintf(`
		WITH rank AS (
			SELECT ROW_NUMBER() OVER (ORDER BY atl.stars DESC) AS rank, u.telegram_id
			FROM %s atl
			JOIN %s u ON atl.user_id = u.id
		)
		SELECT u.telegram_id, u.user_name, u.avatar_url, u.stars_balance, rnk.rank
		FROM %s r
		JOIN %s u ON r.referral = u.id
		JOIN rank rnk ON rnk.telegram_id = u.telegram_id
		WHERE r.referrer = (SELECT id FROM %s WHERE telegram_id = $1)
		LIMIT $2 OFFSET $3;

	`, allTimeLeadersTable, usersTable, referralsTable, usersTable, usersTable)

	rows, err := r.db.Query(query, telegramID, pageSize, offset)
	if err != nil {
		return []UserData{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var referral UserData
		if err := rows.Scan(&referral.TelegramID, &referral.Username, &referral.Avatar,&referral.Stars, &referral.Rank); err != nil {
			return []UserData{}, err
		}
		referrals = append(referrals, referral)
	}

	if err := rows.Err(); err != nil {
		return []UserData{}, err
	}
	return referrals, nil

}

func (r *UsersPostgres) AddPayment(telegramID int64, amount int, paymentType string) error {

	var err error

	tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    defer tx.Rollback()
	// Обновляем stars_balance
	updateStarsBalanceQuery := fmt.Sprintf("UPDATE %s SET stars_balance = stars_balance + $1 WHERE telegram_id = $2", usersTable)
	_, err = tx.Exec(updateStarsBalanceQuery, amount, telegramID)
	if err != nil {
		return err
	}
   
	// Обновляем All_Time_Leaders
	updateAllTimeLeadersQuery := fmt.Sprintf("INSERT INTO %s (user_id, stars) VALUES ((SELECT id FROM %s WHERE telegram_id = $1), $2) ON CONFLICT (user_id) DO UPDATE SET stars = %s.stars + $2", allTimeLeadersTable, usersTable, allTimeLeadersTable)
	_, err = tx.Exec(updateAllTimeLeadersQuery, telegramID, amount)
	if err != nil {
	   fmt.Println("error - 1")
	   return err
	}
   
	dataTimeStart, dataTimeEnd := StartAndEndOfWeek(time.Now())
   
	// Обновляем Current_Week_Leaders
	updateWeekLeadersQuery := fmt.Sprintf("INSERT INTO %s (user_id, stars, week_start, week_end) VALUES ((SELECT id FROM %s WHERE telegram_id = $1), $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET stars = %s.stars + $2", weekLeaderboardTable, usersTable, weekLeaderboardTable)
	_, err = tx.Exec(updateWeekLeadersQuery, telegramID, amount, dataTimeStart, dataTimeEnd)
	if err != nil {
		fmt.Println("error - 2")
		return err
	}
   
	dataTimeStart, dataTimeEnd = StartAndEndOfMonth(time.Now())
   
	// Обновляем Current_Month_Leaders
	updateMonthLeadersQuery := fmt.Sprintf("INSERT INTO %s (user_id, stars, month_start, month_end) VALUES ((SELECT id FROM %s WHERE telegram_id = $1), $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET stars = %s.stars + $2", monthLeaderboardTable, usersTable, monthLeaderboardTable)
	_, err = tx.Exec(updateMonthLeadersQuery, telegramID, amount, dataTimeStart, dataTimeEnd)
	if err != nil {
		fmt.Println("error - 3")
		return err
	}

	// Вставляем запись в Transactions
	insertTransactionQuery := fmt.Sprintf(`INSERT INTO %s (type, user_id, date, amount, currency, status) 
	VALUES ($1, (SELECT id FROM %s WHERE telegram_id = $2), CURRENT_TIMESTAMP, $3, 'stars', 'completed')`, transactionsTable, usersTable)
	_, err = tx.Exec(insertTransactionQuery, paymentType, telegramID, amount)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return  err
	}
	
	return nil
}

//Create user from tg bot with start 
func (r *UsersPostgres) CreateUser(message tgbotapi.Message, refKey string) error {

	referrerKey := uuid.New().String()

	var id int
	currentTime := time.Now()
	userCreated := false
	
	// Вставка пользователя
	query := fmt.Sprintf(
		"INSERT INTO %s (telegram_id, user_name, lang, registration_date, avatar_url, stars_balance, referrer_key) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (telegram_id) DO NOTHING RETURNING id", 
		usersTable)

	err := r.db.QueryRow(query, 
		message.From.ID, 
		message.From.UserName, 
		message.From.LanguageCode, 
		currentTime, 
		"",
		0,
		referrerKey,
		).Scan(&id)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователь уже существует, получаем его id
			query = fmt.Sprintf("SELECT id FROM %s WHERE telegram_id = $1", usersTable)
			err = r.db.QueryRow(query, message.From.ID).Scan(&id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		userCreated = true
	}

	referralBonus := 100

	// Создание связи с реферером, если refKey не пустой
	if refKey != "" && userCreated{
		query = fmt.Sprintf(`
			WITH referrer AS (
				SELECT id FROM %s WHERE referrer_key = $1
			)
			INSERT INTO %s (referrer, referral, bonus) 
			SELECT id, $2, $3 FROM referrer
			ON CONFLICT (referrer, referral) DO NOTHING`, usersTable, referralsTable)

		_, err = r.db.Exec(query, refKey, id, referralBonus)
		if err != nil {
			return err
		}
	}

	err = r.AddPayment(message.From.ID, referralBonus, "referrer")

	if err != nil{
		return err
	}

	return nil
}

// StartAndEndOfWeek возвращает дату и время начала и конца текущей недели
func StartAndEndOfWeek(t time.Time) (time.Time, time.Time) {

	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}

	// Начало недели (понедельник)
	startOfWeek := time.Date(t.Year(), t.Month(), t.Day()+offset, 0, 0, 0, 0, t.Location())

	// Конец недели (воскресенье)
	endOfWeek := startOfWeek.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	return startOfWeek, endOfWeek
}

// StartAndEndOfMonth возвращает дату и время начала и конца текущего месяца
func StartAndEndOfMonth(t time.Time) (time.Time, time.Time) {
	// Начало месяца
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	// Конец месяца
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	return startOfMonth, endOfMonth
}

