package chat

import (
	"context"

	"github.com/Timofey335/chat-server/internal/model"
)

func (s *serv) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	id, err := s.chatRepository.CreateChat(ctx, chat)
	if err != nil {
		return 0, err
	}

	return id, nil
}
