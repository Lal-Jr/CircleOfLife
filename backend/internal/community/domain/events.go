package domain

// Event is something that happened in the Community context that other parts of the system
// (such as the Realtime context's live feed) may react to.
type Event interface {
	eventName() string
}

// PostPublished: a neighbour pinned a new post at a location.
type PostPublished struct {
	Post PostID
	At   Location
}

func (PostPublished) eventName() string { return "post_published" }
