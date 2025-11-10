package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lance4117/gofuse/fileio"
	"github.com/lance4117/gofuse/times"
)

func TestXMLWriter(t *testing.T) {
	filename := fmt.Sprintf("test-xml-writer-%d.xml", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewXMLWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create XML writer: %v", err)
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

	t.Log("XML writer test passed")
}

func TestXMLWriterArray(t *testing.T) {
	filename := fmt.Sprintf("test-xml-array-%d.xml", times.NowMilli())
	defer os.Remove(filename)

	writer, err := fileio.NewXMLWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create XML writer: %v", err)
	}
	defer writer.Close()

	persons := []any{
		Person{Name: "Alice", Age: 25, Email: "alice@example.com"},
		Person{Name: "Bob", Age: 30, Email: "bob@example.com"},
	}

	err = writer.WriteArray(persons)
	if err != nil {
		t.Fatalf("Failed to write array: %v", err)
	}

	t.Log("XML array writer test passed")
}

func TestXMLReader(t *testing.T) {
	filename := fmt.Sprintf("test-xml-reader-%d.xml", times.NowMilli())
	defer os.Remove(filename)

	// 先写入数据
	writer, _ := fileio.NewXMLWriter(filename)
	person := Person{Name: "Alice", Age: 25, Email: "alice@example.com"}
	writer.WriteObject(person)
	writer.Close()

	// 读取数据
	reader, err := fileio.NewXMLReader(filename)
	if err != nil {
		t.Fatalf("Failed to create XML reader: %v", err)
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

	t.Logf("XML reader test passed: %+v", readPerson)
}

func TestXMLStructuredWriterInterface(t *testing.T) {
	filename := fmt.Sprintf("test-xml-interface-%d.xml", times.NowMilli())
	defer os.Remove(filename)

	var writer fileio.StructuredWriter
	var err error
	writer, err = fileio.NewXMLWriter(filename)
	if err != nil {
		t.Fatalf("Failed to create StructuredWriter: %v", err)
	}
	defer writer.Close()

	person := Person{Name: "Test", Age: 20, Email: "test@example.com"}
	err = writer.WriteObject(person)
	if err != nil {
		t.Fatalf("Failed to write object: %v", err)
	}

	t.Log("XML StructuredWriter interface test passed")
}
