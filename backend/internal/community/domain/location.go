// Package domain is the heart of the Community context: neighbours publishing help requests
// and meetups at a place, commenting on them and marking them helpful. It has no knowledge of
// HTTP, SQL or Redis; those live in the application and infrastructure layers around it.
package domain

import "fmt"

// Location is where a post is pinned, as WGS84 latitude and longitude. It is a value object:
// two locations with the same coordinates are the same location.
type Location struct {
	lat, lng float64
}

func NewLocation(lat, lng float64) (Location, error) {
	if lat < -90 || lat > 90 {
		return Location{}, fmt.Errorf("%w: latitude %v is outside -90..90", ErrInvalidLocation, lat)
	}
	if lng < -180 || lng > 180 {
		return Location{}, fmt.Errorf("%w: longitude %v is outside -180..180", ErrInvalidLocation, lng)
	}
	return Location{lat: lat, lng: lng}, nil
}

func (l Location) Lat() float64 { return l.lat }
func (l Location) Lng() float64 { return l.lng }
