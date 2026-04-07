package utils

import "sync"

func SplitMessage(message string, chunkSize int) []string {
	var result []string
	for i := 0; i < len(message); i += chunkSize {
		end := i + chunkSize
		if end > len(message) {
			end = len(message)
		}
		result = append(result, message[i:end])
	}
	return result
}

// Batchify allows to make processing of an array use go concurrency.
// the threads specifies in practice how many threads is used, how many gouroutines
// run at the same time.
func Batchify[A any](array []A, threads int, batchConsumer func(batch []A, batchNumber int)) {
	if threads <= 0 {
		panic("batchify MUST have threads > 0")
	}
	batches := [][]A{}
	chunkSize := (len(array) + threads - 1) / threads
	for i := 0; i < len(array); i += chunkSize {
		end := i + chunkSize
		if end > len(array) {
			end = len(array)
		}
		batches = append(batches, array[i:end])
	}

	wg := sync.WaitGroup{}
	for i, batch := range batches {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batchConsumer(batch, i)
		}()
	}
	wg.Wait()
}
