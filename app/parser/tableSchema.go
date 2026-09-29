package parser

type TableSchema struct {
	Name     string
	Columns []string
	Rootpage int
}

func (ts *TableSchema) GetColumnIndex(name string) int {
	var index int = -1
	for i, column := range ts.Columns {
		if column == name {
			index = i
			break
		}
	}

	return index
}