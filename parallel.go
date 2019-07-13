package parallel

import "sync"

func Execute(tasks []func() error, workersCnt, maxErrorCnt int) {
	taskC := make(chan func() error)
	errOut := make(chan struct{})
	finishWorkers := make(chan struct{})

	// start workers
	var waitgroup sync.WaitGroup
	for i := 0; i < workersCnt; i++ {
		go worker(taskC, errOut, finishWorkers, &waitgroup)
	}

	// count errors
	finishErrorCount := make(chan struct{})
	doneError := make(chan struct{})
	go func() {
		errorCnt := 0
		for {
			select {
			case <-errOut:
				errorCnt++
				if errorCnt == maxErrorCnt {
					close(doneError)
					return
				}
			case <-finishErrorCount:
				return
			}
		}
	}()

out:
	for _, task := range tasks {
		select {
		case taskC <- task:
			waitgroup.Add(1)
		case <-doneError:
			break out
		}
	}

	// wait workers
	doneWorkers := make(chan struct{})
	go func() {
		waitgroup.Wait()
		doneWorkers <- struct{}{}
	}()

	for {
		select {
		case <-doneError:
			close(finishWorkers)
			return
		case <-doneWorkers:
			close(finishErrorCount)
			close(finishWorkers)
			return
		}
	}

}

func worker(task <-chan func() error, errC chan struct{}, finish <-chan struct{}, waitgroup *sync.WaitGroup) {
	for {
		select {
		case task := <-task:
			err := task()
			waitgroup.Done()
			if err != nil {
				errC <- struct{}{}
			}
		case <-finish:
			return
		}
	}

}
