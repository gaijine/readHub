CREATE TABLE users
(
    id              BIGSERIAL PRIMARY KEY,
	telegram_id     BIGINT NOT NULL UNIQUE,
	username        TEXT NOT NULL DEFAULT '',
	created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE books
(
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL,
    open_library_id TEXT NOT NULL,
    title           TEXT NOT NULL,
    author          TEXT NOT NULL DEFAULT '',

    total_pages     INT CHECK (total_pages > 0),
    current_page    INT NOT NULL DEFAULT 0 
        CHECK(current_page >= 0),

    status          TEXT NOT NULL
        CHECK (status IN ('want', 'reading', 'completed')),
    
    cover_url       TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, open_library_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE reading_sessions
(
    id              BIGSERIAL PRIMARY KEY,
    book_id         BIGINT NOT NULL,
    user_id         BIGINT NOT NULL,
    started_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    finished_at     TIMESTAMP,
    start_page      INT,
    end_page        INT,

    FOREIGN KEY (book_id) REFERENCES books(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)