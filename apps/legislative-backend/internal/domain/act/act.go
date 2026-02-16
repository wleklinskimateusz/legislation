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

// Act is a legislative act with identity and internal paragraph storage.
type Act struct {
	id         string
	status     Status
	paragraphs []string
}

// NewDraftAct creates an Act with the given identity in Draft status.
func NewDraftAct(id string) *Act {
	return &Act{
		id:     id,
		status: StatusDraft,
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

// ParagraphCount returns the number of paragraphs stored in the Act.
func (a *Act) ParagraphCount() int {
	return len(a.paragraphs)
}

// AddParagraph appends a paragraph and stores it internally.
func (a *Act) AddParagraph(text string) {
	a.paragraphs = append(a.paragraphs, text)
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
func (a *Act) Accept() error {
	if a.status != StatusVoting {
		return errors.New("Accept only allowed from Voting")
	}
	a.status = StatusAccepted
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
