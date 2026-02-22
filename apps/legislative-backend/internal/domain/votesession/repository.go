package votesession

import "errors"

// ErrVoteSessionNotFound is returned when a VoteSession is not found by Act ID.
var ErrVoteSessionNotFound = errors.New("vote session not found")

// VoteSessionRepository persists and retrieves VoteSessions.
type VoteSessionRepository interface {
	GetByActID(actID string) (*VoteSession, error)
	Save(vs *VoteSession) error
}
