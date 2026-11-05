package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	Remember bool   `json:"remember_me"`
}

type authUser struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	EmployeeNumber string `json:"employee_number"`
	PhoneNumber    string `json:"phone_number"`
}

type contextKey string

const authUserKey contextKey = "auth_user"
const authCookieName = "auth_token"

// jwtSecret reads the secret from the environment.
func jwtSecret() ([]byte, error) {
	secret := strings.TrimSpace(env.GetString("JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}
	return []byte(secret), nil
}

// login verifies the credentials and sets a signed authentication cookie.
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
		slog.Error("login database query failed", "error", err)
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
		slog.Error("login authentication is not configured", "error", err)
		http.Error(w, "Authentication is not configured", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	tokenLifetime := 30 * time.Minute
	if input.Remember {
		tokenLifetime = 30 * 24 * time.Hour
	}
	claims := authClaims{
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	authCookie := &http.Cookie{
		Name:     authCookieName,
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil || strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true"),
		SameSite: http.SameSiteLaxMode,
	}
	if input.Remember {
		authCookie.MaxAge = int(tokenLifetime.Seconds())
		authCookie.Expires = now.Add(tokenLifetime)
	}
	http.SetCookie(w, authCookie)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	})
}

// logout clears the authentication cookie.
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil || strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true"),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	w.WriteHeader(http.StatusNoContent)
}

// authenticate checks the JWT before allowing access to a protected route.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""
		if cookie, err := r.Cookie(authCookieName); err == nil {
			tokenString = cookie.Value
		}
		if tokenString == "" {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}
		if tokenString == "" {
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
			tokenString,
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

func (app *application) me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(authUserKey).(*authClaims)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	user, err := app.getAuthUser(r, claims.UserID)
	if err != nil {
		http.Error(w, "Could not load user profile", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

type updateProfileRequest struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	EmployeeNumber string `json:"employee_number"`
	PhoneNumber    string `json:"phone_number"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (app *application) getAuthUser(r *http.Request, userID string) (authUser, error) {
	var user authUser
	err := app.db.QueryRow(r.Context(), `
		SELECT id::text, full_name, email, role,
		       COALESCE(employee_number, ''), COALESCE(phone_number, '')
		FROM users
		WHERE id::text = $1
	`, userID).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Role, &user.EmployeeNumber, &user.PhoneNumber,
	)
	return user, err
}

func (app *application) updateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(authUserKey).(*authClaims)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var input updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(input.Email)
	input.EmployeeNumber = strings.TrimSpace(input.EmployeeNumber)
	input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)
	if input.FullName == "" || input.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}
	parsedEmail, err := mail.ParseAddress(input.Email)
	if err != nil || parsedEmail.Address != input.Email {
		http.Error(w, "Enter a valid email address", http.StatusBadRequest)
		return
	}

	_, err = app.db.Exec(r.Context(), `
		UPDATE users
		SET full_name = $1,
		    email = $2,
		    employee_number = NULLIF($3, ''),
		    phone_number = NULLIF($4, '')
		WHERE id::text = $5
	`, input.FullName, input.Email, input.EmployeeNumber, input.PhoneNumber, claims.UserID)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			http.Error(w, "That email address or employee number is already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Could not update user profile", http.StatusInternalServerError)
		return
	}

	user, err := app.getAuthUser(r, claims.UserID)
	if err != nil {
		http.Error(w, "Could not load updated profile", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(authUserKey).(*authClaims)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var input changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.CurrentPassword == "" || input.NewPassword == "" || input.ConfirmPassword == "" {
		http.Error(w, "Current password, new password, and confirmation are required", http.StatusBadRequest)
		return
	}
	if input.NewPassword != input.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	if len([]byte(input.NewPassword)) > 72 {
		http.Error(w, "New password must be no more than 72 bytes", http.StatusBadRequest)
		return
	}

	var passwordHash string
	if err := app.db.QueryRow(r.Context(), "SELECT password_hash FROM users WHERE id::text = $1", claims.UserID).Scan(&passwordHash); err != nil {
		http.Error(w, "Could not verify current password", http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.CurrentPassword)); err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not secure new password", http.StatusInternalServerError)
		return
	}
	if _, err := app.db.Exec(r.Context(), "UPDATE users SET password_hash = $1 WHERE id::text = $2", string(newHash), claims.UserID); err != nil {
		http.Error(w, "Could not update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
