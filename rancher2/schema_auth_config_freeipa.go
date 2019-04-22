package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const AuthConfigFreeIpaName = "freeipa"

//Schemas

func authConfigFreeIpaFields() map[string]*schema.Schema {
	return authConfigLdapFields()
}
