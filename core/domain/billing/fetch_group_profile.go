package billing

import (
	"context"

	"gitlab.com/bloom42/bloom/core/api"
	"gitlab.com/bloom42/bloom/core/messages"
	"gitlab.com/bloom42/gobox/graphql"
)

func FetchGroupProfile(params messages.FetchGroupProfileParams) (messages.GroupBillingProfile, error) {
	client := api.Client()

	var resp messages.GroupBillingProfile
	req := graphql.NewRequest(`
	query($groupId: ID!) {
		group(id: $groupId) {
			id
			subscription {
				updatedAt
				usedStorage
				plan {
					id
					product
					price
					name
					storage
				}
			}
			paymentMethods {
				edges {
					node {
						id
						createdAt
						cardLast4
						cardExpirationMonth
						cardExpirationYear
						isDefault
					}
				}
			}
			invoices {
				edges {
					node {
						id
						createdAt
						amount
						stripeId
						stripeHostedUrl
						stripePdfUrl
						paidAt
					}
				}
			}
		}
		billingPlans {
			edges {
				node {
					id
					product
					price
					name
					description
					storage
				}
			}
		}
		stripePublicKey
	}
	`)
	req.Var("groupId", params.ID)

	err := client.Do(context.Background(), req, &resp)

	return resp, err
}
