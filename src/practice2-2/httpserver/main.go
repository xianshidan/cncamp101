package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func sayhello(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析参数，默认是不会解析的

	// 2.设置version
	os.Setenv("VERSION", "v0.0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	fmt.Printf("os version: %s \n", version)

	// 1.将requst中的header 设置到 reponse中
	for k, vArry := range r.Header {
		for _, v := range vArry {
			fmt.Printf("Header key: %s, Header value: %s \n", k, v)
			w.Header().Set(k, v)
		}
	}

	// 3.记录日志并输出
	clientip := ClientIP(r)
	sugarLogger.Infof("Success! Response code: %d", 200)
	sugarLogger.Infof("Success! clientip: %s", clientip)

	fmt.Fprintln(os.Stdout, "Response code: 200")
	fmt.Fprintln(os.Stdout, "clientip: ", clientip)
}

// 4.健康检查的路由
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "i am alive")
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()

	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/", sayhello)

	err := http.ListenAndServe(":8290", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
