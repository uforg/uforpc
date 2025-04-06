package ast

type LineDiff struct {
	// The difference in lines between the start of the first position and the start of the second position.
	StartToStart int
	// The difference in lines between the start of the first position and the end of the second position.
	StartToEnd int
	// The difference in lines between the end of the first position and the start of the second position.
	EndToStart int
	// The difference in lines between the end of the first position and the end of the second position.
	EndToEnd int
}

// GetLineDiff returns the line diff between two positions.
func GetLineDiff(from, to WithPositions) LineDiff {
	return LineDiff{
		StartToStart: from.GetPositions().Pos.Line - to.GetPositions().Pos.Line,
		StartToEnd:   from.GetPositions().Pos.Line - to.GetPositions().EndPos.Line,
		EndToStart:   from.GetPositions().EndPos.Line - to.GetPositions().Pos.Line,
		EndToEnd:     from.GetPositions().EndPos.Line - to.GetPositions().EndPos.Line,
	}
}
