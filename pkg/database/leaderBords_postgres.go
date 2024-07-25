package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type LeaderBordsPostgres struct {
	db *sqlx.DB
}


func NewLeaderBordsPostgres(db *sqlx.DB) *LeaderBordsPostgres {
	return &LeaderBordsPostgres{db: db}
}

type LeaderboardEntry struct {
    TelegramID int64  `json:"telegram_id"`
    UserName   string `json:"user_name"`
    AvatarURL  string `json:"avatar_url"`
    Stars      int    `json:"stars"`
}

func (r *LeaderBordsPostgres) GetAllTimeLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error) {
	query := fmt.Sprintf(`
        SELECT u.telegram_id, u.user_name, u.avatar_url, atl.stars
        FROM %s atl
        JOIN Users u ON atl.user_id = u.id
        ORDER BY atl.stars DESC
        LIMIT 100;
    `, allTimeLeadersTable)

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var leaderboard []LeaderboardEntry
    for rows.Next() {
        var entry LeaderboardEntry
        if err := rows.Scan(&entry.TelegramID, &entry.UserName, &entry.AvatarURL, &entry.Stars); err != nil {
            return nil, 0, err
        }
        leaderboard = append(leaderboard, entry)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, err
    }

	var rank int = 1

	query = fmt.Sprintf("SELECT rank FROM ( SELECT u.telegram_id, ROW_NUMBER() OVER (ORDER BY atl.stars DESC) AS rank FROM %s atl JOIN %s u ON atl.user_id = u.id ) subquery WHERE subquery.telegram_id = $1;", allTimeLeadersTable, usersTable)
	err = r.db.QueryRow(query, telegramID).Scan(&rank)
	if err != nil {
		return nil, 0, err
	}

    return leaderboard, rank, nil
}

func (r *LeaderBordsPostgres) GetCurrentMonthLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error) {
	query := fmt.Sprintf(`
	SELECT u.telegram_id, u.user_name, u.avatar_url, atl.stars
	FROM %s atl
	JOIN %s u ON atl.user_id = u.id
	ORDER BY atl.stars DESC
	LIMIT 100;
	`, monthLeaderboardTable, usersTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.TelegramID, &entry.UserName, &entry.AvatarURL, &entry.Stars); err != nil {
			return nil, 0, err
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var rank int

	query = fmt.Sprintf("SELECT rank FROM ( SELECT u.telegram_id, ROW_NUMBER() OVER (ORDER BY atl.stars DESC) AS rank FROM %s atl JOIN %s u ON atl.user_id = u.id ) subquery WHERE subquery.telegram_id = $1;", monthLeaderboardTable, usersTable)
	err = r.db.QueryRow(query, telegramID).Scan(&rank)

	if err != nil {
		if err == sql.ErrNoRows {
			// Если не найден ранг для данного пользователя
			return leaderboard, 0, nil
		}
		return nil, 0, err
	}

	return leaderboard, rank, nil
}

func (r *LeaderBordsPostgres) GetCurrentWeekLeaderbord(telegramID int64) ([]LeaderboardEntry, int, error) {
	query := fmt.Sprintf(`
	SELECT u.telegram_id, u.user_name, u.avatar_url, atl.stars
	FROM %s atl
	JOIN %s u ON atl.user_id = u.id
	ORDER BY atl.stars DESC
	LIMIT 100;
	`, weekLeaderboardTable, usersTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.TelegramID, &entry.UserName, &entry.AvatarURL, &entry.Stars); err != nil {
			return nil, 0, err
		}
		leaderboard = append(leaderboard, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var rank int

	query = fmt.Sprintf("SELECT rank FROM ( SELECT u.telegram_id, ROW_NUMBER() OVER (ORDER BY atl.stars DESC) AS rank FROM %s atl JOIN %s u ON atl.user_id = u.id ) subquery WHERE subquery.telegram_id = $1;", weekLeaderboardTable, usersTable)
	err = r.db.QueryRow(query, telegramID).Scan(&rank)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// Если не найден ранг для данного пользователя
			return leaderboard, 0, nil
		}
		return nil, 0, err
	}

	return leaderboard, rank, nil
}