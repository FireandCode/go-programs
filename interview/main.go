package main

type Task struct {
	taskID   string
	userID   string
	priority int
}

type TaskManager struct {
	tasks         map[string]Task
	priorityQueue map[string]Task
}

func (tm *TaskManager) CreateTask(userId string, taskId string, priority int) {

	// create task in tasks
	task := Task{
		userID:   userId,
		taskID:   taskId,
		priority: priority,
	}

	tm.tasks[taskId] = task
	//priority queue
	tm.priorityQueue[taskId] = task
}

/*
create task - user_id, priority and task_id
update priority - priority, task_Id
removal task -> task_id
execute task -> pick highest priority and remove it from tasks
*/

func (tm *TaskManager) UpdatePriority(priority int, taskID string) {

	//tasks
	task := tm.tasks[taskID]
	task.priority = priority
	tm.tasks[taskID] = task

	//priority queue
	tm.priorityQueue[taskID] = task
}

func (tm *TaskManager) RemoveTask(taskID string) {
	//remove tasks
	delete(tm.tasks, taskID)
}

func (tm *TaskManager) ExecuteTask() Task {
	//priority queue , task_id = 112
	taskID := "123"
	isExecuted = false 
	for i := 0; i < count; i++ {
		//check in the tasks and priority as well
		if tm.tasks[taskID] != Task{} {
			
		}
	}
	task := tm.priorityQueue[taskID]


}

func main() {

}

/*

There is a task management system that allows users to manage their tasks, each associated with a priority. The system should efficiently handle adding, modifying, executing, and removing tasks.

We can add any task in the system with some priority and associated userId at any point of time.

We can update the priority of the task at any time.

We can remove any task at any time.

There should be a dummy execution function that will pick the task with the highest priority and execute, it should thenremove that task and return the userId and taskId associated with the task.

create task - user_id, priority and task_id
update priority - priority, task_Id
removal task -> task_id
execute task -> pick highest priority and remove it from tasks

Task {
task_id
priority
user_id
}

TaskManager {
 map[string]Task -> task_id , Task
 priority_queue<int, task>
}


create task {
}
*/