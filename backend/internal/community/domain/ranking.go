package domain

import "time"

// The feed ranks posts by how much they matter to someone standing at a point: close, recent
// help requests first. These constants are the whole policy; the Postgres feed reader builds
// its ORDER BY from them so the database ranks exactly as Score does.
const (
	// Posts closer than this are urgent: someone within a few minutes' walk.
	UrgentWithinMeters = 500
	// Posts closer than this are nearby.
	NearbyWithinMeters = 2000

	UrgentBoost      = 0.5
	NearbyBoost      = 0.2
	HelpRequestBoost = 0.2
)

// Proximity is how close a post is to the viewer, shown as its priority in the feed.
type Proximity string

const (
	Urgent  Proximity = "high"
	Nearby  Proximity = "medium"
	Distant Proximity = "normal"
)

func ProximityOf(distanceMeters float64) Proximity {
	switch {
	case distanceMeters < UrgentWithinMeters:
		return Urgent
	case distanceMeters < NearbyWithinMeters:
		return Nearby
	}
	return Distant
}

// Score ranks a post for a viewer: it decays with age (per hour) and distance (per km), with
// boosts for urgent or nearby posts and for help requests. Higher ranks first.
func Score(age time.Duration, distanceMeters float64, kind Kind) float64 {
	s := 1/(age.Hours()+1) + 1/(distanceMeters/1000+1)
	switch ProximityOf(distanceMeters) {
	case Urgent:
		s += UrgentBoost
	case Nearby:
		s += NearbyBoost
	}
	if kind == HelpRequest {
		s += HelpRequestBoost
	}
	return s
}
