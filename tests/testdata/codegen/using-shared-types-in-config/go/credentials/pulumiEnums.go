// Code generated by test DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package credentials

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type HashKind string

const (
	// Adler32 implements the Adler-32 checksum.
	HashKindAdler32 = HashKind("Adler32")
	// CRC32 implements the 32-bit cyclic redundancy check, or CRC-32, checksum.
	HashKindCRC32 = HashKind("CRC32")
)

func (HashKind) ElementType() reflect.Type {
	return reflect.TypeOf((*HashKind)(nil)).Elem()
}

func (e HashKind) ToHashKindOutput() HashKindOutput {
	return pulumi.ToOutput(e).(HashKindOutput)
}

func (e HashKind) ToHashKindOutputWithContext(ctx context.Context) HashKindOutput {
	return pulumi.ToOutputWithContext(ctx, e).(HashKindOutput)
}

func (e HashKind) ToHashKindPtrOutput() HashKindPtrOutput {
	return e.ToHashKindPtrOutputWithContext(context.Background())
}

func (e HashKind) ToHashKindPtrOutputWithContext(ctx context.Context) HashKindPtrOutput {
	return HashKind(e).ToHashKindOutputWithContext(ctx).ToHashKindPtrOutputWithContext(ctx)
}

func (e HashKind) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e HashKind) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e HashKind) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e HashKind) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type HashKindOutput struct{ *pulumi.OutputState }

func (HashKindOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*HashKind)(nil)).Elem()
}

func (o HashKindOutput) ToHashKindOutput() HashKindOutput {
	return o
}

func (o HashKindOutput) ToHashKindOutputWithContext(ctx context.Context) HashKindOutput {
	return o
}

func (o HashKindOutput) ToHashKindPtrOutput() HashKindPtrOutput {
	return o.ToHashKindPtrOutputWithContext(context.Background())
}

func (o HashKindOutput) ToHashKindPtrOutputWithContext(ctx context.Context) HashKindPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v HashKind) *HashKind {
		return &v
	}).(HashKindPtrOutput)
}

func (o HashKindOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o HashKindOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e HashKind) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o HashKindOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o HashKindOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e HashKind) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type HashKindPtrOutput struct{ *pulumi.OutputState }

func (HashKindPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**HashKind)(nil)).Elem()
}

func (o HashKindPtrOutput) ToHashKindPtrOutput() HashKindPtrOutput {
	return o
}

func (o HashKindPtrOutput) ToHashKindPtrOutputWithContext(ctx context.Context) HashKindPtrOutput {
	return o
}

func (o HashKindPtrOutput) Elem() HashKindOutput {
	return o.ApplyT(func(v *HashKind) HashKind {
		if v != nil {
			return *v
		}
		var ret HashKind
		return ret
	}).(HashKindOutput)
}

func (o HashKindPtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o HashKindPtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *HashKind) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// HashKindInput is an input type that accepts values of the HashKind enum
// A concrete instance of `HashKindInput` can be one of the following:
//
//	HashKindAdler32
//	HashKindCRC32
type HashKindInput interface {
	pulumi.Input

	ToHashKindOutput() HashKindOutput
	ToHashKindOutputWithContext(context.Context) HashKindOutput
}

var hashKindPtrType = reflect.TypeOf((**HashKind)(nil)).Elem()

type HashKindPtrInput interface {
	pulumi.Input

	ToHashKindPtrOutput() HashKindPtrOutput
	ToHashKindPtrOutputWithContext(context.Context) HashKindPtrOutput
}

type hashKindPtr string

func HashKindPtr(v string) HashKindPtrInput {
	return (*hashKindPtr)(&v)
}

func (*hashKindPtr) ElementType() reflect.Type {
	return hashKindPtrType
}

func (in *hashKindPtr) ToHashKindPtrOutput() HashKindPtrOutput {
	return pulumi.ToOutput(in).(HashKindPtrOutput)
}

func (in *hashKindPtr) ToHashKindPtrOutputWithContext(ctx context.Context) HashKindPtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(HashKindPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*HashKindInput)(nil)).Elem(), HashKind("Adler32"))
	pulumi.RegisterInputType(reflect.TypeOf((*HashKindPtrInput)(nil)).Elem(), HashKind("Adler32"))
	pulumi.RegisterOutputType(HashKindOutput{})
	pulumi.RegisterOutputType(HashKindPtrOutput{})
}
