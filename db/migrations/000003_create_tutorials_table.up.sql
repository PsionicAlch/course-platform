-- Keywords table is used to store all the keywords that are used by individual tutorials.
CREATE TABLE IF NOT EXISTS keywords (
    id TEXT PRIMARY KEY,                                                                        -- The ID of the individual keyword.

    keyword TEXT NOT NULL                                                                       -- The actual keyword.
);

-- An index to enforce uniqueness of each keyword as well as speed up searching of keywords
-- without the use of their index.
CREATE UNIQUE INDEX idx_keywords_keyword ON keywords(keyword);

-- Tutorials table holds the content and meta data for each individual tutorial to remove
-- the need for a RAM based cache of all the compiled tutorials (viable option though).
CREATE TABLE IF NOT EXISTS tutorials (
    id TEXT PRIMARY KEY,                                                                        -- The ID of the individual tutorial.

    title TEXT NOT NULL,                                                                        -- The title of the tutorial.
    slug TEXT NOT NULL,                                                                         -- URL friendly slug for the tutorial.
    description TEXT NOT NULL,                                                                  -- A short description of the tutorial.
    thumbnail_url TEXT NOT NULL,                                                                -- URL for the thumbnail image of the tutorial.
    banner_url TEXT NOT NULL,                                                                   -- URL for the banner image of the tutorial.
    content TEXT NOT NULL,                                                                      -- HTML based contents of the tutorial.
    published INTEGER DEFAULT 0,                                                                -- BOOLEAN to represent whether or not the tutorial has been published.

    author_id TEXT,                                                                             -- The user ID who published the tutorial (for when I have multiple authors).

    file_checksum TEXT NOT NULL,                                                                -- A SHA256 checksum to speed up the process of checking if a file has changed.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- Timestamp for when the tutorial was initially created.
    update_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                               -- Timestamp for when the tutorial was updated.

    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL                             -- Foreign key for users table. Set to null if the author is deleted.
);

-- A unique index to help speed up searching via slugs.
CREATE UNIQUE INDEX idx_tutorials_slug ON tutorials(slug);

-- Tutorials Keywords is a pivot table since a tutorial can have multiple keywords whilst
-- a keyword can belong to multiple tutorials.
CREATE TABLE IF NOT EXISTS tutorials_keywords (
    id TEXT PRIMARY KEY,                                                                        -- The ID for each tutorials keywords pair.

    tutorial_id TEXT NOT NULL,                                                                  -- A reference to which tutorial uses this keyword.
    keyword_id TEXT NOT NULL,                                                                   -- A reference to the keyword.

    FOREIGN KEY (tutorial_id) REFERENCES tutorials(id) ON DELETE CASCADE,                       -- If the tutorial gets deleted so should this row.
    FOREIGN KEY (keyword_id) REFERENCES keywords(id) ON DELETE CASCADE                          -- If a keyword gets deleted so should this row.
);

-- Ensure that there is ever only one unique pair of tutorial and keyword pairing. It doesn't
-- make sense for a tutorial to have 3 keywords that are all the same.
CREATE UNIQUE INDEX idx_tutorials_keywords ON tutorials_keywords(tutorial_id, keyword_id);

-- Tutorials Likes is a pivot table to keep track of which users like which tutorials. Not
-- really necessary but could be useful in the future.
CREATE TABLE IF NOT EXISTS tutorials_likes (
    id TEXT PRIMARY KEY,                                                                        -- The ID for each tutorials users pair.

    user_id TEXT NOT NULL,                                                                      -- A reference to which user liked the tutorial.
    tutorial_id TEXT NOT NULL,                                                                  -- A reference to which tutorial was liked.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- When did the user like this tutorial.

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,                               -- If the user ever decides to delete their account then there is no need to keep the like information.
    FOREIGN KEY (tutorial_id) REFERENCES tutorials(id) ON DELETE CASCADE                        -- If the tutorial is ever deleted then their is no need to keep the like information.
);

-- Ensure that the user can only ever like a tutorial once. It wouldn't make sense for a user
-- to like the same tutorial multiple times.
CREATE UNIQUE INDEX idx_tutorials_likes_user_id ON tutorials_likes(user_id, tutorial_id);

-- Tutorials Favorites is a pivot table to keep track of which tutorials a user has favored.
CREATE TABLE IF NOT EXISTS tutorials_favorites (
    id TEXT PRIMARY KEY,                                                                        -- The ID for each tutorials users pair.

    user_id TEXT NOT NULL,                                                                      -- A reference to the user who favored this tutorial.
    tutorial_id TEXT NOT NULL,                                                                  -- A reference to which tutorial was favored.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- When did the user favor this tutorial.

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,                               -- If the user ever decides to delete their account then delete this row.
    FOREIGN KEY (tutorial_id) REFERENCES tutorials(id) ON DELETE CASCADE                        -- If the tutorial ever gets deleted then delete this row.
);

-- Ensure that the user can only favor this tutorial once. It wouldn't make a lot of sense if
-- the user could favor the same tutorial multiple times.
CREATE UNIQUE INDEX idx_tutorials_favorites_user_id ON tutorials_favorites(user_id, tutorial_id);
