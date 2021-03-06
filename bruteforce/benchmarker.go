package bruteforce

import (
	"github.com/ngirot/BruteForce/bruteforce/conf"
	"github.com/ngirot/BruteForce/bruteforce/hashs"
	"github.com/ngirot/BruteForce/bruteforce/words"
	"math"
	"time"
)

func BenchHasherOneCpu(hasherCreator func() hashs.Hasher) int {
	var buildActionFunc = getBuildActionFuncForHasher(hasherCreator)
	return bench(buildActionFunc, 1)
}

func BenchHasherMultiCpu(hasherCreator func() hashs.Hasher) int {
	var buildActionFunc = getBuildActionFuncForHasher(hasherCreator)
	return bench(buildActionFunc, conf.BestNumberOfGoRoutine())
}

func BenchWorderOneCpu() int {
	var buildActionFunc = getBuildActionFuncForWorder()
	return bench(buildActionFunc, 1)
}

func BenchWorderMultiCpu() int {
	var buildActionFunc = getBuildActionFuncForWorder()
	return bench(buildActionFunc, conf.BestNumberOfGoRoutine())
}

func bench(buildActionFunc func() func(), cpus int) int {
	var chrono = NewChrono()
	chrono.Start()

	var count = 0
	var oneDone = func() {
		count++
	}

	var quit = make(chan bool)
	for i := 0; i < cpus; i++ {
		go actionLoop(buildActionFunc(), oneDone, quit)
	}

	time.Sleep(time.Second * 5)

	chrono.End()
	for i := 0; i < cpus; i++ {
		quit <- true
	}

	return int(math.Floor(float64(count) / chrono.DurationInSeconds()))
}

func actionLoop(action func(), oneDone func(), quit chan bool) {
	for {
		action()

		select {
		case <-quit:
			return
		default:
			oneDone()
		}
	}
}
func getBuildActionFuncForHasher(hasherCreator func() hashs.Hasher) func() func() {
	return func() func() {
		var hasher = hasherCreator()
		return func() {
			hasher.Hash("1234567890")
		}
	}
}

func getBuildActionFuncForWorder() func() func() {
	return func() func() {
		var alphabet = words.DefaultAlphabet()
		var worder = words.NewWorderAlphabet(alphabet, 1, 0)
		return func() {
			worder.Next()
		}
	}
}
