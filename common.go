package main

func containsStringPointer(strSlice []*string, searchStr *string) bool {
	for _, value := range strSlice {
		if *value == *searchStr {
			return true
		}
	}
	return false
}

func containsString(strSlice []string, searchStr string) bool {
	for _, value := range strSlice {
		if value == searchStr {
			return true
		}
	}
	return false
}

func dedupeStringPointer(strSlice []*string) []*string {
	var returnSlice []*string
	for _, value := range strSlice {
		if !containsStringPointer(returnSlice, value) {
			returnSlice = append(returnSlice, value)
		}
	}
	return returnSlice
}

func dedupeString(strSlice []string) []string {
	var returnSlice []string
	for _, value := range strSlice {
		if !containsString(returnSlice, value) {
			returnSlice = append(returnSlice, value)
		}
	}
	return returnSlice
}

// makeBatchesStringPointer takes a slice of string pointers and returns them as
// a slice of string pointer slices in batch size of batchSize. Useful for splitting
// up work into batches for parallel operations.
func makeBatchesStringPointer(strSlice []*string, batchSize int) (batches [][]*string) {
	numBatches, remainder := len(strSlice)/batchSize, len(strSlice)%batchSize
	// build full batches
	for i := 1; i <= numBatches; i++ {
		var startIndex int
		endIndex := i * batchSize
		if i == 1 {
			startIndex = 0
		} else {
			startIndex = batchSize * (i - 1)
		}
		var b []*string
		b = strSlice[startIndex:endIndex]
		batches = append(batches, b)
	}
	if remainder > 0 {
		// build last partial batch
		startIndex := (len(strSlice) - remainder)
		endIndex := len(strSlice)
		var b []*string
		b = strSlice[startIndex:endIndex]
		batches = append(batches, b)
	}
	return batches
}
