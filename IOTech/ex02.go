//Bryn Jones
//25/01/2023
//IOTech Exercise 02

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Stores an array of all devices read in from input json
type Devices struct {
	Devices []Device `json:"Devices"`
}

// Struct for storing information about a device.
// Input json unmarshals to this
type Device struct {
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Info      string `json:"Info"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

// Unmarshal read in data to Devices struct.
func UnmarshalDevices(data []byte) (Devices, error) {
	var devices Devices
	err := json.Unmarshal(data, &devices)
	return devices, err
}

// Struct that mirrors provided schema.json
// Parsed data will be marshalled to this
type Output struct {
	Uuids      []string `json:"UUIDS"`
	ValueTotal int64    `json:"ValueTotal"`
}

// Marshal to Output struct
func (output *Output) Marshal() ([]byte, error) {
	return json.Marshal(output)
}

func run() {
	//Get the current unix timestamp
	var currentTime = int(time.Now().Unix())
	//fmt.Println("Current time: ", currentTime)
	//Stores all devices read in from input json
	var allDevices Devices
	//Stores parsed output for marshalling to output json
	var output Output

	//regex for parsing uuid, assuming uuids are uniformly 36 characters
	//Matches the 36 characters following the string "uuid:"
	regex := regexp.MustCompile("uuid:(.{36})")

	var inputPath string
	var outputPath string

	inputPath = os.Args[1]
	outputPath = os.Args[2]

	//fmt.Println("Input file: " + inputPath + "\n")
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//read opened file as a byte array
	byteValue, _ := io.ReadAll(file)

	//Unmarshal byte array into the array of devices
	allDevices, err = UnmarshalDevices(byteValue)
	file.Close()

	var valueTotal int = 0
	for i := 0; i < len(allDevices.Devices); i++ {
		//Convert timestamp read in from data.json to an int
		deviceTime, _ := strconv.Atoi(allDevices.Devices[i].Timestamp)
		//Check if timestamp is before current time before parsing
		//i.e. discards timestamps from the past
		if deviceTime > currentTime {
			//Decoding base64 to an array of bytes
			rawBytes, err := base64.StdEncoding.DecodeString(allDevices.Devices[i].Value)
			if err != nil {
				fmt.Println(err)
			}
			//Convert array of bytes to an integer
			intValue, _ := strconv.Atoi(string(rawBytes))
			valueTotal += intValue
			//fmt.Printf("Decoded Value: %d\n", intValue)

			//Parse uuid from Info field using regex
			uuid := regex.FindStringSubmatch(allDevices.Devices[i].Info)[1]
			output.Uuids = append(output.Uuids, uuid)
			//fmt.Println("UUID: " + uuid + "\n")
		}

	}
	//fmt.Println("Total val: ", valueTotal)
	output.ValueTotal = int64(valueTotal)

	//Output to json

	var outputBytes []byte
	outputBytes, err = output.Marshal()
	//Open output file. 0644 permissions = -rw-rw-r--
	err = os.WriteFile(outputPath, outputBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Output to file: " + outputPath)
}
func main() {
	run()
}
