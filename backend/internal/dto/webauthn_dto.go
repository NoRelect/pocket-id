package dto

import (
	"github.com/go-webauthn/webauthn/protocol"
	datatype "github.com/pocket-id/pocket-id/backend/internal/model/types"
)

type MdsAuthenticatorDto struct {
	AAGUID                 string   `json:"aaguid"`
	Description            string   `json:"description"`
	IsCompromised          bool     `json:"isCompromised"`
	FidoCertificationLevel string   `json:"fidoCertificationLevel"`
	AttestationTypes       []string `json:"attestationTypes"`
	KeyProtection          []string `json:"keyProtection"`
	AttachmentHint         []string `json:"attachmentHint"`
	Icon                   string   `json:"icon"`
	IconDark               string   `json:"iconDark"`
}

type WebauthnCredentialDto struct {
	ID              string                            `json:"id"`
	Name            string                            `json:"name"`
	CredentialID    string                            `json:"credentialID"`
	AttestationType string                            `json:"attestationType"`
	Transport       []protocol.AuthenticatorTransport `json:"transport" swaggertype:"array,string"`

	AAGUID              string `json:"aaguid"`
	AttestationVerified bool   `json:"attestationVerified"`
	IsCompromised       bool   `json:"isCompromised"`

	BackupEligible bool `json:"backupEligible"`
	BackupState    bool `json:"backupState"`

	CreatedAt datatype.DateTime `json:"createdAt"`
}

type WebauthnCredentialUpdateDto struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}
