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

// CountVotes returns correct yes, no, abstain counts from a vote map.
func TestCountVotes_ReturnsCorrectCounts(t *testing.T) {
	votes := map[string]VoteChoice{
		"a": VoteYes, "b": VoteYes, "c": VoteNo, "d": VoteAbstain, "e": VoteAbstain,
	}
	yes, no, abstain := CountVotes(votes)
	if yes != 2 || no != 1 || abstain != 2 {
		t.Errorf("CountVotes() = yes=%d no=%d abstain=%d, want 2, 1, 2", yes, no, abstain)
	}
}

func TestCountVotes_EmptyMapReturnsZeros(t *testing.T) {
	yes, no, abstain := CountVotes(nil)
	if yes != 0 || no != 0 || abstain != 0 {
		t.Errorf("CountVotes(nil) = yes=%d no=%d abstain=%d, want 0, 0, 0", yes, no, abstain)
	}
}

// SimpleMajorityPolicy: Passed is true when Yes > No.
func TestSimpleMajorityPolicy_PassedWhenYesGreaterThanNo(t *testing.T) {
	p := SimpleMajorityPolicy{}
	if !p.Passed(3, 2, 0) {
		t.Error("Passed(3, 2, 0): want true, got false")
	}
	if !p.Passed(1, 0, 10) {
		t.Error("Passed(1, 0, 10): want true (abstain ignored), got false")
	}
}

// SimpleMajorityPolicy: Passed is false when Yes <= No.
func TestSimpleMajorityPolicy_NotPassedWhenYesLessOrEqualNo(t *testing.T) {
	p := SimpleMajorityPolicy{}
	if p.Passed(2, 3, 0) {
		t.Error("Passed(2, 3, 0): want false, got true")
	}
	if p.Passed(2, 2, 1) {
		t.Error("Passed(2, 2, 1): want false (tie), got true")
	}
}

// CloseWithResult on an open session sets status Closed and stores result with correct counts and Passed.
func TestVoteSession_CloseWithResult_SetsClosedAndResult(t *testing.T) {
	vs := NewVoteSession("act-1", []EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}, {ID: "m3", Name: "C"}})
	_ = vs.CastVote("m1", VoteYes)
	_ = vs.CastVote("m2", VoteYes)
	_ = vs.CastVote("m3", VoteNo)
	policy := SimpleMajorityPolicy{}

	err := vs.CloseWithResult(policy)

	if err != nil {
		t.Fatalf("CloseWithResult: %v", err)
	}
	if vs.Status() != StatusClosed {
		t.Errorf("Status() = %v, want Closed", vs.Status())
	}
	res := vs.Result()
	if res == nil {
		t.Fatal("Result() = nil, want non-nil")
	}
	if res.YesCount != 2 || res.NoCount != 1 || res.AbstainCount != 0 {
		t.Errorf("Result() counts = Yes=%d No=%d Abstain=%d, want 2, 1, 0", res.YesCount, res.NoCount, res.AbstainCount)
	}
	if !res.Passed {
		t.Error("Result().Passed = false, want true (simple majority)")
	}
}

// CloseWithResult when session already closed returns error.
func TestVoteSession_CloseWithResult_WhenAlreadyClosedReturnsError(t *testing.T) {
	vs := NewVoteSession("act-1", nil)
	vs.Close()

	err := vs.CloseWithResult(SimpleMajorityPolicy{})

	if err == nil {
		t.Fatal("CloseWithResult on closed session: want error, got nil")
	}
}

// Result returns nil when session has not been closed with result.
func TestVoteSession_Result_NilWhenOpen(t *testing.T) {
	vs := NewVoteSession("act-1", nil)

	if vs.Result() != nil {
		t.Errorf("Result() = %v, want nil when open", vs.Result())
	}
}
