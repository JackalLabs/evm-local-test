// Code generated by "stringer -type=mainContent"; DO NOT EDIT.

package tui

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[testCasesMain-0]
	_ = x[cosmosSummaryMain-1]
}

const _mainContent_name = "testCasesMaincosmosSummaryMain"

var _mainContent_index = [...]uint8{0, 13, 30}

func (i mainContent) String() string {
	if i < 0 || i >= mainContent(len(_mainContent_index)-1) {
		return "mainContent(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _mainContent_name[_mainContent_index[i]:_mainContent_index[i+1]]
}
