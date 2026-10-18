// Dominio de negocio para representar uma conta no sistema.
package entity

import (
	"time"
)

// Account representa uma conta no sistema.
type Account struct {
	ID        int
	Document  string
	CreatedAt time.Time
}