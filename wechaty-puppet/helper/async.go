package helper

import (
	"math"
	"runtime"
	"sync"
)

// DefaultWorkerNum default number of goroutines is twice the number of GOMAXPROCS
var DefaultWorkerNum = runtime.GOMAXPROCS(0) * 2

type (
	// IAsync interface
	IAsync interface {
		AddTask(task Task)
		Result() []AsyncResult
	}

	// AsyncResult result struct
	AsyncResult struct {
		Value interface{}
		Err   error
	}

	async struct {
		tasks        []Task
		wg           sync.WaitGroup
		maxWorkerNum int
	}

	// Task task func
	Task func() (interface{}, error)
)

// NewAsync ...
func NewAsync(maxWorkerNum int) IAsync {
	if maxWorkerNum <= 0 {
		maxWorkerNum = DefaultWorkerNum
	}
	return &async{
		maxWorkerNum: maxWorkerNum,
		wg:           sync.WaitGroup{},
	}
}

func (a *async) AddTask(task Task) {
	a.tasks = append(a.tasks, task)
}

func (a *async) Result() []AsyncResult {
	taskChan := make(chan Task)
	resultChan := make(chan AsyncResult)

	taskNum := len(a.tasks)
	workerNum := int(math.Min(float64(taskNum), float64(a.maxWorkerNum)))
	a.wg.Add(taskNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range taskChan {
				result := AsyncResult{}
				result.Value, result.Err = task()
				resultChan <- result
				a.wg.Done()
			}
		}()
	}

	go func() {
		for _, v := range a.tasks {
			taskChan <- v
		}
		a.wg.Wait()
		close(resultChan)
		close(taskChan)
		a.tasks = nil
	}()

	result := make([]AsyncResult, 0, taskNum)
	for v := range resultChan {
		result = append(result, v)
	}
	return result
}
