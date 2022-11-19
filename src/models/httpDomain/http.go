package httpDomain

import "github.com/markwallsgrove/muzz_devops/src/models/domain"

type UserResult struct {
	Result User `json:"result"`
}

type UserResults struct {
	Results []User `json:"results"`
}

type UserProfileResults struct {
	Results []domain.UserProfile `json:"results"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

// UserToHTTPDomain convert a user struct to a http user domain struct
func UserToHTTPDomain(user *domain.User, password string) User {
	return User{
		ID:       user.ID,
		Email:    user.Email,
		Password: password,
		Name:     user.Name,
		Gender:   user.Gender.String(),
		Age:      user.Age,
	}
}

type SwipeResults struct {
	Results SwipeResult `json:"results"`
}

type SwipeResult struct {
	ID      int  `json:"matchID,omitempty"`
	Matched bool `json:"matched"`
}
