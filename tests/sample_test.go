package tests

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

//TODO: Rename Test
func TestModuleDeployment(t *testing.T) {
	t.Parallel()
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
}
