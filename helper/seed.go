package helper

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
	"github.com/michaelyusak/xyz-kredit-plus/repository"
)

type entry struct {
	account      entity.Account
	consumer     entity.Consumer
	accountLimit entity.AccountLimit
}

func accountSeed(ctx context.Context, repo repository.AccountRepository, account entity.Account) (int64, error) {
	accountId, err := repo.InsertAccount(ctx, account)
	if err != nil {
		return 0, fmt.Errorf("[helper][accountSeed][repo.InsertAccount] Error: %w | email: %s", err, account.Email)
	}

	return accountId, nil
}

func consumerSeed(ctx context.Context, repo repository.ConsumerRepository, consumer entity.Consumer) error {
	err := repo.InsertConsumer(ctx, consumer)
	if err != nil {
		return fmt.Errorf("[helper][consumerSeed][repo.InsertConsumer] Error: %w | account_id: %v", err, consumer.AccountId)
	}

	return nil
}

func accountLimitSeed(ctx context.Context, repo repository.AccountLimitRepository, accountLimit entity.AccountLimit) error {
	err := repo.InsertLimit(ctx, accountLimit)
	if err != nil {
		return fmt.Errorf("[helper][accountLimitSeed][repo.InsertLimit] Error: %w | account_id: %v", err, accountLimit.AccountId)
	}

	return nil
}

func Seed(db *sql.DB) error {
	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("[helper][Seed][db.Begin] Error: %w", err)
	}

	accountTx := repository.NewAccountRepositoryMysql(tx)
	consumerTx := repository.NewConsumerRepositoryMysql(tx)
	accountLimitTx := repository.NewAccountLimitRepositoryMysql(tx)

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	seeds := []entry{
		{
			account: entity.Account{
				Email:    "budi@example.com",
				Password: "$2a$10$XjYCovsfMpJfehkoQc6F..ZwHPuSv2aV8LWWPiDPdLJetSLEZ2S5e",
			},
			consumer: entity.Consumer{
				IdentityNumber: "123456789123456789",
				FullName:       "nama lengkap",
				LegalName:      "nama asli",
				PlaceOfBirth:   "jakarta",
				DateOfBirth:    "12 Sep 2024",
				Salary:         1000000000,
			},
			accountLimit: entity.AccountLimit{
				Limit1M: 100000,
				Limit2M: 200000,
				Limit3M: 500000,
				Limit4M: 700000,
			},
		},
		{
			account: entity.Account{
				Email:    "annisa@example.com",
				Password: "$2a$10$XjYCovsfMpJfehkoQc6F..ZwHPuSv2aV8LWWPiDPdLJetSLEZ2S5e",
			},
			consumer: entity.Consumer{
				IdentityNumber: "123456789123456789",
				FullName:       "nama lengkap",
				LegalName:      "nama asli",
				PlaceOfBirth:   "jakarta",
				DateOfBirth:    "12 Sep 2024",
				Salary:         1000000000,
			},
			accountLimit: entity.AccountLimit{
				Limit1M: 1000000,
				Limit2M: 1200000,
				Limit3M: 1500000,
				Limit4M: 2000000,
			},
		},
	}

	for _, seed := range seeds {
		existingAcc, err := accountTx.GetAccountByEmail(ctx, seed.account.Email, false)
		if err != nil {
			return fmt.Errorf("[helper][Seed][accountTx.GetAccountByEmail] Error: %w", err)
		}

		if existingAcc == nil {
			accId, err := accountSeed(ctx, accountTx, seed.account)
			if err != nil {
				return fmt.Errorf("[helper][Seed][accountSeed] Error: %w", err)
			}

			seed.account.Id = accId
			seed.consumer.AccountId = accId
			seed.accountLimit.AccountId = accId
		} else {
			seed.account.Id = existingAcc.Id
			seed.consumer.AccountId = existingAcc.Id
			seed.accountLimit.AccountId = existingAcc.Id
		}

		existingCon, err := consumerTx.GetConsumerByAccountId(ctx, seed.consumer.AccountId, false)
		if err != nil {
			return fmt.Errorf("[helper][Seed][consumerTx.GetConsumerByAccountId] Error: %w", err)
		}

		if existingCon == nil {
			err := consumerSeed(ctx, consumerTx, seed.consumer)
			if err != nil {
				return fmt.Errorf("[helper][Seed][accountSeed] Error: %w", err)
			}
		}

		existingAccLim, err := accountLimitTx.GetAccountLimitByAccountId(ctx, seed.accountLimit.AccountId, false)
		if err != nil {
			return fmt.Errorf("[helper][Seed][accountLimitTx.GetAccountLimitByAccountId] Error: %w", err)
		}

		if existingAccLim == nil {
			err := accountLimitSeed(ctx, accountLimitTx, seed.accountLimit)
			if err != nil {
				return fmt.Errorf("[helper][Seed][accountLimitSeed] Error: %w", err)
			}
		}
	}

	return nil
}
