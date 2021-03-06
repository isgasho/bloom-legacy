package users

import (
	"context"

	"gitlab.com/bloom42/bloom/core/api"
	"gitlab.com/bloom42/bloom/core/api/model"
	"gitlab.com/bloom42/gobox/graphql"
)

func FetchMySessions() (model.User, error) {
	client := api.Client()

	var resp struct {
		Me model.User `json:"me"`
	}
	req := graphql.NewRequest(`
	query {
		me {
			sessions {
				nodes {
					id
					createdAt
					device {
						os
						type
					}
				}
			}
		}
	}
	`)

	err := client.Do(context.Background(), req, &resp)

	return resp.Me, err
}
