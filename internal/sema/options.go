package sema

// AnalysisOptions configures optional sema passes.
type AnalysisOptions struct {
	WarnUnused      bool
	BeginnerLint    bool
	WarnUnreachable bool
}
