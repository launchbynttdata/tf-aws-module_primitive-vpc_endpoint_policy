package testimpl

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComposableComplete verifies the deployed VPC endpoint policy.
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	verifyEndpointPolicy(t, ctx)
}

// TestComposableCompleteReadOnly verifies the deployed VPC endpoint policy using read-only AWS API calls.
func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	verifyEndpointPolicy(t, ctx)
}

func verifyEndpointPolicy(t *testing.T, ctx types.TestContext) {
	opts := ctx.TerratestTerraformOptions()
	region := terraform.Output(t, opts, "region")
	endpointID := terraform.Output(t, opts, "vpc_endpoint_id")
	policy := terraform.Output(t, opts, "policy")

	require.NotEqual(t, "", endpointID)
	assert.Equal(t, terraform.Output(t, opts, "expected_vpc_endpoint_id"), endpointID)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err)

	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeVpcEndpoints(context.Background(), &ec2.DescribeVpcEndpointsInput{VpcEndpointIds: []string{endpointID}})
	require.NoError(t, err)
	require.Len(t, output.VpcEndpoints, 1)

	assert.JSONEq(t, canonicalJSON(t, policy), canonicalJSON(t, aws.ToString(output.VpcEndpoints[0].PolicyDocument)))
}

func canonicalJSON(t *testing.T, value string) string {
	t.Helper()

	var decoded any
	require.NoError(t, json.Unmarshal([]byte(value), &decoded))

	encoded, err := json.Marshal(decoded)
	require.NoError(t, err)

	return string(encoded)
}
