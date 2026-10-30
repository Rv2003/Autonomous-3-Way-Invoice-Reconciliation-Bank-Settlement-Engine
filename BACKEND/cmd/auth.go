package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// authClaims stores the user's identity and role in the JWT.
type authClaims struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authUser struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type contextKey string

const authUserKey contextKey = "auth_user"

// jwtSecret reads the secret from the environment.
func jwtSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}
	return []byte(secret), nil
}

// login verifies the credentials and returns a signed token.
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user authUser
	var passwordHash string

	err := app.db.QueryRow(
		r.Context(),
		`SELECT id::text, full_name, email, password_hash, role
 FROM users WHERE lower(email) = lower($1)`,
		input.Email,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&passwordHash,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash), []byte(input.Password),
	); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	secret, err := jwtSecret()
	if err != nil {
		http.Error(w, "Authentication is not configured", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	claims := authClaims{
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token": tokenString,
		"user":  user,
	})
}

// authenticate checks the JWT before allowing access to a protected route.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		secret, err := jwtSecret()
		if err != nil {
			http.Error(w, "Authentication is not configured", http.StatusInternalServerError)
			return
		}

		claims := &authClaims{}
		token, err := jwt.ParseWithClaims(
			parts[1],
			claims,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}
				return secret, nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// me returns the identity and role from the validated token.
func (app *application) me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(authUserKey).(*authClaims)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authUser{
		ID:       claims.UserID,
		FullName: claims.FullName,
		Email:    claims.Email,
		Role:     claims.Role,
	})
}

// requireRole restricts a route to the specified roles.
func (app *application) requireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return app.authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(authUserKey).(*authClaims)
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Permission denied", http.StatusForbidden)
		}))
	}
}
