package checkadguard

import (
	// "bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// ExitCode contain default text prefix for Nagios compatible exit codes
var ExitCode = map[int]string{
	0: "OK",
	1: "WARNING",
	2: "CRITICAL",
	3: "UNKNOWN",
}

// AdGuardStatus defines structure for json output from AdGuard API for its status
type AdGuardStatus struct {
	dns_addresses []string
	dns_port int
	http_port int
	language string
	protection_enabled bool
	running bool
	version string
}

// CheckAdGuard will perform actual checks and generate status information
func CheckAdGuard(host, port, username, password, mode string, timeout int, ssl, insecure, metrics bool) {
	realHost := makeURL(host, port, ssl)
	authHeader := generateAuth(username, password)

	client := http.Client{
		Timeout: time.Duration(timeout * time.Second),
	}

	req, err := http.Request("GET", realHost + "/status")
	req.Header.Set("Authorization", authHeader)
	exitMessage(err, 3)

	resp, err := client.Do(req)
	exitMessage(err, 3)

	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	exitMessage(err, 3)

	var adguardStatus AdGuardStatus
	err = json.Unmarshal(bodyByte, &adguardStatus)
	exitMessage(err, 3)

	if adguardStatus.running {
		if adguardStatus.protection_enabled {
			exitMessage("AdGuard is running and protection is enabled", 0)
		} else {
			exitMessage("AdGuard is running but protection is disabled", 1)
		}
	} else {
		exitMessage("AdGuard is not running", 2)
	}
}

func exitMessage(err error, exitCode int) {
	if err != nil {
		fmt.Printf("%s: %v\n", ExitCode[exitCode], err)
        os.Exit(exitCode)
	}
}

func generateAuth(username, password string) string {
	authPlain := username + ":" + password
	authBase64 := base64.StdEncoding.EncodeToString([]byte(authPlain))
	authHeader := "Authorization: Basic " + authBase64

	return authHeader
}

func makeURL(host, port string, ssl bool) string {
	prefix := "http://"
	if ssl {
		prefix = "https://"
	}

	hostList := []string{prefix, host, ":", port, "/control"}
	hostString := strings.Join(hostList, "")

	return hostString
}
