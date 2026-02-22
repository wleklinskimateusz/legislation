package legislativebody

import (
	"testing"
)

// A LegislativeBody can be created with an identity.
func TestNewLegislativeBody_HasIdentity(t *testing.T) {
	lb := NewLegislativeBody("lb-1")

	if got := lb.ID(); got != "lb-1" {
		t.Errorf("ID() = %q, want lb-1", got)
	}
}

// AddMember adds a member; ActiveMembers returns only active members with ID, Name, Active.
func TestLegislativeBody_AddMember_ActiveMembersReturnsMember(t *testing.T) {
	lb := NewLegislativeBody("lb-2")
	lb.AddMember("m1", "Alice")

	active := lb.ActiveMembers()
	if len(active) != 1 {
		t.Fatalf("ActiveMembers() length = %d, want 1", len(active))
	}
	if active[0].ID != "m1" || active[0].Name != "Alice" || !active[0].Active {
		t.Errorf("ActiveMembers()[0] = %+v, want ID m1 Name Alice Active true", active[0])
	}
}

// Deactivate sets the member's Active status to false; the member is no longer in ActiveMembers().
func TestLegislativeBody_Deactivate_RemovesMemberFromActiveMembers(t *testing.T) {
	lb := NewLegislativeBody("lb-3")
	lb.AddMember("m1", "Alice")

	err := lb.Deactivate("m1")

	if err != nil {
		t.Fatalf("Deactivate: %v", err)
	}
	active := lb.ActiveMembers()
	if len(active) != 0 {
		t.Errorf("ActiveMembers() length = %d, want 0 after deactivate", len(active))
	}
}

// Remove removes the member from the list; ActiveMembers() no longer contains them.
func TestLegislativeBody_Remove_MemberGoneFromList(t *testing.T) {
	lb := NewLegislativeBody("lb-4")
	lb.AddMember("m1", "Alice")

	err := lb.Remove("m1")

	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	active := lb.ActiveMembers()
	if len(active) != 0 {
		t.Errorf("ActiveMembers() length = %d, want 0 after remove", len(active))
	}
}

// ActiveMembers returns only active members; deactivated members are excluded.
func TestLegislativeBody_ActiveMembers_ReturnsOnlyActiveMembers(t *testing.T) {
	lb := NewLegislativeBody("lb-5")
	lb.AddMember("m1", "Alice")
	lb.AddMember("m2", "Bob")
	if err := lb.Deactivate("m1"); err != nil {
		t.Fatalf("Deactivate: %v", err)
	}

	active := lb.ActiveMembers()
	if len(active) != 1 {
		t.Fatalf("ActiveMembers() length = %d, want 1", len(active))
	}
	if active[0].ID != "m2" || active[0].Name != "Bob" || !active[0].Active {
		t.Errorf("ActiveMembers()[0] = %+v, want ID m2 Name Bob Active true", active[0])
	}
}
