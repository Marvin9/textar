package pkg

func AttachOffsets(dict []Dictionary) {
	for dictIdx, dictionary := range dict {
		offset := int64(0)

		for idx, index := range dictionary.Indexes {
			dict[dictIdx].Indexes[idx].Offset = offset
			offset += int64(len(index.String))
		}
	}
}
