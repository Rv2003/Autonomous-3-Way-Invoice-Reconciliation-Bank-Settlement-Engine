package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Role            string `json:"role"`
	EmployeeNumber  string `json:"employee_number"`
	PhoneNumber     string `json:"phone_number"`
}

type updateUserRoleRequest struct {
	Role string `json:"role"`
}

type managedUser struct {
	ID             string    `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	EmployeeNumber string    `json:"employee_number"`
	PhoneNumber    string    `json:"phone_number"`
	CreatedAt      time.Time `json:"created_at"`
}

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query(r.Context(), `
		SELECT id::text, full_name, email, role, COALESCE(employee_number, ''), COALESCE(phone_number, ''), created_at
		FROM users
		ORDER BY created_at DESC, lower(email)
	`)
	if err != nil {
		http.Error(w, "Could not load users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := make([]managedUser, 0)
	for rows.Next() {
		var user managedUser
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.EmployeeNumber, &user.PhoneNumber, &user.CreatedAt); err != nil {
			http.Error(w, "Could not load users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Could not load users", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(input.Email)
	input.Role = strings.TrimSpace(input.Role)
	input.EmployeeNumber = strings.TrimSpace(input.EmployeeNumber)
	input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)
	if input.FullName == "" || input.Email == "" || input.Password == "" || input.ConfirmPassword == "" {
		http.Error(w, "Name, email, password, and password confirmation are required", http.StatusBadRequest)
		return
	}
	if input.Password != input.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	parsedEmail, err := mail.ParseAddress(input.Email)
	if err != nil || parsedEmail.Address != input.Email {
		http.Error(w, "Enter a valid email address", http.StatusBadRequest)
		return
	}
	if len([]byte(input.Password)) > 72 {
		http.Error(w, "Initial password must be no more than 72 bytes", http.StatusBadRequest)
		return
	}
	if !validUserRole(input.Role) {
		http.Error(w, "Role must be administrator or clerk", http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not secure the initial password", http.StatusInternalServerError)
		return
	}

	var user managedUser
	err = app.db.QueryRow(r.Context(), `
		INSERT INTO users (full_name, email, password_hash, role, employee_number, phone_number)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''))
		RETURNING id::text, full_name, email, role, COALESCE(employee_number, ''), COALESCE(phone_number, ''), created_at
	`, input.FullName, input.Email, string(passwordHash), input.Role, input.EmployeeNumber, input.PhoneNumber).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Role, &user.EmployeeNumber, &user.PhoneNumber, &user.CreatedAt,
	)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			http.Error(w, "A user with that email address or employee number already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (app *application) updateUserRole(w http.ResponseWriter, r *http.Request) {
	var input updateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	input.Role = strings.TrimSpace(input.Role)
	if !validUserRole(input.Role) {
		http.Error(w, "Role must be administrator or clerk", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(authUserKey).(*authClaims)
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	targetID := chi.URLParam(r, "id")
	if input.Role == "clerk" && targetID == claims.UserID {
		http.Error(w, "You cannot change your own administrator role", http.StatusConflict)
		return
	}

	tx, err := app.db.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		http.Error(w, "Could not update user role", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	if input.Role == "clerk" {
		if _, err := tx.Exec(r.Context(), "LOCK TABLE users IN SHARE ROW EXCLUSIVE MODE"); err != nil {
			http.Error(w, "Could not update user role", http.StatusInternalServerError)
			return
		}
		var existingUser string
		if err := tx.QueryRow(r.Context(), "SELECT id::text FROM users WHERE id::text = $1", targetID).Scan(&existingUser); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Could not update user role", http.StatusInternalServerError)
			return
		}
		var remainingAdministrators int
		if err := tx.QueryRow(r.Context(), `
			SELECT count(*) FROM users
			WHERE role = 'administrator' AND id::text <> $1
		`, targetID).Scan(&remainingAdministrators); err != nil {
			http.Error(w, "Could not update user role", http.StatusInternalServerError)
			return
		}
		if remainingAdministrators == 0 {
			http.Error(w, "The last administrator cannot be changed to a clerk", http.StatusConflict)
			return
		}
	}

	tag, err := tx.Exec(r.Context(), "UPDATE users SET role = $1 WHERE id::text = $2", input.Role, targetID)
	if err != nil {
		http.Error(w, "Could not update user role", http.StatusInternalServerError)
		return
	}
	if tag.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		http.Error(w, "Could not update user role", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": targetID, "role": input.Role})
}

func validUserRole(role string) bool {
	return role == "administrator" || role == "clerk"
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
	}
}
