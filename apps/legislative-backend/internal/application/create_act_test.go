package application

import (
	"testing"

	"github.com/wleklinskimateusz/legislation/backend/internal/domain/act"
)

type fakeActRepository struct {
	getByID func(id string) (*act.Act, error)
	save    func(a *act.Act) error
}

func (f *fakeActRepository) GetByID(id string) (*act.Act, error) {
	if f.getByID != nil {
		return f.getByID(id)
	}
	return nil, act.ErrActNotFound
}

func (f *fakeActRepository) Save(a *act.Act) error {
	if f.save != nil {
		return f.save(a)
	}
	return nil
}

func TestCreateAct_SavesNewDraftAct(t *testing.T) {
	var saveCalledWith *act.Act
	actRepo := &fakeActRepository{
		save: func(a *act.Act) error {
			saveCalledWith = a
			return nil
		},
	}
	svc := &CreateActService{ActRepo: actRepo}

	err := svc.CreateAct("act-001")

	if err != nil {
		t.Fatalf("CreateAct: %v", err)
	}
	if saveCalledWith == nil {
		t.Fatal("Save was not called")
	}
	if saveCalledWith.ID() != "act-001" {
		t.Errorf("Save called with Act ID %q, want act-001", saveCalledWith.ID())
	}
	if saveCalledWith.Status() != act.StatusDraft {
		t.Errorf("Save called with Act status %v, want Draft", saveCalledWith.Status())
	}
}
