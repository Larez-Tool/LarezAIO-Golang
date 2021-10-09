package memory

import "github.com/p3tr0v/chacal/antimem"

func StartAntiMemory() bool {
	return antimem.ByMemWatcher()
}
