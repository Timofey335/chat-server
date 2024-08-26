package chat

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/chat-server/internal/model"
	"github.com/Timofey335/chat-server/internal/repository"
)

const (
	tableChats = "chats"

	idColumn    = "id"
	nameColumn  = "name"
	usersColumn = "users"

	tableMessages = "messages"

	userIdColumn    = "user_id"
	chatIdColumn    = "chat_id"
	textColumn      = "text"
	createdAtColumn = "created_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewChat(db *pgxpool.Pool) repository.ChatRepository {
	return &repo{db: db}
}

// CreateChat - create a new chat
func (r *repo) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	var chatId int64

	builderInsert := sq.Insert(tableChats).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, usersColumn).
		Values(chat.Name, chat.Users).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	if err = r.db.QueryRow(ctx, query, args...).Scan(&chatId); err != nil {
		return 0, err
	}

	log.Println(color.BlueString("created chat: %v, with ctx: %v", chat.Name, ctx))

	return chatId, nil
}

// DeleteChat - delete the chat by id
func (r *repo) DeleteChat(ctx context.Context, chatId int64) (*emptypb.Empty, error) {
	var id int64

	builderDelete := sq.Delete(tableChats).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: chatId}).
		Suffix("RETURNING id")

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Println(color.HiMagentaString("error while deleting the chat: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted chat: id %d, with ctx: %v", id, ctx))

	return &emptypb.Empty{}, nil
}

// SendMessage - send message to the server
func (r *repo) SendMessage(ctx context.Context, message *model.Message) (*emptypb.Empty, error) {
	builderInsert := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(userIdColumn, chatIdColumn, textColumn, createdAtColumn).
		Values(message.UserId, message.ChatId, message.Text, message.CreatedAt)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Println(color.HiMagentaString("error while sending the message: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.BlueString("sent message: %s from user id %s at %s, with ctx: %v", message.Text, message.UserId, message.CreatedAt, ctx))

	return &emptypb.Empty{}, nil
}
