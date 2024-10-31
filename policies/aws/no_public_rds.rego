package no_public_rds

import data.util_functions

policyID := "AWSSEC-0003"

has_public_attribute(resource) {
	util_functions.has_key(resource, "publicly_accessible")
}

violation[{"policyId": policyID, "msg": msg}] {
	resource := input.resource.aws_db_instance
	a_resource := resource[name]
	has_public_attribute(a_resource)
	a_resource.publicly_accessible != false

	msg := sprintf("RDS instances must not be publicly exposed. Set `publicly_accessible` to `false` on aws_db_instance.`%s`", [name])
}
