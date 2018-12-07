package rancher2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// CLIConfig used to store data from file.
type CLIConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	CACerts   string `json:"caCerts"`
	URL       string `json:"url"`
	Project   string `json:"project"`
	Path      string `json:"path,omitempty"`
	Insecure  bool   `json:"insecure,omitempty"`
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", ""),
				Description: descriptions["api_url"],
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},
			"token_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_TOKEN_KEY", ""),
				Description: descriptions["token_key"],
			},
			"ca_certs": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CA_CERTS", ""),
				Description: descriptions["ca_certs"],
			},
			"config": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CLIENT_CONFIG", ""),
				Description: descriptions["config"],
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_INSECURE", false),
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher2_auth_config_activedirectory":   resourceRancher2AuthConfigActiveDirectory(),
			"rancher2_auth_config_adfs":              resourceRancher2AuthConfigADFS(),
			"rancher2_auth_config_azuread":           resourceRancher2AuthConfigAzureAD(),
			"rancher2_auth_config_freeipa":           resourceRancher2AuthConfigFreeIpa(),
			"rancher2_auth_config_github":            resourceRancher2AuthConfigGithub(),
			"rancher2_auth_config_openldap":          resourceRancher2AuthConfigOpenLdap(),
			"rancher2_auth_config_ping":              resourceRancher2AuthConfigPing(),
			"rancher2_catalog":                       resourceRancher2Catalog(),
			"rancher2_cluster":                       resourceRancher2Cluster(),
			"rancher2_cluster_logging":               resourceRancher2ClusterLogging(),
			"rancher2_cluster_role_template_binding": resourceRancher2ClusterRoleTemplateBinding(),
			"rancher2_node_pool":                     resourceRancher2NodePool(),
			"rancher2_project":                       resourceRancher2Project(),
			"rancher2_project_logging":               resourceRancher2ProjectLogging(),
			"rancher2_project_role_template_binding": resourceRancher2ProjectRoleTemplateBinding(),
			"rancher2_namespace":                     resourceRancher2Namespace(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"rancher2_setting": dataSourceRancher2Setting(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "API Key used to authenticate with the rancher server",

		"secret_key": "API secret used to authenticate with the rancher server",

		"token_key": "API token used to authenticate with the rancher server",

		"ca_certs": "CA certificates used to sign rancher server tls certificates. Mandatory if self signed.",

		"api_url": "The URL to the rancher API",

		"config": "Path to the Rancher client cli.json config file",

		"insecure": "Allow insecure server connections when using SSL",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiURL := d.Get("api_url").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	tokenKey := d.Get("token_key").(string)
	caCerts := d.Get("ca_certs").(string)
	insecure := d.Get("insecure").(bool)

	if configFile := d.Get("config").(string); configFile != "" {
		config, err := loadConfig(configFile)
		if err != nil {
			return config, err
		}

		if apiURL == "" && config.URL != "" {
			u, err := url.Parse(config.URL)
			if err != nil {
				return config, err
			}
			apiURL = u.Scheme + "://" + u.Host
		}

		if accessKey == "" {
			accessKey = config.AccessKey
		}

		if secretKey == "" {
			secretKey = config.SecretKey
		}

		if tokenKey == "" {
			tokenKey = config.TokenKey
		}

		if caCerts == "" {
			caCerts = config.CACerts
		}

		if insecure == false {
			insecure = config.Insecure
		}
	}

	if apiURL == "" {
		return &Config{}, fmt.Errorf("No api_url provided")
	}

	config := &Config{
		URL:       apiURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
		TokenKey:  tokenKey,
		CACerts:   caCerts,
		Insecure:  insecure,
	}

	_, err := config.ManagementClient()
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}

func loadConfig(path string) (CLIConfig, error) {
	config := CLIConfig{
		Path: path,
	}

	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	config.Path = path

	return config, err
}
