package mds

import "testing"

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

func TestCertificationLevelRankOrdering(t *testing.T) {
	if CertificationLevelRank("FIDO_CERTIFIED_L2") <= CertificationLevelRank("FIDO_CERTIFIED_L1") {
		t.Error("L2 should rank higher than L1")
	}
	if CertificationLevelRank("FIDO_CERTIFIED_L2plus") <= CertificationLevelRank("FIDO_CERTIFIED_L2") {
		t.Error("L2plus should rank higher than L2")
	}
}
