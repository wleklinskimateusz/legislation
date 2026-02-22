package application

import (
	"testing"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

func TestCloseVoteSession_ClosesAndSavesSessionWithResult(t *testing.T) {
	actID := "act-1"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}})
	_ = vs.CastVote("m1", votesession.VoteYes)
	_ = vs.CastVote("m2", votesession.VoteNo)

	var saveCalledWith *votesession.VoteSession
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(s *votesession.VoteSession) error {
			saveCalledWith = s
			return nil
		},
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}}

	result, err := svc.CloseVoteSession(actID)

	if err != nil {
		t.Fatalf("CloseVoteSession: %v", err)
	}
	if saveCalledWith == nil {
		t.Fatal("Save was not called on session repo")
	}
	if saveCalledWith.Status() != votesession.StatusClosed {
		t.Errorf("saved session Status() = %v, want Closed", saveCalledWith.Status())
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if result.YesCount != 1 || result.NoCount != 1 || result.AbstainCount != 0 {
		t.Errorf("result counts = Yes=%d No=%d Abstain=%d, want 1, 1, 0", result.YesCount, result.NoCount, result.AbstainCount)
	}
	if result.Passed {
		t.Error("result.Passed = true, want false (1 Yes, 1 No => tie, simple majority)")
	}
}

func TestCloseVoteSession_WhenSessionNotFound_ReturnsError(t *testing.T) {
	saveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) {
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { saveCalled = true; return nil },
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}}

	result, err := svc.CloseVoteSession("act-none")

	if err != votesession.ErrVoteSessionNotFound {
		t.Errorf("CloseVoteSession err = %v, want ErrVoteSessionNotFound", err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
	if saveCalled {
		t.Error("Save must not be called when session not found")
	}
}

func TestCloseVoteSession_WhenSessionAlreadyClosed_ReturnsError(t *testing.T) {
	actID := "act-1"
	vs := votesession.NewVoteSession(actID, nil)
	_ = vs.CloseWithResult(votesession.SimpleMajorityPolicy{})

	saveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { saveCalled = true; return nil },
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}}

	result, err := svc.CloseVoteSession(actID)

	if err != ErrVoteSessionNotOpen {
		t.Errorf("CloseVoteSession err = %v, want ErrVoteSessionNotOpen", err)
	}
	if result != nil {
		t.Errorf("result = %v, want nil", result)
	}
	if saveCalled {
		t.Error("Save must not be called when session already closed")
	}
}

// When vote passes, application loads the Act, calls Accept(), and saves it.
func TestCloseVoteSession_WhenVotePasses_AcceptsActAndSaves(t *testing.T) {
	actID := "act-pass"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}, {ID: "m3", Name: "C"}})
	_ = vs.CastVote("m1", votesession.VoteYes)
	_ = vs.CastVote("m2", votesession.VoteYes)
	_ = vs.CastVote("m3", votesession.VoteNo)

	a := act.NewDraftAct(actID)
	_ = a.StartVoting()

	var saveActCalledWith *act.Act
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { return nil },
	}
	actRepo := &fakeActRepository{
		getByID: func(id string) (*act.Act, error) {
			if id == actID {
				return a, nil
			}
			return nil, act.ErrActNotFound
		},
		save: func(saved *act.Act) error {
			saveActCalledWith = saved
			return nil
		},
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}, ActRepo: actRepo}

	result, err := svc.CloseVoteSession(actID)

	if err != nil {
		t.Fatalf("CloseVoteSession: %v", err)
	}
	if result == nil || !result.Passed {
		t.Fatalf("result.Passed = false, want true")
	}
	if saveActCalledWith == nil {
		t.Fatal("ActRepo.Save was not called")
	}
	if saveActCalledWith.Status() != act.StatusAccepted {
		t.Errorf("saved act Status() = %v, want Accepted", saveActCalledWith.Status())
	}
}

// When vote does not pass, Act is not modified (ActRepo.Save for act not called).
func TestCloseVoteSession_WhenVoteFails_DoesNotAcceptAct(t *testing.T) {
	actID := "act-fail"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}, {ID: "m3", Name: "C"}})
	_ = vs.CastVote("m1", votesession.VoteYes)
	_ = vs.CastVote("m2", votesession.VoteNo)
	_ = vs.CastVote("m3", votesession.VoteNo)

	actSaveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { return nil },
	}
	actRepo := &fakeActRepository{
		save: func(*act.Act) error { actSaveCalled = true; return nil },
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}, ActRepo: actRepo}

	_, err := svc.CloseVoteSession(actID)

	if err != nil {
		t.Fatalf("CloseVoteSession: %v", err)
	}
	if actSaveCalled {
		t.Error("ActRepo.Save must not be called when vote does not pass")
	}
}

// When vote passes but Act is not in Voting (e.g. Draft), Accept() fails and error is returned; session was already saved.
func TestCloseVoteSession_WhenVotePassesButActNotInVoting_ReturnsError(t *testing.T) {
	actID := "act-draft"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}})
	_ = vs.CastVote("m1", votesession.VoteYes)
	_ = vs.CastVote("m2", votesession.VoteYes)

	a := act.NewDraftAct(actID)
	actSaveCalled := false
	sessionSaveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { sessionSaveCalled = true; return nil },
	}
	actRepo := &fakeActRepository{
		getByID: func(id string) (*act.Act, error) {
			if id == actID {
				return a, nil
			}
			return nil, act.ErrActNotFound
		},
		save: func(*act.Act) error { actSaveCalled = true; return nil },
	}

	svc := &CloseVoteSessionService{SessionRepo: sessionRepo, Policy: votesession.SimpleMajorityPolicy{}, ActRepo: actRepo}

	result, err := svc.CloseVoteSession(actID)

	if err == nil {
		t.Fatal("CloseVoteSession: want error when act not in Voting, got nil")
	}
	if result != nil {
		t.Errorf("result = %v, want nil when error", result)
	}
	if !sessionSaveCalled {
		t.Error("session was not saved; vote result should be persisted before accepting act")
	}
	if actSaveCalled {
		t.Error("ActRepo.Save must not be called when Accept() fails")
	}
}
