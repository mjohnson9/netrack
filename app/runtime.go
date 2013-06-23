package app

import (
	"runtime"
)

func init() {
	// We want to run the number of CPUs + 1 processes. I pulled this number out of nowhere.
	// TODO: Actually do some research about how many processes we want, or if we should even be changing this.
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
}
