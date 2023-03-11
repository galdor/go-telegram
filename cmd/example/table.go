package main

import (
	"fmt"
	"os"
	"strings"
)

type Table struct {
	Header []string
	Rows   [][]interface{}
}

func NewTable(header ...string) *Table {
	return &Table{
		Header: header,
		Rows:   make([][]interface{}, 0),
	}
}

func (t *Table) AddRow(row ...interface{}) {
	t.Rows = append(t.Rows, row)
}

func (t *Table) Write() {
	rows := t.Render()
	widths := t.ColumnWidths(rows)

	if len(t.Header) > 0 {
		for i, label := range t.Header {
			if i > 0 {
				fmt.Fprintf(os.Stderr, "  ")
			}

			label = fmt.Sprintf("%-*s", widths[i], strings.ToUpper(label))
			fmt.Fprintf(os.Stderr, label)
		}

		fmt.Fprintln(os.Stderr)
	}

	for _, row := range rows {
		for j, s := range row {
			if j > 0 {
				fmt.Printf("  ")
			}

			fmt.Printf("%-*s", widths[j], s)
		}

		fmt.Println("")
	}
}

func (t *Table) Render() [][]string {
	rows := make([][]string, len(t.Rows))

	for i, row := range t.Rows {
		rows[i] = make([]string, len(row))

		for j, value := range row {
			rows[i][j] = t.RenderValue(value)
		}
	}

	return rows
}

func (t *Table) RenderValue(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func (t *Table) ColumnWidths(rows [][]string) []int {
	widths := make([]int, t.NbColumns())

	for i, label := range t.Header {
		widths[i] = len(label)
	}

	for _, row := range rows {
		for j, value := range row {
			if len(value) > widths[j] {
				widths[j] = len(value)
			}
		}
	}

	return widths
}

func (t *Table) NbColumns() int {
	n := len(t.Header)

	for _, row := range t.Rows {
		if len(row) > n {
			n = len(row)
		}
	}

	return n
}
