package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	quit := make(chan bool)

	stat, _ := os.Stdin.Stat()
	var getStdinRepeatedly func() io.Reader
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		getStdinRepeatedly = func() io.Reader {
			return bytes.NewBuffer(stdin)
		}
	} else {
		getStdinRepeatedly = func() io.Reader {
			return bytes.NewBufferString("{}")
		}
	}

	var prev = ""
	var backlog = []string{". as $l0"}
	var varname = "$l0"
	var num = 0

loop:
	for {
		go func() {
			reader := bufio.NewReader(os.Stdin)

			for {
				fmt.Print("\u001B[36m(jq)\u001B[m| ")

				bline, _, err := reader.ReadLine()
				if err != nil {
					quit <- true
					return
				}
				line := string(bytes.TrimSpace(bline))
				if line[len(line)-1] == ';' {
					line += "."
				}

				o := new(bytes.Buffer)
				e := new(bytes.Buffer)

				filter := strings.Join(backlog, "|") + "|" + varname + "|" + line

				cmd := exec.Command("jq", "-C", filter)
				cmd.Stdin = getStdinRepeatedly()
				cmd.Env = make([]string, 0)
				cmd.Stdout = o
				cmd.Stderr = e

				err = cmd.Run()
				if err != nil {
					log.Print(e.String())
					continue
				}

				res := strings.TrimSuffix(o.String(), "\n")
				if line == "." || line[0] == '$' && len(line) < 5 || res == prev {
					fmt.Println(res)
				} else {
					// assign to a variable
					num++
					varname = "$l" + strconv.Itoa(num)
					fmt.Println(res + "\u001B[90m as " + varname + "\u001B[m")
					backlog = append(backlog, "("+line+") as "+varname)
				}

				prev = res
			}
		}()

		select {
		case <-quit:
			fmt.Println("[ijq] terminated")
			break loop
		}
	}

	return
}
