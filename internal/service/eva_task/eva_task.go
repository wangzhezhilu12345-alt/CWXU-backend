package eva_task

import (
	context "context"
	eva_task2 "edu-evaluation-backed/api/v1/eva_task"
	"edu-evaluation-backed/internal/biz/eva_task"
	"edu-evaluation-backed/internal/data/dal"
	"strconv"
)

type EvaTaskService struct {
	taskDal *dal.TaskDal
	taskUC  *eva_task.EvaTaskUseCase
}

func (e EvaTaskService) CreateTask(ctx context.Context, req *eva_task2.CreateTaskReq) (*eva_task2.CreateTaskResp, error) {
	id, err := e.taskUC.CreateEvaTask(req.Name, req.CourseIds)
	if err != nil {
		return nil, err
	}
	resp := &eva_task2.CreateTaskResp{
		Id:      strconv.Itoa(int(id)),
		Message: "创建成功",
	}
	return resp, nil
}

func (e EvaTaskService) Detail(ctx context.Context, req *eva_task2.GetTaskReq) (*eva_task2.TaskInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (e EvaTaskService) List(ctx context.Context, req *eva_task2.GetTaskReq) (*eva_task2.GetTaskListResp, error) {
	//TODO implement me
	panic("implement me")
}

func NewEvaTaskService(taskDal *dal.TaskDal, taskUC *eva_task.EvaTaskUseCase) *EvaTaskService {
	return &EvaTaskService{
		taskDal: taskDal,
		taskUC:  taskUC,
	}
}
