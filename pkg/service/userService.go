package service

import (
	back "project-2x"
	"project-2x/pkg/database"
	tgbotapi "project-2x/pkg/telegramBot/all"
)

type UsersDataService struct {
	db database.UsersData
}

type DailyBonusResponse struct {
    CurrentDay int            `json:"current_day"`
	IsAvailable bool	`json:"is_available"`
    Bonuses    []database.BonusInfo    `json:"bonuses"`
}

func NewUsersService(db database.UsersData) *UsersDataService{
	return &UsersDataService{db: db}
}

func (s *UsersDataService) GetUserProfil(telegramID int64) (back.User, error){
	userData, err := s.db.GetUserData(telegramID)

	if err != nil{
		return back.User{}, err
	}
	return userData, nil
}

func (s *UsersDataService) GetDailyBonusInfo(telegramID int64) (DailyBonusResponse, error){
	var dailyBonus DailyBonusResponse

	info, currentDay, available, err := s.db.GetUserDailyBonusData(int64(telegramID))

	if err != nil{
		return DailyBonusResponse{}, err
	}
	dailyBonus.CurrentDay = currentDay
	dailyBonus.Bonuses = info
	dailyBonus.IsAvailable = available

	return dailyBonus, nil
}

func (s *UsersDataService) ClaimDailyBonus(telegramID int64) (int, error){
	value, err := s.db.ClaimDailyBonys(telegramID)

	if err != nil{
		return 0, err
	}
	return value, nil
}

func (s *UsersDataService) AddPayment(telegramID int64, amount int, currency string) error{
	err := s.db.AddPayment(telegramID, amount, currency)

	if err != nil{
		return err
	}
	return nil
}

func (s *UsersDataService) GetReferralCode(telegramID int64) (string, error){
	refKey, err := s.db.GetOrCreateReferralCode(telegramID)

	if err != nil{
		return "", err
	}
	return refKey, nil
}

func (s *UsersDataService) GetUserReferrals(telegramID int64, offset int, pageSize int) ([]database.UserData, error){
	referrals, err := s.db.GetUserReferrals(telegramID, offset, pageSize)

	if err != nil{
		return []database.UserData{}, err
	}
	return referrals, nil
}

func (s *UsersDataService) CreateUser(message tgbotapi.Message, refKey string) error{
	err := s.db.CreateUser(message, refKey)

	if err != nil{
		return err
	}
	return nil
}
