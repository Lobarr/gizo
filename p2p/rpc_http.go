package p2p

//RPC exposes rpc functions
func (d Dispatcher) RPC() {
	d.GetRPC().AddFunction("Version", d.Version)
	d.GetRPC().AddFunction("PeerCount", d.PeerCount)
	d.GetRPC().AddFunction("BlockByHash", d.BlockByHash)
	d.GetRPC().AddFunction("BlockByHeight", d.BlockByHeight)
	d.GetRPC().AddFunction("Latest15Blocks", d.Latest15Blocks)
	d.GetRPC().AddFunction("LatestBlock", d.LatestBlock)
	d.GetRPC().AddFunction("PendingCount", d.PendingCount)
	d.GetRPC().AddFunction("Score", d.Score)
	d.GetRPC().AddFunction("PublicKey", d.PublicKey)
	d.GetRPC().AddFunction("NewJob", d.NewJob)
	d.GetRPC().AddFunction("NewExec", d.NewExec)
	d.GetRPC().AddFunction("WorkersCount", d.WorkersCount)
	d.GetRPC().AddFunction("WorkersCountBusy", d.WorkersCountBusy)
	d.GetRPC().AddFunction("WorkersCountNotBusy", d.WorkersCountNotBusy)
	d.GetRPC().AddFunction("ExecStatus", d.ExecStatus)
	d.GetRPC().AddFunction("CancelExec", d.CancelExec)
	d.GetRPC().AddFunction("ExecTimestamp", d.ExecTimestamp)
	d.GetRPC().AddFunction("ExecTimestampString", d.ExecTimestampString)
	d.GetRPC().AddFunction("ExecDurationNanoseconds", d.ExecDurationNanoseconds)
	d.GetRPC().AddFunction("ExecDurationSeconds", d.ExecDurationSeconds)
	d.GetRPC().AddFunction("ExecDurationMinutes", d.ExecDurationMinutes)
	d.GetRPC().AddFunction("ExecDurationString", d.ExecDurationString)
	d.GetRPC().AddFunction("ExecArgs", d.ExecArgs)
	d.GetRPC().AddFunction("ExecErr", d.ExecErr)
	d.GetRPC().AddFunction("ExecPriority", d.ExecPriority)
	d.GetRPC().AddFunction("ExecResult", d.ExecResult)
	d.GetRPC().AddFunction("ExecRetries", d.ExecRetries)
	d.GetRPC().AddFunction("ExecBackoff", d.ExecBackoff)
	d.GetRPC().AddFunction("ExecExecutionTime", d.ExecExecutionTime)
	d.GetRPC().AddFunction("ExecExecutionTimeString", d.ExecExecutionTimeString)
	d.GetRPC().AddFunction("ExecInterval", d.ExecInterval)
	d.GetRPC().AddFunction("ExecBy", d.ExecBy)
	d.GetRPC().AddFunction("ExecTTLNanoseconds", d.ExecTTLNanoseconds)
	d.GetRPC().AddFunction("ExecTTLSeconds", d.ExecTTLSeconds)
	d.GetRPC().AddFunction("ExecTTLMinutes", d.ExecTTLMinutes)
	d.GetRPC().AddFunction("ExecTTLHours", d.ExecTTLHours)
	d.GetRPC().AddFunction("ExecTTLString", d.ExecTTLString)
	d.GetRPC().AddFunction("JobQueueCount", d.JobQueueCount)
	d.GetRPC().AddFunction("LatestBlockHeight", d.LatestBlockHeight)
	d.GetRPC().AddFunction("Job", d.Job)
	d.GetRPC().AddFunction("JobSubmissionTimeUnix", d.JobSubmisstionTimeUnix)
	d.GetRPC().AddFunction("JobSubmissionTimeString", d.JobSubmisstionTimeString)
	d.GetRPC().AddFunction("IsJobPrivate", d.IsJobPrivate)
	d.GetRPC().AddFunction("JobName", d.JobName)
	d.GetRPC().AddFunction("JobLatestExec", d.JobLatestExec)
	d.GetRPC().AddFunction("JobExecs", d.JobExecs)
	d.GetRPC().AddFunction("BlockHashesHex", d.BlockHashesHex)
	d.GetRPC().AddFunction("KeyPair", d.KeyPair)
	d.GetRPC().AddFunction("Solo", d.Solo)
	d.GetRPC().AddFunction("Chord", d.Chord)
	d.GetRPC().AddFunction("Chain", d.Chain)
	d.GetRPC().AddFunction("Batch", d.Batch)
	d.GetRPC().AddFunction("GetUptime", d.GetUptime)
	d.GetRPC().AddFunction("GetUptimeString", d.GetUptimeString)
}
