CREATE TABLE public.users (
    user_id serial PRIMARY KEY,
    user_name VARCHAR(255)
)

CREATE TABLE public.threads (
    thread_id serial PRIMARY KEY,
    thread_title VARCHAR(255),
    thread_desc TEXT,
    created_by integer,
    created_at timestamp without time zone default now(),
    updated_at timestamp without time zone default now(),
    CONSTRAINT fk_users
        FOREIGN KEY(created_by)
            REFERENCES public.users(user_id)
)

CREATE TABLE public.comments (
   comment_id serial PRIMARY KEY,
   thread_id integer,
   comment_desc TEXT,
   created_by integer,
   created_at timestamp without time zone default now(),
   updated_at timestamp without time zone default now(),
   CONSTRAINT fk_threads
       FOREIGN KEY(thread_id)
           REFERENCES public.threads(thread_id)
)

CREATE TABLE public.replies (
    reply_id serial PRIMARY KEY,
    comment_id integer,
    reply_desc TEXT,
    created_by integer,
    created_at timestamp without time zone default now(),
    updated_at timestamp without time zone default now(),
    CONSTRAINT fk_comments
     FOREIGN KEY(comment_id)
         REFERENCES public.comments(comment_id)
)

CREATE TABLE public.likes (
    thread_id integer,
    created_by integer,
    created_at timestamp without time zone default now(),
    CONSTRAINT fk_likes,
        FOREIGN KEY(thread_id)
            REFERENCES public.threads(thread_id)
)