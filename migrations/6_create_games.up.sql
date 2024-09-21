CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses,
    start_at TIMESTAMP NOT NULL DEFAULT NOW(),
    end_at TIMESTAMP
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

-- comment
CREATE TABLE games_questions (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games,
    question_id INTEGER REFERENCES questions,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL
);

CREATE TABLE questions_answers (
    id SERIAL PRIMARY KEY,
    game_question_id INTEGER REFERENCES games_questions,
    answer_id INTEGER REFERENCES answers,
    user_id UUID REFERENCES users
)