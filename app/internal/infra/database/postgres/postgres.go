package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
)

// Config contém as configurações necessárias para conectar ao banco de dados PostgreSQL.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB cria uma nova conexão com o banco de dados PostgreSQL usando as configurações fornecidas.
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := otelsql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão com o banco de dados: %w", err)
	}

	// Configura o pool de conexões do banco de dados.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Testa a conexão com o banco de dados.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	return db, nil
}
