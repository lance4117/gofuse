package fileio

import (
	"fmt"
	"os"
	"testing"

	"github.com/lance4117/gofuse/times"
)

func TestCSVWriter(t *testing.T) {
	filename := fmt.Sprintf("test-writer-%d.csv", times.NowMilli())
	defer os.Remove(filename) // 清理测试文件

	// 创建CSV写入器
	writer, err := NewCSVWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create CSV writer: %v", err)
	}
	defer writer.Close()

	// 写入表头
	err = writer.WriteHeader([]string{"Name", "Age", "City"})
	if err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	// 写入数据行
	err = writer.WriteRow([]string{"Alice", "25", "Beijing"})
	if err != nil {
		t.Fatalf("Failed to write row: %v", err)
	}

	// 使用通用 Write 接口
	err = writer.Write([]string{"Bob", "30", "Shanghai"})
	if err != nil {
		t.Fatalf("Failed to write using Write(): %v", err)
	}

	t.Log("CSV writer test passed")
}

func TestCSVReader(t *testing.T) {
	filename := fmt.Sprintf("test-reader-%d.csv", times.NowMilli())
	defer os.Remove(filename) // 清理测试文件

	// 先创建测试文件
	writer, err := NewCSVWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create CSV writer: %v", err)
	}

	err = writer.WriteHeader([]string{"Name", "Age", "City"})
	if err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	err = writer.Write([][]string{
		{"Alice", "25", "Beijing"},
		{"Bob", "30", "Shanghai"},
		{"Charlie", "35", "Guangzhou"},
	})
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}
	writer.Close()

	// 读取文件
	reader, err := NewCSVReader(filename)
	if err != nil {
		t.Fatalf("Failed to create CSV reader: %v", err)
	}
	defer reader.Close()

	// 读取所有数据
	data, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read all data: %v", err)
	}

	if len(data) != 4 { // 1 header + 3 rows
		t.Fatalf("Expected 4 rows, got %d", len(data))
	}

	t.Logf("CSV reader test passed, read %d rows", len(data))
}

func TestTableWriterInterface(t *testing.T) {
	filename := fmt.Sprintf("test-interface-%d.csv", times.NowMilli())
	defer os.Remove(filename)

	// 测试接口实现
	var writer TableWriter
	var err error
	writer, err = NewCSVWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create TableWriter: %v", err)
	}
	defer writer.Close()

	err = writer.WriteHeader([]string{"ID", "Value"})
	if err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}

	err = writer.WriteRow([]string{"1", "Test"})
	if err != nil {
		t.Fatalf("Failed to write row: %v", err)
	}

	t.Log("TableWriter interface test passed")
}

func TestTableReaderInterface(t *testing.T) {
	filename := fmt.Sprintf("test-reader-interface-%d.csv", times.NowMilli())
	defer os.Remove(filename)

	// 创建测试文件
	writer, _ := NewCSVWriter(filename)
	writer.Write([]string{"Header1", "Header2"})
	writer.Write([]string{"Value1", "Value2"})
	writer.Close()

	// 测试接口实现
	var reader TableReader
	var err error
	reader, err = NewCSVReader(filename)
	if err != nil {
		t.Fatalf("Failed to create TableReader: %v", err)
	}
	defer reader.Close()

	data, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read all: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("Expected 2 rows, got %d", len(data))
	}

	t.Log("TableReader interface test passed")
}
