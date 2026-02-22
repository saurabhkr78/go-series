package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"learned/domain"
)

type userRepository struct {
	db *pgxpool.Pool
}

// constructor returns INTERFACE
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
	).Scan(&user.ID)
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
func (r *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}
