package mds

import "testing"

func TestCertificationLevels(t *testing.T) {
	levels := CertificationLevels()
	if len(levels) != len(fidoCertStatusOrder) {
		t.Fatalf("expected %d levels, got %d", len(fidoCertStatusOrder), len(levels))
	}
	if levels[0] != "FIDO_CERTIFIED" {
		t.Errorf("expected first level FIDO_CERTIFIED, got %q", levels[0])
	}
	if levels[len(levels)-1] != "FIDO_CERTIFIED_L3plus" {
		t.Errorf("expected last level FIDO_CERTIFIED_L3plus, got %q", levels[len(levels)-1])
	}
}

func TestCertificationLevelRank(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  int
	}{
		{"empty", "", 0},
		{"unknown", "NOT_A_LEVEL", 0},
		{"not certified", "NOT_FIDO_CERTIFIED", 0},
		{"base", "FIDO_CERTIFIED", 1},
		{"l1", "FIDO_CERTIFIED_L1", 2},
		{"l3plus", "FIDO_CERTIFIED_L3plus", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CertificationLevelRank(tt.level); got != tt.want {
				t.Errorf("CertificationLevelRank(%q) = %d, want %d", tt.level, got, tt.want)
			}
		})
	}
}

func TestDistinctKeyProtection(t *testing.T) {
	s := &Service{}

	if got := s.DistinctKeyProtection(); got != nil {
		t.Errorf("expected nil for empty service, got %v", got)
	}

	e1 := Entry{}
	e1.MetadataStatement.KeyProtection = []string{"HARDWARE", "tee"}
	e2 := Entry{}
	e2.MetadataStatement.KeyProtection = []string{"tee", "secure_element", ""}
	s.SetEntriesForTest(map[string]Entry{"a": e1, "b": e2})

	got := s.DistinctKeyProtection()
	want := []string{"hardware", "secure_element", "tee"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestEntryKeyProtection(t *testing.T) {
	e := Entry{}
	e.MetadataStatement.KeyProtection = []string{"HARDWARE", " secure_element "}
	got := EntryKeyProtection(e)
	want := []string{"hardware", "secure_element"}
	if len(got) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCertificationLevelRankOrdering(t *testing.T) {
	if CertificationLevelRank("FIDO_CERTIFIED_L2") <= CertificationLevelRank("FIDO_CERTIFIED_L1") {
		t.Error("L2 should rank higher than L1")
	}
	if CertificationLevelRank("FIDO_CERTIFIED_L2plus") <= CertificationLevelRank("FIDO_CERTIFIED_L2") {
		t.Error("L2plus should rank higher than L2")
	}
}
