package act

import (
	"testing"
)

// An Act in Draft status can add a paragraph and stores it internally.
func TestActInDraft_AddParagraph_StoresItInternally(t *testing.T) {
	// Act must have identity
	id := "act-001"
	a := NewDraftAct(id)

	// Act must start in Draft status
	if a.Status() != StatusDraft {
		t.Fatalf("new act status = %v, want Draft", a.Status())
	}

	// Adding paragraph must increase paragraph count
	before := a.ParagraphCount()
	a.AddParagraph("First paragraph text")
	after := a.ParagraphCount()

	if after != before+1 {
		t.Errorf("paragraph count = %d, want %d", after, before+1)
	}
}

// An Act in Draft status can transition to Voting status using StartVoting().
func TestActInDraft_StartVoting_ChangesStatusToVoting(t *testing.T) {
	a := NewDraftAct("act-002")

	// Act must start in Draft
	if a.Status() != StatusDraft {
		t.Fatalf("new act status = %v, want Draft", a.Status())
	}

	// Calling StartVoting() changes status to Voting
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting from Draft: %v", err)
	}

	if a.Status() != StatusVoting {
		t.Errorf("after StartVoting status = %v, want Voting", a.Status())
	}
}

// An Act in Voting status can transition to Accepted using Accept().
func TestActInVoting_Accept_ChangesStatusToAccepted(t *testing.T) {
	a := NewDraftAct("act-003")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if a.Status() != StatusVoting {
		t.Fatalf("act status = %v, want Voting", a.Status())
	}

	if err := a.Accept(); err != nil {
		t.Fatalf("Accept from Voting: %v", err)
	}

	if a.Status() != StatusAccepted {
		t.Errorf("after Accept status = %v, want Accepted", a.Status())
	}
}

// An Act in Accepted status can transition to Published using Publish().
func TestActInAccepted_Publish_ChangesStatusToPublished(t *testing.T) {
	a := NewDraftAct("act-004")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if err := a.Accept(); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if a.Status() != StatusAccepted {
		t.Fatalf("act status = %v, want Accepted", a.Status())
	}

	if err := a.Publish(); err != nil {
		t.Fatalf("Publish from Accepted: %v", err)
	}

	if a.Status() != StatusPublished {
		t.Errorf("after Publish status = %v, want Published", a.Status())
	}
}

// StartVoting() from non-Draft is rejected; status remains unchanged.
func TestActInVoting_StartVoting_ReturnsErrorAndStatusUnchanged(t *testing.T) {
	a := NewDraftAct("act-005")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if a.Status() != StatusVoting {
		t.Fatalf("act status = %v, want Voting", a.Status())
	}

	err := a.StartVoting()

	if err == nil {
		t.Error("StartVoting from Voting: want error, got nil")
	}
	if a.Status() != StatusVoting {
		t.Errorf("status changed to %v, want Voting unchanged", a.Status())
	}
}

// Accept() from non-Voting is rejected; status remains unchanged.
func TestActInDraft_Accept_ReturnsErrorAndStatusUnchanged(t *testing.T) {
	a := NewDraftAct("act-006")
	if a.Status() != StatusDraft {
		t.Fatalf("act status = %v, want Draft", a.Status())
	}

	err := a.Accept()

	if err == nil {
		t.Error("Accept from Draft: want error, got nil")
	}
	if a.Status() != StatusDraft {
		t.Errorf("status changed to %v, want Draft unchanged", a.Status())
	}
}

// Publish() from non-Accepted is rejected; status remains unchanged.
func TestActInDraft_Publish_ReturnsErrorAndStatusUnchanged(t *testing.T) {
	a := NewDraftAct("act-007")
	if a.Status() != StatusDraft {
		t.Fatalf("act status = %v, want Draft", a.Status())
	}

	err := a.Publish()

	if err == nil {
		t.Error("Publish from Draft: want error, got nil")
	}
	if a.Status() != StatusDraft {
		t.Errorf("status changed to %v, want Draft unchanged", a.Status())
	}
}

// The aggregate exposes its identity via ID().
func TestNewDraftAct_ExposesIdentityViaID(t *testing.T) {
	a := NewDraftAct("id-123")

	if got := a.ID(); got != "id-123" {
		t.Errorf("ID() = %q, want id-123", got)
	}
}
