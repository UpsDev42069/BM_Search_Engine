DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  password_reset_required BOOLEAN DEFAULT FALSE
);

-- Original pages table without the last_updated column
CREATE TABLE IF NOT EXISTS pages (
    title TEXT PRIMARY KEY UNIQUE,
    url TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en', -- How you define ENUM type in SQLite
    last_updated TIMESTAMP,
    content TEXT NOT NULL
);

-- Create FTS5 table for full-text searching
CREATE VIRTUAL TABLE IF NOT EXISTS pages_fts USING fts5(
  title,
  content,
  language UNINDEXED
  );

  -- Insert trigger to keep FTS table up to date
  CREATE TRIGGER IF NOT EXISTS pages_ai AFTER INSERT ON pages BEGIN
    INSERT INTO pages_fts (title, content, language) VALUES (new.title, new.content, new.language);
    END;
 
  -- Update trigger to keep FTS table up to date
  CREATE TRIGGER IF NOT EXISTS pages_au AFTER UPDATE ON pages BEGIN
    UPDATE pages_fts 
    SET content = NEW.content,
        language = NEW.language
    WHERE title = OLD.title;
    END;
  
  -- Delete trigger to keep FTS table up to date
  CREATE TRIGGER IF NOT EXISTS pages_ad AFTER DELETE ON pages BEGIN
    DELETE FROM pages_fts WHERE title = OLD.title;
    END;
  
  -- Automatically update `last_updated` upon any update to the pages table
  CREATE TRIGGER IF NOT EXISTS pages_last_updated_trigger AFTER UPDATE ON pages
    BEGIN
      UPDATE pages
      SET last_updated = CURRENT_TIMESTAMP
      WHERE title = OLD.title;
    END;