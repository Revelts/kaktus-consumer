package threads

import (
	"database/sql"
	"kaktus-task/Connections"
	"kaktus-task/Controllers/Dto/request"
	"kaktus-task/Controllers/Dto/response"
	"log"
)

func (t thread) CreateThread(params request.CreateThread) (response response.CreateThread, err error) {
	query := `INSERT INTO public.threads (thread_title, thread_desc, created_by) VALUES ($1, $2, $3)
					RETURNING thread_id, thread_title, thread_desc, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS')`
	err = Connections.PostgreConn.QueryRow(query, params.ThreadTitle, params.ThreadDesc, params.CreatedBy).Scan(
		&response.Id, &response.Title, &response.Desc, &response.CreatedAt)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	return
}

func (t thread) ViewThreads() (threadsList []response.ViewAllThread, err error) {
	var threadsRow *sql.Rows
	query := `SELECT thread_title, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS') FROM public.threads`
	threadsRow, err = Connections.PostgreConn.Query(query)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	for threadsRow.Next() {
		var thread response.ViewAllThread
		err = threadsRow.Scan(&thread.Title, &thread.CreatedAt)
		if err != nil {
			log.Panicf("[ERROR]: %s", err)
			return
		}
		threadsList = append(threadsList, thread)
	}
	return
}

func (t thread) LikeThread(params request.LikeThread) (executedRow int64, err error) {
	query := `INSERT INTO public.likes (thread_id, created_by) 
				SELECT $1, $2 WHERE NOT EXISTS (
				    SELECT 1 
			  			FROM public.likes 
				WHERE thread_id = $1 AND created_by = $2
	)`
	result, err := Connections.PostgreConn.Exec(query, params.ThreadId, params.UserId)
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	executedRow, err = result.RowsAffected()
	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	return
}

func (t thread) ViewThreadDetail(threadId int) (threadDetail response.ViewThreadDetail, err error) {
	query := `
		SELECT
		    t.thread_id as thread_id,
			t.thread_title AS title,
			t.thread_desc AS description,
			u.user_name AS created_by,
			to_char(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			(SELECT COUNT(*) FROM public.likes l
				WHERE l.thread_id  = t.thread_id
			) AS total_likes,
			(SELECT COUNT(*) FROM public.comments c
				WHERE c.thread_id  = t.thread_id
			) AS total_comments
		FROM public.threads t
			INNER JOIN users u ON t.created_by = u.user_id
		WHERE t.thread_id = $1
	`
	err = Connections.PostgreConn.QueryRow(query, threadId).Scan(
		&threadDetail.ThreadInfo.ThreadId,
		&threadDetail.Title,
		&threadDetail.Description,
		&threadDetail.ThreadInfo.CreatedBy,
		&threadDetail.ThreadInfo.CreatedAt,
		&threadDetail.ThreadInfo.TotalLikes,
		&threadDetail.ThreadInfo.TotalComments,
	)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	query = `SELECT 
    			c.comment_id,
				c.comment_desc, 
				u.user_name , 
				to_char(c.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
			FROM public.comments c
			LEFT JOIN public.users u on u.user_id  = c.created_by 
			WHERE c.thread_id = $1 ORDER BY c.comment_id ASC`

	rows, err := Connections.PostgreConn.Query(query, threadId)

	rowIndex := 0
	for rows.Next() {
		var thread response.CommentInfo
		err = rows.Scan(&thread.CommentId, &thread.CommentDesc, &thread.CreatedBy, &thread.CreatedAt)
		if err != nil {
			log.Panicf("[ERROR]: %s", err)
			return
		}
		threadDetail.Comments = append(threadDetail.Comments, thread)

		queryReply := `SELECT 
    			r.reply_desc,
				u.user_name , 
				to_char(r.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at
			FROM public.replies r
			LEFT JOIN public.users u on u.user_id  = r.created_by 
			WHERE r.comment_id = $1`

		var replyRows *sql.Rows
		replyRows, err = Connections.PostgreConn.Query(queryReply, thread.CommentId)

		if err != nil {
			log.Panicf("[ERROR]: %s", err)
			return
		}

		for replyRows.Next() {
			var replies response.RepliesInfo
			err = replyRows.Scan(&replies.ReplyDesc, &replies.CreatedBy, &replies.CreatedAt)
			if err != nil {
				log.Panicf("[ERROR]: %s", err)
				return
			}

			threadDetail.Comments[rowIndex].Replies = append(
				threadDetail.Comments[rowIndex].Replies,
				replies)
		}
		threadDetail.Comments[rowIndex].TotalReplies = len(threadDetail.Comments[rowIndex].Replies)

		if err != nil {
			log.Panicf("[ERROR]: %s", err)
			return
		}
		rowIndex++
	}
	return
}

func (t thread) CommentThread(params request.CommentThread) (response response.CommentThread, err error) {
	query := `INSERT INTO public.comments (thread_id, comment_desc, created_by) VALUES ($1, $2, $3)
					RETURNING thread_id, comment_desc, created_by, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS')`

	var userId int
	err = Connections.PostgreConn.QueryRow(query, params.ThreadId, params.CommentDesc, params.CreatedBy).Scan(
		&response.ThreadId, &response.CommentDesc, &userId, &response.CreatedAt)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	query = `SELECT user_name FROM public.users WHERE user_id = $1`
	err = Connections.PostgreConn.QueryRow(query, userId).Scan(&response.CreatedBy)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	return
}

func (t thread) ReplyComment(params request.ReplyComment) (response response.ReplyComment, err error) {
	query := `INSERT INTO public.replies (comment_id, reply_desc, created_by) VALUES ($1, $2, $3)
					RETURNING comment_id, reply_desc, created_by, TO_CHAR(created_at,'YYYY-MM-DD HH24:MI:SS')`

	var userId int
	err = Connections.PostgreConn.QueryRow(query, params.CommentId, params.ReplyDesc, params.CreatedBy).Scan(
		&response.CommentId, &response.ReplyDesc, &userId, &response.CreatedAt)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	query = `SELECT thread_id FROM public.comments WHERE comment_id = $1`
	err = Connections.PostgreConn.QueryRow(query, response.CommentId).Scan(&response.ThreadId)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}

	query = `SELECT user_name FROM public.users WHERE user_id = $1`
	err = Connections.PostgreConn.QueryRow(query, userId).Scan(&response.CreatedBy)

	if err != nil {
		log.Panicf("[ERROR]: %s", err)
		return
	}
	return
}
