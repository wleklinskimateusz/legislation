package application

import (
	"errors"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

// ErrVotingRecordNotAvailable is returned when the voting record cannot be exposed (e.g. session not closed).
var ErrVotingRecordNotAvailable = errors.New("voting record not available: session not closed")

// VotingRecord is the public record of a closed vote session (per-member final votes and aggregated result).
type VotingRecord struct {
	FinalVotes map[string]votesession.VoteChoice
	Result     *votesession.VotingResult
}

// GetVotingRecordService returns the public voting record for a vote session.
type GetVotingRecordService struct {
	SessionRepo votesession.VoteSessionRepository
}

// GetVotingRecord loads the session for the act and returns the public record (final votes and result) if the session is closed.
// Returns ErrVoteSessionNotFound if no session exists for the act, or ErrVotingRecordNotAvailable if the session is not closed.
func (s *GetVotingRecordService) GetVotingRecord(actID string) (*VotingRecord, error) {
	vs, err := s.SessionRepo.GetByActID(actID)
	if err != nil {
		return nil, err
	}
	if vs == nil {
		return nil, votesession.ErrVoteSessionNotFound
	}
	if vs.Status() != votesession.StatusClosed {
		return nil, ErrVotingRecordNotAvailable
	}
	return &VotingRecord{
		FinalVotes: vs.FinalVotes(),
		Result:     vs.Result(),
	}, nil
}
