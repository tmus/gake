package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmus/gake"
)

func main() {
	r := gake.Runner()
	t1 := gake.Rule("build_world").Recipe(func(ctx context.Context) error {
		_, cancel := context.WithCancel(ctx)
		fmt.Println("Building world")
		cancel()
		return nil
	}).Phony(true)
	t2 := gake.Rule("hello_world").Recipe(func(ctx context.Context) error {
		fmt.Println("Hello world")
		return nil
	}).Dependencies(t1).Phony(true)

	r.DefaultGoal(t2)

	r.Add(t1, t2)

	r.Run(os.Args)
}
