-- Tutorials Favorites is a pivot table to keep track of which tutorials a user has favored.
CREATE TABLE IF NOT EXISTS tutorials_bookmarks (
    id TEXT PRIMARY KEY,                                                                        -- The ID for each tutorials users pair.

    user_id TEXT NOT NULL,                                                                      -- A reference to the user who favored this tutorial.
    tutorial_id TEXT NOT NULL,                                                                  -- A reference to which tutorial was favored.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,                                              -- When did the user favor this tutorial.

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,                               -- If the user ever decides to delete their account then delete this row.
    FOREIGN KEY (tutorial_id) REFERENCES tutorials(id) ON DELETE CASCADE                        -- If the tutorial ever gets deleted then delete this row.
);

-- Ensure that the user can only favor this tutorial once. It wouldn't make a lot of sense if
-- the user could favor the same tutorial multiple times.
CREATE UNIQUE INDEX idx_tutorials_bookmarks_user_id ON tutorials_bookmarks(user_id, tutorial_id);
