// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sampling

import (
	"sync"
	"time"

	"go.opentelemetry.io/collector/model/pdata"
)

// TraceData stores the sampling related trace data.
type TraceData struct {
	sync.Mutex
	// Decisions gives the current status of the sampling decision for each policy.
	Decisions []Decision
	// Arrival time the first span for the trace was received.
	ArrivalTime time.Time
	// Decisiontime time when sampling decision was taken.
	DecisionTime time.Time
	// SpanCount track the number of spans on the trace.
	SpanCount int64
	// ReceivedBatches stores all the batches received for the trace.
	ReceivedBatches []pdata.Traces
}

// Decision gives the status of sampling decision.
type Decision int32

const (
	// Unspecified indicates that the status of the decision was not set yet.
	Unspecified Decision = iota
	// Pending indicates that the policy was not evaluated yet.
	Pending
	// Sampled is used to indicate that the decision was already taken
	// to sample the data.
	Sampled
	// NotSampled is used to indicate that the decision was already taken
	// to not sample the data.
	NotSampled
	// Dropped is used when data needs to be purged before the sampling policy
	// had a chance to evaluate it.
	Dropped
	// Error is used to indicate that policy evaluation was not succeeded.
	Error
	// InvertSampled is used on the invert match flow and indicates to sample
	// the data.
	InvertSampled
	// InvertNotSampled is used on the invert match flow and indicates to not
	// sample the data.
	InvertNotSampled
)

// PolicyEvaluator implements a tail-based sampling policy evaluator,
// which makes a sampling decision for a given trace when requested.
type PolicyEvaluator interface {
	// OnLateArrivingSpans notifies the evaluator that the given list of spans arrived
	// after the sampling decision was already taken for the trace.
	// This gives the evaluator a chance to log any message/metrics and/or update any
	// related internal state.
	OnLateArrivingSpans(earlyDecision Decision, spans []*pdata.Span) error

	// Evaluate looks at the trace data and returns a corresponding SamplingDecision.
	Evaluate(traceID pdata.TraceID, trace *TraceData) (Decision, error)
}
