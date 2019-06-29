package parallel

import "sync"

func Execute(tasks []func() error, workersCnt, maxErrorCnt int) {
	taskC := make(chan func() error)
	errC := make(chan struct{})
	closeC := make(chan struct{})
	var waitgroup sync.WaitGroup

	// start workers
	for i := 0; i < workersCnt; i++ {
		go worker(taskC, errC, closeC, &waitgroup)
	}

	errorCnt := 0
out:
	for _, task := range tasks {
		select {
		case <-errC:
			errorCnt++
			if errorCnt == maxErrorCnt {
				break out
			}
		case taskC <- task:
		}
	}
	close(closeC)
	waitgroup.Wait()
}

func worker(task <-chan func() error, errC chan struct{}, close <-chan struct{}, waitgroup *sync.WaitGroup) {
	waitgroup.Add(1)
	for {
		select {
		case task := <-task:
			err := task()
			if err != nil && !closed(close) {
				errC <- struct{}{}
			}
		case <-close:
			waitgroup.Done()
			return
		}
	}

}

func closed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
