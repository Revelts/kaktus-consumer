package threads

import (
	"kaktus-task/Controllers/Dto/request"
	"kaktus-task/Controllers/Dto/response"
)

type ThreadInterface interface {
	CreateThread(params request.CreateThread) (response response.CreateThread, err error)
	ViewThreads() (threadsList []response.ViewAllThread, err error)
	LikeThread(params request.LikeThread) (executedRow int64, err error)
	ViewThreadDetail(threadId int) (threadDetail response.ViewThreadDetail, err error)
	CommentThread(params request.CommentThread) (resp response.CommentThread, err error)
	ReplyComment(params request.ReplyComment) (resp response.ReplyComment, err error)
}

type thread struct{}

func InitThreads() ThreadInterface {
	return &thread{}
}
