package models

type User struct {
	PublicIdentifier string `json:"publicIdentifier"`
	FirebaseUID      string `json:"firebaseUID"`
	Email            string `json:"email"`
}
