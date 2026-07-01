# Complete Example

This example creates a complete VPC endpoint policy deployment with the dependencies required to exercise the primitive module.

## Usage

```hcl
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

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

module "resource_names" {
  source  = "terraform.registry.launch.nttdata.com/module_library/resource_name/launch"
  version = "~> 2.0"

  for_each = var.resource_names_map

  logical_product_family  = var.logical_product_family
  logical_product_service = var.logical_product_service
  class_env               = var.class_env
  instance_env            = var.instance_env
  instance_resource       = var.instance_resource
  cloud_resource_type     = each.value.name
  maximum_length          = each.value.max_length

  region = join("", split("-", data.aws_region.current.region))
}

resource "aws_vpc" "example" {
  cidr_block           = "10.42.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags                 = merge(var.tags, { Name = module.resource_names["vpc"].standard })
}

resource "aws_route_table" "example" {
  vpc_id = aws_vpc.example.id
  tags   = merge(var.tags, { Name = module.resource_names["route_table"].standard })
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = aws_vpc.example.id
  service_name      = "com.amazonaws.${data.aws_region.current.region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [aws_route_table.example.id]
  tags              = merge(var.tags, { Name = module.resource_names["endpoint"].standard })
}

data "aws_iam_policy_document" "endpoint" {
  statement {
    effect    = "Allow"
    actions   = ["s3:ListAllMyBuckets"]
    resources = ["*"]
    principals { type = "*" identifiers = ["*"] }
  }
}

module "endpoint_policy" {
  source = "../.."

  vpc_endpoint_id = aws_vpc_endpoint.s3.id
  policy          = data.aws_iam_policy_document.endpoint.json
}
```

<!-- BEGIN_TF_DOCS -->
<!-- END_TF_DOCS -->
