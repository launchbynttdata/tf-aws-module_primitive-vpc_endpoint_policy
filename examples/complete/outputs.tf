// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

output "id" {
  description = "The VPC endpoint policy ID."
  value       = module.endpoint_policy.id
}
output "policy" {
  description = "The configured VPC endpoint policy JSON."
  value       = module.endpoint_policy.policy
}
output "vpc_endpoint_id" {
  description = "The VPC endpoint ID."
  value       = module.endpoint_policy.vpc_endpoint_id
}
output "expected_vpc_endpoint_id" {
  description = "Expected VPC endpoint ID."
  value       = aws_vpc_endpoint.s3.id
}

output "region" {
  description = "The AWS Region where the example resources are deployed."
  value       = data.aws_region.current.region
}
