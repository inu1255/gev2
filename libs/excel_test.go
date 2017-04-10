package libs

import (
	"fmt"
	"os"

	"testing"
)

func TestWrite(t *testing.T) {
	filename := os.TempDir() + "a.xlsx"
	fmt.Println(filename)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Log(err)
	}
	SimpleWriteExcel(f, [][]string{[]string{"1", "2"}, []string{"3", "4"}})
}
