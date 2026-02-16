package act

// Status represents the lifecycle state of an Act.
type Status string

const (
	StatusDraft  Status = "Draft"
	StatusVoting Status = "Voting"
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
func (a *Act) StartVoting() {
	a.status = StatusVoting
}
