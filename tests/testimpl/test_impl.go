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

// TestComposableComplete verifies the deployed VPC endpoint policy and exercises a reversible policy update.
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	client, endpointID, policy := verifyEndpointPolicy(t, ctx)
	exerciseEndpointPolicyWrite(t, client, endpointID, policy)
}

// TestComposableCompleteReadOnly verifies the deployed VPC endpoint policy using read-only AWS API calls.
func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	verifyEndpointPolicy(t, ctx)
}

func verifyEndpointPolicy(t *testing.T, ctx types.TestContext) (*ec2.Client, string, string) {
	opts := ctx.TerratestTerraformOptions()
	region := terraform.OutputContext(t, context.Background(), opts, "region")
	endpointID := terraform.OutputContext(t, context.Background(), opts, "vpc_endpoint_id")
	policy := terraform.OutputContext(t, context.Background(), opts, "policy")

	require.NotEqual(t, "", endpointID)
	assert.Equal(t, terraform.OutputContext(t, context.Background(), opts, "expected_vpc_endpoint_id"), endpointID)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	require.NoError(t, err)

	client := ec2.NewFromConfig(cfg)
	output, err := client.DescribeVpcEndpoints(context.Background(), &ec2.DescribeVpcEndpointsInput{VpcEndpointIds: []string{endpointID}})
	require.NoError(t, err)
	require.Len(t, output.VpcEndpoints, 1)

	assert.JSONEq(t, canonicalJSON(t, policy), canonicalJSON(t, aws.ToString(output.VpcEndpoints[0].PolicyDocument)))

	return client, endpointID, policy
}

func exerciseEndpointPolicyWrite(t *testing.T, client *ec2.Client, endpointID string, originalPolicy string) {
	t.Helper()

	mutatedPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:ListAllMyBuckets","Resource":"*"},{"Effect":"Deny","Principal":"*","Action":"s3:DeleteBucket","Resource":"*"}]}`
	defer func() {
		_, restoreErr := client.ModifyVpcEndpoint(context.Background(), &ec2.ModifyVpcEndpointInput{
			VpcEndpointId:  aws.String(endpointID),
			PolicyDocument: aws.String(originalPolicy),
		})
		assert.NoError(t, restoreErr)
	}()

	_, err := client.ModifyVpcEndpoint(context.Background(), &ec2.ModifyVpcEndpointInput{
		VpcEndpointId:  aws.String(endpointID),
		PolicyDocument: aws.String(mutatedPolicy),
	})
	require.NoError(t, err)

	_, err = client.ModifyVpcEndpoint(context.Background(), &ec2.ModifyVpcEndpointInput{
		VpcEndpointId:  aws.String(endpointID),
		PolicyDocument: aws.String(originalPolicy),
	})
	require.NoError(t, err)
}

func canonicalJSON(t *testing.T, value string) string {
	t.Helper()

	var decoded any
	require.NoError(t, json.Unmarshal([]byte(value), &decoded))

	encoded, err := json.Marshal(decoded)
	require.NoError(t, err)

	return string(encoded)
}
