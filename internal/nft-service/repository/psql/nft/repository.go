package nft

import (
	"car-sell-buy-system/internal/nft-service/domain/nft"
	"car-sell-buy-system/pkg/postgres"
	"car-sell-buy-system/pkg/sqlutil"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"math/rand"
)

const (
	nftTableName = "nfts"
)

type Repository struct {
	*postgres.Postgres
}

func NewRepository(pg *postgres.Postgres) *Repository {
	return &Repository{
		pg,
	}
}

func (r *Repository) StoreNft(ctx context.Context, dto nft.StoreNftDTO) (nft.Nft, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	sql, args, err := r.Builder.
		Insert(nftTableName).
		Columns(
			"token_id",
			"vin",
			"metadata_url",
			"is_minted",
		).
		Values(
			rand.Int(),
			dto.Vin,
			dto.MetadataUrl,
			true,
		).
		Suffix("RETURNING token_id").
		ToSql()
	if err != nil {
		return nft.Nft{}, fmt.Errorf("nft - Repository - StoreNft - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	nftToken := nft.Nft{
		Vin:         dto.Vin,
		MetadataUrl: dto.MetadataUrl,
		IsMinted:    false,
	}
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&nftToken.TokenId)
	if err != nil {
		return nft.Nft{}, fmt.Errorf("nft - Repository - StoreNft - row.Scan: %w", err)
	}

	return nftToken, nil
}

func (r *Repository) GetNftByVin(ctx context.Context, vin string) (nft.Nft, error) {
	userId := ctx.Value("userId")
	fmt.Println("userId: ", userId)

	sql, args, err := r.Builder.
		Select(
			sqlutil.TableColumn(nftTableName, "token_id"),
			sqlutil.TableColumn(nftTableName, "vin"),
			sqlutil.TableColumn(nftTableName, "metadata_url"),
			sqlutil.TableColumn(nftTableName, "is_minted"),
			sqlutil.TableColumn(nftTableName, "created_at"),
		).
		From(nftTableName).
		Where(squirrel.Eq{sqlutil.TableColumn(nftTableName, "vin"): vin}).
		ToSql()
	if err != nil {
		return nft.Nft{}, fmt.Errorf("nft - Repository - GetByTransactionId - r.Builder: %w", err)
	}
	fmt.Println("sql: ", sql)
	fmt.Println("args: ", args)

	var nftToken nft.Nft
	err = r.Pool.
		QueryRow(ctx, sql, args...).
		Scan(
			&nftToken.TokenId,
			&nftToken.Vin,
			&nftToken.MetadataUrl,
			&nftToken.IsMinted,
			&nftToken.CreatedAt,
		)
	if err != nil {
		return nft.Nft{}, fmt.Errorf("nft - Repository - GetByTransactionId - row.Scan: %w", err)
	}

	return nftToken, nil
}
