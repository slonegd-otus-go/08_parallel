package parallel_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	parallel "github.com/slonegd-otus-go/08_parallel"
)

var log string
var logmtx sync.Mutex

func makeFoo(id int, duration time.Duration, isError bool) func() error {
	return func() error {
		time.Sleep(duration)
		logmtx.Lock()
		defer logmtx.Unlock()
		if isError {
			log = fmt.Sprintf("%verror %v\n", log, id)
			return errors.New("")
		}
		log = fmt.Sprintf("%vexecute %v\n", log, id)
		return nil
	}
}

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
		wantLog     string
	}{
		{
			name: "1 goroutine 1 error",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false),
				makeFoo(2, 10*time.Millisecond, true),
				makeFoo(3, 1*time.Millisecond, false),
			},
			workersCnt:  1,
			maxErrorCnt: 1,
			wantLog: `execute 1
error 2
`,
		},
		{
			name: "1 goroutine 2 error",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false),
				makeFoo(2, 20*time.Millisecond, true),
				makeFoo(3, 30*time.Millisecond, false),
				makeFoo(4, 40*time.Millisecond, true),
				makeFoo(5, 50*time.Millisecond, false),
			},
			workersCnt:  1,
			maxErrorCnt: 2,
			wantLog: `execute 1
error 2
execute 3
error 4
`,
		},
		{
			name: "2 goroutine 1 error",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false),
				makeFoo(2, 30*time.Millisecond, true),
				makeFoo(3, 10*time.Millisecond, false),
				makeFoo(4, 50*time.Millisecond, false),
				makeFoo(5, 50*time.Millisecond, false),
			},
			workersCnt:  2,
			maxErrorCnt: 1,
			wantLog: `execute 1
execute 3
error 2
`,
		},
		{
			name: "2 goroutine 2 error",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false),
				makeFoo(2, 30*time.Millisecond, true),
				makeFoo(3, 10*time.Millisecond, false),
				makeFoo(4, 20*time.Millisecond, true),
				makeFoo(5, 50*time.Millisecond, false),
				makeFoo(6, 50*time.Millisecond, false),
			},
			workersCnt:  2,
			maxErrorCnt: 2,
			wantLog: `execute 1
execute 3
error 2
error 4
`,
		},
		{
			name: "2 goroutine 2 error out",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false), // 1 до 10
				makeFoo(2, 30*time.Millisecond, true),  // 2 до 30
				makeFoo(3, 10*time.Millisecond, false), // 1 до 20
				makeFoo(4, 30*time.Millisecond, true),  // 1 до 50 - больше не работает
				makeFoo(5, 10*time.Millisecond, false), // 2 до 40
				makeFoo(6, 30*time.Millisecond, false), // 2 до 70
			},
			workersCnt:  2,
			maxErrorCnt: 2,
			wantLog: `execute 1
execute 3
error 2
execute 5
error 4
`,
		},
		{
			name: "2 goroutine execute all",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, false), // 1 до 10
				makeFoo(2, 30*time.Millisecond, true),  // 2 до 30
				makeFoo(3, 10*time.Millisecond, false), // 1 до 20
				makeFoo(4, 30*time.Millisecond, false), // 1 до 50
				makeFoo(5, 10*time.Millisecond, false), // 2 до 40
				makeFoo(6, 30*time.Millisecond, false), // 2 до 70
			},
			workersCnt:  2,
			maxErrorCnt: 2,
			wantLog: `execute 1
execute 3
error 2
execute 5
execute 4
execute 6
`,
		},
		{
			name: "2 goroutine all error",
			tasks: []func() error{
				makeFoo(1, 10*time.Millisecond, true), // 1 до 10
				makeFoo(2, 30*time.Millisecond, true), // 2 до 30
				makeFoo(3, 10*time.Millisecond, true), // 1 до 20
				makeFoo(4, 30*time.Millisecond, true), // 1 до 50
				makeFoo(5, 10*time.Millisecond, true), // 2 до 40
				makeFoo(6, 30*time.Millisecond, true), // 2 до 70
			},
			workersCnt:  2,
			maxErrorCnt: 2,
			wantLog: `error 1
error 3
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logmtx.Lock()
			log = ""
			logmtx.Unlock()
			parallel.Execute(tt.tasks, tt.workersCnt, tt.maxErrorCnt)
			logmtx.Lock()
			assert.Equal(t, tt.wantLog, log)
			logmtx.Unlock()
			time.Sleep(100 * time.Millisecond)
		})
	}
}
