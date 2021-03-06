package mutation

import (
	"context"

	"gitlab.com/bloom42/bloom/server/api"
	"gitlab.com/bloom42/bloom/server/api/graphql/model"
	"gitlab.com/bloom42/bloom/server/domain/groups"
	"gitlab.com/bloom42/bloom/server/domain/sync"
)

// AcceptGroupInvitation accepts a group invitaiton
func (resolver *Resolver) AcceptGroupInvitation(ctx context.Context, input model.AcceptGroupInvitationInput) (ret *model.Group, err error) {
	params := groups.AcceptInvitationParams{
		InvitationID:       input.InvitationID,
		EncryptedMasterKey: input.EncryptedMasterKey,
		MasterKeyNonce:     input.MasterKeyNonce,
	}
	group, err := resolver.groupsService.AcceptInvitation(ctx, params)
	if err != nil {
		err = api.NewError(err)
		return
	}

	state := sync.EncodeState(group.State)
	ret = &model.Group{
		ID:          &group.ID,
		CreatedAt:   &group.CreatedAt,
		Name:        group.Name,
		Description: group.Description,
		AvatarURL:   nil,
		State:       &state,
	}
	return ret, nil
}
