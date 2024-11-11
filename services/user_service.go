package services

import (
	"context"
	"errors"
	"go-webserver-performance-test/models/data"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	DB *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{DB: db}
}

func (m *UserService) GetAllUsers() (*[]data.User, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT id, username, email, password_hash FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []data.User
	for rows.Next() {
		var user data.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}
func (m *UserService) GetUserByID(userID string) (*data.User, error) {
	searchID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var user data.User

	err = m.DB.QueryRow(context.Background(), "SELECT id, username, email, password_hash FROM users WHERE id = $1", searchID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserService) GetUserByUsername(username string) (*data.User, error) {
	var user data.User

	err := m.DB.QueryRow(context.Background(), "SELECT id, username, email, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (m *UserService) GetUserByEmail(email string) (*data.User, error) {
	var user data.User

	err := m.DB.QueryRow(context.Background(), "SELECT id, username, email, password_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user an returns the ID of the new user or an error
// If the error isn't nil, the ID will be 0
func (m *UserService) CreateUser(ctx context.Context, user *data.User) (*uuid.UUID, error) {
	var id uuid.UUID
	// now need to generate a UUID for the user, because postgres will create a UUID for us
	query := `INSERT INTO users (id, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`
	err := m.DB.QueryRow(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (m *UserService) UpdateUser(ctx context.Context, user *data.User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	commandTag, err := m.DB.Exec(ctx, query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("No rows affected, user not found")
	}
	return nil
}
func (m *UserService) DeleteUser(userID uuid.UUID) error {
	commandTag, err := m.DB.Exec(context.TODO(), "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("Now rows affected, user not found")
	}
	return nil
}
