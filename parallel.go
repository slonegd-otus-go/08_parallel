package parallel

func Execute(tasks []func() error, workersCnt, maxErrorCnt int) {
	taskC := make(chan func() error)
	errC := make(chan struct{})
	closeC := make(chan struct{})
	defer func() {
		close(closeC)
	}()

	// start workers
	for i := 0; i < workersCnt; i++ {
		go worker(taskC, errC, closeC)
	}

	errorCnt := 0
	for _, task := range tasks {
		select {
		case <-errC:
			errorCnt++
			if errorCnt == maxErrorCnt {
				return
			}
		case taskC <- task:
		}
	}
}

func worker(task <-chan func() error, errC chan struct{}, close <-chan struct{}) {
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
