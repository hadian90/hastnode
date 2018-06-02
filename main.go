package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(location.Default())

	r.POST("/hook/:directory", hook)

	r.Run(":5000")
}

func hook(c *gin.Context) {

	// check if url has hook command
	directory := c.Param("directory")

	// change all | to /
	directory = strings.Replace(directory, "|", "/", -1)
	// add / infront
	directory = "/" + directory

	// go to directory and find hastehook file
	fileHandle, _ := os.Open(directory + "/deployscript")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		// run all command inside the hastehook file inside the directory
		lsCmd := exec.Command("ls")
		cmd := strings.Split(fileScanner.Text(), " ")
		if len(cmd) < 2 {
			lsCmd = exec.Command(cmd[0])
		} else {
			lsCmd = exec.Command(cmd[0], cmd[1:]...)
		}
		lsCmd.Dir = directory
		out, err := lsCmd.Output()

		check(err)

		fmt.Printf(string(out))

	}
}
