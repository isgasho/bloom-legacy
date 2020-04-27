package users

import (
	"context"
	"time"

	"gitlab.com/bloom42/bloom/cmd/bloom/server/db"
	"gitlab.com/bloom42/lily/crypto"
	"gitlab.com/bloom42/lily/rz"
	"gitlab.com/bloom42/lily/uuid"
)

type VerifyPendingUserParams struct {
	PendingUserID uuid.UUID
	Code          string
}

func VerifyPendingUser(ctx context.Context, params VerifyPendingUserParams) (err error) {
	logger := rz.FromCtx(ctx)
	var pendingUser PendingUser

	// verify pending user
	tx, err := db.DB.Beginx()
	if err != nil {
		logger.Error("users.VerifyPendingUser: Starting transaction", rz.Err(err))
		err = NewError(ErrorVerifyingPendingUser)
		return
	}

	err = tx.Get(&pendingUser, "SELECT * FROM pending_users WHERE id = $1 FOR UPDATE", params.PendingUserID)
	if err != nil {
		tx.Rollback()
		logger.Error("users.VerifyPendingUser: getting pending user", rz.Err(err),
			rz.String("pending_user.id", params.PendingUserID.String()))
		err = NewError(ErrorVerifyingPendingUser)
		return
	}

	if pendingUser.FailedAttempts+1 >= 5 {
		tx.Rollback()
		return NewError(ErrorMaximumVerificationTrialsReached)
	}

	now := time.Now().UTC()
	since := now.Sub(pendingUser.UpdatedAt)
	if since >= 30*time.Minute {
		tx.Rollback()
		return NewError(ErrorRegistrationCodeExpired)
	}

	if !crypto.VerifyPasswordHash([]byte(params.Code), pendingUser.VerificationCodeHash) {
		tx.Rollback()
		err = NewError(ErrorRegistrationCodeIsNotValid)
		tx, _ := db.DB.Beginx()
		if tx != nil {
			err2 := failPendingUserVerification(ctx, tx, &pendingUser)
			if err2 != nil {
				tx.Rollback()
				err = NewError(ErrorVerifyingPendingUser)
				return
			}
			tx.Commit()
		}
		return
	}

	now = time.Now().UTC()
	pendingUser.VerifiedAt = &now
	pendingUser.UpdatedAt = now

	_, err = tx.Exec("UPDATE pending_users SET verified_at = $1, updated_at = $1 WHERE id = $2",
		now, pendingUser.ID)
	if err != nil {
		tx.Rollback()
		logger.Error("users.VerifyPendingUser: error verifying pending user", rz.Err(err),
			rz.String("pending_user.id", pendingUser.ID.String()))
		return NewError(ErrorVerifyingPendingUser)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("users.VerifyPendingUser: Committing transaction", rz.Err(err))
		err = NewError(ErrorVerifyingPendingUser)
		return
	}

	return
}
