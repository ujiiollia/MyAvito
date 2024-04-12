package postgresql

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type users struct {
	dbu *postgres
}

func NewUser(pgl *postgres) *users {
	return &users{dbu: pgl}
}

func (u *users) GetRole(userId string, ctx *gin.Context) (string, error) {
	const el = "postgresql.users.GetRole"

	conn, err := u.dbu.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", el, err)
	}
	defer conn.Release()
	var role string
	err = conn.QueryRow(ctx,
		`SELECT is_admin FROM "user" WHERE id = $1;`,
		userId).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("%s: %w", el, err)
	}
	return role, nil
}
