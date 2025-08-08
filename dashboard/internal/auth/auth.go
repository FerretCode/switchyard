package auth

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"sync"

	"github.com/ferretcode/switchyard/dashboard/internal/types"
)

var (
	sessionsMu sync.Mutex
	sessions   = make(map[string]bool) // map[sessionToken]isAuthenticated
)

type AuthService struct {
	Config *types.Config
}

func NewAuthService(config *types.Config) AuthService {
	return AuthService{
		Config: config,
	}
}

func (a *AuthService) RenderLogin(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	cookie, err := r.Cookie(a.Config.SessionsCookieName)
	if err == nil {
		sessionsMu.Lock()
		valid := sessions[cookie.Value]
		sessionsMu.Unlock()

		if valid {
			http.Redirect(w, r, "/dashboard/home", http.StatusSeeOther)
			return nil
		}
	}

	return templates.ExecuteTemplate(w, "login.html", nil)
}

func (a *AuthService) Login(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != a.Config.AdminUsername || password != a.Config.AdminPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return nil
	}

	token, err := generateToken()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil
	}

	sessionsMu.Lock()
	sessions[token] = true
	sessionsMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     a.Config.SessionsCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	http.Redirect(w, r, "/dashboard/home", http.StatusSeeOther)

	return nil
}

func (a *AuthService) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(a.Config.SessionsCookieName)
	if err == nil {
		sessionsMu.Lock()
		delete(sessions, cookie.Value)
		sessionsMu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:   a.Config.SessionsCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/dashboard/home", http.StatusSeeOther)

	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *AuthService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(a.Config.SessionsCookieName)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		sessionsMu.Lock()
		valid := sessions[cookie.Value]
		sessionsMu.Unlock()

		if !valid {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
