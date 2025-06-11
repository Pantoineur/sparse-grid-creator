package main

func GenerateCells(size int) []Cell {
	cells := make([]Cell, size*size)

	for y := range size {
		for x := range size {
			cells[y*size+x] = Cell{
				X: x,
				Y: y,
			}
		}
	}

	return cells
}

func SetClosestHorizontal(m model, negative bool) map[Cell]bool {
	count := 0

	for {
		if negative {
			count--
		} else {
			count++
		}

		newCell := Cell{X: m.cursor.X + count, Y: m.cursor.Y}
		if (newCell.X > m.grid.Size-1 && !negative) || (newCell.X < 0 && negative) {

			return m.additionalCursors
		}

		if _, ok := m.additionalCursors[newCell]; !ok {
			m.additionalCursors[newCell] = true
			return m.additionalCursors
		}
	}
}

func SetClosestVertical(m model, negative bool) map[Cell]bool {
	count := 0

	for {
		if negative {
			count--
		} else {
			count++
		}

		newCell := Cell{X: m.cursor.X, Y: m.cursor.Y + count}
		if (newCell.Y > m.grid.Size-1 && !negative) || (newCell.Y < 0 && negative) {
			return m.additionalCursors
		}

		if _, ok := m.additionalCursors[newCell]; !ok {
			m.additionalCursors[newCell] = true
			return m.additionalCursors
		}
	}
}
