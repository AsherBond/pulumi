// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package component

import (
	"context"
	"reflect"

	"errors"
	"example.com/pulumi-component/sdk/go/v13/component/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A component resource that accepts a reference to a custom resource. The input resource's `value` is used to create a child custom resource inside the component, before a reference to this child is returned.
type ComponentCustomRefInputOutput struct {
	pulumi.ResourceState

	InputRef  CustomOutput `pulumi:"inputRef"`
	OutputRef CustomOutput `pulumi:"outputRef"`
}

// NewComponentCustomRefInputOutput registers a new resource with the given unique name, arguments, and options.
func NewComponentCustomRefInputOutput(ctx *pulumi.Context,
	name string, args *ComponentCustomRefInputOutputArgs, opts ...pulumi.ResourceOption) (*ComponentCustomRefInputOutput, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.InputRef == nil {
		return nil, errors.New("invalid value for required argument 'InputRef'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource ComponentCustomRefInputOutput
	err := ctx.RegisterRemoteComponentResource("component:index:ComponentCustomRefInputOutput", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

type componentCustomRefInputOutputArgs struct {
	InputRef *Custom `pulumi:"inputRef"`
}

// The set of arguments for constructing a ComponentCustomRefInputOutput resource.
type ComponentCustomRefInputOutputArgs struct {
	InputRef CustomInput
}

func (ComponentCustomRefInputOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*componentCustomRefInputOutputArgs)(nil)).Elem()
}

type ComponentCustomRefInputOutputInput interface {
	pulumi.Input

	ToComponentCustomRefInputOutputOutput() ComponentCustomRefInputOutputOutput
	ToComponentCustomRefInputOutputOutputWithContext(ctx context.Context) ComponentCustomRefInputOutputOutput
}

func (*ComponentCustomRefInputOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ComponentCustomRefInputOutput)(nil)).Elem()
}

func (i *ComponentCustomRefInputOutput) ToComponentCustomRefInputOutputOutput() ComponentCustomRefInputOutputOutput {
	return i.ToComponentCustomRefInputOutputOutputWithContext(context.Background())
}

func (i *ComponentCustomRefInputOutput) ToComponentCustomRefInputOutputOutputWithContext(ctx context.Context) ComponentCustomRefInputOutputOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ComponentCustomRefInputOutputOutput)
}

type ComponentCustomRefInputOutputOutput struct{ *pulumi.OutputState }

func (ComponentCustomRefInputOutputOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ComponentCustomRefInputOutput)(nil)).Elem()
}

func (o ComponentCustomRefInputOutputOutput) ToComponentCustomRefInputOutputOutput() ComponentCustomRefInputOutputOutput {
	return o
}

func (o ComponentCustomRefInputOutputOutput) ToComponentCustomRefInputOutputOutputWithContext(ctx context.Context) ComponentCustomRefInputOutputOutput {
	return o
}

func (o ComponentCustomRefInputOutputOutput) InputRef() CustomOutput {
	return o.ApplyT(func(v *ComponentCustomRefInputOutput) CustomOutput { return v.InputRef }).(CustomOutput)
}

func (o ComponentCustomRefInputOutputOutput) OutputRef() CustomOutput {
	return o.ApplyT(func(v *ComponentCustomRefInputOutput) CustomOutput { return v.OutputRef }).(CustomOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ComponentCustomRefInputOutputInput)(nil)).Elem(), &ComponentCustomRefInputOutput{})
	pulumi.RegisterOutputType(ComponentCustomRefInputOutputOutput{})
}
