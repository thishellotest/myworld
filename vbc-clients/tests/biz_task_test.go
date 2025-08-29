package tests

import (
	"fmt"
	"testing"
	"vbc/internal/biz"
	"vbc/lib"
)

func Test_TaskUsecasek(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Take(&entity)
	a := UT.TaskUsecase.ExecTask(&entity)
	lib.DPrintln(a)
}

func Test_TaskUsecase_Task(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=9").Take(&entity)
	//err := UT.TaskUsecase.CreateEnvelope(&entity)
	//lib.DPrintln(err)
}

func Test_TaskUsecase_ExecTask(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=12905").Take(&entity)
	err := UT.TaskUsecase.ExecTask(&entity)
	lib.DPrintln(err)
}

func Test_TaskUsecase_TestInvoke(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=61").Take(&entity)
	err := UT.TaskUsecase.Invoke(&entity)
	fmt.Println(err)
}

func Test_TaskUsecase_CreateTask(t *testing.T) {
	a := make(lib.TypeMap)
	a.Set("envelope_id", "641eafbd-1d03-409b-b092-37219af0ae41")
	err := UT.TaskCreateUsecase.CreateTask(0, a, biz.Task_Dag_GetEnvelopeDocuments, 0, "", "")
	fmt.Println(err)
}

func Test_TaskUsecase_Invoke_ContractNonResponsive(t *testing.T) {
	var entity biz.TaskEntity
	UT.CommonUsecase.DB().Where("id=13").Take(&entity)
	err := UT.TaskUsecase.Invoke(&entity)
	fmt.Println(err)
}

func Test_TaskUsecase_WaitingTasks(t *testing.T) {

	sql, _ := UT.TaskUsecase.WaitingTasks()
	a, b, c := lib.SqlRowsTrans(sql)
	lib.DPrintln(a, b, c)
}
