CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(200),
    photo_url VARCHAR(200),
    background_url VARCHAR(200) DEFAULT ''
);

CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR(200),
    min_rating INTEGER NOT NULL DEFAULT 0,
    optimal_rating INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories ON DELETE CASCADE,
    photo_url VARCHAR(200)
);

CREATE TABLE ratings (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users,
    category_id INTEGER REFERENCES categories ON DELETE CASCADE,
    rating INTEGER,
    UNIQUE(user_id, category_id)
);