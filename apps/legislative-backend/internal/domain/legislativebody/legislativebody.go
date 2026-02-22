package legislativebody

import "errors"

// Member has a unique ID, name, and active status.
type Member struct {
	ID     string
	Name   string
	Active bool
}

// LegislativeBody represents a body with identity and dynamic membership.
type LegislativeBody struct {
	id      string
	members []Member
}

// NewLegislativeBody creates a legislative body with the given identity and no members.
func NewLegislativeBody(id string) *LegislativeBody {
	return &LegislativeBody{
		id: id,
	}
}

// ID returns the unique identifier of the legislative body.
func (lb *LegislativeBody) ID() string {
	return lb.id
}

// AddMember adds a member with the given id and name; the member is active by default.
func (lb *LegislativeBody) AddMember(id, name string) {
	lb.members = append(lb.members, Member{ID: id, Name: name, Active: true})
}

// ActiveMembers returns a copy of members that are active.
func (lb *LegislativeBody) ActiveMembers() []Member {
	var out []Member
	for _, m := range lb.members {
		if m.Active {
			out = append(out, m)
		}
	}
	return out
}

// Deactivate sets the member's Active status to false. Returns error if member not found.
func (lb *LegislativeBody) Deactivate(memberID string) error {
	for i := range lb.members {
		if lb.members[i].ID == memberID {
			lb.members[i].Active = false
			return nil
		}
	}
	return errors.New("member not found")
}

// Remove removes the member with the given ID from the list. Returns error if member not found.
func (lb *LegislativeBody) Remove(memberID string) error {
	for i := range lb.members {
		if lb.members[i].ID == memberID {
			lb.members = append(lb.members[:i], lb.members[i+1:]...)
			return nil
		}
	}
	return errors.New("member not found")
}
