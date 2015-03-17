package cludo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/oxtopus/cludo/systemd"
	"github.com/oxtopus/cludo/unit"
)

// Build Unit file from CLI params
func Run(wd string, cliArgs []string) {
	unit := cludo.MakeUnit()

	userCommand := strings.Join(cliArgs, " ")
	command := fmt.Sprintf("/usr/bin/docker run cludo-base /bin/bash -c \"%s\"", userCommand)

	unit.AddSection("Unit")
	unit.AddItem("Unit", "Description", userCommand)
	unit.AddSection("Service")
	unit.AddItem("Service", "ExecStart", command)

	d, err := ioutil.TempDir(wd, ".cludo-")
	if err != nil {
		panic(err)
	}
	//defer os.RemoveAll(d)

	fname := systemd.Escape(userCommand) + ".service"

	outp, err := os.Create(path.Join(d, fname))
	if err != nil {
		panic(err)
	}

	unit.Export(outp)

	defer os.Chdir(wd)
	os.Chdir(d)

	// Load unit file
	cmd := exec.Command("fleetctl", "load", fname)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())

	// Start unit
	cmd = exec.Command("fleetctl", "start", fname)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())
}
