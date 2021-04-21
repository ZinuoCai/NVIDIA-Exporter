package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type processInfo struct {
	gpuId     string
	processId string
	smUtil    int
	memUtil   int
	command   string
}

// 指标结构体
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(metricName, docString, labels, nil)
}

/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics() *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"sm_util":  newGlobalMetric("nvidia_sm_util", "sm_util", []string{"gpu_id", "process_id", "command"}),
			"mem_util": newGlobalMetric("nvidia_mem_util", "mem_util", []string{"gpu_id", "process_id", "command"}),
		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()

	processInfo := c.GetProcessInfo()
	for _, info := range processInfo {
		ch <- prometheus.MustNewConstMetric(c.metrics["sm_util"], prometheus.GaugeValue, float64(info.smUtil), info.gpuId, info.processId, info.command)
		ch <- prometheus.MustNewConstMetric(c.metrics["mem_util"], prometheus.GaugeValue, float64(info.memUtil), info.gpuId, info.processId, info.command)
	}
}

/**
 * 函数：GetProcessInfo
 * 功能：生成数据
 */
func (c *Metrics) GetProcessInfo() []processInfo {
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
			gpuId := strarray[0]
			processId := strarray[1]
			smUtil, _ := strconv.Atoi(strarray[3])
			memUtil, _ := strconv.Atoi(strarray[4])
			command := strarray[7]

			info = append(info, processInfo{gpuId: gpuId, processId: processId, smUtil: smUtil, memUtil: memUtil, command: command})
		}
	}
	return info
}
