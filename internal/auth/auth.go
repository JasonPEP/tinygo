package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"tinygo/internal/config"
	"tinygo/internal/logger"

	"github.com/gorilla/sessions"
)

// Store holds the session store
var Store *sessions.CookieStore

// Init initializes the authentication system
func Init(cfg config.AuthConfig) {
	// Generate a random session key if not provided
	sessionKey := cfg.SessionKey
	if sessionKey == "" {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			logger.Log.Fatalf("generate session key: %v", err)
		}
		sessionKey = base64.StdEncoding.EncodeToString(key)
		logger.Log.Warnf("generated random session key: %s", sessionKey)
	}

	Store = sessions.NewCookieStore([]byte(sessionKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	logger.Log.Info("authentication system initialized")
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	session, err := Store.Get(r, "auth")
	if err != nil {
		logger.Log.Errorf("get session: %v", err)
		return false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	return ok && authenticated
}

// SetAuthenticated sets the authentication status in the session
func SetAuthenticated(w http.ResponseWriter, r *http.Request, authenticated bool) error {
	session, err := Store.Get(r, "auth")
	if err != nil {
		return fmt.Errorf("get session: %v", err)
	}

	session.Values["authenticated"] = authenticated
	if authenticated {
		session.Values["login_time"] = time.Now().Unix() // Store as Unix timestamp instead of time.Time
	} else {
		delete(session.Values, "login_time")
	}

	return session.Save(r, w)
}

// GetLoginTime returns the login time from the session
func GetLoginTime(r *http.Request) (time.Time, bool) {
	session, err := Store.Get(r, "auth")
	if err != nil {
		return time.Time{}, false
	}

	loginTimeUnix, ok := session.Values["login_time"].(int64)
	if !ok {
		return time.Time{}, false
	}

	return time.Unix(loginTimeUnix, 0), true
}

// RequireAuth is a middleware that requires authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			// Redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoginRequired is a middleware that requires authentication for API endpoints
func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "authentication required"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
