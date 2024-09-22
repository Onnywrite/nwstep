CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses ON DELETE CASCADE,
    start_at TIMESTAMP NOT NULL DEFAULT NOW(),
    end_at TIMESTAMP,
    last_question_number INTEGER NOT NULL DEFAULT 0,
    last_question_time TIMESTAMP
);

CREATE TABLE places (
    place_id SERIAL PRIMARY KEY,
    reward INTEGER
);

CREATE TABLE games_users (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games ON DELETE CASCADE,
    user_id UUID REFERENCES users,
    place_id INTEGER REFERENCES places ON DELETE CASCADE,
    health INTEGER NOT NULL,
    UNIQUE(game_id, user_id)
);

CREATE TABLE games_questions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions ON DELETE CASCADE,
    number INTEGER NOT NULL,
    UNIQUE(game_id, question_id)
);

CREATE TABLE questions_answers (
    id SERIAL PRIMARY KEY,
    game_question_id INTEGER REFERENCES games_questions ON DELETE CASCADE,
    answer_id INTEGER REFERENCES answers ON DELETE CASCADE,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    answered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(answer_id, user_id, game_question_id)
);