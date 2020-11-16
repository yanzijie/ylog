package ylog

import "testing"

func TestConsoleLevel(t *testing.T){
	t.Logf("%T %v\n",DebugLevel,DebugLevel)
	t.Logf("%T %v\n",InfoLevel,InfoLevel)
	t.Logf("%T %v\n",WarningLevel,WarningLevel)
	t.Logf("%T %v\n",ErrorLevel,ErrorLevel)
	t.Logf("%T %v\n",FatalLevel,FatalLevel)
}
