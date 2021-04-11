package postgres

import (
	"github.com/toundaherve/sad-api/user"
)

type PostgresDB struct{}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (p *PostgresDB) CreateUser(u *user.User) error {
	return nil
}

func (p *PostgresDB) GetUserByEmail(email string) (*user.User, error) {
	return nil, nil
}
