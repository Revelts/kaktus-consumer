package Repository

import (
	"kaktus-task/Repository/public"
	"kaktus-task/Repository/threads"
)

type repository struct {
	Threads threads.ThreadInterface
	Public  public.PublicInterface
}

var AllRepository = repository{
	Threads: threads.InitThreads(),
	Public:  public.InitPublic(),
}
