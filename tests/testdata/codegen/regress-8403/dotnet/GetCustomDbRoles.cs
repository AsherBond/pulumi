// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Mongodbatlas
{
    public static class GetCustomDbRoles
    {
        public static Task<GetCustomDbRolesResult> InvokeAsync(GetCustomDbRolesArgs? args = null, InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.InvokeAsync<GetCustomDbRolesResult>("mongodbatlas::getCustomDbRoles", args ?? new GetCustomDbRolesArgs(), options.WithDefaults());

        public static Output<GetCustomDbRolesResult> Invoke(InvokeOptions? options = null)
            => global::Pulumi.Deployment.Instance.Invoke<GetCustomDbRolesResult>("mongodbatlas::getCustomDbRoles", InvokeArgs.Empty, options.WithDefaults());

        public static Output<GetCustomDbRolesResult> Invoke(InvokeOutputOptions options)
            => global::Pulumi.Deployment.Instance.Invoke<GetCustomDbRolesResult>("mongodbatlas::getCustomDbRoles", InvokeArgs.Empty, options.WithDefaults());
    }


    public sealed class GetCustomDbRolesArgs : global::Pulumi.InvokeArgs
    {
        public GetCustomDbRolesArgs()
        {
        }
        public static new GetCustomDbRolesArgs Empty => new GetCustomDbRolesArgs();
    }


    [OutputType]
    public sealed class GetCustomDbRolesResult
    {
        public readonly Outputs.GetCustomDbRolesResult? Result;

        [OutputConstructor]
        private GetCustomDbRolesResult(Outputs.GetCustomDbRolesResult? result)
        {
            Result = result;
        }
    }
}
