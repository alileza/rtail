package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	os.Exit(Main())
}

var (
	colorList = []func(string, ...interface{}) string{
		color.BlueString,
		color.CyanString,
		color.GreenString,
		color.MagentaString,
		color.YellowString,
	}
	colorCounter int
)

type tconfig struct {
	Name    string   `json:"name"`
	Servers []string `json:"server_addresses"`
	File    string   `json:"file"`
	Options []string `json:"options"`
}

var (
	configFile = flag.String("config.file", "", "")

	configs []tconfig
)

func Main() int {

	if configFile := *configFile; configFile != "" {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			println(err.Error())
			return 1
		}
		err = json.Unmarshal(b, &configs)
		if err != nil {
			println(err.Error())
		}
		return MainConfig()
	}

	if len(os.Args) < 2 {
		fmt.Println("usage: rtail [<servers>] [-F | -f | -r] [-q] [-b # | -c # | -n #] <file>")
		return 0
	}
	servers := os.Args[1]

	var grCount int
	errChan := make(chan error)

	for _, server := range strings.Split(servers, ",") {
		go run("", server, flag.Args()[1:], errChan)
		grCount++
	}
	for grCount != 0 {
		err := <-errChan
		if err != nil {
			println(err.Error())
		}
		grCount--
	}

	return 0
}

func MainConfig() int {
	var grCount int
	errChan := make(chan error)
	for _, cfg := range configs {
		for _, server := range cfg.Servers {
			go run(cfg.Name, server, append(cfg.Options, cfg.File), errChan)
			grCount++
		}

	}
	for grCount != 0 {
		err := <-errChan
		if err != nil {
			println(err.Error())
		}
		grCount--
	}

	return 0
}

func run(name, server string, command []string, err chan error) {
	cmds := []string{"ssh", server, "tail", strings.Join(command, " ")}
	fmt.Println("Executing : ", cmds)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmd.Stdout = &writer{prefix: randColor(fmt.Sprintf("[%s:%s] ", name, server))}
	cmd.Stderr = &writer{prefix: errColor(fmt.Sprintf("[%s:%s] ERR : ", name, server))}

	if errno := cmd.Run(); errno != nil {
		err <- fmt.Errorf("[%s:%s] %s", name, server, errno.Error())
		return
	}

	err <- fmt.Errorf("[%s:%s] %s", name, server, "session closed")
}

func randColor(s string) string {
	colorCounter++
	if colorCounter == len(colorList) {
		colorCounter = 0
	}
	return colorList[colorCounter](s)
}
func errColor(s string) string {
	return color.RedString(s)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type writer struct {
	prefix string
}

func (c *writer) Write(b []byte) (int, error) {
	fmt.Print(c.prefix + string(b))
	return len(b), nil
}
