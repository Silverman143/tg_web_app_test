package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)
const (
	usersTable      = "users"
	allTimeLeadersTable = "All_Time_Leaders"
	referralsTable = "Referrals"
	weekLeaderboardTable = "Current_Week_Leaders"
	historicWeekLeaderboardTable  = "Historic_Week_Leaders"
	monthLeaderboardTable = "Current_Month_Leaders"
	historicMonthLeaderboardTable = "Historic_Month_Leaders"
	usersDailyBonusesTable = "Users_daily_bonuses"
	dailyBonusesInfoTable = "Daily_bonuses_info"
	tasksTable = "Tasks"
	tasksUpdatesTable = "Tasks_updates"
	transactionsTable = "Transactions"
)

type Config struct {
	URL      	string
	Host 		string
	Port 		string
	Username 	string
	Password 	string
	DBname 		string 
	SSLmode 	string
}

func NewPostgresDB(cfg *Config) (*sqlx.DB, error){
	var connStr string

	if cfg.URL != "" {

		connStr = cfg.URL

	} else {

		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode= %s", cfg.Host, cfg.Port,cfg.Username, cfg.DBname, cfg.Password, cfg.SSLmode)

	}

	db, err := sqlx.Open("postgres", connStr)

	if err != nil{
		return nil, err
	}
	err = db.Ping()
	if err != nil{
		return nil, err
	}
	return db, nil
}