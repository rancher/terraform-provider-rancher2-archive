package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Types

type azureCredentialConfig struct {
	ClientID       string `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret   string `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	SubscriptionID string `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
}

//Schemas

func cloudCredentialAzureFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account ID",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account password",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Azure Subscription ID",
		},
	}

	return s
}
