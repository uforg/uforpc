package ast

type LineDiff struct {
	StartToStart int
	StartToEnd   int
	EndToStart   int
	EndToEnd     int
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
