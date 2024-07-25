package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConfigPostgres struct {
	db *sqlx.DB
}

type BonusInfo struct {
    Day   int `json:"day"`
    Price  int `json:"price"`
}

func NewConfigPostgres(db *sqlx.DB) *ConfigPostgres {
	return &ConfigPostgres{db: db}
}

// Returns bonus days and their prices 
func (r *UsersPostgres) GetDailyBonusesInfo() ([]BonusInfo, error) {
	var bonusesInfo []BonusInfo
	query := fmt.Sprintf("SELECT day, price FROM %s ORDER BY day", dailyBonusesInfoTable)

	err := r.db.Select(&bonusesInfo, query)
    if err != nil {
        return nil, err
    }

	return bonusesInfo, nil
}