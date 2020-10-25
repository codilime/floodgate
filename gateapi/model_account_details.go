/*
 * Spinnaker API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type AccountDetails struct {
	AccountType string `json:"accountType,omitempty"`
	AccountId string `json:"accountId,omitempty"`
	PrimaryAccount bool `json:"primaryAccount,omitempty"`
	ChallengeDestructiveActions bool `json:"challengeDestructiveActions,omitempty"`
	Environment string `json:"environment,omitempty"`
	CloudProvider string `json:"cloudProvider,omitempty"`
	Name string `json:"name,omitempty"`
	Permissions map[string][]string `json:"permissions,omitempty"`
	Type_ string `json:"type,omitempty"`
	RequiredGroupMembership []string `json:"requiredGroupMembership,omitempty"`
}
