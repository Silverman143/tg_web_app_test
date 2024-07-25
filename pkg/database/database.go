package database

import (
	back "project-2x"
	tgbotapi "project-2x/pkg/telegramBot/all"

	"github.com/jmoiron/sqlx"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type Authorization interface {
	GetOrCreateUser(user initdata.InitData) (back.User, error)
}

type UsersData interface{
	GetUserData(telegramID int64) (back.User, error) 
	GetUserDailyBonusData(telegramID int64) ([]BonusInfo, int, bool, error)
	ClaimDailyBonys(telegramID int64) (int, error)
	AddPayment(telegramID int64, amount int, paymentType string) error
	GetOrCreateReferralCode(telegramID int64) (string, error)
	GetUserReferrals(telegramID int64, offset int, pageSize int ) ([]UserData, error)
	CreateUser(message tgbotapi.Message, refKey string) error
}

type LeaderBords interface{
	GetAllTimeLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error)
	GetCurrentMonthLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error)
	GetCurrentWeekLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error)
}

type Wallet interface{
	GetAllTransactions(telegramID int64) ([]TransactionInfo, error)
	GetPositiveTransactions(telegramID int64) ([]TransactionInfo, error)
	GetNegativeTransactions(telegramID int64) ([]TransactionInfo, error)
	GetBalance(telegramID int64) (int, error)
}

type System interface{
	UpdateDailyBonuses(dailyBonuses map[string]string) error
}

type Payment interface {

}

type Database struct {
	Authorization
	UsersData
	LeaderBords
	Wallet
	System
	Payment
}


func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		Authorization: NewAuthPostgres(db), 
		UsersData: NewUserPostgres(db),
		LeaderBords: NewLeaderBordsPostgres(db),
		Wallet: NewWalletPostgres(db),
		System: NewSystemPostgres(db),
		Payment: NewPaymentPostgres(db),
	}
}
