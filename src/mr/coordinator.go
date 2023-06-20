package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Coordinator struct {
	// Your definitions here.
	mutex sync.Mutex
	// MapReduce job status [map, reduce, done]
	currentTaskType int
	// Number of map tasks
	nMap int
	// Number of reduce tasks
	nReduce int
	// Maptask queue
	mapTasks []Task
	// Reducetask queue
	reduceTasks []Task
}

// Your code here -- RPC handlers for the worker to call.
// TODO: Worker call RPC to get a task
func (c *Coordinator) ApplyTask(args *ApplyTaskArgs, reply *ApplyTaskReply) error {
	return nil
}

// TODO: Worker call RPC to report task completion
func (c *Coordinator) Report(args *ReportArgs, reply *ReportReply) error {
	return nil
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	if c.currentTaskType == Done {
		ret = true
	}

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	c := Coordinator{
		currentTaskType: MapTask,
		nMap:            len(files),
		nReduce:         nReduce,
		mapTasks:        make([]Task, max(len(files), nReduce)),
		reduceTasks:     make([]Task, max(len(files), nReduce)),
	}

	// Your code here.
	for idx, fileName := range files {
		task := Task{
			TaskType: MapTask,
			FileName: fileName,
			NMap:     -1,
			NReduce:  idx,
		}
		c.mapTasks[idx] = task
	}

	go func() {
		for {
			time.Sleep(time.Microsecond * 500)
			c.mutex.Lock()
			defer c.mutex.Unlock()
			switch c.currentTaskType {
			case MapTask:
				for _, task := range c.mapTasks {
					if task.NMap != -1 && time.Now().After(task.Deadline) {
						log.Printf("Map task %v timeout\n", task)
						task.NMap = -1
						// TODO: modify task queue to channel
						c.mapTasks = append(c.mapTasks, task)
					}
				}
			}
		}
	}()

	c.server()
	return &c
}
