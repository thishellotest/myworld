package tests

import (
	"context"
	"sync"
	"testing"
	"time"
	"vbc/lib"
	. "vbc/lib/builder"
)

func Test_ZoomRecordingFileJobUsecase_HandleExec(t *testing.T) {

	//
	//[82.167ms] [rows:1] UPDATE `zoom_recording_files` SET `play_url`='https://us06web.zoom.us/rec/play/kQ7QA6pZcR3uYmthU5lQJz-RmwHZgqenb09QYQjsvgi8fCWhGXHf3GlXzsSrhuqlpfLdERDOuv_1tkpS.ADzmKaPeDunpco8P',`download_url`='https://us06web.zoom.us/rec/download/kQ7QA6pZcR3uYmthU5lQJz-RmwHZgqenb09QYQjsvgi8fCWhGXHf3GlXzsSrhuqlpfLdERDOuv_1tkpS.ADzmKaPeDunpco8P'
	entity, _ := UT.ZoomRecordingFileUsecase.GetByCond(Eq{"id": 3})
	err := UT.ZoomRecordingFileJobUsecase.HandleExec(context.TODO(), entity)
	lib.DPrintln(err)
}

func Test_ZoomRecordingFileJobUsecase_WaitingTasks(t *testing.T) {
	a, err := UT.ZoomRecordingFileJobUsecase.WaitingTasks(context.TODO())
	lib.DPrintln(err)
	a1, a2, err := lib.SqlRowsTrans(a)
	lib.DPrintln(a1, a2, err)
}

func Test_ZoomRecordingFileJobUsecase_RunHandleCustomJob(t *testing.T) {
	var w sync.WaitGroup
	w.Add(1)
	err := UT.ZoomRecordingFileJobUsecase.RunHandleCustomJob(context.TODO(),
		1,
		time.Second*10,
		UT.ZoomRecordingFileJobUsecase.WaitingTasks,
		UT.ZoomRecordingFileJobUsecase.Handle)
	if err != nil {
		panic(err)
	}
	w.Wait()
}

func Test_ZoomRecordingFileJobUsecase_ExecuteCrontabHandleProcessingRecording(t *testing.T) {
	err := UT.ZoomRecordingFileJobUsecase.ExecuteCrontabHandleProcessingRecording()
	lib.DPrintln(err)
}
