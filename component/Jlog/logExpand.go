package Jlog

type LogIoWrite struct {
	Msg string
	Flag string
}

func (log *LogIoWrite) Write(p []byte) (n int, err error){
	processSugared.Infow(log.Msg, log.Flag, string(p))
	return
}

