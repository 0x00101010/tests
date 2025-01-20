## op-node-derivers

Key abstractions

```golang
type Deriver interface {
    OnEvent(ev Event) bool
}

DeriverMux // DeriverMux takes an event-signal as deriver, and synchronously fans it out to all contained Deriver ends.
EngDeriver
PipelineDeriver
SyncDeriver
EngineResetDeriver
L2Verifier
DisabledSequencer
L1OriginSelector
Sequencer
L1Tracker
StatusTracker

// debug / tests
DebugDeriver
NoopDeriver
DeriverFunc

// potential?
AttributesHandler
CLSync
StepSchedulingDeriver
AltDAFinalizer
Finalizer
InteropDeriver
systemActor
```

```golang
type Emitter interface {
	Emit(ev Event)
}

EmitterFunc
NoopEmitter
MockEmitter
Limiter
systemActor

emitter seems to all be handled by Sys
```

```golang
type Event interface {
	// String returns the name of the event.
	// The name must be simple and identify the event type, not the event content.
	// This name is used for metric-labeling.
	String() string
}

// error:
CriticalErrorEvent

// forkchoice:
ForkchoiceRequestEvent
ForkchoiceUpdateEvent
PromoteUnsafeEvent

// interop:
RequestCrossUnsafeEvent
PromoteCrossUnsafeEvent
CrossUnsafeUpdateEvent
InteropPendingSafeChangedEvent
RequestCrossSafeEvent
CrossSafeUpdateEvent
PromoteSafeEvent
CrossUpdateRequestEvent

// derivation:
UnsafeUpdateEvent
PendingSafeUpdateEvent
PromotePendingSafeEvent
PromoteLocalSafeEvent
LocalSafeUpdateEvent
SafeDerivedEvent
ProcessAttributesEvent
PendingSafeRequestEvent
ProcessUnsafePayloadEvent
TryBackupUnsafeReorgEvent
ForceEngineResetEvent
EngineResetConfirmedEvent
PromoteFinalizedEvent
FinalizedUpdateEvent
RequestFinalizedUpdateEvent
```

```golang
type Drainer interface {
	// Drain processes all events.
	Drain() error
	// DrainUntil processes all events until a condition is hit.
	// If excl, the event that matches the condition is not processed yet.
	// If not excl, the event that matches is processed.
	DrainUntil(fn func(ev Event) bool, excl bool) error
}
```

```golang
type EmitterDrainer interface {
	Emitter
	Drainer
}
```