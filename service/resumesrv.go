package service

import "sync"

var ResumeController = &resumeController{
	resumeFlag: true,
	mutex:      &sync.Mutex{},
}

type resumeController struct {
	resumeFlag bool
	mutex      *sync.Mutex
}

func (r *resumeController) EnResumed() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.resumeFlag = true
}

func (r *resumeController) UnResumed() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.resumeFlag = false
}
func (r *resumeController) IsResumed() bool {
	return r.resumeFlag
}
