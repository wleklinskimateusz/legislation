package application

import (
	"testing"

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
