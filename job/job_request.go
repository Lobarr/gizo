package job

//Request requests job execs
type Request struct {
	ID    string
	Execs []*Exec
}

//NewRequest initializes job request
func NewRequest(id string, exec ...*Exec) *Request {
	return &Request{
		ID:    id,
		Execs: exec,
	}
}

//AppendExec adds exec to job request
func (jr *Request) AppendExec(exec *Exec) {
	jr.Execs = append(jr.Execs, exec)
}

//SetID sets job id
func (jr *Request) SetID(id string) {
	jr.ID = id
}

//GetID returns job reqiest id
func (jr Request) GetID() string {
	return jr.ID
}

//GetExec returns job execs
func (jr Request) GetExec() []*Exec {
	return jr.Execs
}
