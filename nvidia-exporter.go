package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type processInfo struct {
	gpuId     int
	processId int
	smUtil    int
	memUtil   int
	command   string
}

func getProcessInfo() []processInfo {
	command := exec.Command("nvidia-smi", "pmon", "-c", "1")
	out, err := command.Output()

	if err != nil {
		log.Fatal("Error: ", err)
	}

	s := string(out)
	arr := strings.Split(s, "\n")
	var info []processInfo

	for _, line := range arr {
		strarray := strings.Fields(strings.TrimSpace(line))

		// omit the title line
		// omit idle gpu cards
		// omit the process do not use gpu currently
		if len(strarray) == 8 && strarray[1] != "-" && strarray[3] != "-" && strarray[4] != "-" {
			// gpu_id, process, sm, mem, command
			// fmt.Println(strarray[0], strarray[1], strarray[3], strarray[4], strarray[7])
			gpuId, _ := strconv.Atoi(strarray[0])
			processId, _ := strconv.Atoi(strarray[1])
			smUtil, _ := strconv.Atoi(strarray[3])
			memUtil, _ := strconv.Atoi(strarray[4])
			command := strarray[7]

			info = append(info, processInfo{gpuId: gpuId, processId: processId, smUtil: smUtil, memUtil: memUtil, command: command})
		}
	}
	return info
}

func metrics(response http.ResponseWriter, request *http.Request) {
	processInfos := getProcessInfo()
	result := ""

	for _, info := range processInfos {
		result = fmt.Sprintf("%s%s{gpu=\"%d\", prosess=\"%d\", command=\"%s\"} %d\n",
			result, "sm_util", info.gpuId, info.processId, info.command, info.smUtil)
		result = fmt.Sprintf("%s%s{gpu=\"%d\", prosess=\"%d\", command=\"%s\"} %d\n",
			result, "mem_util", info.gpuId, info.processId, info.command, info.memUtil)
	}
	_, _ = fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}

var (
	listenAddress string
	metricsPath   string
)

func init() {
	flag.StringVar(&listenAddress, "web.listen-address", ":9114", "Address to listen on")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics/", "Path under which to expose metrics.")
	flag.Parse()
}

func main() {
	http.HandleFunc(metricsPath, metrics)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
