package votesession

// EligibleVoter is a snapshot of one eligible voter at session creation time.
type EligibleVoter struct {
	ID   string
	Name string
}

// Status represents the state of a VoteSession.
type Status string

const (
	StatusOpen   Status = "Open"
	StatusClosed Status = "Closed"
)

// VoteSession stores ActID, status, and an immutable snapshot of eligible voters.
type VoteSession struct {
	actID          string
	status         Status
	eligibleVoters []EligibleVoter
}

// NewVoteSession creates a session for the given act ID with status Open.
// The eligibleVoters slice is copied; the session stores an immutable snapshot.
func NewVoteSession(actID string, eligibleVoters []EligibleVoter) *VoteSession {
	var voters []EligibleVoter
	if len(eligibleVoters) > 0 {
		voters = make([]EligibleVoter, len(eligibleVoters))
		copy(voters, eligibleVoters)
	}
	return &VoteSession{
		actID:          actID,
		status:         StatusOpen,
		eligibleVoters: voters,
	}
}

// ActID returns the Act ID this session is for.
func (vs *VoteSession) ActID() string {
	return vs.actID
}

// Status returns the current status of the session.
func (vs *VoteSession) Status() Status {
	return vs.status
}

// EligibleVoters returns a copy of the eligible voters snapshot.
func (vs *VoteSession) EligibleVoters() []EligibleVoter {
	if len(vs.eligibleVoters) == 0 {
		return nil
	}
	out := make([]EligibleVoter, len(vs.eligibleVoters))
	copy(out, vs.eligibleVoters)
	return out
}

// Close sets the session status to Closed.
func (vs *VoteSession) Close() {
	vs.status = StatusClosed
}
