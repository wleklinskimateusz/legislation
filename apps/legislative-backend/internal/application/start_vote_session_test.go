package application

import (
	"testing"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/legislativebody"
	"github.com/wleklinskimateusz/legislation/backend/internal/domain/votesession"
)

type fakeLegislativeBodyRepository struct {
	getByID func(id string) (*legislativebody.LegislativeBody, error)
	save    func(lb *legislativebody.LegislativeBody) error
}

func (f *fakeLegislativeBodyRepository) GetByID(id string) (*legislativebody.LegislativeBody, error) {
	if f.getByID != nil {
		return f.getByID(id)
	}
	return nil, legislativebody.ErrLegislativeBodyNotFound
}

func (f *fakeLegislativeBodyRepository) Save(lb *legislativebody.LegislativeBody) error {
	if f.save != nil {
		return f.save(lb)
	}
	return nil
}

type fakeVoteSessionRepository struct {
	getByActID func(actID string) (*votesession.VoteSession, error)
	save       func(vs *votesession.VoteSession) error
}

func (f *fakeVoteSessionRepository) GetByActID(actID string) (*votesession.VoteSession, error) {
	if f.getByActID != nil {
		return f.getByActID(actID)
	}
	return nil, votesession.ErrVoteSessionNotFound
}

func (f *fakeVoteSessionRepository) Save(vs *votesession.VoteSession) error {
	if f.save != nil {
		return f.save(vs)
	}
	return nil
}

func TestStartVoteSession_SavesSessionWithEligibleVotersFromBody(t *testing.T) {
	// Act in Voting
	a := act.NewDraftAct("act-vote")
	_ = a.StartVoting()
	actRepo := &fakeActRepository{
		getByID: func(id string) (*act.Act, error) {
			if id == "act-vote" {
				return a, nil
			}
			return nil, act.ErrActNotFound
		},
	}
	// Body with two active members
	lb := legislativebody.NewLegislativeBody("body-1")
	lb.AddMember("m1", "Alice")
	lb.AddMember("m2", "Bob")
	bodyRepo := &fakeLegislativeBodyRepository{
		getByID: func(id string) (*legislativebody.LegislativeBody, error) {
			if id == "body-1" {
				return lb, nil
			}
			return nil, legislativebody.ErrLegislativeBodyNotFound
		},
	}
	// No existing session
	var saveCalledWith *votesession.VoteSession
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) {
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(vs *votesession.VoteSession) error {
			saveCalledWith = vs
			return nil
		},
	}
	svc := &StartVoteSessionService{
		ActRepo:     actRepo,
		BodyRepo:    bodyRepo,
		SessionRepo: sessionRepo,
	}

	err := svc.StartVoteSession("act-vote", "body-1")

	if err != nil {
		t.Fatalf("StartVoteSession: %v", err)
	}
	if saveCalledWith == nil {
		t.Fatal("Save was not called on session repo")
	}
	if saveCalledWith.ActID() != "act-vote" {
		t.Errorf("saved session ActID = %q, want act-vote", saveCalledWith.ActID())
	}
	ev := saveCalledWith.EligibleVoters()
	if len(ev) != 2 {
		t.Fatalf("EligibleVoters length = %d, want 2", len(ev))
	}
	ids := map[string]bool{ev[0].ID: true, ev[1].ID: true}
	if !ids["m1"] || !ids["m2"] {
		t.Errorf("EligibleVoters IDs = %v, want m1 and m2", ids)
	}
}

func TestStartVoteSession_WhenSessionAlreadyExistsForAct_ReturnsError(t *testing.T) {
	a := act.NewDraftAct("act-dup")
	_ = a.StartVoting()
	actRepo := &fakeActRepository{getByID: func(id string) (*act.Act, error) {
		if id == "act-dup" {
			return a, nil
		}
		return nil, act.ErrActNotFound
	}}
	lb := legislativebody.NewLegislativeBody("body-1")
	bodyRepo := &fakeLegislativeBodyRepository{getByID: func(id string) (*legislativebody.LegislativeBody, error) {
		if id == "body-1" {
			return lb, nil
		}
		return nil, legislativebody.ErrLegislativeBodyNotFound
	}}
	existingSession := votesession.NewVoteSession("act-dup", nil)
	saveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) {
			if actID == "act-dup" {
				return existingSession, nil
			}
			return nil, votesession.ErrVoteSessionNotFound
		},
		save: func(vs *votesession.VoteSession) error { saveCalled = true; return nil },
	}
	svc := &StartVoteSessionService{ActRepo: actRepo, BodyRepo: bodyRepo, SessionRepo: sessionRepo}

	err := svc.StartVoteSession("act-dup", "body-1")

	if err == nil {
		t.Error("StartVoteSession when session exists: want error, got nil")
	}
	if saveCalled {
		t.Error("Save must not be called when session already exists")
	}
}

func TestStartVoteSession_WhenActNotInVoting_ReturnsError(t *testing.T) {
	// Act in Draft, not Voting
	a := act.NewDraftAct("act-draft")
	actRepo := &fakeActRepository{getByID: func(id string) (*act.Act, error) {
		if id == "act-draft" {
			return a, nil
		}
		return nil, act.ErrActNotFound
	}}
	lb := legislativebody.NewLegislativeBody("body-1")
	bodyRepo := &fakeLegislativeBodyRepository{getByID: func(id string) (*legislativebody.LegislativeBody, error) {
		if id == "body-1" {
			return lb, nil
		}
		return nil, legislativebody.ErrLegislativeBodyNotFound
	}}
	saveCalled := false
	sessionRepo := &fakeVoteSessionRepository{
		getByActID: func(actID string) (*votesession.VoteSession, error) { return nil, votesession.ErrVoteSessionNotFound },
		save:       func(vs *votesession.VoteSession) error { saveCalled = true; return nil },
	}
	svc := &StartVoteSessionService{ActRepo: actRepo, BodyRepo: bodyRepo, SessionRepo: sessionRepo}

	err := svc.StartVoteSession("act-draft", "body-1")

	if err == nil {
		t.Error("StartVoteSession when act not in Voting: want error, got nil")
	}
	if saveCalled {
		t.Error("Save must not be called when act is not in Voting")
	}
}
