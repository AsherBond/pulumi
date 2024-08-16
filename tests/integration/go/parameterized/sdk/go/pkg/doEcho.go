// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package pkg

import (
	"context"
	"reflect"

	"example.com/pulumi-pkg/sdk/go/pkg/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A test invoke that echoes its input.
func DoEcho(ctx *pulumi.Context, args *DoEchoArgs, opts ...pulumi.InvokeOption) (*DoEchoResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	ref, err := internal.PkgGetPackageRef(ctx)
	if err != nil {
		return nil, err
	}
	var rv DoEchoResult
	err = ctx.InvokePackage("pkg:index:doEcho", args, &rv, ref, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

type DoEchoArgs struct {
	Echo *string `pulumi:"echo"`
}

type DoEchoResult struct {
	Echo *string `pulumi:"echo"`
}

func DoEchoOutput(ctx *pulumi.Context, args DoEchoOutputArgs, opts ...pulumi.InvokeOption) DoEchoResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (DoEchoResult, error) {
			args := v.(DoEchoArgs)
			r, err := DoEcho(ctx, &args, opts...)
			var s DoEchoResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(DoEchoResultOutput)
}

type DoEchoOutputArgs struct {
	Echo pulumi.StringPtrInput `pulumi:"echo"`
}

func (DoEchoOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*DoEchoArgs)(nil)).Elem()
}

type DoEchoResultOutput struct{ *pulumi.OutputState }

func (DoEchoResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*DoEchoResult)(nil)).Elem()
}

func (o DoEchoResultOutput) ToDoEchoResultOutput() DoEchoResultOutput {
	return o
}

func (o DoEchoResultOutput) ToDoEchoResultOutputWithContext(ctx context.Context) DoEchoResultOutput {
	return o
}

func (o DoEchoResultOutput) Echo() pulumi.StringPtrOutput {
	return o.ApplyT(func(v DoEchoResult) *string { return v.Echo }).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterOutputType(DoEchoResultOutput{})
}