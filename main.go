package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var exitFn = os.Exit

func main() { exitFn(run()) }

func run() int {
	data := dataFromPipe()
	port, err := getFreePort()
	if err != nil {
		panic(err)
	}

	syncGroup := &sync.WaitGroup{}
	syncGroup.Add(1)

	startHttpServer(data, port, syncGroup)

	openBrowser("http://127.0.0.1:" + port)

	syncGroup.Wait()

	return 0
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func startHttpServer(data string, port string, wg *sync.WaitGroup) *http.Server {
	m := http.NewServeMux()
	srv := &http.Server{Addr: ":" + port, Handler: m}

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			return
		}

		_, _ = fmt.Fprintf(w, data)

		go func() {
			time.Sleep(100 * time.Millisecond)
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	})

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return srv
}

func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}

	defer func () {
		_ = l.Close()
	}()

	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}


func dataFromPipe() string {
	out, _ := readUnixPipe()
	return strings.Join(out[:], "\n")
}

func readUnixPipe() ([]string, error) {
	a := make([]string, 0)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := read(reader)
		if err != nil {
			return a, err
		}

		line = strings.TrimSpace(line)
		line = strings.Trim(line, "\n")

		a = append(a, line)
	}
}

func read(r *bufio.Reader) (string, error) {
	line, _ := r.ReadString('\n')

	if line == "" {
		return line, fmt.Errorf("")
	}

	return line, nil
}
