package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmus/gake"
)

var contextKey struct{}

func main() {
	r := gake.Runner()

	t1 := gake.Rule("build_world").
		Predicate(func(ctx context.Context) (context.Context, bool) {
			return ctx, false
		}).
		Recipe(func(ctx context.Context) (context.Context, error) {
			ctx = context.WithValue(ctx, contextKey, "asd")
			fmt.Println("Building world")
			return ctx, nil
		})

	t2 := gake.Rule("hello_world").
		Dependencies(t1).
		Predicate(func(ctx context.Context) (context.Context, bool) {
			ctx = context.WithValue(ctx, contextKey, "modifying context from predicate...")
			return ctx, true
		}).
		Recipe(func(ctx context.Context) (context.Context, error) {
			fmt.Println(ctx.Value(contextKey))
			fmt.Println("Hello world")
			return ctx, nil
		})

	r.DefaultGoal(t2)

	r.Add(t1, t2)

	r.Run(os.Args)
}
