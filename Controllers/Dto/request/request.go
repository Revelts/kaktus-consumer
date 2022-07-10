package request

type CreateThread struct {
	ThreadTitle string `json:"thread_title"`
	ThreadDesc  string `json:"thread_desc"`
	CreatedBy   int    `json:"created_by"`
}

type ViewThreadDetail struct {
	ThreadId int `json:"thread_id"`
}

type LikeThread struct {
	ThreadId int `json:"thread_id"`
	UserId   int `json:"user_id"`
}

type CommentThread struct {
	ThreadId    int    `json:"thread_id"`
	CommentDesc string `json:"comment_desc"`
	CreatedBy   int    `json:"created_by"`
}

type ReplyComment struct {
	CommentId int    `json:"comment_id"`
	ReplyDesc string `json:"reply_desc"`
	CreatedBy int    `json:"created_by"`
}
