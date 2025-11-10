package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/times"
)

// 测试用的结构体
type Person struct {
	Name  string `json:"name" xml:"name"`
	Age   int    `json:"age" xml:"age"`
	Email string `json:"email" xml:"email"`
}

type PersonList struct {
	Persons []Person `xml:"person"`
}

func TestJSONWriter(t *testing.T) {
	filename := fmt.Sprintf("test-json-writer-%d.json", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewJSONWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create JSON writer: %v", err)
	}
	defer writer.Close()

	person := Person{
		Name:  "Alice",
		Age:   25,
		Email: "alice@example.com",
	}

	err = writer.WriteObject(person)
	if err != nil {
		t.Fatalf("Failed to write object: %v", err)
	}

	t.Log("JSON writer test passed")
}

func TestJSONWriterArray(t *testing.T) {
	filename := fmt.Sprintf("test-json-array-%d.json", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewJSONWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create JSON writer: %v", err)
	}
	defer writer.Close()

	persons := []any{
		Person{Name: "Alice", Age: 25, Email: "alice@example.com"},
		Person{Name: "Bob", Age: 30, Email: "bob@example.com"},
		Person{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
	}

	err = writer.WriteArray(persons)
	if err != nil {
		t.Fatalf("Failed to write array: %v", err)
	}

	t.Log("JSON array writer test passed")
}

func TestJSONReader(t *testing.T) {
	filename := fmt.Sprintf("test-json-reader-%d.json", times.NowMilli())
	defer os.Remove(filename)

	// 先写入数据
	writer, _ := fileio.NewJSONWriter(filename)
	person := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
	writer.WriteObject(person)
	writer.Close()

	// 读取数据
	reader, err := fileio.NewJSONReader(filename)
	if err != nil {
		t.Fatalf("Failed to create JSON reader: %v", err)
	}
	defer reader.Close()

	var readPerson Person
	err = reader.ReadObject(&readPerson)
	if err != nil {
		t.Fatalf("Failed to read object: %v", err)
	}

	if readPerson.Name != "Alice" || readPerson.Age != 25 {
		t.Fatalf("Data mismatch: got %+v", readPerson)
	}

	t.Logf("JSON reader test passed: %+v", readPerson)
}

func TestJSONReaderArray(t *testing.T) {
	filename := fmt.Sprintf("test-json-reader-array-%d.json", times.NowMilli())
	defer os.Remove(filename)

	// 先写入数组
	writer, _ := fileio.NewJSONWriter(filename)
	persons := []Person{
		{Name: "Alice", Age: 25, Email: "alice@example.com"},
		{Name: "Bob", Age: 30, Email: "bob@example.com"},
	}
	writer.Write(persons)
	writer.Close()

	// 读取数组
	reader, err := fileio.NewJSONReader(filename)
	if err != nil {
		t.Fatalf("Failed to create JSON reader: %v", err)
	}
	defer reader.Close()

	var readPersons []Person
	err = reader.ReadArray(&readPersons)
	if err != nil {
		t.Fatalf("Failed to read array: %v", err)
	}

	if len(readPersons) != 2 {
		t.Fatalf("Expected 2 persons, got %d", len(readPersons))
	}

	t.Logf("JSON array reader test passed: read %d persons", len(readPersons))
}

func TestJSONStructuredWriterInterface(t *testing.T) {
	filename := fmt.Sprintf("test-json-interface-%d.json", times.NowMilli())
	defer os.Remove(filename)

	var writer fileio.StructuredWriter
	var err error
	writer, err = fileio.NewJSONWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create StructuredWriter: %v", err)
	}
	defer writer.Close()

	person := Person{Name: "Test", Age: 20, Email: "test@example.com"}
	err = writer.WriteObject(person)
	if err != nil {
		t.Fatalf("Failed to write object: %v", err)
	}

	t.Log("JSON StructuredWriter interface test passed")
}

// ==================== 通用接口测试 ====================

func TestStructuredWriterPolymorphism(t *testing.T) {
	person := Person{Name: "Polymorphism", Age: 99, Email: "poly@example.com"}

	// 测试 JSON
	jsonFile := fmt.Sprintf("test-poly-json-%d.json", times.NowMilli())
	defer os.Remove(jsonFile)

	var jsonWriter fileio.StructuredWriter
	jsonWriter, _ = fileio.NewJSONWriter(jsonFile)
	jsonWriter.WriteObject(person)
	jsonWriter.Close()

	// 测试 XML
	xmlFile := fmt.Sprintf("test-poly-xml-%d.xml", times.NowMilli())
	defer os.Remove(xmlFile)

	var xmlWriter fileio.StructuredWriter
	xmlWriter, _ = fileio.NewXMLWriter(xmlFile)
	xmlWriter.WriteObject(person)
	xmlWriter.Close()

	t.Log("Polymorphism test passed: same interface works for JSON and XML")
}

func TestInvalidFile(t *testing.T) {
	// 测试读取不存在的 JSON 文件
	_, err := fileio.NewJSONReader("nonexistent.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent JSON file")
	}

	// 测试读取不存在的 XML 文件
	_, err = fileio.NewXMLReader("nonexistent.xml")
	if err == nil {
		t.Fatal("Expected error for nonexistent XML file")
	}

	t.Log("Invalid file test passed")
}
