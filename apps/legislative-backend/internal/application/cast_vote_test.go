package application

import (
	"testing"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

func TestCastVote_SavesSessionWithVote(t *testing.T) {
	actID := "act-1"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "Alice"}})

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

	svc := &CastVoteService{SessionRepo: sessionRepo}

	err := svc.CastVote(actID, "m1", votesession.VoteYes)

	if err != nil {
		t.Fatalf("CastVote: %v", err)
	}
	if saveCalledWith == nil {
		t.Fatal("Save was not called on session repo")
	}
	votes := saveCalledWith.Votes()
	if votes["m1"] != votesession.VoteYes {
		t.Errorf("saved session Votes()[m1] = %v, want Yes", votes["m1"])
	}
}

func TestCastVote_WhenSessionNotFound_ReturnsError(t *testing.T) {
	saveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) {
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(*votesession.VoteSession) error { saveCalled = true; return nil },
	}

	svc := &CastVoteService{SessionRepo: sessionRepo}

	err := svc.CastVote("act-none", "m1", votesession.VoteYes)

	if err != votesession.ErrVoteSessionNotFound {
		t.Errorf("CastVote: err = %v, want ErrVoteSessionNotFound", err)
	}
	if saveCalled {
		t.Error("Save must not be called when session not found")
	}
}
