package repository

import (
	"context"
	"go-project-template/internal/domain"
	"go-project-template/internal/utils"

	"github.com/google/uuid"
)

type postgresUserRepository struct {
	conn Connection
}

// NewUserRepository returns a new [UserRepository].
func NewUserRepository(conn Connection) domain.UserRepository {
	return &postgresUserRepository{conn: conn}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uu []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.UUID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
		); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return uu, nil
}

func (p *postgresUserRepository) GetByID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	query := `
		SELECT *
		FROM users
		WHERE uuid = $1`

	srs, err := p.fetch(ctx, query, uuid)

	if err != nil {
		return domain.User{}, err
	}
	if len(srs) == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return srs[0], nil
}

func (p *postgresUserRepository) CreateOrUpdate(ctx context.Context, u *domain.User) (*domain.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users (uuid, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(uuid) DO UPDATE
		SET first_name = $2, last_name = $3, email = $4;
		`

	_, err := p.conn.Exec(
		ctx,
		query,
		u.UUID,
		u.NormalizedFirstName(),
		u.NormalizedLastName(),
		&u.Email,
	)

	return u, err
}

func (p *postgresUserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM users WHERE uuid = $1`
	_, err := p.conn.Exec(ctx, query, uuid)
	return err
}

func (p *postgresUserRepository) GetList(ctx context.Context, pq *utils.PaginationQuery) (*utils.PaginationResponse[domain.User], error) {
	var count int
	countQuery := "SELECT count(uuid) FROM users;"
	if err := p.conn.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		return nil, err
	}

	if count == 0 {
		return utils.DefaultPaginationResponse[domain.User](pq), nil
	}

	query := "SELECT * FROM users ORDER BY $1::text OFFSET $2 LIMIT $3"
	rows, err := p.conn.Query(ctx, query, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uu []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.UUID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
		); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}

	return utils.PaginatedResponse(count, pq, uu), nil
}
