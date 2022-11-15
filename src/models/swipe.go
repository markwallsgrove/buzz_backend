package models

type Swipe struct {
	ID               int  `json:"id"`
	FirstUserID      int  `json:"firstUserID"`
	SecondUserID     int  `json:"secondUserID"`
	FirstUserSwiped  bool `json:"firstUserSwiped"`
	SecondUserSwiped bool `json:"secondUserSwiped"`
}

type SwipeResults struct {
	Results SwipeResult
}

type SwipeResult struct {
	ID      int  `json:"matchID,omitempty"`
	Matched bool `json:"matched"`
}
