module l3-component-simple

go 1.20

require (
	github.com/pulumi/pulumi/sdk/v3 v3.30.0
	example.com/pulumi-simple/sdk/go/v2 v2.0.0
)

replace github.com/pulumi/pulumi/sdk/v3 => ../../artifacts/github.com_pulumi_pulumi_sdk_v3

replace example.com/pulumi-simple/sdk/go/v2 => ../../artifacts/example.com_pulumi-simple_sdk_go_v2
