package contacts

import (
	"context"

	"gitlab.com/bloom42/bloom/core/domain/kernel"
	"gitlab.com/bloom42/bloom/core/domain/objects"
	"gitlab.com/bloom42/bloom/core/messages"
)

func FindContacts(params messages.ContactsFindParams) (Contacts, error) {
	ret := Contacts{Contacts: []objects.Object{}}

	objects, err := objects.FindObjectsByType(context.Background(), nil, kernel.OBJECT_TYPE_CONTACT, params.GroupID)
	if err != nil {
		return ret, err
	}

	ret.Contacts = objects

	return ret, nil
}
