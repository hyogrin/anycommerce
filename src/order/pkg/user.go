package pkg

import (
	"time"
)

// User Struct
type User struct {
	ID             string     `json:"id" yaml:"id"`
	Username       string     `json:"username" yaml:"username"`
	Email          string     `json:"email" yaml:"email"`
	FirstName      string     `json:"first_name" yaml:"first_name"`
	LastName       string     `json:"last_name" yaml:"last_name"`
	Age            int        `json:"age" yaml:"age"`
	Gender         string     `json:"gender" yaml:"gender"`
	Segment        string     `json:"segment" yaml:"segment"`
	SignUpDate     *time.Time `json:"sign_up_date,omitempty"`
	SelectableUser bool       `json:"selectable_user" yaml:"selectable_user"`
	LastSignInDate *time.Time `json:"last_sign_in_date,omitempty"`
	IdentityId     string     `json:"identity_id,omitempty"`
	PhoneNumber    string     `json:"phone_number,omitempty"`
}

// Users Array
type Users []User
