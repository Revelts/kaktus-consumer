package response

type CreateThread struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	CreatedAt string `json:"created_at"`
}

type ViewAllThread struct {
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type CommentThread struct {
	ThreadId    int    `json:"thread_id"`
	CommentDesc string `json:"comment_desc"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}

type ReplyComment struct {
	ThreadId  int    `json:"thread_id"`
	CommentId int    `json:"comment_id"`
	ReplyDesc string `json:"reply_desc"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type ViewThreadDetail struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Comments    []CommentInfo `json:"comments"`
	ThreadInfo  threadInfo    `json:"thread_info"`
}

type CommentInfo struct {
	CommentId    int           `json:"comment_id"`
	CommentDesc  string        `json:"comment_desc"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    string        `json:"created_at"`
	Replies      []RepliesInfo `json:"replies"`
	TotalReplies int           `json:"total_replies"`
}

type RepliesInfo struct {
	ReplyDesc string `json:"reply_desc"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

type threadInfo struct {
	ThreadId      int    `json:"thread_id"`
	CreatedBy     string `json:"created_by"`
	TotalLikes    int    `json:"total_likes"`
	TotalComments int    `json:"total_comments"`
	CreatedAt     string `json:"created_at"`
}
