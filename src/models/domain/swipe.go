package domain

// Swipe a indication that one or both users are interested in matching.
//
// The first user id should be set to the lowest user id value. This is
// to provide predictable ordering (else we might create duplicated records).
type Swipe struct {
	ID               int  `json:"id"`
	FirstUserID      int  `json:"firstUserID"`
	SecondUserID     int  `json:"secondUserID"`
	FirstUserSwiped  bool `json:"firstUserSwiped"`
	SecondUserSwiped bool `json:"secondUserSwiped"`
}
