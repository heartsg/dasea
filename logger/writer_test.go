package logger

import (
	"testing"
	"strconv"
	"os"
)

func TestNoRotate(t *testing.T) {
	filepath := "test.txt"
	
	//writer object creation
	writer := NewWriter(filepath, false, 100, 0)
	if writer == nil {
		t.Errorf("Writer failed to create %s.", filepath)
	}
	
	//check whether correctly written or not
	writedata := []byte("test\n")
	len, err := writer.Write(writedata)
	if len != 5 || err != nil {
		t.Errorf("Failed to write \"%s\": %s", filepath, err)
	}
	file, err := os.Open(filepath)
	if err != nil {
		t.Errorf("Failed to open read \"%s\": %s", filepath, err)
	}
	data := make([]byte, 5)
	len, err = file.Read(data)
	if len != 5 || err != nil {
		t.Errorf("Failed to read \"%s\": %s", filepath, err)
	}
	for i, b := range data {
		if b != writedata[i] {
			t.Errorf("\"%s\" has wrong data.")
		}
	}
	file.Close()
	
	//check size and whether backup will be created etc.
	for i := 0; i < 30; i++ {
		writer.Write(writedata)
	}
	
	stat, err := os.Stat(filepath + ".bak")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath+".bak")
	}
	if stat.Size() != 100 {
		t.Errorf("\"%s\" file size is: %d (should be 100)", filepath+".bak", stat.Size())
	}
	
	stat, err = os.Stat(filepath)
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath)
	}
	if stat.Size() != 55 {
		t.Errorf("\"%s\" file size is: %d (should be 55)", filepath, stat.Size())
	}
	
		
	//13 bytes
	writedata = []byte("i am testing\n")
	for i := 0; i < 10; i++ {
		writer.Write(writedata)
	}

	stat, err = os.Stat(filepath + ".bak")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath+".bak")
	}
	if stat.Size() != 94 {
		t.Errorf("\"%s\" file size is: %d (should be 94)", filepath+".bak", stat.Size())
	}
	
	stat, err = os.Stat(filepath)
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath)
	}
	if stat.Size() != 91 {
		t.Errorf("\"%s\" file size is: %d (should be 91)", filepath, stat.Size())
	}
		
	
	writer.Close()
	os.Remove(filepath)
	os.Remove(filepath + ".bak")
}



func TestRotate(t *testing.T) {
	filepath := "test.txt"
	
	//writer object creation
	writer := NewWriter(filepath, true, 100, 5)
	if writer == nil {
		t.Errorf("Writer failed to create %s.", filepath)
	}
	
	//check whether correctly written or not
	writedata := []byte("test\n")
	len, err := writer.Write(writedata)
	if len != 5 || err != nil {
		t.Errorf("Failed to write \"%s\": %s", filepath, err)
	}
	file, err := os.Open(filepath + ".0")
	if err != nil {
		t.Errorf("Failed to open read \"%s\": %s", filepath+".0", err)
	}
	data := make([]byte, 5)
	len, err = file.Read(data)
	if len != 5 || err != nil {
		t.Errorf("Failed to read \"%s\": %s", filepath+".0", err)
	}
	for i, b := range data {
		if b != writedata[i] {
			t.Errorf("\"%s\" has wrong data.")
		}
	}
	file.Close()
	
	//check size and whether backup will be created etc.
	for i := 0; i < 30; i++ {
		writer.Write(writedata)
	}
	
	stat, err := os.Stat(filepath + ".0")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath+".0")
	}
	if stat.Size() != 100 {
		t.Errorf("\"%s\" file size is: %d (should be 100)", filepath+".0", stat.Size())
	}
	
	stat, err = os.Stat(filepath+".1")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath+".1")
	}
	if stat.Size() != 55 {
		t.Errorf("\"%s\" file size is: %d (should be 55)", filepath+".1", stat.Size())
	}
	
		
	//13 bytes
	writedata = []byte("i am testing\n")
	for i := 0; i < 20; i++ {
		writer.Write(writedata)
	}

	stat, err = os.Stat(filepath + ".1")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath+".1")
	}
	if stat.Size() != 94 {
		t.Errorf("\"%s\" file size is: %d (should be 94)", filepath+".1", stat.Size())
	}
	
	stat, err = os.Stat(filepath + ".2")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".2")
	}
	if stat.Size() != 91 {
		t.Errorf("\"%s\" file size is: %d (should be 91)", filepath + ".2", stat.Size())
	}
	
	stat, err = os.Stat(filepath + ".3")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".3")
	}
	if stat.Size() != 91 {
		t.Errorf("\"%s\" file size is: %d (should be 91)", filepath + ".3", stat.Size())
	}

	stat, err = os.Stat(filepath + ".4")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".4")
	}
	if stat.Size() != 39 {
		t.Errorf("\"%s\" file size is: %d (should be 39)", filepath + ".4", stat.Size())
	}
	
	for i := 0; i < 10; i++ {
		writer.Write(writedata)
	}
	
	stat, err = os.Stat(filepath + ".4")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".4")
	}
	if stat.Size() != 91 {
		t.Errorf("\"%s\" file size is: %d (should be 91)", filepath + ".4", stat.Size())
	}
	
	
	stat, err = os.Stat(filepath + ".5")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".5")
	}
	if stat.Size() != 78 {
		t.Errorf("\"%s\" file size is: %d (should be 78)", filepath + ".5", stat.Size())
	}
	
	stat, err = os.Stat(filepath + ".0")
	if err == nil {
		t.Errorf("\"%s\" file should be deleted.", filepath + ".0")
	}

	for i := 0; i < 10; i++ {
		writer.Write(writedata)
	}
	
	stat, err = os.Stat(filepath + ".1")
	if err == nil {
		t.Errorf("\"%s\" file should be deleted.", filepath + ".1")
	}
	
	stat, err = os.Stat(filepath + ".2")
	if err == nil {
		t.Errorf("\"%s\" file should be deleted.", filepath + ".2")
	}
		
	
	writer.Close()
	os.Remove(filepath + ".3")
	os.Remove(filepath + ".4")
	os.Remove(filepath + ".5")
	os.Remove(filepath + ".6")
	os.Remove(filepath + ".7")
}


func TestSmallFileLimit(t *testing.T) {
	filepath := "test.txt"
	
	//writer object creation
	writer := NewWriter(filepath, true, 10, 0)
	if writer == nil {
		t.Errorf("Writer failed to create %s.", filepath)
	}
	
	//check whether correctly written or not
	writedata := []byte("test\n")
	len, err := writer.Write(writedata)
	if len != 5 || err != nil {
		t.Errorf("Failed to write \"%s\": %s", filepath, err)
	}
	file, err := os.Open(filepath + ".0")
	if err != nil {
		t.Errorf("Failed to open read \"%s\": %s", filepath+".0", err)
	}
	data := make([]byte, 5)
	len, err = file.Read(data)
	if len != 5 || err != nil {
		t.Errorf("Failed to read \"%s\": %s", filepath+".0", err)
	}
	for i, b := range data {
		if b != writedata[i] {
			t.Errorf("\"%s\" has wrong data.")
		}
	}
	file.Close()
	
	//check size and whether backup will be created etc.
	writedata = []byte("i am testing\n")
	for i := 0; i < 30; i++ {
		size, err := writer.Write(writedata)
		if err != nil {
			t.Errorf("Write failed")
		}
		if size != 13 {
			t.Errorf("Returned write size incorrect")
		}
	}
	//totally 30*13 + 5 bytes = 395
	//there should be 40 files and the last one is 5 bytes
	stat, err := os.Stat(filepath + ".39")
	if err != nil {
		t.Errorf("\"%s\" file does not exist.", filepath + ".39")
	}
	if stat.Size() != 5 {
		t.Errorf("\"%s\" file size is: %d (should be 5)", filepath + ".39", stat.Size())
	}
	
	
	writer.Close()
	for i := 0; i < 40; i++ {
		os.Remove(filepath + "." + strconv.Itoa(i))
	}
}