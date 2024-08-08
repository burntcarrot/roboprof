package storage

type Storage interface {
	Write(b []byte) error
}
