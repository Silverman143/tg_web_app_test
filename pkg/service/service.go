package service

import (
	back "project-2x"
	"project-2x/pkg/database"
	"project-2x/pkg/telegramBot"
	tgbotapi "project-2x/pkg/telegramBot/all"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type Authorization interface {
	CreateUser(user initdata.InitData) (back.User, error)
	GenerateAccessToken(userID int64) (string, error)
	ParseToken(token string) (int64, error)
}

type UsersData interface {
	GetUserProfil(telegramID int64) (back.User, error)
	GetDailyBonusInfo(telegramID int64) (DailyBonusResponse, error)
	ClaimDailyBonus(telegramID int64) (int, error)
	AddPayment(telegramID int64, amount int, currency string) error
	GetReferralCode(telegramID int64) (string, error)
	GetUserReferrals(telegramID int64, offset int, pageSize int) ([]database.UserData, error)
	CreateUser(message tgbotapi.Message, refKey string) error
}

type LeaderBords interface {
	GetAllTimeLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error)
	GetCurrentMonthLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error)
	GetCurrentWeekLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error)
}

type Wallet interface {
	GetAllTransactions(telegramID int64) ([]database.TransactionInfo, error)
	GetPositiveTransactions(telegramID int64) ([]database.TransactionInfo, error)
	GetNegativeTransactions(telegramID int64) ([]database.TransactionInfo, error)
	GetBalance(telegramID int64) (int, error)
}

type Payment interface{
	CreateStarsInvoice(amount int) (string, error)
}

type Service struct {
	Authorization
	UsersData
	LeaderBords
	Wallet
	Payment
}

func NewService(db *database.Database, telegramBot *telegramBot.Bot) *Service {
	return &Service{
		Authorization: NewAuthService( db.Authorization),
		UsersData: NewUsersService( db.UsersData ),
		LeaderBords: NewLeaderBoardsService(db.LeaderBords),
		Wallet: NewWalletService(db),
		Payment: NewPaymentService(db, *telegramBot),
	}
}