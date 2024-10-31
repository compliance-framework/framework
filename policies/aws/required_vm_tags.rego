package required_vm_tags

policyID := "AWSSEC-0005"

minimum_required_tags = [
	"owner_application",
	"owner_domain",
	"owner_team",
	"owner_iam_role",
]

tags_contain_proper_keys(tags) {
	keys := {key | tags[key]}
	minimum_tags_set := {x | x := minimum_required_tags[i]}
	leftover := minimum_tags_set - keys

	# If all minimum_tags exist in keys, the leftover set should be empty - equal to a new set()
	leftover == set()
}

warn[msg] {
	resource := input.resource[resource_type]
	tags := resource[name].tags

	# Create an array of resources, only if they are missing the minimum tags
	resources := [sprintf("%v.%v", [resource_type, name]) | not tags_contain_proper_keys(tags)]

	resources != []
	msg := sprintf("%s: Invalid tags (missing minimum required tags) for the following resource(s): `%v`. Required tags: `%v`", [policyID, resources, minimum_required_tags])
}
