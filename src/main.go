package main

import (
	p "RaGodahn/pslistwin"
	"fmt"
	"os"
	"slices"
)

func getProcess() int {
	processList, err := p.Processes()
	var list []int
	var i int
	if err != nil {
		fmt.Println("cannot read processes")
		os.Exit(3)
	}

	for x := range processList {
		var process p.Process = processList[x]
		fmt.Printf("%d\t%s\n", process.Pid(), process.Executable())
		list = append(list, process.Pid())
	}

	fmt.Printf("Please input the PID of the process you wish to inject:")
	fmt.Scanln(&i)

	if !slices.Contains(list, i) {
		fmt.Println("PID inputed does not match any running process")
		os.Exit(3)
	}
	return i

}

func main() {
	getProcess()
}
