// Package coopdatadog starts and stops the Datadog integration: tracer,
// profiler and Dogstatsd metrics.
//
// Deprecated: The whole module is deprecated. Use Start from
// github.com/coopnorge/go-telemetry-lib and OpenTelemetry instead. The migration
// guide is docs/migration-to-opentelemetry.md in this repository, and the
// decision is recorded in the ADR
// https://inventory.internal.coop/catalog/default/component/go-telemetry-lib/adrs?record=0001-deprecate-go-datadog-lib.md.
// The module will be archived on 2027-03-01. It receives non-breaking dependency
// updates only until then.
package coopdatadog
