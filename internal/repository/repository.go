package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/chat-server/internal/model"
)

// ChatRepository - интерфейс repo слоя
type ChatRepository interface {
	CreateChat(ctx context.Context, chat *model.Chat) (int64, error)
	DeleteChat(ctx context.Context, chatId int64) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error)
}
