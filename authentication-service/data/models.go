package data

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var dbTimeout = 5 * time.Second // seconds

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password,omitempty"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, active, created_at, updated_at FROM users Order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `Select * from users where email = $1`
	row := db.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `Select * from users where id = $1`
	row := db.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `Update users set email = $1, first_name = $2, last_name = $3, active = $4, updated_at = $5 where id = $6`
	_, err := db.ExecContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `Delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newId int

	query := `Insert into users (email, first_name, last_name, password, active, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7) returning id`
	err := db.QueryRowContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
		u.Active,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return err
	}
	u.ID = newId
	return nil
}

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `Update users set password = $1, updated_at = $2 where id = $3`
	_, err = db.ExecContext(ctx, query, HashedPassword, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PasswordMatches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(text))
	if err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
