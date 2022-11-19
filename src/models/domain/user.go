package domain

type Gender int

const (
	UnknownGender = iota
	Male
	Female
)

var stringToGender map[string]Gender = map[string]Gender{
	"Unknown": UnknownGender,
	"Male":    Male,
	"Female":  Female,
}

var Genders []Gender = []Gender{Male, Female}

func (g Gender) String() string {
	return [...]string{"Unknown", "Male", "Female"}[g]
}

// StringToGender convert a string into a gender enum
func StringToGender(gender string) Gender {
	g, ok := stringToGender[gender]
	if !ok {
		return UnknownGender
	}

	return g
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
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}
