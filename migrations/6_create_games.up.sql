CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses,
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
    game_id INTEGER REFERENCES games,
    user_id UUID REFERENCES users,
    place_id INTEGER REFERENCES places,
    health INTEGER NOT NULL,
    UNIQUE(game_id, user_id)
);

CREATE TABLE games_questions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games,
    question_id INTEGER REFERENCES questions,
    number INTEGER NOT NULL,
    UNIQUE(game_id, question_id)
);

CREATE TABLE questions_answers (
    id SERIAL PRIMARY KEY,
    game_question_id INTEGER REFERENCES games_questions,
    answer_id INTEGER REFERENCES answers,
    user_id UUID REFERENCES users,
    answered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(answer_id, user_id, game_question_id)
);