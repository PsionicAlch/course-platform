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
CREATE UNIQUE INDEX IF NOT EXISTS idx_tutorials_likes_user_id ON tutorials_likes(user_id, tutorial_id);
