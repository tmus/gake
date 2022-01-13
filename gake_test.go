package gake

import (
	"context"
	"testing"
)

func TestEmptyRunnerAborts(t *testing.T) {
	r := Runner()

	_, err := r.Run([]string{"filename"})
	if err == nil {
		t.Fatalf("expected err to have value")
	}
}

func TestRunnerUsesFirstRuleWhenNoArgsPassed(t *testing.T) {
	var called bool
	var notCalled bool
	rule1 := Rule("first").Recipe(func(ctx context.Context) (context.Context, error) {
		called = true
		return ctx, nil
	})

	rule2 := Rule("second").Recipe(func(ctx context.Context) (context.Context, error) {
		notCalled = true
		return ctx, nil
	})
	r := Runner()
	r.Add(rule1, rule2)

	r.Run([]string{"filename"})
	if !called || notCalled {
		t.Fatalf("expected called to be true and notCalled to be false, got: %v and %v", called, notCalled)
	}
}

func TestRunnerUsesDefaultGoal(t *testing.T) {
	var called bool
	var notCalled bool
	rule1 := Rule("first").Recipe(func(ctx context.Context) (context.Context, error) {
		notCalled = true
		return ctx, nil
	})

	rule2 := Rule("second").Recipe(func(ctx context.Context) (context.Context, error) {
		called = true
		return ctx, nil
	})
	r := Runner()
	r.Add(rule1, rule2)
	r.DefaultGoal(rule2)

	r.Run([]string{"filename"})
	if !called || notCalled {
		t.Fatalf("expected called to be true and notCalled to be false, got: %v and %v", called, notCalled)
	}
}

func TestReturnsErrorIfRuleNotFound(t *testing.T) {
	rule1 := Rule("first").Recipe(func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	})

	r := Runner()
	r.Add(rule1)
	_, err := r.Run([]string{"filename", "invalid_rule"})

	if err == nil || err.Error() != "target `invalid_rule` not found" {
		t.Fatalf("expected error to be returned")
	}
}
