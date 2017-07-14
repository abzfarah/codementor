// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.

package main

import (
	"fmt"
	"os"
)

func CommandPrintln(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func CommandPrint(a ...interface{}) (int, error) {
	return fmt.Print(a...)
}

func CommandPrintErrorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func CommandPrintError(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, a...)
}

func CommandPrettyPrintln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func CommandPrettyPrint(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, a...)
}
