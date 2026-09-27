package db

import (
	"context"
	"fmt"
)

// Sample neighbourhood around central Bangalore, so a reviewer logging in
// with the demo account sees a lived-in community: several neighbours,
// help requests and meetups at different distances, comment threads and
// Helpful votes.

type seedUser struct {
	key, name, email string
}

type seedPost struct {
	key, author, title, description, kind string
	lat, lng                              float64
	ageMinutes                            int
	// Minutes from now; 0 means no meetup time.
	meetupInMinutes int
}

type seedComment struct {
	post, author, content string
	ageMinutes            int
}

var neighbourUsers = []seedUser{
	{"priya", "Priya Raman", "priya@neighbours.circleoflife.app"},
	{"arjun", "Arjun Mehta", "arjun@neighbours.circleoflife.app"},
	{"meera", "Meera Iyer", "meera@neighbours.circleoflife.app"},
	{"rahul", "Rahul Nair", "rahul@neighbours.circleoflife.app"},
	{"ananya", "Ananya Das", "ananya@neighbours.circleoflife.app"},
}

var neighbourPosts = []seedPost{
	{"jumper", "arjun", "Car battery died - anyone have jumper cables?",
		"Stuck in the basement parking of Brigade Road complex. Just need a jump start, will take 5 minutes.",
		"help", 12.9730, 77.6070, 18, 0},
	{"medicine", "meera", "Need someone to pick up medicine for my mother",
		"She's unwell and I'm stuck at work until 8pm. The pharmacy on Residency Road has the order ready and paid for.",
		"help", 12.9665, 77.6010, 45, 0},
	{"dog", "rahul", "Lost dog near Cubbon Park - brown indie, answers to Bruno",
		"He slipped his leash near the Queen's statue gate this morning. Red collar with a tag. Please call if you see him!",
		"help", 12.9763, 77.5929, 90, 0},
	{"laptop", "ananya", "Can someone help set up a laptop for my grandfather?",
		"He wants to video call family abroad. Needs basic setup and a quick walkthrough - happy to pay in chai and snacks.",
		"help", 12.9580, 77.5990, 240, 0},
	{"ladder", "priya", "Borrowing a ladder for the weekend",
		"Need to change a ceiling fan. Will return it Sunday evening, promise!",
		"help", 12.9810, 77.6110, 360, 0},
	{"run", "priya", "Sunday morning 5K run at Cubbon Park",
		"Easy pace, all levels welcome. Meeting at the bandstand, coffee afterwards.",
		"meetup", 12.9760, 77.5930, 120, 2 * 24 * 60},
	{"books", "meera", "Book swap at the corner cafe",
		"Bring 2-3 books you've finished, leave with new ones. Fiction, non-fiction, kids' books - all welcome.",
		"meetup", 12.9700, 77.6080, 300, 3 * 24 * 60},
	{"garden", "rahul", "Community garden planting day",
		"We're planting tomatoes and herbs in the apartment terrace garden. Gloves provided, bring water.",
		"meetup", 12.9640, 77.5880, 600, 4 * 24 * 60},
	{"cricket", "arjun", "Evening cricket in the ground behind the school",
		"Short 6-over matches, tennis ball. We're 2 players short!",
		"meetup", 12.9850, 77.5980, 30, 5 * 60},
	{"coding", "ananya", "Beginner coding study group",
		"Working through a Python course together every Thursday. Laptops needed, questions encouraged.",
		"meetup", 12.9550, 77.6150, 1440, 5 * 24 * 60},
}

var neighbourComments = []seedComment{
	{"jumper", "rahul", "I'm 5 minutes away, heading over with cables now.", 12},
	{"jumper", "arjun", "Legend, thank you! I'm on level B2.", 10},
	{"medicine", "priya", "I can pick it up on my way home at 6, send me the order details.", 40},
	{"medicine", "meera", "Thank you so much Priya, messaging you now.", 38},
	{"dog", "ananya", "Shared this in our building group. Will keep an eye out on my evening walk.", 80},
	{"dog", "demo", "Saw a brown dog with a red collar near the library entrance about 20 minutes ago!", 60},
	{"dog", "rahul", "Heading there now, thank you!!", 58},
	{"laptop", "arjun", "I can come over Saturday afternoon if that works?", 200},
	{"run", "meera", "Count me in! First 5K in a while, so I'll be at the back.", 100},
	{"run", "demo", "Sounds fun, I'll be there.", 90},
	{"books", "ananya", "Bringing a few kids' books my niece has outgrown.", 250},
	{"cricket", "rahul", "I'm in, can bring an extra bat.", 20},
	{"groceries", "arjun", "I'm in the next building, can help in 10 minutes.", 5},
	{"cleanup", "priya", "Great idea! I'll bring gloves and trash bags.", 4},
	{"cleanup", "meera", "Joining with my kids.", 3},
}

// Helpful votes: post key -> voters.
var neighbourLikes = map[string][]string{
	"jumper":    {"priya", "meera"},
	"medicine":  {"ananya", "rahul", "demo"},
	"dog":       {"priya", "arjun", "meera", "ananya", "demo"},
	"laptop":    {"meera"},
	"run":       {"arjun", "ananya", "rahul"},
	"books":     {"priya", "demo"},
	"garden":    {"meera", "ananya"},
	"cricket":   {"priya"},
	"groceries": {"priya", "rahul"},
	"cleanup":   {"arjun", "ananya", "rahul"},
}

// seedNeighbourhood populates the sample community once, in one transaction,
// so a failed run leaves nothing behind and is retried on the next startup.
// Neighbour accounts get an unusable password hash: only the demo account
// can log in.
func seedNeighbourhood(ctx context.Context, demoUserID string) error {
	var exists bool
	if err := Pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		neighbourUsers[0].email).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	tx, err := Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	users := map[string]string{"demo": demoUserID}
	for _, u := range neighbourUsers {
		var id string
		if err := tx.QueryRow(ctx,
			`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, '!seeded-no-login') RETURNING id`,
			u.name, u.email,
		).Scan(&id); err != nil {
			return fmt.Errorf("creating %s: %v", u.email, err)
		}
		users[u.key] = id
	}

	posts := map[string]string{}
	for _, p := range neighbourPosts {
		var id string
		if err := tx.QueryRow(ctx, `
			INSERT INTO posts (user_id, title, description, type, location, meetup_time, created_at)
			VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326),
				CASE WHEN $7::int > 0 THEN now() + make_interval(mins => $7::int) END,
				now() - make_interval(mins => $8::int))
			RETURNING id`,
			users[p.author], p.title, p.description, p.kind, p.lng, p.lat, p.meetupInMinutes, p.ageMinutes,
		).Scan(&id); err != nil {
			return fmt.Errorf("creating post %q: %v", p.title, err)
		}
		posts[p.key] = id
	}

	// The demo account's own posts, so neighbours can reply to them too.
	for key, title := range map[string]string{
		"groceries": "Need help carrying groceries upstairs",
		"cleanup":   "Weekend park cleanup meetup",
	} {
		var id string
		if err := tx.QueryRow(ctx, `SELECT id FROM posts WHERE user_id = $1 AND title = $2 LIMIT 1`,
			demoUserID, title).Scan(&id); err == nil {
			posts[key] = id
		}
	}

	for _, c := range neighbourComments {
		postID, ok := posts[c.post]
		if !ok {
			continue
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO comments (post_id, user_id, content, created_at)
			VALUES ($1, $2, $3, now() - make_interval(mins => $4::int))`,
			postID, users[c.author], c.content, c.ageMinutes,
		); err != nil {
			return fmt.Errorf("creating comment: %v", err)
		}
	}

	for postKey, voters := range neighbourLikes {
		postID, ok := posts[postKey]
		if !ok {
			continue
		}
		for _, v := range voters {
			if _, err := tx.Exec(ctx,
				`INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
				postID, users[v],
			); err != nil {
				return fmt.Errorf("creating vote: %v", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	fmt.Printf("Seeded sample neighbourhood: %d neighbours, %d posts.\n", len(neighbourUsers), len(neighbourPosts))
	return nil
}
