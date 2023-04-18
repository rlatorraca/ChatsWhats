package entity

import (
	"errors"
	"github.com/google/uuid"
)

type ChatConfig struct {
	Model           *Model
	Temperature     float32  // 0.0 a 1.0 - Mais preciso ou menos preciso na pergunta (0 = mais preciso, 1 = menos preciso)
	TopP            float32  // 0.0 a 1.0 - Mai conservador na escolha dos tokens (0 = mais conservador, 1 = menos conservador)
	N               int      // numero de mensagens geradas
	Stop            []string // a cadeia de palavras usada para parar o chat
	MaxTokens       int      // numero de tokens que chamada no chat pode aceitar
	PresencePenalty float32  // -2.0 a 2.0 - isso afeta como o modelo penaliza novos tokens com base no fato de terem aparecido ou não no texto até o momento.
	// Valores positivos aumentarão a probabilidade do modelo falar sobre novos tópicos ao penalizar novos tokens que já foram usados.
	FrequencyPenalty float32 // -2.0 a 2.0 - e afeta como o modelo penaliza novos tokens com base em sua frequência existente no texto.
	// valores positivos diminuirão a probabilidade de o modelo repetir a mesma linha textualmente, penalizando novos tokens que já foram usados com frequência.
}

type Chat struct {
	ID              string
	UserId          string
	Status          string
	TokenUsage      int
	Config          *ChatConfig
	InitialMessage  *Message
	Messages        []*Message
	DeletedMessages []*Message
}

func NewChat(userID string, initialSystemMessage *Message, chatConfig *ChatConfig) (*Chat, error) {
	chat := &Chat{
		ID:             uuid.New().String(),
		UserId:         userID,
		InitialMessage: initialSystemMessage,
		Status:         "active",
		Config:         chatConfig,
		TokenUsage:     0,
	}
	err := chat.AddMessage(initialSystemMessage)
	if err != nil {
		return nil, err
	}

	if err := chat.Validate(); err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *Chat) Validate() error {
	if c.UserId == "" {
		return errors.New("user id is empty")
	}
	if c.Status != "active" && c.Status != "ended" {
		return errors.New("invalid status")
	}
	if c.Config.Temperature < 0 || c.Config.Temperature > 2 {
		return errors.New("invalid temperature")
	}
	// ... more validations for config
	return nil
}

func (c *Chat) AddMessage(m *Message) error {
	if c.Status == "ended" {
		return errors.New("chat is ended. no more messages allowed")
	}
	for {
		if c.Config.Model.GetMaxTokens() >= m.GetQuantityOfTokens()+c.TokenUsage {
			c.Messages = append(c.Messages, m)
			c.RefreshTokenUsage()
			break
		}
		c.DeletedMessages = append(c.DeletedMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokenUsage()
	}
	return nil
}

func (c *Chat) RefreshTokenUsage() {
	c.TokenUsage = 0
	for m := range c.Messages {
		c.TokenUsage += c.Messages[m].GetQuantityOfTokens()
	}
}

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessages() int {
	return len(c.Messages)
}

func (c *Chat) End() {
	c.Status = "ended"
}
