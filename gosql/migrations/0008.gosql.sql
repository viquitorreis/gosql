-- gosql Up
CREATE TABLE testando (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);

-- gosql Down 
DROP TABLE "testando";
