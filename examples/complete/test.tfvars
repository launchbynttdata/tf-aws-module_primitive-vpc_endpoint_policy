logical_product_family  = "launch"
logical_product_service = "vpcepol"
class_env               = "dev"
instance_env            = 19
instance_resource       = 1
resource_names_map      = { vpc = { name = "vpc", max_length = 64 }, route_table = { name = "routetable", max_length = 64 }, endpoint = { name = "endpoint", max_length = 64 } }
tags                    = { environment = "test", module = "vpc_endpoint_policy" }
