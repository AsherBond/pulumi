// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"context"
	"fmt"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	surveycore "github.com/AlecAivazis/survey/v2/core"

	"github.com/pulumi/pulumi/pkg/v3/backend/display"
	sdkDisplay "github.com/pulumi/pulumi/pkg/v3/display"
	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag/colors"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/result"
)

// ApplierOptions is a bag of configuration settings for an Applier.
type ApplierOptions struct {
	// DryRun indicates if the update should not change any resource state and instead just preview changes.
	DryRun bool
	// ShowLink indicates if a link to the update persisted result can be displayed.
	ShowLink bool
}

// Applier applies the changes specified by this update operation against the target stack.
type Applier func(ctx context.Context, kind apitype.UpdateKind, stack Stack, op UpdateOperation,
	opts ApplierOptions, events chan<- engine.Event) (*deploy.Plan, sdkDisplay.ResourceChanges, error)

// Explainer provides a contract for explaining changes that will be made to the stack.
// For Pulumi Cloud, this is a Copilot explainer.
type Explainer interface {
	// Explain returns a human-readable explanation of the changes that will be made to the stack.
	Explain(ctx context.Context, stackRef StackReference, kind apitype.UpdateKind, op UpdateOperation,
		events []engine.Event) (string, error)
	// IsExplainPreviewEnabled returns whether the explainer is enabled for the given project.
	IsExplainPreviewEnabled(ctx context.Context, opts display.Options) bool
}

func ActionLabel(kind apitype.UpdateKind, dryRun bool) string {
	v := updateTextMap[kind]
	contract.Assertf(v.previewText != "", "preview text for %q cannot be empty", kind)
	contract.Assertf(v.text != "", "text for %q cannot be empty", kind)

	if dryRun {
		return "Previewing " + v.previewText
	}

	return v.text
}

var updateTextMap = map[apitype.UpdateKind]struct {
	previewText string
	text        string
}{
	apitype.PreviewUpdate:        {"update", "Previewing"},
	apitype.UpdateUpdate:         {"update", "Updating"},
	apitype.RefreshUpdate:        {"refresh", "Refreshing"},
	apitype.DestroyUpdate:        {"destroy", "Destroying"},
	apitype.StackImportUpdate:    {"stack import", "Importing"},
	apitype.ResourceImportUpdate: {"import", "Importing"},
}

type response string

const (
	yes     response = "yes"
	no      response = "no"
	details response = "details"
)

func PreviewThenPrompt(ctx context.Context, kind apitype.UpdateKind, stack Stack,
	op UpdateOperation, apply Applier, explainer Explainer,
) (*deploy.Plan, sdkDisplay.ResourceChanges, error) {
	// create a channel to hear about the update events from the engine. this will be used so that
	// we can build up the diff display in case the user asks to see the details of the diff

	// Note that eventsChannel is not closed in a `defer`. It is generally unsafe to do so, since defers run during
	// panics and we can't know whether or not we were in the middle of writing to this channel when the panic occurred.
	//
	// Instead of using a `defer`, we manually close `eventsChannel` on every exit of this function.
	eventsChannel := make(chan engine.Event)

	var events []engine.Event
	go func() {
		// Pull out relevant events we will want to display in the confirmation below.
		for e := range eventsChannel {
			// Don't include internal events in the confirmation stats.
			if e.Internal() {
				continue
			}
			if e.Type == engine.ResourcePreEvent ||
				e.Type == engine.ResourceOutputsEvent ||
				e.Type == engine.PolicyRemediationEvent ||
				e.Type == engine.SummaryEvent {
				events = append(events, e)
			}
		}
	}()

	// Perform the update operations, passing true for dryRun, so that we get a preview.
	// We perform the preview (DryRun), but don't display the cloud link since the
	// thing the user cares about would be the link to the actual update if they
	// confirm the prompt.
	opts := ApplierOptions{
		DryRun:   true,
		ShowLink: true,
	}

	plan, changes, err := apply(ctx, kind, stack, op, opts, eventsChannel)
	if err != nil {
		close(eventsChannel)
		return plan, changes, err
	}

	// If there are no changes, or we're auto-approving or just previewing, we can skip the confirmation prompt.
	if op.Opts.AutoApprove || kind == apitype.PreviewUpdate {
		close(eventsChannel)
		// If we're running in experimental mode then return the plan generated, else discard it. The user may
		// be explicitly setting a plan but that's handled higher up the call stack.
		if !op.Opts.Engine.Experimental {
			plan = nil
		}
		return plan, changes, nil
	}

	stats := computeUpdateStats(events)

	infoPrefix := "\b" + op.Opts.Display.Color.Colorize(colors.SpecWarning+"info: "+colors.Reset)
	if kind != apitype.UpdateUpdate {
		// If not an update, we can skip displaying warnings
	} else if stats.numNonStackResources == 0 {
		// This is an update and there are no resources being CREATED
		fmt.Print(infoPrefix, "There are no resources in your stack (other than the stack resource).\n\n")
	}

	// Warn user if an update is going to leave untracked resources in the environment.
	if (kind == apitype.UpdateUpdate || kind == apitype.PreviewUpdate || kind == apitype.DestroyUpdate) &&
		len(stats.retainedResources) != 0 {
		fmt.Printf(
			"%sThis update will leave %d resource(s) untracked in your environment:\n",
			infoPrefix, len(stats.retainedResources))
		for _, res := range stats.retainedResources {
			urn := res.URN
			fmt.Printf("    - %s %s\n", urn.Type().DisplayName(), urn.Name())
		}
		fmt.Print("\n")
	}

	// Otherwise, ensure the user wants to proceed.
	plan, err = confirmBeforeUpdating(ctx, kind, stack.Ref(), op, events, plan, explainer)
	close(eventsChannel)
	return plan, changes, err
}

// confirmBeforeUpdating asks the user whether to proceed. A nil error means yes.
func confirmBeforeUpdating(ctx context.Context, kind apitype.UpdateKind, stackRef StackReference, op UpdateOperation,
	events []engine.Event, plan *deploy.Plan, explainer Explainer,
	// askOpts is used for testing to control stdin/stdout
	askOpts ...survey.AskOpt,
) (*deploy.Plan, error) {
	for {
		opts := op.Opts
		var response string

		surveycore.DisableColor = true
		surveyIcons := survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question = survey.Icon{}
			icons.SelectFocus = survey.Icon{Text: opts.Display.Color.Colorize(colors.BrightGreen + ">" + colors.Reset)}
		})

		choices := []string{string(yes), string(no)}

		// explain is down here instead of with the other choices so we can optionally emit an emoji
		explainChoice := "explain " + cmdutil.EmojiOr("✨", "")

		// For non-previews, we can also offer a detailed summary.
		if !opts.SkipPreview {
			choices = append(choices, string(details))

			// If we have an explainer (pulumi-cloud) we can offer to explain the changes.
			if explainer != nil && explainer.IsExplainPreviewEnabled(ctx, opts.Display) {
				choices = append(choices, explainChoice)
			}
		}

		var previewWarning string
		if opts.SkipPreview {
			previewWarning = colors.SpecWarning + " without a preview" + colors.Bold
		}

		// Create a prompt. If this is a refresh, we'll add some extra text so it's clear we aren't updating resources.
		prompt := "\b" + opts.Display.Color.Colorize(
			colors.SpecPrompt+fmt.Sprintf("Do you want to perform this %s%s?",
				updateTextMap[kind].previewText, previewWarning)+colors.Reset)
		if kind == apitype.RefreshUpdate {
			prompt += "\n" +
				opts.Display.Color.Colorize(colors.SpecImportant+
					"No resources will be modified as part of this refresh; just your stack's state will be.\n"+
					colors.Reset)
		}

		// Now prompt the user for a yes, no, or details, and then proceed accordingly.
		allAskOpts := append([]survey.AskOpt{surveyIcons}, askOpts...)
		if err := survey.AskOne(&survey.Select{
			Message: prompt,
			Options: choices,
			Default: string(no),
		}, &response, allAskOpts...); err != nil {
			return nil, fmt.Errorf("confirmation cancelled, not proceeding with the %s: %w", kind, err)
		}

		if response == string(no) {
			return nil, result.FprintBailf(os.Stdout, "confirmation declined, not proceeding with the %s", kind)
		}

		if response == string(yes) {
			// If we're in experimental mode always use the plan
			if opts.Engine.Experimental {
				return plan, nil
			}
			return nil, nil
		}

		if response == string(details) {
			diff, err := display.CreateDiff(events, opts.Display)
			if err != nil {
				return nil, err
			}
			_, err = os.Stdout.WriteString(diff + "\n")
			contract.IgnoreError(err)
			continue
		}

		if response == explainChoice {
			contract.Assertf(explainer != nil, "explainer must be present if explain option was selected")

			explanation, err := explainer.Explain(ctx, stackRef, kind, op, events)

			stdout := op.Opts.Display.Stdout
			if stdout == nil {
				stdout = os.Stdout
			}
			if err != nil {
				_, err = stdout.Write([]byte(
					opts.Display.Color.Colorize(
						"An error occurred while explaining the changes:\n" +
							colors.BrightRed + err.Error() + colors.Reset + "\n\n")))
				contract.IgnoreError(err)
				continue
			}
			_, err = stdout.Write([]byte(explanation + "\n"))
			contract.IgnoreError(err)
			continue
		}
	}
}

func PreviewThenPromptThenExecute(ctx context.Context, kind apitype.UpdateKind, stack Stack,
	op UpdateOperation, apply Applier, explainer Explainer, events chan<- engine.Event,
) (sdkDisplay.ResourceChanges, error) {
	// Preview the operation to the user and ask them if they want to proceed.
	if !op.Opts.SkipPreview {
		// We want to run the preview with the given plan and then run the full update with the initial plan as well,
		// but because plans are mutated as they're checked we need to clone it here.
		// We want to use the original plan because a program could be non-deterministic and have a plan of
		// operations P0, the update preview could return P1, and then the actual update could run P2, were P1 < P2 < P0.
		var originalPlan *deploy.Plan
		if op.Opts.Engine.Plan != nil {
			originalPlan = op.Opts.Engine.Plan.Clone()
		}

		plan, changes, err := PreviewThenPrompt(ctx, kind, stack, op, apply, explainer)
		if err != nil || kind == apitype.PreviewUpdate {
			return changes, err
		}

		// If we had an original plan use it, else if prompt said to use the plan from Preview then use the
		// newly generated plan
		if originalPlan != nil {
			op.Opts.Engine.Plan = originalPlan
		} else if plan != nil {
			op.Opts.Engine.Plan = plan
		} else {
			op.Opts.Engine.Plan = nil
		}
	}

	// Perform the change (!DryRun) and show the cloud link to the result.
	// We don't care about the events it issues, so just pass a nil channel along.
	opts := ApplierOptions{
		DryRun:   false,
		ShowLink: true,
	}
	// No need to generate a plan at this stage, there's no way for the system or user to extract the plan
	// after here.
	op.Opts.Engine.GeneratePlan = false
	_, changes, res := apply(ctx, kind, stack, op, opts, events)
	return changes, res
}

type updateStats struct {
	numNonStackResources int
	retainedResources    []engine.StepEventMetadata
}

func computeUpdateStats(events []engine.Event) updateStats {
	var stats updateStats

	for _, e := range events {
		if e.Type != engine.ResourcePreEvent {
			continue
		}
		p, ok := e.Payload().(engine.ResourcePreEventPayload)

		if !ok {
			continue
		}

		if p.Metadata.Type.String() != "pulumi:pulumi:Stack" {
			stats.numNonStackResources++
		}

		// Track deleted resources that are retained.
		switch p.Metadata.Op {
		case deploy.OpDelete, deploy.OpReplace:
			if old := p.Metadata.Old; old != nil && old.State != nil && old.State.RetainOnDelete {
				stats.retainedResources = append(stats.retainedResources, p.Metadata)
			}
		}
	}
	return stats
}
