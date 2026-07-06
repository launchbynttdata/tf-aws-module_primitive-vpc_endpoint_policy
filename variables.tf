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

# -----------------------------------------------------------------------------
# Required
# -----------------------------------------------------------------------------

variable "vpc_endpoint_id" {
  description = "ID of the VPC endpoint to attach the policy to."
  type        = string

  validation {
    condition     = can(regex("^vpce-[a-z0-9]+$", var.vpc_endpoint_id))
    error_message = "vpc_endpoint_id must start with vpce-."
  }
}

# -----------------------------------------------------------------------------
# Optional
# -----------------------------------------------------------------------------

variable "policy" {
  description = "JSON policy document for the VPC endpoint. Defaults to the provider-managed policy when omitted."
  type        = string
  default     = null

  validation {
    condition     = var.policy == null ? true : can(jsondecode(var.policy))
    error_message = "policy must be valid JSON."
  }
}
variable "timeouts" {
  description = "Timeouts for creating or deleting the endpoint policy."
  type = object({
    create = optional(string)
    delete = optional(string)
  })
  default = null
}
variable "region" {
  description = "AWS Region where this resource is managed. Defaults to the provider-configured Region."
  type        = string
  default     = null
}
