package fileio

import (
	"fmt"
	"testing"

	"gitee.com/lance4117/GoFuse/times"
)

func TestCSV(t *testing.T) {
	csv := NewCSVFileIO(fmt.Sprintf("test-%d", times.NowMilli()))
	err := csv.Create([]string{"1", "2"})
	if err != nil {
		t.Fatal(err)
	}
	err = csv.Write([]string{"data1", "data2"})
	if err != nil {
		t.Fatal(err)
	}

}
