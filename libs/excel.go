package libs

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

func SimpleReadExcel(r io.Reader) ([][]string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// xls 文件
	xlFile, err := xls.OpenReader(bytes.NewReader(bs), "utf-8")
	if err == nil {
		return xlFile.ReadAllCells(10000), nil
	}
	// xlsx 文件
	file, err := xlsx.OpenBinary(bs)
	if err != nil {
		return nil, err
	}
	table, err := file.ToSlice()
	if err != nil {
		return nil, err
	}
	if len(table) < 1 {
		return nil, errors.New("没有Sheet")
	}
	for _, item := range table[1:] {
		table[0] = append(table[0], item...)
	}
	return table[0], nil
}

func SimpleWriteExcel(f io.Writer, table [][]string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	for _, item := range table {
		row = sheet.AddRow()
		for _, col := range item {
			cell = row.AddCell()
			cell.Value = col
		}
	}
	if err = file.Write(f); err != nil {
		return err
	}
	return nil
}

func SimpleWriteCsv(f io.Writer, table [][]string) error {
	_, err := f.Write([]byte("\xEF\xBB\xBF")) // 写入UTF-8 BOM
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	for _, item := range table {
		err = w.Write(item)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
