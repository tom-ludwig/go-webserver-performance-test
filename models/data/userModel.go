package data

import (
	"github.com/google/uuid"
)

// User-Roles: A user has a default role of User, which gives him access to the normal public stuff and his own profile
// Other roles include: Admin: Can manage certain user groups. Super Admin: Has access to everything. Developer: Has read access but no write access.
// More User roles: User.ReadWrite, User.Upload, User.Hide, User.Maintain
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`     // omit from JSON
	Roles        string    `json:"roles"` // TODO: Check if roles should be a array of strings or just an array
}

// FIXME: Doesn't work with postgres(Private and PublicUser can't be split but then in the same tabel)
type PrivateUser struct {
	User                // Embedded user type
	PasswordHash string `json:"passwordHash"`
}

type PublicUser struct {
	User
}
