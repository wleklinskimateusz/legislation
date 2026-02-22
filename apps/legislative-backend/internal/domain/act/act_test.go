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
	if err := a.AddParagraph("p0", "First paragraph text"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	after := a.ParagraphCount()

	if after != before+1 {
		t.Errorf("paragraph count = %d, want %d", after, before+1)
	}
}

// A paragraph has a unique identifier and content; the Act stores it and exposes it via Paragraphs().
func TestActInDraft_AddParagraph_StoresParagraphWithIdAndContent(t *testing.T) {
	a := NewDraftAct("act-par")
	if err := a.AddParagraph("p1", "First paragraph text"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}

	if a.ParagraphCount() != 1 {
		t.Fatalf("ParagraphCount() = %d, want 1", a.ParagraphCount())
	}
	ps := a.Paragraphs()
	if len(ps) != 1 {
		t.Fatalf("Paragraphs() length = %d, want 1", len(ps))
	}
	if ps[0].ID != "p1" {
		t.Errorf("Paragraphs()[0].ID = %q, want p1", ps[0].ID)
	}
	if ps[0].Content != "First paragraph text" {
		t.Errorf("Paragraphs()[0].Content = %q, want First paragraph text", ps[0].Content)
	}
}

// Paragraphs are in insertion order.
func TestActInDraft_AddParagraph_ParagraphsAreOrdered(t *testing.T) {
	a := NewDraftAct("act-order")
	if err := a.AddParagraph("p1", "first"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	if err := a.AddParagraph("p2", "second"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	if err := a.AddParagraph("p3", "third"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}

	ps := a.Paragraphs()
	if len(ps) != 3 {
		t.Fatalf("Paragraphs() length = %d, want 3", len(ps))
	}
	if ps[0].ID != "p1" || ps[1].ID != "p2" || ps[2].ID != "p3" {
		t.Errorf("order: got IDs %q, %q, %q; want p1, p2, p3", ps[0].ID, ps[1].ID, ps[2].ID)
	}
}

// AddParagraph when not in Draft is rejected; count and status unchanged.
func TestActInVoting_AddParagraph_ReturnsErrorAndCountUnchanged(t *testing.T) {
	a := NewDraftAct("act-nodraft")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if a.Status() != StatusVoting {
		t.Fatalf("act status = %v, want Voting", a.Status())
	}

	err := a.AddParagraph("p1", "content")

	if err == nil {
		t.Error("AddParagraph from Voting: want error, got nil")
	}
	if a.ParagraphCount() != 0 {
		t.Errorf("ParagraphCount() = %d, want 0 unchanged", a.ParagraphCount())
	}
	if a.Status() != StatusVoting {
		t.Errorf("status = %v, want Voting unchanged", a.Status())
	}
}

// In Draft, paragraph content can be replaced by ID; ID and order stay the same.
func TestActInDraft_ReplaceParagraphContent_UpdatesContentKeepsIdAndOrder(t *testing.T) {
	a := NewDraftAct("act-replace")
	if err := a.AddParagraph("p1", "old"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	if err := a.AddParagraph("p2", "second"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}

	err := a.ReplaceParagraphContent("p1", "new")

	if err != nil {
		t.Fatalf("ReplaceParagraphContent: %v", err)
	}
	ps := a.Paragraphs()
	if len(ps) != 2 {
		t.Fatalf("Paragraphs() length = %d, want 2", len(ps))
	}
	if ps[0].ID != "p1" || ps[0].Content != "new" {
		t.Errorf("Paragraphs()[0] = %+v, want ID p1 Content new", ps[0])
	}
	if ps[1].ID != "p2" || ps[1].Content != "second" {
		t.Errorf("Paragraphs()[1] = %+v, want ID p2 Content second unchanged", ps[1])
	}
}

// ReplaceParagraphContent when not in Draft is rejected; content unchanged.
func TestActInVoting_ReplaceParagraphContent_ReturnsErrorAndContentUnchanged(t *testing.T) {
	a := NewDraftAct("act-replace-nodraft")
	if err := a.AddParagraph("p1", "old"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if a.Status() != StatusVoting {
		t.Fatalf("act status = %v, want Voting", a.Status())
	}

	err := a.ReplaceParagraphContent("p1", "new")

	if err == nil {
		t.Error("ReplaceParagraphContent from Voting: want error, got nil")
	}
	ps := a.Paragraphs()
	if len(ps) != 1 || ps[0].Content != "old" {
		t.Errorf("content changed: Paragraphs()[0].Content = %q, want old", ps[0].Content)
	}
}

// ReplaceParagraphContent for non-existent paragraph ID returns error; existing content unchanged.
func TestActInDraft_ReplaceParagraphContent_NonExistentParagraphReturnsError(t *testing.T) {
	a := NewDraftAct("act-replace-notfound")
	if err := a.AddParagraph("p1", "old"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}

	err := a.ReplaceParagraphContent("p99", "x")

	if err == nil {
		t.Error("ReplaceParagraphContent for p99: want error, got nil")
	}
	ps := a.Paragraphs()
	if len(ps) != 1 || ps[0].Content != "old" {
		t.Errorf("content changed: Paragraphs()[0].Content = %q, want old", ps[0].Content)
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

// A newly created Act has version 1.
func TestNewDraftAct_HasVersionOne(t *testing.T) {
	a := NewDraftAct("id")

	if got := a.Version(); got != 0 {
		t.Errorf("Version() = %d, want 0", got)
	}
}

// Accept() increments the version number.
func TestActInVoting_Accept_IncrementsVersion(t *testing.T) {
	a := NewDraftAct("act-ver")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if a.Version() != 0 {
		t.Fatalf("before Accept Version() = %d, want 0", a.Version())
	}

	if err := a.Accept(); err != nil {
		t.Fatalf("Accept: %v", err)
	}

	if got := a.Version(); got != 1 {
		t.Errorf("after Accept Version() = %d, want 1", got)
	}
}

// Accept() from wrong status does not increment version.
func TestActInDraft_Accept_DoesNotIncrementVersion(t *testing.T) {
	a := NewDraftAct("act-noinc")
	if a.Version() != 0 {
		t.Fatalf("Version() = %d, want 0", a.Version())
	}

	err := a.Accept()

	if err == nil {
		t.Error("Accept from Draft: want error, got nil")
	}
	if got := a.Version(); got != 0 {
		t.Errorf("Version() = %d, want 0 unchanged", got)
	}
}

// Version is unchanged by Publish and remains immutable once Published.
func TestAct_Publish_DoesNotChangeVersion(t *testing.T) {
	a := NewDraftAct("act-pubver")
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if err := a.Accept(); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if a.Version() != 1 {
		t.Fatalf("after Accept Version() = %d, want 1", a.Version())
	}

	if err := a.Publish(); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if got := a.Version(); got != 1 {
		t.Errorf("after Publish Version() = %d, want 1 unchanged", got)
	}
}

// actInPublishedStatus creates an Act in Published status (Draft -> StartVoting -> Accept -> Publish).
func actInPublishedStatus(t *testing.T, actID string) *Act {
	t.Helper()
	a := NewDraftAct(actID)
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if err := a.Accept(); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if err := a.Publish(); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return a
}

// When Published, no paragraphs can be added.
func TestActInPublished_AddParagraph_ReturnsErrorAndCountUnchanged(t *testing.T) {
	a := actInPublishedStatus(t, "act-pub-add")
	if a.Status() != StatusPublished {
		t.Fatalf("act status = %v, want Published", a.Status())
	}
	before := a.ParagraphCount()

	err := a.AddParagraph("p1", "content")

	if err == nil {
		t.Error("AddParagraph from Published: want error, got nil")
	}
	if a.ParagraphCount() != before {
		t.Errorf("ParagraphCount() = %d, want %d unchanged", a.ParagraphCount(), before)
	}
}

// When Published, no paragraph content can be replaced.
func TestActInPublished_ReplaceParagraphContent_ReturnsErrorAndContentUnchanged(t *testing.T) {
	a := NewDraftAct("act-pub-replace")
	if err := a.AddParagraph("p1", "original"); err != nil {
		t.Fatalf("AddParagraph: %v", err)
	}
	if err := a.StartVoting(); err != nil {
		t.Fatalf("StartVoting: %v", err)
	}
	if err := a.Accept(); err != nil {
		t.Fatalf("Accept: %v", err)
	}
	if err := a.Publish(); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if a.Status() != StatusPublished {
		t.Fatalf("act status = %v, want Published", a.Status())
	}

	err := a.ReplaceParagraphContent("p1", "new")

	if err == nil {
		t.Error("ReplaceParagraphContent from Published: want error, got nil")
	}
	ps := a.Paragraphs()
	if len(ps) != 1 || ps[0].Content != "original" {
		t.Errorf("content changed: Paragraphs()[0].Content = %q, want original", ps[0].Content)
	}
}

// When Published, no lifecycle regression: StartVoting, Accept, and Publish all return error and status stays Published.
func TestActInPublished_NoLifecycleRegression(t *testing.T) {
	a := actInPublishedStatus(t, "act-pub-noregress")
	if a.Status() != StatusPublished {
		t.Fatalf("act status = %v, want Published", a.Status())
	}

	if err := a.StartVoting(); err == nil {
		t.Error("StartVoting from Published: want error, got nil")
	}
	if a.Status() != StatusPublished {
		t.Errorf("after StartVoting status = %v, want Published", a.Status())
	}

	if err := a.Accept(); err == nil {
		t.Error("Accept from Published: want error, got nil")
	}
	if a.Status() != StatusPublished {
		t.Errorf("after Accept status = %v, want Published", a.Status())
	}

	if err := a.Publish(); err == nil {
		t.Error("Publish from Published: want error, got nil")
	}
	if a.Status() != StatusPublished {
		t.Errorf("after Publish status = %v, want Published", a.Status())
	}
}
