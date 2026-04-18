package apikeys

import (
	"context"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ApiKeySqrcRepository struct{}

func NewApiKeySqrcRepository() ApiKeySqrcRepository {
	return ApiKeySqrcRepository{}
}

func (repo ApiKeySqrcRepository) countWithTx(ctx context.Context, tx pgx.Tx, userId uuid.UUID) int {
	q := db.New(tx)

	pgUserID := pgtype.UUID{Bytes: userId, Valid: true}
	count, err := q.CountApiKeyByUserID(ctx, pgUserID)
	if err != nil {
		return 0
	}
	return int(count)
}

func (repo ApiKeySqrcRepository) saveWithTx(ctx context.Context, tx pgx.Tx, apiKey hashedApiKey) error {
	q := db.New(tx)

	apikeyID := pgtype.UUID{Bytes: apiKey.apiKeyID, Valid: true}
	userID := pgtype.UUID{Bytes: apiKey.userID, Valid: true}

	var expiredAt pgtype.Timestamptz
	if apiKey.expiredAt == nil {
		expiredAt = pgtype.Timestamptz{
			InfinityModifier: pgtype.Infinity,
			Valid:            true,
		}
	} else {
		expiredAt = pgtype.Timestamptz{
			Time:  *apiKey.expiredAt,
			Valid: true,
		}
	}

	return q.InsertApiKey(ctx, db.InsertApiKeyParams{
		ApikeyID:       apikeyID,
		KeyHash:        apiKey.hashedApiKey,
		UserID:         userID,
		KeyPlainSuffix: apiKey.plainApiKeySuffix,
		ExpiredAt:      expiredAt,
	})
}
