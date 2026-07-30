package webauthn

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pocket-id/pocket-id/backend/internal/appconfig"
	"github.com/pocket-id/pocket-id/backend/internal/common"
	"github.com/pocket-id/pocket-id/backend/internal/utils"
	"github.com/pocket-id/pocket-id/backend/internal/webauthn/mds"
)

const (
	AttestationModeDisabled = "disabled"
	AttestationModeOptional = "optional"
	AttestationModeRequired = "required"
)

func (s *Service) checkAttestation(dbConfig *appconfig.AppConfigModel, aaguid []byte, attestationType string, allowedAAGUIDs []string) error {
	mode := strings.ToLower(strings.TrimSpace(dbConfig.PasskeyAttestationMode.String()))

	if mode == "" || mode == AttestationModeDisabled || mode == AttestationModeOptional {
		return nil
	}

	if attestationType == "" || attestationType == "none" {
		return &common.PasskeyAttestationError{
			Reason: "authenticator did not provide attestation information",
		}
	}

	if attestationType != "basic_full" {
		return &common.PasskeyAttestationError{
			Reason: fmt.Sprintf("authenticator attestation type %q is not allowed; only full basic attestation is accepted", attestationType),
		}
	}

	aaguidStr := utils.FormatAAGUID(aaguid)

	if len(allowedAAGUIDs) > 0 && !isAAGUIDAllowed(aaguidStr, allowedAAGUIDs) {
		return &common.PasskeyAttestationError{
			Reason: fmt.Sprintf("authenticator AAGUID %s is not in the allowed list", aaguidStr),
		}
	}

	if err := s.checkCredentialRestrictions(dbConfig, aaguidStr); err != nil {
		return err
	}

	return nil
}

func (s *Service) checkCredentialRestrictions(dbConfig *appconfig.AppConfigModel, aaguidStr string) error {
	if err := s.checkMinCertificationLevel(aaguidStr, strings.TrimSpace(dbConfig.PasskeyMinCertificationLevel.String())); err != nil {
		return err
	}

	return nil
}

func (s *Service) checkMinCertificationLevel(aaguidStr, minLevel string) error {
	if minLevel == "" {
		return nil
	}

	requiredRank := mds.CertificationLevelRank(minLevel)
	if requiredRank == 0 {
		return nil
	}

	if s.mds == nil {
		return &common.PasskeyAttestationError{
			Reason: "authenticator metadata is unavailable, cannot verify certification level",
		}
	}

	entry, found := s.mds.Lookup(aaguidStr)
	if !found {
		return &common.PasskeyAttestationError{
			Reason: fmt.Sprintf("authenticator AAGUID %s has no metadata entry, cannot verify certification level", aaguidStr),
		}
	}

	actualLevel := mds.HighestCertificationLevel(entry)
	if mds.CertificationLevelRank(actualLevel) < requiredRank {
		return &common.PasskeyAttestationError{
			Reason: fmt.Sprintf("authenticator certification level %q is below the required minimum %q", actualLevel, minLevel),
		}
	}

	return nil
}

func isAAGUIDAllowed(aaguid string, allowedAAGUIDs []string) bool {
	lower := strings.ToLower(aaguid)
	return slices.ContainsFunc(allowedAAGUIDs, func(s string) bool {
		return strings.ToLower(s) == lower
	})
}

func parseAllowedAAGUIDs(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		if trimmedLower := strings.ToLower(strings.TrimSpace(p)); trimmedLower != "" {
			result = append(result, trimmedLower)
		}
	}
	return result
}
