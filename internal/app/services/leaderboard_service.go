package services

import (
	"sort"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type LeaderboardService struct {
	LeaderboardRepo *repositories.LeaderboardRepository
	ParticipantRepo repositories.ParticipantRepositoryInterface
	UserRepo        repositories.UserRepositoryInterface
	QuizRepo        repositories.QuizRepositoryInterface
}

func NewLeaderboardService(
	leaderboardRepo *repositories.LeaderboardRepository,
	participantRepo repositories.ParticipantRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	quizRepo repositories.QuizRepositoryInterface,
) *LeaderboardService {
	return &LeaderboardService{
		LeaderboardRepo: leaderboardRepo,
		ParticipantRepo: participantRepo,
		UserRepo:        userRepo,
		QuizRepo:        quizRepo,
	}
}

type LeaderboardEntry struct {
	User          models.User `json:"user"`
	Score         int64       `json:"score"`
	Rank          int         `json:"rank"`
	IsCurrentUser bool        `json:"is_current_user"`
}

type LeaderboardResponse struct {
	Juara1      *LeaderboardEntry  `json:"juara1"`
	Juara2      *LeaderboardEntry  `json:"juara2"`
	Juara3      *LeaderboardEntry  `json:"juara3"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}

func (s *LeaderboardService) GetLeaderboard(moduleID *uint, quizType string, currentUserID *uint) (*LeaderboardResponse, error) {
	// Get all users
	allUsers, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Create map to store user scores
	userScores := make(map[uint]int64)
	userMap := make(map[uint]models.User)

	// Initialize all users with 0 score
	for _, user := range allUsers {
		userScores[user.ID] = 0
		userMap[user.ID] = user
	}

	// Get participants with their scores
	participants, err := s.ParticipantRepo.GetParticipantsWithScores(moduleID, quizType)
	if err != nil {
		return nil, err
	}

	// Update scores for users who have participated
	for _, participant := range participants {
		if participant.User.ID != 0 {
			userScores[participant.UserID] += int64(participant.TotalScore)
		}
	}

	// Convert to slice and sort by score
	var entries []LeaderboardEntry
	for userID, score := range userScores {
		isCurrentUser := false
		if currentUserID != nil && *currentUserID == userID {
			isCurrentUser = true
		}

		entries = append(entries, LeaderboardEntry{
			User:          userMap[userID],
			Score:         score,
			IsCurrentUser: isCurrentUser,
		})
	}

	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	response := &LeaderboardResponse{}

	// Assign top 3 positions
	if len(entries) > 0 {
		response.Juara1 = &entries[0]
	}
	if len(entries) > 1 {
		response.Juara2 = &entries[1]
	}
	if len(entries) > 2 {
		response.Juara3 = &entries[2]
	}

	// Rest go to leaderboard array (from 4th place onwards)
	if len(entries) > 3 {
		response.Leaderboard = entries[3:]
	} else {
		response.Leaderboard = []LeaderboardEntry{}
	}

	return response, nil
}

func (s *LeaderboardService) GetUserRank(userID uint, moduleID *uint, quizType string, currentUserID *uint) (*LeaderboardEntry, error) {
	// Get leaderboard with current user marked
	leaderboard, err := s.GetLeaderboard(moduleID, quizType, currentUserID)
	if err != nil {
		return nil, err
	}

	// Find user in top 3
	if leaderboard.Juara1 != nil && leaderboard.Juara1.User.ID == userID {
		return leaderboard.Juara1, nil
	}
	if leaderboard.Juara2 != nil && leaderboard.Juara2.User.ID == userID {
		return leaderboard.Juara2, nil
	}
	if leaderboard.Juara3 != nil && leaderboard.Juara3.User.ID == userID {
		return leaderboard.Juara3, nil
	}

	// Find user in remaining leaderboard
	for _, entry := range leaderboard.Leaderboard {
		if entry.User.ID == userID {
			return &entry, nil
		}
	}

	// User not found in leaderboard
	user, err := s.UserRepo.GetUserById(int(userID))
	if err != nil {
		return nil, err
	}

	isCurrentUser := false
	if currentUserID != nil && *currentUserID == userID {
		isCurrentUser = true
	}

	return &LeaderboardEntry{
		User:          user,
		Score:         0,
		Rank:          -1, // Not ranked
		IsCurrentUser: isCurrentUser,
	}, nil
}

func (s *LeaderboardService) UpdateUserScore(userID uint, moduleID uint, quizType string, score int) error {
	// This method can be used to update user scores after quiz completion
	// For now, we'll work with participant scores directly
	return nil
}
