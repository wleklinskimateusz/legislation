package votesession

import (
	"testing"
)

// VoteSession can be created for an Act ID and starts in Open state.
func TestNewVoteSession_HasActIDAndStartsOpen(t *testing.T) {
	vs := NewVoteSession("act-1", nil)

	if got := vs.ActID(); got != "act-1" {
		t.Errorf("ActID() = %q, want act-1", got)
	}
	if got := vs.Status(); got != StatusOpen {
		t.Errorf("Status() = %v, want Open", got)
	}
	ev := vs.EligibleVoters()
	if len(ev) != 0 {
		t.Errorf("EligibleVoters() length = %d, want 0", len(ev))
	}
}

// VoteSession stores EligibleVoters snapshot; it is immutable after creation.
func TestNewVoteSession_StoresEligibleVotersSnapshotImmutable(t *testing.T) {
	voters := []EligibleVoter{{ID: "m1", Name: "Alice"}}
	vs := NewVoteSession("act-2", voters)

	ev := vs.EligibleVoters()
	if len(ev) != 1 {
		t.Fatalf("EligibleVoters() length = %d, want 1", len(ev))
	}
	if ev[0].ID != "m1" || ev[0].Name != "Alice" {
		t.Errorf("EligibleVoters()[0] = %+v, want ID m1 Name Alice", ev[0])
	}

	// Mutate the slice passed to NewVoteSession; session snapshot must be unchanged.
	voters[0].Name = "Mutated"
	ev2 := vs.EligibleVoters()
	if len(ev2) != 1 || ev2[0].Name != "Alice" {
		t.Errorf("after mutating caller slice, EligibleVoters()[0].Name = %q, want Alice (unchanged)", ev2[0].Name)
	}
}

// VoteSession can be closed; status becomes Closed.
func TestVoteSession_Close_SetsStatusClosed(t *testing.T) {
	vs := NewVoteSession("act-3", nil)
	if vs.Status() != StatusOpen {
		t.Fatalf("Status() = %v, want Open", vs.Status())
	}

	vs.Close()

	if vs.Status() != StatusClosed {
		t.Errorf("after Close Status() = %v, want Closed", vs.Status())
	}
}

// Domain allows multiple VoteSessions for the same ActID; "only one per Act" is enforced by the application (e.g. repository).
func TestNewVoteSession_DomainAllowsMultipleSessionsPerActID(t *testing.T) {
	vs1 := NewVoteSession("act-same", nil)
	vs2 := NewVoteSession("act-same", nil)

	if vs1.ActID() != "act-same" || vs2.ActID() != "act-same" {
		t.Errorf("ActID: vs1=%q vs2=%q, both want act-same", vs1.ActID(), vs2.ActID())
	}
	if vs1 == vs2 {
		t.Error("two sessions for same ActID are distinct instances; duplicate prevention is application responsibility")
	}
}
