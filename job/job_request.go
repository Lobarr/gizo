package job

//JobRequest requests job exec
type JobRequest struct {
	ID    string
	Execs []*Exec
}

func NewJobRequest(id string, exec ...*Exec) *JobRequest {
	return &JobRequest{
		ID:    id,
		Execs: exec,
	}
}

func (jr *JobRequest) AppendExec(exec *Exec) {
	jr.Execs = append(jr.Execs, exec)
}

func (jr *JobRequest) SetID(id string) {
	jr.ID = id
}

func (jr JobRequest) GetID() string {
	return jr.ID
}

func (jr JobRequest) GetExec() []*Exec {
	return jr.Execs
}
