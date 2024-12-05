package connectivity

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"net"
	"strings"
)

const pidFlag = "d"

var count = int32(1000)

func getNextId() int32 {
	atomic.AddInt32(&count, 1)
	if count > 8000 {
		atomic.StoreInt32(&count, 1000)
	}
	return count
}

func getPid4Trace() string {
	pid := os.Getpid()
	if pid < 0 {
		pid = 0
	} else if pid > 65535 {
		pid = pid % 60000
	}
	return strconv.FormatInt(int64(pid), 16)
}

func getIp16(ip string) string {
	ips := strings.Split(ip, ".")
	buf := bytes.Buffer{}
	for _, v := range ips {
		d, _ := strconv.Atoi(v)
		hex := strconv.FormatInt(int64(d), 16)
		if len(hex) == 1 {
			buf.WriteString("0" + hex)
		} else {
			buf.WriteString(hex)
		}
	}
	return buf.String()
}

func getLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "127.0.0.1"
    }

    for _, addr := range addrs {
        // 检查地址类型并确保它不是 loopback 地址
        ipNet, ok := addr.(*net.IPNet)
        if ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String()
            }
        }
    }

    return "127.0.0.1"
}

func GenerateTraceId() string {
	b := make([]byte, 0, 32)
	buf := bytes.NewBuffer(b)
	buf.WriteString(getIp16(getLocalIP()))
	buf.WriteString(fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	buf.WriteString(fmt.Sprintf("%d", getNextId()))
	buf.WriteString(pidFlag)
	buf.WriteString(getPid4Trace())
	return buf.String()
}

type EagleEye struct {
	TraceId  string
	RpcId    string
	Index     int
}

const DefaultRpcId = "3"

func (eagleeye *EagleEye) GetTraceId() string {
	return eagleeye.TraceId
}

func (eagleeye *EagleEye) GetRpcId() string {
	eagleeye.Index += 1
	return fmt.Sprintf("%s.%d", eagleeye.RpcId, eagleeye.Index)
}
