package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

var (
	testNamespaceResourceQuotaLimitConf      *clusterClient.ResourceQuotaLimit
	testNamespaceResourceQuotaLimitInterface []interface{}
	testNamespaceResourceQuotaConf           *clusterClient.NamespaceResourceQuota
	testNamespaceResourceQuotaInterface      []interface{}
	testNamespaceConf                        *clusterClient.Namespace
	testNamespaceInterface                   map[string]interface{}
)

func init() {
	testNamespaceResourceQuotaLimitConf = &clusterClient.ResourceQuotaLimit{
		ConfigMaps:             "config",
		LimitsCPU:              "cpu",
		LimitsMemory:           "memory",
		PersistentVolumeClaims: "pvc",
		Pods:                   "pods",
		ReplicationControllers: "rc",
		RequestsCPU:            "r_cpu",
		RequestsMemory:         "r_memory",
		RequestsStorage:        "r_storage",
		Secrets:                "secrets",
		Services:               "services",
		ServicesLoadBalancers:  "lb",
		ServicesNodePorts:      "np",
	}
	testNamespaceResourceQuotaLimitInterface = []interface{}{
		map[string]interface{}{
			"config_maps":              "config",
			"limits_cpu":               "cpu",
			"limits_memory":            "memory",
			"persistent_volume_claims": "pvc",
			"pods":                     "pods",
			"replication_controllers":  "rc",
			"requests_cpu":             "r_cpu",
			"requests_memory":          "r_memory",
			"requests_storage":         "r_storage",
			"secrets":                  "secrets",
			"services":                 "services",
			"services_load_balancers":  "lb",
			"services_node_ports":      "np",
		},
	}
	testNamespaceResourceQuotaConf = &clusterClient.NamespaceResourceQuota{
		Limit: testNamespaceResourceQuotaLimitConf,
	}
	testNamespaceResourceQuotaInterface = []interface{}{
		map[string]interface{}{
			"limit": testNamespaceResourceQuotaLimitInterface,
		},
	}
	testNamespaceConf = &clusterClient.Namespace{
		ProjectID:     "project:test",
		Name:          "test",
		Description:   "description",
		ResourceQuota: testNamespaceResourceQuotaConf,
	}
	testNamespaceInterface = map[string]interface{}{
		"project_id":     "project:test",
		"name":           "test",
		"description":    "description",
		"resource_quota": testNamespaceResourceQuotaInterface,
	}
}

func TestFlattenNamespaceResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.ResourceQuotaLimit
		ExpectedOutput []interface{}
	}{
		{
			testNamespaceResourceQuotaLimitConf,
			testNamespaceResourceQuotaLimitInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNamespaceResourceQuotaLimit(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenNamespaceResourceQuota(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.NamespaceResourceQuota
		ExpectedOutput []interface{}
	}{
		{
			testNamespaceResourceQuotaConf,
			testNamespaceResourceQuotaInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNamespaceResourceQuota(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenNamespace(t *testing.T) {

	cases := []struct {
		Input          *clusterClient.Namespace
		ExpectedOutput map[string]interface{}
	}{
		{
			testNamespaceConf,
			testNamespaceInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, namespaceFields(), map[string]interface{}{})
		err := flattenNamespace(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandNamespaceResourceQuotaLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *clusterClient.ResourceQuotaLimit
	}{
		{
			testNamespaceResourceQuotaLimitInterface,
			testNamespaceResourceQuotaLimitConf,
		},
	}

	for _, tc := range cases {
		output := expandNamespaceResourceQuotaLimit(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNamespaceResourceQuota(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *clusterClient.NamespaceResourceQuota
	}{
		{
			testNamespaceResourceQuotaInterface,
			testNamespaceResourceQuotaConf,
		},
	}

	for _, tc := range cases {
		output := expandNamespaceResourceQuota(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNamespace(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *clusterClient.Namespace
	}{
		{
			testNamespaceInterface,
			testNamespaceConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, namespaceFields(), tc.Input)
		output := expandNamespace(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
