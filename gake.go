package gake

import (
	"context"
	"fmt"
)

var (
	errTargetNotFound = "target `%s` not found"
)

type recipe func(ctx context.Context) error

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

func Rule(target string) rule {
	return rule{
		target: target,
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

func (r rule) run(ctx context.Context) error {
	if r.phony {
		for _, d := range r.deps {
			err := d.run(ctx)
			if err != nil {
				return err
			}
		}

		return r.fn(ctx)
	}

	// TODO // Non-phony targets may not always run
	return nil
}

func Runner() runner {
	return runner{}
}

func (r *runner) Add(rules ...rule) {
	r.rules = rules
}

func (r *runner) DefaultGoal(rule rule) {
	r.defaultGoal = rule
}

func (r runner) Run(args []string) {
	var rule rule
	if len(args) == 1 {
		if r.defaultGoal.target == "" {
			r.rules[0].run(context.Background())
		} else {
			r.defaultGoal.run(context.Background())
		}

		return
	}

	rule, err := r.findRule(args[1])
	if err != nil {
		panic(err)
	}

	rule.run(context.Background())
}

func (r runner) findRule(target string) (rule, error) {
	for _, t := range r.rules {
		if t.target == target {
			return t, nil
		}
	}

	return rule{}, fmt.Errorf(errTargetNotFound, target)
}