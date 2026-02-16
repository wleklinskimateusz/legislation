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
	a.StartVoting()

	if a.Status() != StatusVoting {
		t.Errorf("after StartVoting status = %v, want Voting", a.Status())
	}
}
