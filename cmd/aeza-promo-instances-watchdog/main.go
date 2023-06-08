package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	ipAddresses             []string // env vars: IP_ADDRESSES
	aezaApiKey              string   //           AEZA_API_KEY
	watchDelay              string   //           WATCH_DELAY
	aezaApiBaseURL          string   //           AEZA_API_BASE_URL
	aezaVmBaseURL           string   //           AEZA_VM_BASE_URL
	ntfyUrl                 string   //           NTFY_URL
	ntfyChannel             string   //           NTFY_CHANNEL
	githubActionModeEnabled bool     //           GITHUB_ACTION_MODE
)

const (
	defaultAezaApiBaseURL = "https://core.aeza.net"
	defaultAezaVmBaseURL  = "https://vm.aeza.net"
	defaultWatchDelay     = "5m"
	defaultNtfyUrl        = "https://ntfy.sh"
	emptyValue            = ""
	errorCode             = 1
	passedCode            = 0
)

func init() {
	ipsEnvValue := getEnvValue("IP_ADDRESSES", true, "")
	ipAddresses = strings.Split(ipsEnvValue, ",")
	aezaApiKey = getEnvValue("AEZA_API_KEY", true, "")
	watchDelay = getEnvValue("WATCH_DELAY", false, defaultWatchDelay)
	aezaApiBaseURL = getEnvValue("AEZA_API_BASE_URL", false, defaultAezaApiBaseURL)
	aezaVmBaseURL = getEnvValue("AEZA_VM_BASE_URL", false, defaultAezaVmBaseURL)

	ntfyUrl = getEnvValue("NTFY_URL", false, defaultNtfyUrl)
	ntfyChannel = getEnvValue("NTFY_CHANNEL", false, emptyValue)

	_, githubActionModeEnabled = os.LookupEnv("GITHUB_ACTION_MODE")
}

func main() {
	if githubActionModeEnabled {
		exitCode, message := runDog()
		sendNtfyIfEnabled(exitCode, message)
		exit(exitCode, message)
	} else {
		runCron()
	}
}

func runCron() {
	tickDuration, err := time.ParseDuration(watchDelay)
	if err != nil {
		exit(1, "parse delay. invalid duration: %s", watchDelay)
	}

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	runJob := func() {
		exitCode, message := runDog()
		sendNtfyIfEnabled(exitCode, message)
	}
	runJob()
	for {
		select {
		case <-ticker.C:
			runJob()
		case <-ctx.Done():
			fmt.Println("stop.")
			return
		}
	}
}

func runDog() (int, string) {
	aezaClient := NewAezaAPIClient(0, aezaApiKey, aezaApiBaseURL)
	fmt.Println("get all services...")
	services, err := aezaClient.GetAllServices()
	if err != nil {
		return errorCode, fmt.Sprintf("ERROR: %s", err.Error())
	}

	if len(services.Data.Items) == 0 {
		return passedCode, "services not find"
	}

	serviceID := services.Data.Items[0].Id
	fmt.Println("service found")

	fmt.Println("get vm token...")
	token, err := aezaClient.GetVMToken(serviceID)
	if err != nil {
		return errorCode, fmt.Sprintf("VM token not find: %s", err.Error())
	}
	fmt.Println("token found")

	vmClient := NewAezaVMClient(0, aezaVmBaseURL, token)
	fmt.Println("vm auth...")
	if err := vmClient.Auth(); err != nil {
		return errorCode, fmt.Sprintf("VM auth error: %s", err.Error())
	}
	fmt.Println("logged.")

	fmt.Println("get all instances...")
	allInstances, err := vmClient.GetAllInstances()
	if err != nil {
		return errorCode, fmt.Sprintf("get instances: %s", err.Error())
	}

	if len(allInstances.List) == 0 {
		return passedCode, "not find active instances"
	}
	fmt.Println("instances received")

	fmt.Println("find dead services...")
	var msgBuffer bytes.Buffer
	exitCode := passedCode
	for _, i := range allInstances.List {
		if i.State != "active" {
			for _, ip := range i.Ip4 {
				if isWatchIp(ip.Ip) {
					if err := vmClient.StartInstance(i.Id); err != nil {
						exitCode = errorCode
						msgBuffer.WriteString(fmt.Sprintf("error start instance: %s", err.Error()))
						msgBuffer.WriteString("\n")
					} else {
						msgBuffer.WriteString(fmt.Sprintf("start instance: %s", ip.Ip))
						msgBuffer.WriteString("\n")
					}
				}
			}
		}
	}
	fmt.Println("done.")
	return exitCode, msgBuffer.String()
}

func isWatchIp(currentIp string) bool {
	for _, ip := range ipAddresses {
		if ip == currentIp {
			return true
		}
	}
	return false
}

func getEnvValue(varName string, required bool, defaultValue string) string {
	value, exist := os.LookupEnv(varName)
	if exist {
		if value == "" && required {
			exit(1, "ERROR: env var `%s` is empty", varName)
		}
		if value == "" {
			return defaultValue
		}
		return value
	}
	if required {
		exit(1, "ERROR: env var `%s` is empty", varName)
	}
	return defaultValue
}

func exit(code int, msg string, args ...interface{}) {
	if code == errorCode && msg != "" {
		fmt.Printf(msg+"\n", args...)
	}
	os.Exit(code)
}
