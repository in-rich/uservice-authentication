package models

type UpdateUser struct {
	PublicIdentifier string `json:"publicIdentifier" validate:"required,max=255"`
}
