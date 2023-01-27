// File comparison functionality from: https://go.dev/play/p/xDjugk2Yzz

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// Computes hashes of each input file and compares them.
// Two different hashes means the files are definitely NOT the same
// https://go.dev/play/p/xDjugk2Yzz
func eq(r1, r2 io.Reader) (bool, error) {
	w1 := sha256.New()
	w2 := sha256.New()
	n1, err1 := io.Copy(w1, r1)
	if err1 != nil {
		return false, err1
	}
	n2, err2 := io.Copy(w2, r2)
	if err2 != nil {
		return false, err2
	}

	var b1, b2 [sha256.Size]byte
	copy(b1[:], w1.Sum(nil))
	copy(b2[:], w2.Sum(nil))

	return n1 != n2 || b1 == b2, nil
}

func checkEquality(f1 string, f2 string) bool {
	r1 := strings.NewReader(f1)
	r2 := strings.NewReader(f2)

	result, err := eq(r1, r2)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

// Default test: Checks the program outputs the correct json file given the default data.json provided
func TestDefault(t *testing.T) {
	os.Args[1] = "exercise-02/data/data.json"
	os.Args[2] = "exercise-02/data/output.json"
	run()

	r1 := strings.NewReader("exercise-02/data/output.json")
	r2 := strings.NewReader("exercise-02/testing/defaultOutput.json")

	result, err := eq(r1, r2)
	if err != nil {
		fmt.Println(err)
	}
	got := result
	want := true

	if got != want {
		t.Errorf("TestDefault fail")
	} else {
		fmt.Println("TestDefault pass")
	}

}

func TestReadIn(t *testing.T) {
	var allDevices Devices
	file, err := os.Open("exercise-02/data/data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//read opened file as a byte array
	byteValue, _ := io.ReadAll(file)

	//Unmarshal byte array into the array of devices
	allDevices, err = UnmarshalDevices(byteValue)
	file.Close()

	var outputBytes []byte
	outputBytes, err = allDevices.Marshal()
	//Open output file. 0644 permissions = -rw-rw-r--
	err = os.WriteFile("exercise-02/testing/readInOutput.json", outputBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	got := checkEquality("exercise-02/testing/readInOutput.json", "exercise-02/data/data.json")
	want := true

	if got != want {
		t.Errorf("TestReadIn fail")
	} else {
		fmt.Println("TestReadIn pass")
	}
}
func (devices *Devices) Marshal() ([]byte, error) {
	return json.Marshal(devices)
}

func TestUUIDs(t *testing.T) {

}
