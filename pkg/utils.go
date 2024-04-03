package pkg

import "time"

func AttachOffsets(dict []Dictionary) {
	for dictIdx, dictionary := range dict {
		offset := int64(0)

		for idx, index := range dictionary.Indexes {
			dict[dictIdx].Indexes[idx].Offset = offset
			offset += int64(len(index.String))
		}
	}
}

func MeasureTime(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}
