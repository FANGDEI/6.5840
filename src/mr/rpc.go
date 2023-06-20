package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
	"time"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
const (
	// Task types
	MapTask = iota
	ReduceTask
	Done
)

type (
	Task struct {
		TaskType int
		FileName string
		// TODO: modify NMap to WorkerId
		NMap     int
		NReduce  int
		Deadline time.Time
	}

	ApplyTaskArgs struct {
		WorkerId int
	}

	ApplyTaskReply struct {
		Task Task
	}

	ReportArgs struct {
		TaskType int
	}

	ReportReply struct {
		Ok bool
	}
)

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
