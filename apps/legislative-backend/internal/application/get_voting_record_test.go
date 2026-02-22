package application

import (
	"testing"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

func TestGetVotingRecord_ReturnsFinalVotesAndResultWhenClosed(t *testing.T) {
	actID := "act-1"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}, {ID: "m2", Name: "B"}})
	_ = vs.CastVote("m1", votesession.VoteYes)
	_ = vs.CastVote("m2", votesession.VoteNo)
	_ = vs.CloseWithResult(votesession.SimpleMajorityPolicy{})

	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
	}

	svc := &GetVotingRecordService{SessionRepo: sessionRepo}

	rec, err := svc.GetVotingRecord(actID)

	if err != nil {
		t.Fatalf("GetVotingRecord: %v", err)
	}
	if rec == nil {
		t.Fatal("record = nil, want non-nil")
	}
	if rec.FinalVotes["m1"] != votesession.VoteYes || rec.FinalVotes["m2"] != votesession.VoteNo {
		t.Errorf("FinalVotes = %v, want m1=Yes m2=No", rec.FinalVotes)
	}
	if rec.Result == nil {
		t.Fatal("Result = nil, want non-nil")
	}
	if rec.Result.YesCount != 1 || rec.Result.NoCount != 1 || rec.Result.AbstainCount != 0 {
		t.Errorf("Result counts = Yes=%d No=%d Abstain=%d, want 1, 1, 0", rec.Result.YesCount, rec.Result.NoCount, rec.Result.AbstainCount)
	}
}

func TestGetVotingRecord_WhenSessionNotFound_ReturnsError(t *testing.T) {
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) {
			return nil, votesession.ErrVoteSessionNotFound
		},
	}

	svc := &GetVotingRecordService{SessionRepo: sessionRepo}

	rec, err := svc.GetVotingRecord("act-none")

	if err != votesession.ErrVoteSessionNotFound {
		t.Errorf("GetVotingRecord err = %v, want ErrVoteSessionNotFound", err)
	}
	if rec != nil {
		t.Errorf("record = %v, want nil", rec)
	}
}

func TestGetVotingRecord_WhenSessionOpen_ReturnsError(t *testing.T) {
	actID := "act-1"
	vs := votesession.NewVoteSession(actID, []votesession.EligibleVoter{{ID: "m1", Name: "A"}})
	_ = vs.CastVote("m1", votesession.VoteYes)

	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(id string) (*votesession.VoteSession, error) {
			if id == actID {
				return vs, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
	}

	svc := &GetVotingRecordService{SessionRepo: sessionRepo}

	rec, err := svc.GetVotingRecord(actID)

	if err != ErrVotingRecordNotAvailable {
		t.Errorf("GetVotingRecord err = %v, want ErrVotingRecordNotAvailable", err)
	}
	if rec != nil {
		t.Errorf("record = %v, want nil", rec)
	}
}
