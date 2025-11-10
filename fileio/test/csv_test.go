package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/times"
)

func TestCSVWriter(t *testing.T) {
	filename := fmt.Sprintf("test-writer-%d.csv", times.NowMilli())
	defer os.Remove(filename) // 清理测试文件

	// 创建CSV写入器
	writer, err := fileio.NewCSVWriter(filename)
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
	writer, err := fileio.NewCSVWriter(filename)
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
	reader, err := fileio.NewCSVReader(filename)
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
	var writer fileio.TableWriter
	var err error
	writer, err = fileio.NewCSVWriter(filename)
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
	writer, _ := fileio.NewCSVWriter(filename)
	writer.Write([]string{"Header1", "Header2"})
	writer.Write([]string{"Value1", "Value2"})
	writer.Close()

	// 测试接口实现
	var reader fileio.TableReader
	var err error
	reader, err = fileio.NewCSVReader(filename)
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

// TestCSVWriterBatchPerformance 测试批量写入性能（手动控制刷新）
func TestCSVWriterBatchPerformance(t *testing.T) {
	filename := fmt.Sprintf("test-batch-%d.csv", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewCSVWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create CSV writer: %v", err)
	}
	defer writer.Close()

	// 写入表头
	writer.WriteHeader([]string{"ID", "Name", "Score"})

	// 批量写入 1000 行，最后统一刷新
	for i := 0; i < 1000; i++ {
		err = writer.WriteRow([]string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("User%d", i),
			fmt.Sprintf("%d", 60+i%40),
		})
		if err != nil {
			t.Fatalf("Failed to write row %d: %v", i, err)
		}
	}

	// 手动刷新
	if err := writer.Flush(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}

	t.Logf("Batch performance test passed: wrote 1000 rows")
}

// TestCSVReaderRowByRow 测试逐行读取
func TestCSVReaderRowByRow(t *testing.T) {
	filename := fmt.Sprintf("test-rowbyrow-%d.csv", times.NowMilli())
	defer os.Remove(filename)

	// 创建测试文件
	writer, _ := fileio.NewCSVWriter(filename)
	writer.Write([][]string{
		{"Name", "Age"},
		{"Alice", "25"},
		{"Bob", "30"},
		{"Charlie", "35"},
	})
	writer.Close()

	// 逐行读取
	reader, err := fileio.NewCSVReader(filename)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
	defer reader.Close()

	// 读取表头
	header, err := reader.ReadHeader()
	if err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}
	if len(header) != 2 || header[0] != "Name" || header[1] != "Age" {
		t.Fatalf("Invalid header: %v", header)
	}

	// 逐行读取数据
	rowCount := 0
	for {
		row, err := reader.ReadRow()
		if err != nil {
			// 到达文件末尾
			break
		}
		rowCount++
		if len(row) != 2 {
			t.Fatalf("Invalid row: %v", row)
		}
	}

	if rowCount != 3 {
		t.Fatalf("Expected 3 data rows, got %d", rowCount)
	}

	t.Logf("Row-by-row reading test passed: read %d rows", rowCount)
}

// TestCSVReaderInvalidFile 测试读取不存在的文件
func TestCSVReaderInvalidFile(t *testing.T) {
	_, err := fileio.NewCSVReader("nonexistent-file.csv")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
	t.Logf("Invalid file test passed: %v", err)
}

// TestCSVWriterFlushControl 测试刷新控制
func TestCSVWriterFlushControl(t *testing.T) {
	filename := fmt.Sprintf("test-flush-%d.csv", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewCSVWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()

	// 写入数据但不刷新
	writer.WriteRow([]string{"Header1", "Header2"})
	writer.WriteRow([]string{"Value1", "Value2"})

	// 手动刷新
	if err := writer.Flush(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}

	// 继续写入
	writer.WriteRow([]string{"Value3", "Value4"})

	// Close 会自动刷新
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}

	// 验证文件内容
	reader, _ := fileio.NewCSVReader(filename)
	defer reader.Close()
	data, _ := reader.ReadAll()

	if len(data) != 3 {
		t.Fatalf("Expected 3 rows, got %d", len(data))
	}

	t.Log("Flush control test passed")
}
