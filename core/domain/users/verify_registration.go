package users

import (
	"context"
	"strings"

	"gitlab.com/bloom42/bloom/core/api"
	"gitlab.com/bloom42/bloom/core/api/model"
	"gitlab.com/bloom42/gobox/graphql"
)

func VerifyUser(params VerifyRegistrationParams) (bool, error) {
	client := api.Client()
	ret := false

	code := strings.ToLower(params.Code)

	input := model.VerifyUserInput{
		ID:   params.ID,
		Code: code,
	}
	var resp struct {
		VerifyUser *bool `json:"verifyUser"`
	}
	req := graphql.NewRequest(`
        mutation ($input: VerifyUserInput!) {
			verifyUser(input: $input)
		}
	`)
	req.Var("input", input)

	err := client.Do(context.Background(), req, &resp)
	if err != nil {
		return ret, err
	}

	if resp.VerifyUser != nil {
		ret = *resp.VerifyUser
	}

	return ret, err
}
