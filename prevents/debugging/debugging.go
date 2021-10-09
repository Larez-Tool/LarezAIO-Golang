package debugging

import (
	"github.com/p3tr0v/chacal/antidebug"
)

func StartAntiDebugging() bool {
	return antidebug.ByProcessWatcher()
}
