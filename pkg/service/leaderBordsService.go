package service

import "project-2x/pkg/database"

type LeaderBoardsService struct {
	db database.LeaderBords
}


func NewLeaderBoardsService(db database.LeaderBords) *LeaderBoardsService{
	return &LeaderBoardsService{db: db}
}

func (s *LeaderBoardsService) GetAllTimeLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error){
	leaderBoard, rank, err := s.db.GetAllTimeLeaderbord(telegram_ID)

	if err != nil{
		return leaderBoard, rank, err
	}
	return leaderBoard, rank, nil
}

func (s *LeaderBoardsService) GetCurrentMonthLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error){
	leaderBoard, rank, err := s.db.GetCurrentMonthLeaderbord(telegram_ID)

	if err != nil{
		return leaderBoard, rank, err
	}
	return leaderBoard, rank, nil
}

func (s *LeaderBoardsService) GetCurrentWeekLeaderboard(telegram_ID int64) ([]database.LeaderboardEntry, int, error){
	leaderBoard, rank, err := s.db.GetCurrentWeekLeaderbord(telegram_ID)

	if err != nil{
		return leaderBoard, rank, err
	}
	return leaderBoard, rank, nil
}