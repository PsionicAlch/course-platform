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
