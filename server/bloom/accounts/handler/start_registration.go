package handler

import (
	"context"
	"time"

	"github.com/twitchtv/twirp"
	rpc "gitlab.com/bloom42/bloom/common/rpc/accounts"
	"gitlab.com/bloom42/bloom/server/api/apictx"
	"gitlab.com/bloom42/bloom/server/bloom/accounts"
	"gitlab.com/bloom42/bloom/server/db"
	"gitlab.com/bloom42/bloom/server/services/util"
	"gitlab.com/bloom42/libs/rz-go"
)

func (s Handler) StartRegistration(ctx context.Context, params *rpc.StartRegistrationParams) (*rpc.RegistrationStarted, error) {
	logger := rz.FromCtx(ctx)
	apiCtx, ok := ctx.Value(apictx.Key).(*apictx.Context)
	if !ok {
		logger.Error("accounts.StartRegistration: error getting apiCtx from context")
		return &rpc.RegistrationStarted{}, twirp.InternalError("error creating account")
	}
	if apiCtx.AuthenticatedAccount != nil {
		twerr := twirp.NewError(twirp.PermissionDenied, "Must not be authenticated")
		return &rpc.RegistrationStarted{}, twerr
	}

	// sleep to prevent spam and bruteforce
	time.Sleep(time.Duration(util.InsecureRandomInt(500, 800)) * time.Millisecond)

	// created pending account
	tx, err := db.DB.Beginx()
	if err != nil {
		logger.Error("accounts.StartRegistration: Starting transaction", rz.Err(err))
		return &rpc.RegistrationStarted{}, twirp.InternalError("Error creating account. Please try again.")
	}

	newPendingAccount, verificationCode, twerr := accounts.CreatePendingAccount(ctx, tx, params.DisplayName, params.Email)
	if twerr != nil {
		tx.Rollback()
		return &rpc.RegistrationStarted{}, twerr
	}

	err = accounts.SendAccountVerificationCode(params.Email, params.DisplayName, verificationCode)
	if err != nil {
		tx.Rollback()
		logger.Error("accounts.StartRegistration: Sending confirmation email", rz.Err(err))
		return &rpc.RegistrationStarted{}, twirp.InternalError("Error sending confirmation email. Please try again.")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("accounts.StartRegistration: Committing transaction", rz.Err(err))
		return &rpc.RegistrationStarted{}, twirp.InternalError("Error creating account. Please try again.")
	}

	ret := rpc.RegistrationStarted{
		Id: newPendingAccount.ID,
	}
	return &ret, nil
}
