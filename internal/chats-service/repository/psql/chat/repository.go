package chat

import (
	"car-sell-buy-system/internal/chats-service/domain/chat"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
)

const (
	messageTableName = "messages"
	chatTableName    = "chats"
)

type Repository struct {
	*postgres.Postgres
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg,
	}
}

func (r *Repository) StoreMessage(ctx context.Context, chatId int64, dto chat.StoreMessageDTO) (chat.Message, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	sql, args, err := r.Builder.
		Insert(messageTableName).
		Columns(
			"chat_id",
			"sender_id",
			"text",
			"is_read",
			"image_url",
		).
		Values(
			chatId,
			userId,
			dto.Text,
			false,
			dto.CurrentImageUrl,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return chat.Message{}, fmt.Errorf("chat - Repository - StoreMessage - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	msg := chat.Message{
		ChatId:   chatId,
		SenderId: userId.(int64),
		Text:     dto.Text,
		IsRead:   false,
		Mine:     true,
	}
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&msg.Id)
	if err != nil {
		return chat.Message{}, fmt.Errorf("chat - Repository - StoreMessage - row.Scan: %w", err)
	}

	return msg, nil
}

func (r *Repository) StoreChat(ctx context.Context, dto chat.StoreChatDTO) (chat.Chat, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)
	if userId == nil {
		return chat.Chat{}, fmt.Errorf("user is not loggined")
	}

	sql, args, err := r.Builder.
		Insert(chatTableName).
		Columns(
			"buyer_id",
			"seller_id",
			"ad_id",
		).
		Values(
			userId,
			dto.SellerId,
			dto.AdId,
		).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return chat.Chat{}, fmt.Errorf("chat - Repository - StoreChat - r.Builder: %w", err)
	}

	cht := chat.Chat{
		BuyerId:  userId.(int64),
		SellerId: dto.SellerId,
		AdId:     dto.AdId,
	}
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&cht.Id)
	if err != nil {
		return chat.Chat{}, fmt.Errorf("chat - Repository - StoreChat - row.Scan: %w", err)
	}

	return cht, nil
}

func (r *Repository) ListChats(ctx context.Context) ([]chat.Chat, int64, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	builder := r.Builder.
		Select(
			sqlutil.TableColumn(chatTableName, "id"),
			sqlutil.TableColumn(chatTableName, "buyer_id"),
			sqlutil.TableColumn(chatTableName, "seller_id"),
			sqlutil.TableColumn(chatTableName, "ad_id"),
			sqlutil.TableColumn(chatTableName, "created_at"),
			"ads.title",
			"buyers.name",
			"sellers.name",
		).
		From(chatTableName).
		InnerJoin("ads on ads.id = chats.ad_id").
		InnerJoin("users as buyers on buyers.id = chats.buyer_id").
		InnerJoin("users as sellers on sellers.id = chats.seller_id").
		Where(squirrel.Or{
			squirrel.Eq{sqlutil.TableColumn(chatTableName, "buyer_id"): userId},
			squirrel.Eq{sqlutil.TableColumn(chatTableName, "seller_id"): userId},
		}).
		OrderBy("created_at DESC")

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("chat - Repository - List - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("chat - Repository - List - r.Pool.Query: %w", err)
	}

	var messages []chat.Chat
	var count int64
	for rows.Next() {
		var msg chat.Chat
		err = rows.Scan(
			&msg.Id,
			&msg.BuyerId,
			&msg.SellerId,
			&msg.AdId,
			&msg.CreatedAt,
			&msg.AdTitle,
			&msg.BuyerName,
			&msg.SellerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("chat - Repository - List - rows.Scan: %w", err)
		}
		messages = append(messages, msg)
		count++
	}

	return messages, count, nil
}

func (r *Repository) ListMessagesByChatId(ctx context.Context, chatId int64) ([]chat.Message, int64, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	builder := r.Builder.
		Select(
			sqlutil.TableColumn(messageTableName, "id"),
			sqlutil.TableColumn(messageTableName, "chat_id"),
			sqlutil.TableColumn(messageTableName, "sender_id"),
			sqlutil.TableColumn(messageTableName, "text"),
			sqlutil.TableColumn(messageTableName, "is_read"),
			sqlutil.TableColumn(messageTableName, "created_at"),
			sqlutil.TableColumn(messageTableName, "image_url"),
		).
		Column(squirrel.Alias(squirrel.Case().When(squirrel.Expr("messages.sender_id = ?", userId), "true").Else("false"), "mine")).
		From(messageTableName).
		Where(squirrel.Eq{"chat_id": chatId}).
		OrderBy("created_at ASC").
		Limit(200)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("AdRepository - List - r.Pool.Query: %w", err)
	}

	var messages []chat.Message
	var count int64
	for rows.Next() {
		var msg chat.Message
		err = rows.Scan(
			&msg.Id,
			&msg.ChatId,
			&msg.SenderId,
			&msg.Text,
			&msg.IsRead,
			&msg.CreatedAt,
			&msg.ImageUrl,
			&msg.Mine,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("ChatRepository - List - rows.Scan: %w", err)
		}
		messages = append(messages, msg)
		count++
	}

	return messages, count, nil
}
