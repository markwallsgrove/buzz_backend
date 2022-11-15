package models

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
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   Gender `json:"gender"`
	Age      int    `json:"age"`
}

type Results struct {
	Result *User `json:"result"`
}
