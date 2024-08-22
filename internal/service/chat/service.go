package chat

import (
	"github.com/Timofey335/chat-server/internal/repository"
	def "github.com/Timofey335/chat-server/internal/service"
)

var _ def.ChatService = (*serv)(nil)

type serv struct {
	chatRepository repository.ChatRepository
}

func NewService(chatRepository repository.ChatRepository) *serv {
	return &serv{
		chatRepository: chatRepository,
	}
}
