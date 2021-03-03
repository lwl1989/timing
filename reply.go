package timer

var successReply = Reply{
	Code: "200",
	Msg:  "执行成功",
	Err:  nil,
}

func getDefaultSuccessReply(task TaskInterface) Reply {
	return Reply{
		Code: "200",
		Msg:  "执行成功",
		Err:  nil,
		Ts:   task,
	}
}

func getDefaultErrorReply(task TaskInterface, err error) Reply {
	return Reply{
		Code: "-1",
		Msg:  "执行失败",
		Err:  err,
		Ts:   task,
	}
}

func GetReply(task TaskInterface, Code, Msg string, err error) Reply {
	return Reply{
		Code: Code,
		Msg:  "执行失败",
		Err:  err,
		Ts:   task,
	}
}
