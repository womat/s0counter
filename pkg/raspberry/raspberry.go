package raspberry

// Edge represents the change in P level that triggers an interrupt.
type Edge string

const (
	// EdgeNone indicates no level transitions will trigger an interrupt
	EdgeNone Edge = "none"

	// EdgeRising indicates an interrupt is triggered when the P transitions from low to high.
	EdgeRising Edge = "rising"

	// EdgeFalling indicates an interrupt is triggered when the P transitions from high to low.
	EdgeFalling Edge = "falling"

	// EdgeBoth indicates an interrupt is triggered when the P changes level.
	EdgeBoth Edge = "both"
)
