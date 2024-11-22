-- Keywords table is used to store all the keywords that are used by individual tutorials.
CREATE TABLE IF NOT EXISTS keywords (
    id TEXT PRIMARY KEY,                                                                        -- The ID of the individual keyword.

    keyword TEXT NOT NULL                                                                       -- The actual keyword.
);

-- An index to enforce uniqueness of each keyword as well as speed up searching of keywords
-- without the use of their index.
CREATE UNIQUE INDEX idx_keywords_keyword ON keywords(keyword);
