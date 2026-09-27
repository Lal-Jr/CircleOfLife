// Matches the demo account seeded by the backend (backend/internal/db/db.go)
// so recruiters/reviewers can try the app without signing up.
export const DEMO_EMAIL = "demo@circleoflife.app";
export const DEMO_PASSWORD = "CircleDemo123!";

// Centre of the seeded sample neighbourhood (central Bangalore). The demo
// account is always "placed" here, whatever the visitor's real location, so
// every reviewer sees the same populated feed and map.
export const DEMO_LOCATION = { lat: 12.9716, lng: 77.5946 };
