CREATE TABLE members (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    power       INTEGER DEFAULT 0,    -- 位运算权限（预留，第一阶段全为 0 或 1）
    nick        TEXT NOT NULL UNIQUE, -- 登录用户名，唯一
    email       TEXT NOT NULL,
    passwd      TEXT NOT NULL,        -- 存储 bcrypt 哈希值（60字节）
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_delete   INTEGER DEFAULT 0     -- 0=正常，1=已删除/禁用
);
CREATE TABLE posts (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    title            TEXT NOT NULL,
    author           TEXT,                    -- 可空
    overview         TEXT NOT NULL,           -- 概览
    content          TEXT NOT NULL,           -- 原始内容， 依靠 format 区分格式
    format           INTEGER NOT NULL,        -- content 的格式
    status           INTEGER DEFAULT 1,       -- 1=待审核，0=已发布，-1=审核不通过
    view_count       INTEGER DEFAULT 0,       -- 阅读量（原 click_time 更名）
    create_time      DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time      DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE tags (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE,
    description TEXT
);
CREATE TABLE post_tag_map (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_id  INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
