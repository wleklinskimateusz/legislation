package application

import (
	"errors"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

// ErrVoteSessionNotOpen is returned when closing a session that is not Open.
var ErrVoteSessionNotOpen = errors.New("vote session is not open")

// CloseVoteSessionService closes a vote session and computes the result using the configured policy.
type CloseVoteSessionService struct {
	SessionRepo votesession.VoteSessionRepository
	Policy      votesession.OutcomePolicy
	ActRepo     act.ActRepository
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
	result := vs.Result()
	if result != nil && result.Passed && s.ActRepo != nil {
		if err := s.acceptActForPassedVote(actID); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// acceptActForPassedVote loads the act, calls Accept(), and saves it. Call only when result.Passed and ActRepo is set.
func (s *CloseVoteSessionService) acceptActForPassedVote(actID string) error {
	a, err := s.ActRepo.GetByID(actID)
	if err != nil {
		return err
	}
	if a == nil {
		return act.ErrActNotFound
	}
	if err := a.Accept(); err != nil {
		return err
	}
	return s.ActRepo.Save(a)
}
