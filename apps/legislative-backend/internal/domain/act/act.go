package act

import "errors"

// Status represents the lifecycle state of an Act.
type Status string

const (
	StatusDraft     Status = "Draft"
	StatusVoting    Status = "Voting"
	StatusAccepted  Status = "Accepted"
	StatusPublished Status = "Published"
)

// Paragraph is a single paragraph with a unique identifier and content.
type Paragraph struct {
	ID      string
	Content string
}

// Act is a legislative act with identity and internal paragraph storage.
type Act struct {
	id         string
	status     Status
	version    int
	paragraphs []Paragraph
}

// NewDraftAct creates an Act with the given identity in Draft status.
func NewDraftAct(id string) *Act {
	return &Act{
		id:      id,
		status:  StatusDraft,
		version: 0,
	}
}

// ID returns the unique identifier of the Act.
func (a *Act) ID() string {
	return a.id
}

// Status returns the current status of the Act.
func (a *Act) Status() Status {
	return a.status
}

// Version returns the version number of the Act.
func (a *Act) Version() int {
	return a.version
}

// ParagraphCount returns the number of paragraphs stored in the Act.
func (a *Act) ParagraphCount() int {
	return len(a.paragraphs)
}

// Paragraphs returns a copy of the paragraphs in insertion order.
func (a *Act) Paragraphs() []Paragraph {
	if len(a.paragraphs) == 0 {
		return nil
	}
	out := make([]Paragraph, len(a.paragraphs))
	copy(out, a.paragraphs)
	return out
}

// AddParagraph appends a paragraph with the given id and content and stores it internally.
// Returns an error if the Act is not in Draft status.
func (a *Act) AddParagraph(id string, content string) error {
	if a.status != StatusDraft {
		return errors.New("paragraphs can only be added in Draft status")
	}
	a.paragraphs = append(a.paragraphs, Paragraph{ID: id, Content: content})
	return nil
}

// ReplaceParagraphContent replaces the content of the paragraph with the given ID.
// Returns an error if the Act is not in Draft or the paragraph is not found.
// Paragraph identity and order remain unchanged.
func (a *Act) ReplaceParagraphContent(paragraphID string, newContent string) error {
	if a.status != StatusDraft {
		return errors.New("paragraph content can only be replaced in Draft status")
	}
	for i := range a.paragraphs {
		if a.paragraphs[i].ID == paragraphID {
			a.paragraphs[i].Content = newContent
			return nil
		}
	}
	return errors.New("paragraph not found")
}

// StartVoting transitions the Act from Draft to Voting status.
// Returns an error if the Act is not in Draft status.
func (a *Act) StartVoting() error {
	if a.status != StatusDraft {
		return errors.New("StartVoting only allowed from Draft")
	}
	a.status = StatusVoting
	return nil
}

// Accept transitions the Act from Voting to Accepted status.
// Returns an error if the Act is not in Voting status.
// On success, increments the version number.
func (a *Act) Accept() error {
	if a.status != StatusVoting {
		return errors.New("Accept only allowed from Voting")
	}
	a.status = StatusAccepted
	a.version++
	return nil
}

// Publish transitions the Act from Accepted to Published status.
// Returns an error if the Act is not in Accepted status.
func (a *Act) Publish() error {
	if a.status != StatusAccepted {
		return errors.New("Publish only allowed from Accepted")
	}
	a.status = StatusPublished
	return nil
}
