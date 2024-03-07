package database

type Db interface {
	Get(key string) (string, error)
	Set(key string) error
	CloseDb()
}
