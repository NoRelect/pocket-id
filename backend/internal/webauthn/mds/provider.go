package mds

import (
	"context"

	"github.com/go-webauthn/webauthn/metadata"
	"github.com/google/uuid"
)

func (s *Service) GetEntry(ctx context.Context, aaguid uuid.UUID) (*metadata.Entry, error) {
	entry, found := s.Lookup(aaguid.String())
	if !found {
		return nil, nil
	}
	return &entry, nil
}

func (s *Service) GetValidateEntry(ctx context.Context) bool {
	return s.enforce.Load()
}

func (s *Service) GetValidateEntryPermitZeroAAGUID(ctx context.Context) bool {
	return false
}

func (s *Service) GetValidateTrustAnchor(ctx context.Context) bool {
	return s.enforce.Load()
}

func (s *Service) GetValidateStatus(ctx context.Context) bool {
	return s.enforce.Load()
}

func (s *Service) GetValidateAttestationTypes(ctx context.Context) bool {
	return s.enforce.Load()
}

func (s *Service) ValidateStatusReports(ctx context.Context, reports []metadata.StatusReport) error {
	return metadata.ValidateStatusReports(reports, nil, metadata.DefaultUndesiredAuthenticatorStatuses())
}

var _ metadata.Provider = (*Service)(nil)
