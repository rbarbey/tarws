package main

import "github.com/rbarbey/tarws/cmd"

func main() {
	cmd.Execute()

	//err := filepath.Walk("/Users/robert/Development/golang/src/", tar)
	// tarWriter := tar.NewWriter(ioutil.Discard)
	// defer tarWriter.Close()
	//
	// iterate("/Users/robert/Development/golang/src/", tarWriter)
}
