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

	// Convert to slice and sort by TotalXP (from User.TotalXP field)
	var entries []LeaderboardEntry
	for _, user := range allUsers {
		isCurrentUser := false
		if currentUserID != nil && *currentUserID == user.ID {
			isCurrentUser = true
		}

		entries = append(entries, LeaderboardEntry{
			User:          user,
			Score:         int64(user.TotalXP), // Use TotalXP field from User
			IsCurrentUser: isCurrentUser,
		})
	}

	// Sort by TotalXP descending
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

// UpdateUserPosition calculates and updates position change for a specific user
func (s *LeaderboardService) UpdateUserPosition(userID uint) error {
	// Get current leaderboard to find user's current position
	leaderboard, err := s.GetLeaderboard(nil, "", &userID)
	if err != nil {
		return err
	}

	// Find user's current position in leaderboard
	currentPosition := 0
	if leaderboard.Juara1 != nil && leaderboard.Juara1.User.ID == userID {
		currentPosition = 1
	} else if leaderboard.Juara2 != nil && leaderboard.Juara2.User.ID == userID {
		currentPosition = 2
	} else if leaderboard.Juara3 != nil && leaderboard.Juara3.User.ID == userID {
		currentPosition = 3
	} else {
		// Check in remaining leaderboard entries
		for _, entry := range leaderboard.Leaderboard {
			if entry.User.ID == userID {
				currentPosition = entry.Rank
				break
			}
		}
	}

	if currentPosition == 0 {
		// User not found in ranked positions, get total user count and assume they're at the bottom
		allUsers, err := s.UserRepo.GetAllUsers()
		if err == nil {
			currentPosition = len(allUsers) // Last position
		} else {
			currentPosition = 100 // Default fallback
		}
	}

	// Get user data to check previous position (stored in database)
	user, err := s.UserRepo.GetUserById(int(userID))
	if err != nil {
		return err
	}

	// Get the stored previous position 
	previousPosition := user.PreviousPosition
	if previousPosition == 0 {
		// First time tracking - set previous position to current for future tracking
		// and assume they moved up for this display
		previousPosition = currentPosition + 2 // Assume they moved up 2 positions
	}

	// Calculate position change
	positionType := "stable"
	changeAmount := 0
	
	if currentPosition < previousPosition {
		positionType = "increasing" // Lower number = higher rank = increasing
		changeAmount = previousPosition - currentPosition
	} else if currentPosition > previousPosition {
		positionType = "decreasing" // Higher number = lower rank = decreasing
		changeAmount = currentPosition - previousPosition
	}

	// Update user position fields
	user.PositionType = positionType
	user.PositionChange = changeAmount
	user.PreviousPosition = currentPosition // Store current position for next time

	// Save updated user
	err = s.UserRepo.UpdateUser(&user)
	return err
}
