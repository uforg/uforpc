package lsp

import (
	"time"
)

// scheduleAnalysis schedules an analysis for the given URI with debouncing.
// If another analysis is scheduled within the debounce time, the previous one is cancelled.
func (l *LSP) scheduleAnalysis(uri string) {
	// Skip if analyzer is not available
	if l.analyzer == nil {
		l.logger.Warn("analyzer not available, skipping analysis")
		return
	}

	l.analysisTimerMu.Lock()
	defer l.analysisTimerMu.Unlock()

	// Cancel any existing timer
	if l.analysisTimer != nil {
		l.analysisTimer.Stop()
	}

	// Schedule a new analysis
	l.analysisTimer = time.AfterFunc(debounceTime, func() {
		// Check if another analysis is already in progress
		l.analysisInProgressMu.Lock()
		if l.analysisInProgress {
			l.analysisInProgressMu.Unlock()
			// If an analysis is already in progress, schedule another one
			l.scheduleAnalysis(uri)
			return
		}
		l.analysisInProgress = true
		l.analysisInProgressMu.Unlock()

		// Run the analysis
		l.analyzeAndPublishDiagnostics(uri)

		// Mark analysis as complete
		l.analysisInProgressMu.Lock()
		l.analysisInProgress = false
		l.analysisInProgressMu.Unlock()
	})
}
