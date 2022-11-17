package domain

type Swipe struct {
	ID               int  `json:"id"`
	FirstUserID      int  `json:"firstUserID"`
	SecondUserID     int  `json:"secondUserID"`
	FirstUserSwiped  bool `json:"firstUserSwiped"`
	SecondUserSwiped bool `json:"secondUserSwiped"`
}
