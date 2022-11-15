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

type UserProfile struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender Gender `json:"gender"`
	Age    int    `json:"age"`
}

type UserResult struct {
	Result User `json:"result"`
}

type UserResults struct {
	Results []User `json:"results"`
}

type UserProfileResults struct {
	Results []UserProfile `json:"results"`
}
