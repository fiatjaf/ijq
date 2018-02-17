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
	"regexp"
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
	var globals = []string{"module {iqj: true}"}
	var backlog = []string{". as $v0"}
	var varname = "$v0"
	var num = 0

loop:
	for {
		go func() {
			reader := bufio.NewReader(os.Stdin)

			for {
				fmt.Print("\u001B[33m(./jq)\u001B[m| ")

				filterPrefix := strings.Join(globals, ";") + ";" +
					strings.Join(backlog, "|") + "|" +
					varname + "|"

				bline, _, err := reader.ReadLine()
				if err != nil {
					quit <- true
					return
				}
				line := strings.TrimSpace(string(bline))

				// special commands
				switch line {
				case "debug":
					fmt.Println(filterPrefix)
					continue
				}

				// check for globals (def, include, import)
				var addedglobals = 0
				line = GLOBAL.ReplaceAllStringFunc(line, func(global string) string {
					addedglobals++

					if global[len(global)-1] == ';' {
						globals = append(globals, global[0:len(global)-1])
					} else {
						globals = append(globals, global)
					}
					return ""
				})

				o := new(bytes.Buffer)
				e := new(bytes.Buffer)

				line = strings.TrimSpace(line)
				if len(line) == 0 {
					// if this line came with some global statement,
					// just to be sure, but don't print the output
					if addedglobals > 0 {
						cmd := exec.Command("jq", strings.Join(globals, ";")+";.")
						cmd.Stdin = bytes.NewBufferString("{}")
						cmd.Stdout = o
						cmd.Stderr = e

						err = cmd.Run()
						if err != nil {
							// there is an error in some global statement:
							// log the error,
							// remove them and proceed as if nothing happened.
							log.Print(e.String())
							globals = globals[0 : len(globals)-addedglobals]
							continue
						}
					}

					// otherwise it is just an empty line, proceed while doing nothing.
					continue
				}

				filter := filterPrefix + line

				cmd := exec.Command("jq", "-C", filter)
				cmd.Stdin = getStdinRepeatedly()
				cmd.Env = make([]string, 0)
				cmd.Stdout = o
				cmd.Stderr = e

				err = cmd.Run()
				if err != nil {
					log.Print(e.String())
					globals = globals[0 : len(globals)-addedglobals]
					continue
				}

				res := strings.TrimSuffix(o.String(), "\n")
				if line == "." || line[0] == '$' && len(line) < 5 || res == prev {
					fmt.Println(res)
				} else {
					// assign to a variable
					num++
					varname = "$v" + strconv.Itoa(num)
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

var GLOBAL = regexp.MustCompile(`(def|import|include)[^;]+(;)?`)
