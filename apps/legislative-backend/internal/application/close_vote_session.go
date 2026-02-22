package application

import (
	"errors"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

// ErrVoteSessionNotOpen is returned when closing a session that is not Open.
var ErrVoteSessionNotOpen = errors.New("vote session is not open")

// CloseVoteSessionService closes a vote session and computes the result using the configured policy.
type CloseVoteSessionService struct {
	SessionRepo votesession.VoteSessionRepository
	Policy      votesession.OutcomePolicy
}

// CloseVoteSession loads the session for the act, closes it with the policy, saves, and returns the result.
// Returns ErrVoteSessionNotFound if no session exists for the act, or ErrVoteSessionNotOpen if the session is not Open.
func (s *CloseVoteSessionService) CloseVoteSession(actID string) (*votesession.VotingResult, error) {
	vs, err := s.SessionRepo.GetByActID(actID)
	if err != nil {
		return nil, err
	}
	if vs == nil {
		return nil, votesession.ErrVoteSessionNotFound
	}
	if vs.Status() != votesession.StatusOpen {
		return nil, ErrVoteSessionNotOpen
	}
	if err := vs.CloseWithResult(s.Policy); err != nil {
		return nil, err
	}
	if err := s.SessionRepo.Save(vs); err != nil {
		return nil, err
	}
	return vs.Result(), nil
}
