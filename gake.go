package gake

import (
	"context"
	"errors"
	"fmt"
)

var (
	errTargetNotFound = "target `%s` not found"
	errEmptyRunner    = errors.New("runner is empty")
)

type recipe func(ctx context.Context) (context.Context, error)

type rule struct {
	target string
	fn     recipe
	deps   []rule
	phony  bool
}

type runner struct {
	rules       []rule
	defaultGoal rule
}

// Rule is a single task that exists on a Runner.
func Rule(target string) rule {
	return rule{
		target: target,
		// TODO
		phony: true,
	}
}

func (r rule) Phony(val bool) rule {
	r.phony = val
	return r
}

func (r rule) Recipe(fn recipe) rule {
	r.fn = fn
	return r
}

func (r rule) Dependencies(rules ...rule) rule {
	r.deps = rules
	return r
}

func (r rule) run(ctx context.Context) (context.Context, error) {
	var err error
	if r.phony {
		for _, d := range r.deps {
			ctx, err = d.run(ctx)
			if err != nil {
				return ctx, err
			}
		}

		return r.fn(ctx)
	}

	// TODO // Non-phony targets may not always run
	return ctx, nil
}

// Runner contains multiple rules and orchestrates tasks being ran.
func Runner() runner {
	return runner{}
}

func (r *runner) Add(rules ...rule) {
	r.rules = rules
}

func (r *runner) DefaultGoal(rule rule) {
	r.defaultGoal = rule
}

func (r runner) Run(args []string) (context.Context, error) {
	if len(r.rules) == 0 {
		return context.Background(), errEmptyRunner
	}
	if len(args) == 1 {
		if r.defaultGoal.target == "" {
			return r.rules[0].run(context.Background())
		}

		return r.defaultGoal.run(context.Background())
	}

	rule, err := r.findRule(args[1])
	if err != nil {
		return context.Background(), err
	}

	return rule.run(context.Background())
}

func (r runner) findRule(target string) (rule, error) {
	for _, t := range r.rules {
		if t.target == target {
			return t, nil
		}
	}

	return rule{}, fmt.Errorf(errTargetNotFound, target)
}
