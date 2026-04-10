package apikeys

import (
	"context"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type sqrcQueries interface {
	InsertApiKey(ctx context.Context, arg db.InsertApiKeyParams) error
	CountApiKeyByUserID(ctx context.Context, userid pgtype.UUID) (int64, error)
}

type ApiKeySqrcRepository struct {
	queries sqrcQueries
}

func NewApiKeySqrcRepository(sqrc sqrcQueries) ApiKeySqrcRepository {
	return ApiKeySqrcRepository{queries: sqrc}
}

func (repo ApiKeySqrcRepository) count(ctx context.Context, userId uuid.UUID) int {
	pgUserID := pgtype.UUID{Bytes: userId, Valid: true}
	count, err := repo.queries.CountApiKeyByUserID(ctx, pgUserID)
	if err != nil {
		return 0
	}
	return int(count)
}

func (repo ApiKeySqrcRepository) save(ctx context.Context, apiKey hashedApiKey) error {
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

	return repo.queries.InsertApiKey(ctx, db.InsertApiKeyParams{
		ApikeyID:       apikeyID,
		KeyHash:        apiKey.hashedApiKey,
		UserID:         userID,
		KeyPlainSuffix: apiKey.plainApiKeySuffix,
		ExpiredAt:      expiredAt,
	})
}
