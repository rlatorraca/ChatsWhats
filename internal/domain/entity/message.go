package entity

import (
	"errors"
	"github.com/google/uuid"
	tiktoken_go "github.com/j178/tiktoken-go"
	"time"
)

type Message struct {
	ID       string
	Role     string // User, system or assistant
	Content  string
	Tokens   int // quantidade de tokens da mensagem
	Model    *Model
	CreateAt time.Time
}

func NewMessage(role, content string, model *Model) (*Message, error) {
	totalTokens := tiktoken_go.CountTokens(model.GetModelName(), content)
	msg := &Message{
		ID:       uuid.New().String(),
		Role:     role,
		Content:  content,
		Tokens:   totalTokens,
		Model:    model,
		CreateAt: time.Now(),
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (m *Message) Validate() error {
	if m.Role != "user" && m.Role != "system" && m.Role != "assistant" {
		return errors.New("Invalid role")
	}

	if m.Content == "" {
		return errors.New("Invalid content")
	}

	if m.CreateAt.IsZero() {
		return errors.New("Invalid 'Created At'")
	}

	return nil
}

func (m *Message) GetQuantityOfTokens() int {
	return m.Tokens
}
