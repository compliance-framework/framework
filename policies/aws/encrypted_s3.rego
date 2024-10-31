package encrypted_s3

import data.util_functions

policyID := "AWSSEC-0001"

violation[{"policyId": policyID, "msg": msg}] {
	resource := input.resource.aws_s3_bucket
	a_resource := resource[name]
	not util_functions.has_key(a_resource, "server_side_encryption_configuration")

	msg := sprintf("Missing S3 encryption for `%s`. Required flag: `server_side_encryption_configuration`", [name])
}
