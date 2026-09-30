package application

import (
	"context"
	"errors"
	"testing"

	"circleoflife/internal/identity/domain"
)

type fakeUsers struct{ byEmail map[domain.Email]*domain.User }

func (f *fakeUsers) Add(_ context.Context, u *domain.User) error {
	if _, taken := f.byEmail[u.Email()]; taken {
		return domain.ErrEmailTaken
	}
	f.byEmail[u.Email()] = u
	return nil
}
func (f *fakeUsers) ByEmail(_ context.Context, e domain.Email) (*domain.User, error) {
	if u, ok := f.byEmail[e]; ok {
		return u, nil
	}
	return nil, domain.ErrUserNotFound
}
func (f *fakeUsers) ByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	for _, u := range f.byEmail {
		if u.ID() == id {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// reversing is a stand-in hasher: fast and obviously not the plain password.
type reversing struct{}

func (reversing) Hash(p string) (string, error) {
	r := []rune(p)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}
func (h reversing) Matches(p, hash string) bool { got, _ := h.Hash(p); return got == hash }

type tokenForID struct{}

func (tokenForID) Issue(id domain.UserID) (string, error) { return "token:" + string(id), nil }

func newService() (*Service, *fakeUsers) {
	users := &fakeUsers{byEmail: map[domain.Email]*domain.User{}}
	return NewService(users, reversing{}, tokenForID{}), users
}

func TestSignUpStoresAHashAndIssuesAToken(t *testing.T) {
	svc, users := newService()
	token, err := svc.SignUp(context.Background(), "Alice", "alice@example.com", "secret1")
	if err != nil {
		t.Fatal(err)
	}
	u := users.byEmail["alice@example.com"]
	if u == nil || u.PasswordHash() == "secret1" {
		t.Fatalf("stored user = %+v", u)
	}
	if token != "token:"+string(u.ID()) {
		t.Errorf("token = %q", token)
	}
}

func TestSignUpRules(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	if _, err := svc.SignUp(ctx, "Alice", "alice@example.com", "secret1"); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name, email, password string
		want                  error
	}{
		{"Alice", "alice@example.com", "secret1", domain.ErrEmailTaken},
		{"Bob", "not-an-email", "secret1", domain.ErrInvalidEmail},
		{" ", "bob@example.com", "secret1", domain.ErrInvalidName},
		{"Bob", "bob@example.com", "12345", domain.ErrWeakPassword},
	} {
		if _, err := svc.SignUp(ctx, c.name, c.email, c.password); !errors.Is(err, c.want) {
			t.Errorf("SignUp(%q, %q, %q): err = %v, want %v", c.name, c.email, c.password, err, c.want)
		}
	}
}

func TestLogInDoesNotRevealWhichPartWasWrong(t *testing.T) {
	svc, _ := newService()
	ctx := context.Background()
	if _, err := svc.SignUp(ctx, "Alice", "alice@example.com", "secret1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.LogIn(ctx, "alice@example.com", "wrong"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("wrong password: err = %v", err)
	}
	if _, err := svc.LogIn(ctx, "nobody@example.com", "secret1"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("unknown email: err = %v", err)
	}
	if token, err := svc.LogIn(ctx, "alice@example.com", "secret1"); err != nil || token == "" {
		t.Errorf("correct password: token=%q err=%v", token, err)
	}
}

func TestProfileOfUnknownMemberIsNotFound(t *testing.T) {
	svc, _ := newService()
	if _, err := svc.Profile(context.Background(), "missing"); !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("err = %v", err)
	}
}
