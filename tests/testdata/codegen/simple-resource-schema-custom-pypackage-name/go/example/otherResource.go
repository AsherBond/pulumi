// Code generated by test DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package example

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"simple-resource-schema-custom-pypackage-name/example/internal"
)

type OtherResource struct {
	pulumi.ResourceState

	Foo ResourceOutput `pulumi:"foo"`
}

// NewOtherResource registers a new resource with the given unique name, arguments, and options.
func NewOtherResource(ctx *pulumi.Context,
	name string, args *OtherResourceArgs, opts ...pulumi.ResourceOption) (*OtherResource, error) {
	if args == nil {
		args = &OtherResourceArgs{}
	}

	opts = internal.PkgResourceDefaultOpts(opts)
	var resource OtherResource
	err := ctx.RegisterRemoteComponentResource("example::OtherResource", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type otherResourceArgs struct {
	Foo *Resource `pulumi:"foo"`
}

// The set of arguments for constructing a OtherResource resource.
type OtherResourceArgs struct {
	Foo ResourceInput
}

func (OtherResourceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*otherResourceArgs)(nil)).Elem()
}

type OtherResourceInput interface {
	pulumi.Input

	ToOtherResourceOutput() OtherResourceOutput
	ToOtherResourceOutputWithContext(ctx context.Context) OtherResourceOutput
}

func (*OtherResource) ElementType() reflect.Type {
	return reflect.TypeOf((**OtherResource)(nil)).Elem()
}

func (i *OtherResource) ToOtherResourceOutput() OtherResourceOutput {
	return i.ToOtherResourceOutputWithContext(context.Background())
}

func (i *OtherResource) ToOtherResourceOutputWithContext(ctx context.Context) OtherResourceOutput {
	return pulumi.ToOutputWithContext(ctx, i).(OtherResourceOutput)
}

type OtherResourceOutput struct{ *pulumi.OutputState }

func (OtherResourceOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**OtherResource)(nil)).Elem()
}

func (o OtherResourceOutput) ToOtherResourceOutput() OtherResourceOutput {
	return o
}

func (o OtherResourceOutput) ToOtherResourceOutputWithContext(ctx context.Context) OtherResourceOutput {
	return o
}

func (o OtherResourceOutput) Foo() ResourceOutput {
	return o.ApplyT(func(v *OtherResource) ResourceOutput { return v.Foo }).(ResourceOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*OtherResourceInput)(nil)).Elem(), &OtherResource{})
	pulumi.RegisterOutputType(OtherResourceOutput{})
}
