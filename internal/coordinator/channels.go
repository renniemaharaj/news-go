package coordinator

import "github.com/renniemaharaj/news/internal/types"

var PreResultsChannel = make(chan types.Report, 100)

var PreContentChannel = make(chan types.Report, 100)

var PreModelChannel = make(chan types.Report, 100)

var JobsCompleteChannel = make(chan types.Report, 100)
