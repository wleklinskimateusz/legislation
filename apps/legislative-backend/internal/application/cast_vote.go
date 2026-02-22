package application

import (
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

// CastVoteService records a vote for an act's vote session.
type CastVoteService struct {
	SessionRepo votesession.VoteSessionRepository
}

// CastVote loads the vote session for the act, records the voter's choice, and saves the session.
// Returns ErrVoteSessionNotFound if no session exists for the act.
func (s *CastVoteService) CastVote(actID, voterID string, choice votesession.VoteChoice) error {
	vs, err := s.SessionRepo.GetByActID(actID)
	if err != nil {
		return err
	}
	if vs == nil {
		return votesession.ErrVoteSessionNotFound
	}
	if err := vs.CastVote(voterID, choice); err != nil {
		return err
	}
	return s.SessionRepo.Save(vs)
}
