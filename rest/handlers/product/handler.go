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
	}
}

type MessageStore struct {
	mu       sync.RWMutex
	messages map[string]string
}

func NewMessageStore() *MessageStore {
	return &MessageStore{
		messages: make(map[string]string),
	}
}

func (s *MessageStore) Set(barcode, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[barcode] = msg
}

func (s *MessageStore) Get(barcode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msg, ok := s.messages[barcode]
	return msg, ok
}
