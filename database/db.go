package database

type Db interface {
	Get()
	Set()
}

type PostgresDb struct {
}

func NewPostgresDb() *PostgresDb {
	return &PostgresDb{}
}
