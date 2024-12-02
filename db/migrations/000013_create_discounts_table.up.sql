CREATE TABLE IF NOT EXISTS discounts (
    id TEXT PRIMARY KEY,

    title TEXT NOT NULL,
    description TEXT NOT NULL,
    code TEXT NOT NULL,
    discount INTEGER NOT NULL CHECK (discount > 0 AND discount <= 100),
    uses INTEGER NOT NULL CHECK (uses > 0),
    active INTEGER NOT NULL DEFAULT 0 CHECK (active >= 0 AND active <= 1),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX discounts_title_idx ON discounts(title);

CREATE UNIQUE INDEX discounts_title_code ON discounts(code);
