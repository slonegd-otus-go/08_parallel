package parallel

import "sync"

func Execute(tasks []func() error, workersCnt, maxErrorCnt int) {
	taskC := make(chan func() error)
	errC := make(chan struct{})
	closeC := make(chan struct{})
	var waitgroup sync.WaitGroup

	// start workers
	waitgroup.Add(1)
	for i := 0; i < workersCnt; i++ {
		go worker(taskC, errC, closeC, &waitgroup)
	}

	go func() {
		errorCnt := 0
		for {
			select {
			case <-errC:
				errorCnt++
				if errorCnt == maxErrorCnt {
					close(closeC)
					return
				}
			case <-closeC:
				return
			}
		}
	}()

	go func() {
		waitgroup.Wait()
		if closed(closeC) {
			return
		}
		close(closeC)
	}()

out:
	for i := 0; i < len(tasks); {
		select {
		case taskC <- tasks[i]:
			waitgroup.Add(1)
			i++
		case <-closeC:
			break out
		}
	}

	waitgroup.Done()
	<-closeC

}

func worker(task <-chan func() error, errC chan struct{}, close <-chan struct{}, waitgroup *sync.WaitGroup) {
	for {
		select {
		case task := <-task:
			err := task()
			waitgroup.Done()
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
