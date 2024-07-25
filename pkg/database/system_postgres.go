package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Create needed params

type SystemPostgres struct {
	db *sqlx.DB
}

func NewSystemPostgres(db *sqlx.DB) *SystemPostgres {
	return &SystemPostgres{db: db}
}

func (r *SystemPostgres) UpdateDailyBonuses(dailyBonuses map[string]string) error {

	for day, price := range dailyBonuses {
		query := fmt.Sprintf(`
		INSERT INTO %s (day, price) 
		VALUES ($1, $2) 
		ON CONFLICT (day) 
		DO UPDATE SET price = EXCLUDED.price;
		`, dailyBonusesInfoTable)
		_, err := r.db.Exec(query, day, price)
		if err != nil {
			return fmt.Errorf("failed to update daily bonus for day %d: %w", day, err)
		}
	}
	
	logrus.Info("Daily bonuses updated successfully")
	return nil
}