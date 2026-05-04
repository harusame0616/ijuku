ALTER TABLE
    authors
ADD
    COLUMN slug TEXT NOT NULL;

ALTER TABLE
    authors
ADD
    CONSTRAINT uq_authors_slug UNIQUE (slug);
