# 第一阶段设计文档（修订版）

## 1. 阶段目标

- **内部成员管理**：通过服务器端 CLI 工具（非 Web 接口）添加内部管理员账号。
- **内部帖子上传**：管理员通过 Web 接口上传 `.zip` 压缩包（需认证）。
- **帖子展示**：前端支持按标签筛选，以只读方式浏览已发布的帖子。
- **阅读量统计**：记录帖子的点击浏览数量。

### 明确不做

- 普通用户注册 / 公网上传功能。
- 复杂的用户身份认证体系（Session / JWT）。

---

## 2. 技术选型与数据存储

- **后端语言**：Golang（标准库 `net/http` + `database/sql`）
- **数据库**：SQLite（文件型，便于部署）
- **密码存储**：`bcrypt` 哈希（CLI 添加成员时直接生成哈希，避免明文存储）
- **认证方式**：**HTTP Basic Auth**（无状态，每个受保护接口请求头携带 `Authorization: Basic base64(nick:pass)`）

---

## 3. SQLite 数据库表设计（含精确类型）

> **注意**：SQLite 建议使用 `INTEGER` 作为主键（自增），时间字段使用 `DATETIME` 或 `INTEGER (Unix时间戳)`，布尔/标志位使用 `INTEGER`（0/1）。

### 3.1 成员表 (`members`)

用于存放内部管理员账号（仅 CLI 工具写入）。

```sql
CREATE TABLE members (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    power       INTEGER DEFAULT 0,    -- 位运算权限（预留，第一阶段全为 0 或 1）
    nick        TEXT NOT NULL UNIQUE, -- 登录用户名，唯一
    email       TEXT NOT NULL,
    passwd      TEXT NOT NULL,        -- 存储 bcrypt 哈希值（60字节）
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_delete   INTEGER DEFAULT 0     -- 0=正常，1=已删除/禁用
);
```

### 3.2 帖子原始数据表 (`posts`)

```sql
CREATE TABLE posts (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    title            TEXT NOT NULL,
    author           TEXT,                    -- 可空
    overview         TEXT NOT NULL,                    -- 概览
    markdown_content TEXT NOT NULL,           -- 原始 Markdown
    is_pending       INTEGER DEFAULT 1,       -- 1=待审核，0=已发布，2=审核不通过
    view_count       INTEGER DEFAULT 0,       -- 阅读量（原 click_time 更名）
    create_time      DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_time      DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3.3 标签表 (`tags`)

```sql
CREATE TABLE tags (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE,
    description TEXT
);
```

### 3.4 帖子-标签关联表 (`post_tag_map`)

```sql
CREATE TABLE post_tag_map (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    tag_id  INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
```

---

## 4. 成员添加（CLI 工具）

> 不开发 Web 注册页面。在服务器上直接运行一个独立的 Go 二进制或 Python 脚本。

**命令示例**：

```bash
./admin_tool add --nick=zhangsan --email=zs@example.com --password=123456
```

**脚本内部逻辑**：

1. 校验 `nick` 唯一性。
2. 使用 `bcrypt` 对 `password` 进行哈希。
3. 插入 `members` 表。
4. `power` 字段暂设为 `1`（代表拥有全部管理员权限）。

---

## 5. 核心接口规范（RESTful）

### 通用约定

- **公开接口**：无需认证。
- **管理接口**：全部使用 **HTTP Basic Auth**（校验 `nick` 与 `passwd` 是否匹配 `members` 表）。
- **响应格式**：统一 `application/json`。
  - 成功：`{"code": 200, "status": "OK", "data": ...}`
  - 失败：`{"code": 4xx/5xx, "status": "Error", "message": "..."}`

---

### 5.1 公开 - 获取帖子列表（按标签筛选）

`GET /api/posts`

- **认证**：无
- **请求参数（Query String）**：
  - `tags`：可选，逗号分隔的标签名，例如 `?tags=math,travel`。若不传则返回全部。
- **响应示例**：

```json
{
    "code": 200,
    "status": "OK",
    "data": {
        "posts": [
            {
                "id": 10,
                "title": "如何高效刷题",
                "author": "张三",
                "visit_path": "/posts/open/10",  // 前端直接使用此路径跳转
                "view_count": 128
            }
        ]
    }
}
```

- **补充说明**：`/posts/open/{id}` 由后端静态文件服务提供，实际映射到 `./posts/{id}/index.html`。

---

### 5.2 管理 - 上传帖子压缩包

`POST /api/posts/upload`

- **认证**：需要 Basic Auth（管理员）。
- **请求体**：`multipart/form-data`，字段名 `zip_file`。
- **处理流程**：
  1. 解压 `.zip` 到临时目录。
  2. 校验根目录下必须包含 `index.md` 和 `config.toml`（用于提取 Title/Author/Tags）。
  3. 将整个目录移动到 `./posts_origin/{timestamp}/`。
  4. 向 `posts` 表插入记录，`is_pending=1`（待审核），`markdown_content` 存储 `index.md` 内容。
  5. 解析 `config.toml` 中的标签，存入 `post_tag_map`。
- **响应**：

```json
{"code": 200, "status": "OK", "data": {"post_id": 10, "message": "上传成功，等待审核"}}
```

---

### 5.3 管理 - 获取待审核列表

`GET /api/posts/pending/list`

- **认证**：需要 Basic Auth。
- **响应**（仅展示标题与上传时间）：

```json
{
    "code": 200,
    "status": "OK",
    "data": {
        "list": [
            {"id": 5, "title": "待审核文章", "create_time": "2026-08-16 10:00:00"}
        ]
    }
}
```

---

### 5.4 管理 - 查看待审核帖子的详细内容

`GET /api/posts/pending/{id}`

- **认证**：需要 Basic Auth。
- **响应**：返回完整的元数据 + Markdown 内容（用于管理员预览）。

```json
{
    "code": 200,
    "status": "OK",
    "data": {
        "post_id": 5,
        "title": "待审核文章",
        "author": "李四",
        "markdown_content": "# 正文...",
        "tags": ["math", "life"]
    }
}
```

---

### 5.5 管理 - 审核通过（发布）

`POST /api/posts/pending/{id}/approve`

- **认证**：需要 Basic Auth。
- **逻辑**：
  1. 检查 `is_pending` 是否为 `1`。
  2. 使用 `goldmark` 等库将 `markdown_content` 渲染为 `index.html`。
  3. 将 `./posts_origin/{id}` 中的图片等静态资源一并复制到 `./posts/{id}/`。
  4. 更新数据库：`is_pending = 0`, `update_time = CURRENT_TIMESTAMP`。
- **响应**：

```json
{"code": 200, "status": "OK", "message": "发布成功"}
```

---

### 5.6 管理 - 审核不通过（拒绝）

`POST /api/posts/pending/{id}/reject`

- **认证**：需要 Basic Auth。
- **逻辑**：
  1. 更新 `is_pending = 2`（审核不通过）。
  2. 可选择记录拒绝原因（可选，暂不强制）。
- **响应**：

```json
{"code": 200, "status": "OK", "message": "已拒绝"}
```

---

## 6. 静态资源与目录结构

### 6.1 源文件存储目录（原始压缩包解压处）

用于存放待审核的原始文件（含 md 和图片），供审核时预览或重新生成。

```
posts_origin/
├── 5/              # post_id 为 5
│   ├── index.md
│   ├── config.toml
│   └── pic1.png
└── 6/
    └── ...
```

### 6.2 网页发布目录（最终对外服务）

审核通过后，从此目录提供静态文件服务。

```
posts/
├── 10/             # post_id 为 10
│   ├── index.html  # 由 Markdown 生成
│   └── pic1.png
└── 11/
    └── ...
```

### 6.3 路由映射

- `/posts/open/{id}` → 静态文件服务指向 `./posts/{id}/index.html`（或目录，自动寻找 index.html）。

---

## 7. 阅读量（`view_count`）记录逻辑

- 当请求 `GET /posts/open/{id}` 时，后端在处理静态文件前，异步执行 `UPDATE posts SET view_count = view_count + 1 WHERE id = ?`。
- **注意**：为防止爬虫或频繁刷新导致计数失真，可增加简单的 IP 去重（第一阶段可忽略，直接累加即可）。

---

## 8. 开发阶段排期建议（单人两周）

| 天数 | 任务 |
| :--- | :--- |
| Day 1-2 | 搭建 Golang 项目结构，完成 SQLite 初始化与 ORM（或原生 SQL）封装。 |
| Day 3-4 | 实现 CLI 工具（添加成员），实现 Basic Auth 中间件。 |
| Day 5-6 | 实现公开接口（`GET /api/posts`）与静态文件服务（含 `view_count`）。 |
| Day 7-9 | 实现上传压缩包接口、解压与元数据提取逻辑。 |
| Day 10-11 | 实现审核列表、查看详情、审核通过/拒绝接口（含 Markdown 渲染为 HTML）。 |
| Day 12-14 | 端到端联调、编写简单的前端测试页面、修复 Bug、部署上线。 |

---

## 9. 安全与注意事项（简化但重要）

1. **Basic Auth 务必搭配 HTTPS**（内网环境可酌情放宽）。
2. **密码永不明文存储**，CLI 工具必须使用 `bcrypt`。
3. **上传文件限制**：限制压缩包大小（如 20MB），防止恶意攻击。
4. **解压安全**：防范 Zip Slip 漏洞（路径遍历），确保解压文件均在目标子目录内。
