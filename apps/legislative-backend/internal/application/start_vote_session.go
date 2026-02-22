package application

import (
	"errors"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/legislativebody"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

// StartVoteSessionService starts a vote session for an act using a legislative body's active members.
type StartVoteSessionService struct {
	ActRepo     act.ActRepository
	BodyRepo    legislativebody.LegislativeBodyRepository
	SessionRepo votesession.VoteSessionRepository
}

// StartVoteSession starts a VoteSession for the given act, using the given legislative body's active members as eligible voters.
// Returns error if the act is not found, not in Voting status, the body is not found, or a session already exists for this act.
func (s *StartVoteSessionService) StartVoteSession(actID, bodyID string) error {
	a, err := s.ActRepo.GetByID(actID)
	if err != nil {
		return err
	}
	if a == nil {
		return act.ErrActNotFound
	}
	if a.Status() != act.StatusVoting {
		return errors.New("act must be in Voting status to start a vote session")
	}
	existing, _ := s.SessionRepo.GetByActID(actID)
	if existing != nil {
		return errors.New("vote session already exists for this act")
	}
	lb, err := s.BodyRepo.GetByID(bodyID)
	if err != nil {
		return err
	}
	if lb == nil {
		return legislativebody.ErrLegislativeBodyNotFound
	}
	members := lb.ActiveMembers()
	eligible := make([]votesession.EligibleVoter, len(members))
	for i, m := range members {
		eligible[i] = votesession.EligibleVoter{ID: m.ID, Name: m.Name}
	}
	vs := votesession.NewVoteSession(actID, eligible)
	return s.SessionRepo.Save(vs)
}
