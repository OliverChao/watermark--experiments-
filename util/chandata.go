package util

type ChanMetaData struct {
	ID      uint
	Col     string  // column the need to update
	NewData float64 // the new number needed to be assigned to Col
	//	tx.Model(&Student{ID:ID}).Update(Col, new)
}

type VerifyResultData struct {
	TotalCount int
	MatchCount int
}
