package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
)

type Config struct {
	Spawn int `short:"s" long:"spawn" description:"How many children to spawn [sum_(j=0)^a product_(i=0)^(j - 1)(a - i) where a = spawn)]" default:"2"`
}

func main() {
	var conf Config
	if _, err := flags.Parse(&conf); err != nil {
		os.Exit(1)
	}

	binPath, _ := filepath.Abs(os.Args[0])
	for i := 0; i < conf.Spawn; i++ {
		c := exec.Command(binPath, "--spawn", strconv.Itoa(conf.Spawn-1))

		fmt.Println("Spawning proccess", i)
		c.Start()
	}

	fmt.Println("Spawner", os.Getpid(), "startet, waiting...")
	for {
		time.Sleep(10 * time.Second)
	}
}
