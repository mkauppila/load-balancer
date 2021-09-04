package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

type Color string

const (
	Reset  Color = "\x1b[0m"
	Red    Color = "\x1b[31m"
	Green  Color = "\x1b[32m"
	Yellow Color = "\x1b[33m"
	Blue   Color = "\x1b[34m"
)

func GetRandomColor() Color {
	colors := []Color{Red, Green, Yellow, Blue}
	index := rand.Intn(len(colors))
	return colors[index]
}

type Logger struct {
	w   io.Writer
	id  int
	buf *bytes.Buffer
}

func CreateLogger(w io.Writer, id int) *Logger {
	return &Logger{w, id, &bytes.Buffer{}}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	n, err = l.buf.Write(p)
	if err != nil {
		return
	}
	err = l.writeLineByLine()
	return
}

func (l *Logger) formatOutputLine(line string) []byte {
	temp := fmt.Sprintf("%s %d |%s %s", Green, l.id, Reset, line)
	return []byte(temp)
}

func (l *Logger) writeLineByLine() (err error) {
	for {
		line, err := l.buf.ReadString('\n')
		if io.EOF == err {
			break
		}
		if err != nil {
			return err
		}

		if _, err := l.w.Write(l.formatOutputLine(line)); err != nil {
			return err
		}
	}
	return
}

func desiredInstanceCountOrDefault(defaultCount int) int {
	envValue := os.Getenv("COUNT")
	if envValue == "" {
		return defaultCount
	}

	temp, err := strconv.ParseInt(envValue, 10, 32)
	if err != nil {
		return defaultCount
	}
	return int(temp)
}

func launchTestApi(identifier, port int, wg *sync.WaitGroup) {
	cmd := exec.Command("go", "run", "../testApi/api.go")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", port))
	cmd.Stdout = CreateLogger(os.Stdout, identifier)
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintln("Failed to spawn test API: ", err))
	}

	if err := cmd.Wait(); err != nil {
		panic(fmt.Sprintln("Failed to execture test API: ", err))
	}
}

func main() {
	fmt.Println("Set up test APIs")

	count := desiredInstanceCountOrDefault(2)
	defaultPort := 4001
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go launchTestApi(i+1, defaultPort, &wg)
		defaultPort++
	}
	wg.Wait()
}
