package utils

func GetBatchGenerator(baseChan chan uint64, batchSize int) chan []uint64 {
	if batchSize < 1 {
		batchSize = 10
	}
	result := make(chan []uint64)

	go func () {
		batch := make([]uint64, 0, batchSize)
		for value := range baseChan {
			batch = append(batch, value)
			if len(batch) >= batchSize {
				result <- batch
				batch = make([]uint64, 0, batchSize)
			}
		}
		if len(batch) > 0 {
			result <- batch
		}
		close(result)
	}()

	return result
}
