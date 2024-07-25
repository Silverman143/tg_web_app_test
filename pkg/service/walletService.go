package service

import "project-2x/pkg/database"

type WalletService struct {
	db database.Wallet
}

func NewWalletService(db database.Wallet) *WalletService{
	return &WalletService{db: db}
}

func (s *WalletService) GetAllTransactions(telegramID int64) ([]database.TransactionInfo, error){
	transactions, err := s.db.GetAllTransactions(telegramID)

	if err != nil{
		return transactions, err
	}
	return transactions, nil
}

func (s *WalletService) GetPositiveTransactions(telegramID int64) ([]database.TransactionInfo, error){
	transactions, err := s.db.GetPositiveTransactions(telegramID)

	if err != nil{
		return transactions, err
	}
	return transactions, nil
}

func (s *WalletService) GetNegativeTransactions(telegramID int64) ([]database.TransactionInfo, error){
	transactions, err := s.db.GetNegativeTransactions(telegramID)

	if err != nil{
		return transactions, err
	}
	return transactions, nil
}

func (s *WalletService) GetBalance(telegramID int64) (int, error){
	balance, err := s.db.GetBalance(telegramID)

	if err != nil{
		return 0, err
	}
	return balance, nil
}