package billing

import (
	"context"

	"gitlab.com/bloom42/bloom/core/api"
	"gitlab.com/bloom42/bloom/core/messages"
	"gitlab.com/bloom42/gobox/graphql"
)

func FetchMyProfile() (messages.MyBillingProfile, error) {
	client := api.Client()

	var resp messages.MyBillingProfile
	req := graphql.NewRequest(`
	query {
		me {
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

	err := client.Do(context.Background(), req, &resp)

	return resp, err
}
