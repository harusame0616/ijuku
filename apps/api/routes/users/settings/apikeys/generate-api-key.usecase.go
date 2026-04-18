package apikeys

import (
	"context"
	"hash/fnv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type transactionRunner interface {
	RunInTransaction(ctx context.Context, f func(tx pgx.Tx) error) error
	AcquireAdvisoryLock(ctx context.Context, tx pgx.Tx, lockKey int64) error
}

type generateApiKeyRepository interface {
	countWithTx(ctx context.Context, tx pgx.Tx, userId uuid.UUID) int
	saveWithTx(ctx context.Context, tx pgx.Tx, apiKey hashedApiKey) error
}

type generateApiKeyUsecase struct {
	apiKeyRepository generateApiKeyRepository
	txRunner         transactionRunner
}

func NewGenerateApiKeyUsecase(repository generateApiKeyRepository, txRunner transactionRunner) generateApiKeyUsecase {
	return generateApiKeyUsecase{
		apiKeyRepository: repository,
		txRunner:         txRunner,
	}
}

type generateApiKeyExecuteResult struct {
	Apikey string `json:"apikey"`
}

func apiKeyGenLockKey(userID uuid.UUID) int64 {
	h := fnv.New64a()
	h.Write([]byte("APIKEY_GEN_LOCK_" + userID.String()))
	return int64(h.Sum64())
}

func (usecase *generateApiKeyUsecase) Execute(ctx context.Context, userId uuid.UUID, expiredAt *time.Time) (generateApiKeyExecuteResult, error) {
	var result generateApiKeyExecuteResult

	err := usecase.txRunner.RunInTransaction(ctx, func(tx pgx.Tx) error {
		if err := usecase.txRunner.AcquireAdvisoryLock(ctx, tx, apiKeyGenLockKey(userId)); err != nil {
			return err
		}

		if usecase.apiKeyRepository.countWithTx(ctx, tx, userId) >= apiKeyMaxCount {
			return ErrApiKeyCountExceedsLimit
		}

		hashedApiKey, plainApiKey := NewHashedApiKey(NewHashedApiKeyParams{
			UserID:    userId,
			ExpiredAt: expiredAt,
		})

		if err := usecase.apiKeyRepository.saveWithTx(ctx, tx, hashedApiKey); err != nil {
			return err
		}

		result = generateApiKeyExecuteResult{Apikey: plainApiKey}
		return nil
	})

	return result, err
}
