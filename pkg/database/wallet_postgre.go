package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type WalletPostgres struct {
	db *sqlx.DB
}

func NewWalletPostgres(db *sqlx.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

type TransactionInfo struct {
    TransactionType string    `json:"transactionType"`
    Date            time.Time `json:"date"`
    Value           int       `json:"value"`
    Currency        string    `json:"currency"`
    Status          string    `json:"status"`
}

func (r *WalletPostgres) GetAllTransactions(telegramID int64) ([]TransactionInfo, error) {
    var transactions []TransactionInfo
    query := fmt.Sprintf(`
        SELECT type, date, amount, currency, status
        FROM %s
        WHERE user_id = (SELECT id FROM %s WHERE telegram_id = $1)
    `, transactionsTable, usersTable)
    rows, err := r.db.Query(query, telegramID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var transaction TransactionInfo
        if err := rows.Scan(&transaction.TransactionType, &transaction.Date, &transaction.Value, &transaction.Currency, &transaction.Status); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return transactions, nil
}

func (r *WalletPostgres) GetPositiveTransactions(telegramID int64) ([]TransactionInfo, error) {
    var transactions []TransactionInfo
    query := fmt.Sprintf(`
        SELECT type, date, amount, currency, status
        FROM %s
        WHERE user_id = (SELECT id FROM %s WHERE telegram_id = $1) AND amount > 0
    `, transactionsTable, usersTable)
    rows, err := r.db.Query(query, telegramID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var transaction TransactionInfo
        if err := rows.Scan(&transaction.TransactionType, &transaction.Date, &transaction.Value, &transaction.Currency, &transaction.Status); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return transactions, nil
}

func (r *WalletPostgres) GetNegativeTransactions(telegramID int64) ([]TransactionInfo, error) {
    var transactions []TransactionInfo
    query := fmt.Sprintf(`
        SELECT type, date, amount, currency, status
        FROM %s
        WHERE user_id = (SELECT id FROM %s WHERE telegram_id = $1) AND amount < 0
    `, transactionsTable, usersTable)
    rows, err := r.db.Query(query, telegramID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var transaction TransactionInfo
        if err := rows.Scan(&transaction.TransactionType, &transaction.Date, &transaction.Value, &transaction.Currency, &transaction.Status); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return transactions, nil
}

func (r *WalletPostgres) GetBalance(telegramID int64) (int, error){
	var balance int
	query := fmt.Sprintf(`
        SELECT stars_balance
        FROM %s
        WHERE telegram_id = $1
    `, usersTable)
	err := r.db.QueryRow(query, telegramID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows{
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return balance, nil
}