package parallel

import (
	"errors"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	type args struct {
		tasks       []func() error
		workersCnt  int
		maxErrorCnt int
	}
	tests := []struct {
		name        string
		tasks       []func() error
		workersCnt  int
		maxErrorCnt int
	}{
		{
			name: "test",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false),
				makeFoo(2, 20*time.Millisecond, true),
				makeFoo(3, 30*time.Millisecond, false),
				makeFoo(4, 40*time.Millisecond, true),
				makeFoo(5, 50*time.Millisecond, false),
				makeFoo(6, 60*time.Millisecond, true),
				makeFoo(7, 70*time.Millisecond, false),
				makeFoo(8, 80*time.Millisecond, true),
				makeFoo(9, 90*time.Millisecond, false),
				makeFoo(10, 100*time.Millisecond, true),
				makeFoo(11, 110*time.Millisecond, false),
			},
			workersCnt:  1,
			maxErrorCnt: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute(tt.tasks, tt.workersCnt, tt.maxErrorCnt)
		})
	}
}

func makeFoo(id int, duration time.Duration, isError bool) func() error {
	return func() error {
		time.Sleep(duration)
		if isError {
			println("error", id)
			return errors.New("")
		}
		println("execute", id)
		return nil
	}
}
