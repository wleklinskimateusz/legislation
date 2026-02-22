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

// Eligible voter can cast a vote; it is stored and returned by Votes().
func TestVoteSession_CastVote_EligibleVoterStoresVote(t *testing.T) {
	vs := NewVoteSession("act-vote", []EligibleVoter{{ID: "m1", Name: "Alice"}})

	err := vs.CastVote("m1", VoteYes)

	if err != nil {
		t.Fatalf("CastVote: %v", err)
	}
	votes := vs.Votes()
	if votes["m1"] != VoteYes {
		t.Errorf("Votes()[m1] = %v, want Yes", votes["m1"])
	}
}

// Non-eligible voter cannot cast a vote; CastVote returns error and vote is not stored.
func TestVoteSession_CastVote_NonEligibleReturnsError(t *testing.T) {
	vs := NewVoteSession("act-vote", []EligibleVoter{{ID: "m1", Name: "Alice"}})

	err := vs.CastVote("m99", VoteYes)

	if err == nil {
		t.Fatal("CastVote(m99): want error, got nil")
	}
	votes := vs.Votes()
	if votes != nil {
		if _, has := votes["m99"]; has {
			t.Errorf("Votes() must not contain m99, got %v", votes)
		}
	}
}

// Voting is only allowed while session is Open; after Close, CastVote returns error and Votes() stays empty.
func TestVoteSession_CastVote_ClosedSessionReturnsError(t *testing.T) {
	vs := NewVoteSession("act-vote", []EligibleVoter{{ID: "m1", Name: "Alice"}})
	vs.Close()

	err := vs.CastVote("m1", VoteYes)

	if err == nil {
		t.Fatal("CastVote after Close: want error, got nil")
	}
	votes := vs.Votes()
	if votes != nil && len(votes) != 0 {
		t.Errorf("Votes() = %v, want empty", votes)
	}
}

// Re-voting overwrites the previous vote; one active vote per voter.
func TestVoteSession_CastVote_RevoteOverwrites(t *testing.T) {
	vs := NewVoteSession("act-vote", []EligibleVoter{{ID: "m1", Name: "Alice"}})

	_ = vs.CastVote("m1", VoteYes)
	err := vs.CastVote("m1", VoteNo)

	if err != nil {
		t.Fatalf("second CastVote: %v", err)
	}
	votes := vs.Votes()
	if votes["m1"] != VoteNo {
		t.Errorf("Votes()[m1] = %v, want No (overwritten)", votes["m1"])
	}
}

// Yes, No, and Abstain are all valid choices; three voters can cast each.
func TestVoteSession_CastVote_YesNoAbstainAllValid(t *testing.T) {
	vs := NewVoteSession("act-vote", []EligibleVoter{
		{ID: "m1", Name: "Alice"},
		{ID: "m2", Name: "Bob"},
		{ID: "m3", Name: "Carol"},
	})

	_ = vs.CastVote("m1", VoteYes)
	_ = vs.CastVote("m2", VoteNo)
	_ = vs.CastVote("m3", VoteAbstain)

	votes := vs.Votes()
	if votes["m1"] != VoteYes || votes["m2"] != VoteNo || votes["m3"] != VoteAbstain {
		t.Errorf("Votes() = %v, want m1=Yes m2=No m3=Abstain", votes)
	}
}
