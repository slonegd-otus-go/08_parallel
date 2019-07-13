package parallel

import "sync"

func Execute(tasks []func() error, workersCnt, maxErrorCnt int) {
	taskC := make(chan func() error)
	errC := make(chan struct{})
	closeC := make(chan struct{})
	var waitgroup sync.WaitGroup

	// start workers
	for i := 0; i < workersCnt; i++ {
		waitgroup.Add(1)
		go worker(taskC, errC, closeC, &waitgroup)
	}

	go func() {
		errorCnt := 0
		for {
			select {
			case <-errC:
				errorCnt++
				if errorCnt == maxErrorCnt {
					if !closed(closeC) {
						close(closeC)
					}
					return
				}
			case <-closeC:
				return
			}
		}
	}()

out:
	for _, task := range tasks {
		select {
		case taskC <- task:

		case <-closeC:
			break out
		}
	}

	if !closed(closeC) {
		close(closeC)
	}
	waitgroup.Wait()
}

func worker(task <-chan func() error, errC chan struct{}, close <-chan struct{}, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	for {
		select {
		case task := <-task:
			err := task()
			if err != nil && !closed(close) {
				errC <- struct{}{}
			}
		case <-close:
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
