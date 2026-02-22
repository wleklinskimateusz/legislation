package votesession

import "errors"

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

// VoteChoice represents a vote option.
type VoteChoice string

const (
	VoteYes     VoteChoice = "Yes"
	VoteNo      VoteChoice = "No"
	VoteAbstain VoteChoice = "Abstain"
)

// VotingResult holds the tally and outcome of a closed vote session.
type VotingResult struct {
	YesCount     int
	NoCount      int
	AbstainCount int
	Passed       bool
}

// CountVotes returns the number of Yes, No, and Abstain votes from a vote map.
func CountVotes(votes map[string]VoteChoice) (yes, no, abstain int) {
	for _, c := range votes {
		switch c {
		case VoteYes:
			yes++
		case VoteNo:
			no++
		case VoteAbstain:
			abstain++
		}
	}
	return yes, no, abstain
}

// OutcomePolicy determines whether a vote passes given the counts.
type OutcomePolicy interface {
	Passed(yes, no, abstain int) bool
}

// SimpleMajorityPolicy passes when Yes count exceeds No count (abstentions ignored).
type SimpleMajorityPolicy struct{}

// Passed returns true when yes > no.
func (SimpleMajorityPolicy) Passed(yes, no, abstain int) bool {
	return yes > no
}

// VoteSession stores ActID, status, and an immutable snapshot of eligible voters.
type VoteSession struct {
	actID          string
	status         Status
	eligibleVoters []EligibleVoter
	votes          map[string]VoteChoice
	result         *VotingResult
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
		votes:          make(map[string]VoteChoice),
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

// CloseWithResult closes the session and computes and stores the VotingResult using the given policy.
// Returns error if the session is not Open.
func (vs *VoteSession) CloseWithResult(policy OutcomePolicy) error {
	if vs.status != StatusOpen {
		return errors.New("vote session is not open")
	}
	yes, no, abstain := CountVotes(vs.votes)
	vs.status = StatusClosed
	vs.result = &VotingResult{
		YesCount:     yes,
		NoCount:      no,
		AbstainCount: abstain,
		Passed:       policy.Passed(yes, no, abstain),
	}
	return nil
}

// Result returns the voting result after the session has been closed with CloseWithResult; nil otherwise.
func (vs *VoteSession) Result() *VotingResult {
	return vs.result
}

// CastVote records or overwrites the vote for the given voter. Only eligible voters can vote, and only while the session is Open.
func (vs *VoteSession) CastVote(voterID string, choice VoteChoice) error {
	if vs.status != StatusOpen {
		return errors.New("voting only allowed while session is open")
	}
	var eligible bool
	for _, v := range vs.eligibleVoters {
		if v.ID == voterID {
			eligible = true
			break
		}
	}
	if !eligible {
		return errors.New("voter not eligible")
	}
	vs.votes[voterID] = choice
	return nil
}

// Votes returns a copy of the current votes (voterID -> choice).
func (vs *VoteSession) Votes() map[string]VoteChoice {
	if len(vs.votes) == 0 {
		return nil
	}
	out := make(map[string]VoteChoice, len(vs.votes))
	for k, v := range vs.votes {
		out[k] = v
	}
	return out
}
