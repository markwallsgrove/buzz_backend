package domain

type Gender int

const (
	UnknownGender = iota
	Male
	Female
)

func (g Gender) String() string {
	return [...]string{"Unknown", "Male", "Female"}[g]
}

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Name         string `json:"name"`
	Gender       Gender `json:"gender"`
	Age          int    `json:"age"`
}

type UserProfile struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender Gender `json:"gender"`
	Age    int    `json:"age"`
}
