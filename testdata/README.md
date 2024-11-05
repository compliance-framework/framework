# This directory holds testdata for the policy evaluator

## Running the agent locally for a test

```shell
go run main.go agent \
    --policy-bundle=testdata/bundle_local_ssh.tar.gz \
    --policy-bundle=testdata/bundle_local_ssh_clash.tar.gz
```

## Local SSH Bundles

Policy Bundles
```yaml
./bundle_local_ssh.tar.gz
./bundle_local_ssh_clash.tar.gz
```

There are two policies where the top level package `compliance_framework` in them match, and opa generally 
doesn't like that. We need to support cases like this as folks might inject many policy bundles and all
of them will have the same root package.

When executing the agent with multiple policy bundles, it should successfully check all policies contained
regardless of whether they share the same root package.
