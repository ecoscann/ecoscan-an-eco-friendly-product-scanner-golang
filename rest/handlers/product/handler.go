package product

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

type ProductHandler struct {
	DB *sqlx.DB
	Store  *MessageStore
}

func NewProductHandler(db *sqlx.DB) *ProductHandler {
	return &ProductHandler{
		DB: db,
		Store: NewMessageStore(),
	}
}

type MessageStore struct {
    mu   sync.RWMutex
    data map[string]string
}

func NewMessageStore() *MessageStore {
    return &MessageStore{
        data: make(map[string]string),
    }
}

func (s *MessageStore) Set(barcode, msg string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[barcode] = msg
}

func (s *MessageStore) Get(barcode string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    msg, ok := s.data[barcode]
    return msg, ok
}
