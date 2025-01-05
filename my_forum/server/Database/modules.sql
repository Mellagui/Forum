-- 1. Users 👤​
CREATE TABLE IF NOT EXISTS users (
    id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    password_hash NVARCHAR NOT NULL,
    user_image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Posts 📝​
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT,
    content TEXT,
    image_url TEXT, -- URL or path to the post image
    category TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 3. Comments 💭​
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 4. Categories 🏷️​
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    category_name TEXT UNIQUE NOT NULL,
    created_by_user_id INTEGER, -- The user who created this category
    FOREIGN KEY (created_by_user_id) REFERENCES users(id)
);

-- 5. PostCategory 🔗​
CREATE TABLE IF NOT EXISTS postCategory (
    post_id INTEGER UNIQUE NOT NULL,
    category_id INTEGER UNIQUE NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- 6. LikeDislike 👍👎​
CREATE TABLE IF NOT EXISTS likeDislike (
    id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    post_id INTEGER,
    is_like BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
    -- FOREIGN KEY (comment_id) REFERENCES Comments(id)
);

-- 7. Session 🍪​
-- CREATE TABLE IF NOT EXISTS Session (
--     id TEXT UNIQUE PRIMARY KEY AUTOINCREMENT,
--     user_id TEXT UNIQUE NOT NULL,
--     token TEXT UNIQUE NOT NULL,
--     expires_at DATETIME NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES Users(id)
-- );

-- 8. Conversations 💬​
-- CREATE TABLE IF NOT EXISTS Conversations (
--     id TEXT UNIQUE PRIMARY KEY AUTOINCREMENT,
--     user1_id TEXT UNIQUE NOT NULL,
--     user2_id TEXT UNIQUE NOT NULL,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (user1_id) REFERENCES Users(id),
--     FOREIGN KEY (user2_id) REFERENCES Users(id)
-- );

-- 9. Messages ✉️​
-- CREATE TABLE IF NOT EXISTS Messages (
--     id TEXT UNIQUE PRIMARY KEY AUTOINCREMENT,
--     conversation_id TEXT UNIQUE NOT NULL,
--     sender_id TEXT UNIQUE NOT NULL,
--     content TEXT NOT NULL,
--     sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (conversation_id) REFERENCES Conversations(id),
--     FOREIGN KEY (sender_id) REFERENCES Users(id)
-- );
