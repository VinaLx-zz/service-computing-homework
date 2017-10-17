package err

// Err is the error code used in the agenda
type Err int

const (
	// OK means no error
	OK Err = 0
	// TimeConflict happens when trying to add a user to a meeting
	// in whose duration he have no time
	TimeConflict Err = 1
	// NoSuchMeeting happens when looking up a meeting by title but none of
	// the meeting have such title
	NoSuchMeeting Err = 2
	// NoSuchUser happens when looking up a user by username but none of the
	// user have such name
	NoSuchUser Err = 3
	// DuplicateMeeting happens when trying to host a meeting with the title
	// that already exist in the meeting set
	DuplicateMeeting Err = 4
	// DuplicateUser happens when trying to add a user with the username that
	// already exist in the User set
	DuplicateUser Err = 5
	// InvalidTime happens when there're some error with time
	InvalidTime Err = 6
)
