package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	testCloudCredentialConfAmazonec2         *CloudCredential
	testCloudCredentialInterfaceAmazonec2    map[string]interface{}
	testCloudCredentialConfAzure             *CloudCredential
	testCloudCredentialInterfaceAzure        map[string]interface{}
	testCloudCredentialConfDigitalocean      *CloudCredential
	testCloudCredentialInterfaceDigitalocean map[string]interface{}
	testCloudCredentialConfOpenstack         *CloudCredential
	testCloudCredentialInterfaceOpenstack    map[string]interface{}
	testCloudCredentialConfVsphere           *CloudCredential
	testCloudCredentialInterfaceVsphere      map[string]interface{}
)

func init() {
	testCloudCredentialConfAmazonec2 = &CloudCredential{
		Amazonec2CredentialConfig: testCloudCredentialAmazonec2Conf,
	}
	testCloudCredentialConfAmazonec2.Name = "cloudCredential-test"
	testCloudCredentialConfAmazonec2.Description = "description"
	testCloudCredentialInterfaceAmazonec2 = map[string]interface{}{
		"name":                        "cloudCredential-test",
		"description":                 "description",
		"amazonec2_credential_config": testCloudCredentialAmazonec2Interface,
		"driver":                      amazonec2ConfigDriver,
	}
	testCloudCredentialConfAzure = &CloudCredential{
		AzureCredentialConfig: testCloudCredentialAzureConf,
	}
	testCloudCredentialConfAzure.Name = "cloudCredential-test"
	testCloudCredentialConfAzure.Description = "description"
	testCloudCredentialInterfaceAzure = map[string]interface{}{
		"name":                    "cloudCredential-test",
		"description":             "description",
		"azure_credential_config": testCloudCredentialAzureInterface,
		"driver":                  azureConfigDriver,
	}
	testCloudCredentialConfDigitalocean = &CloudCredential{
		DigitaloceanCredentialConfig: testCloudCredentialDigitaloceanConf,
	}
	testCloudCredentialConfDigitalocean.Name = "cloudCredential-test"
	testCloudCredentialConfDigitalocean.Description = "description"
	testCloudCredentialInterfaceDigitalocean = map[string]interface{}{
		"name":                           "cloudCredential-test",
		"description":                    "description",
		"digitalocean_credential_config": testCloudCredentialDigitaloceanInterface,
		"driver":                         digitaloceanConfigDriver,
	}
	testCloudCredentialConfOpenstack = &CloudCredential{
		OpenstackCredentialConfig: testCloudCredentialOpenstackConf,
	}
	testCloudCredentialConfOpenstack.Name = "cloudCredential-test"
	testCloudCredentialConfOpenstack.Description = "description"
	testCloudCredentialInterfaceOpenstack = map[string]interface{}{
		"name":                        "cloudCredential-test",
		"description":                 "description",
		"openstack_credential_config": testCloudCredentialOpenstackInterface,
		"driver":                      openstackConfigDriver,
	}
	testCloudCredentialConfVsphere = &CloudCredential{
		VmwarevsphereCredentialConfig: testCloudCredentialVsphereConf,
	}
	testCloudCredentialConfVsphere.Name = "cloudCredential-test"
	testCloudCredentialConfVsphere.Description = "description"
	testCloudCredentialInterfaceVsphere = map[string]interface{}{
		"name":                      "cloudCredential-test",
		"description":               "description",
		"vsphere_credential_config": testCloudCredentialVsphereInterface,
		"driver":                    vmwarevsphereConfigDriver,
	}
}

func TestFlattenCloudCredential(t *testing.T) {

	cases := []struct {
		Input          *CloudCredential
		ExpectedOutput map[string]interface{}
	}{
		{
			testCloudCredentialConfAmazonec2,
			testCloudCredentialInterfaceAmazonec2,
		},
		{
			testCloudCredentialConfAzure,
			testCloudCredentialInterfaceAzure,
		},
		{
			testCloudCredentialConfDigitalocean,
			testCloudCredentialInterfaceDigitalocean,
		},
		{
			testCloudCredentialConfOpenstack,
			testCloudCredentialInterfaceOpenstack,
		},
		{
			testCloudCredentialConfVsphere,
			testCloudCredentialInterfaceVsphere,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, cloudCredentialFields(), tc.ExpectedOutput)
		err := flattenCloudCredential(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandCloudCredential(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *CloudCredential
	}{
		{
			testCloudCredentialInterfaceAmazonec2,
			testCloudCredentialConfAmazonec2,
		},
		{
			testCloudCredentialInterfaceAzure,
			testCloudCredentialConfAzure,
		},
		{
			testCloudCredentialInterfaceDigitalocean,
			testCloudCredentialConfDigitalocean,
		},
		{
			testCloudCredentialInterfaceOpenstack,
			testCloudCredentialConfOpenstack,
		},
		{
			testCloudCredentialInterfaceVsphere,
			testCloudCredentialConfVsphere,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, cloudCredentialFields(), tc.Input)
		output := expandCloudCredential(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
