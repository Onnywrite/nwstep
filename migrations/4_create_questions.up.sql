CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    question VARCHAR(200),
    course_id INTEGER REFERENCES courses
);

CREATE TABLE answers (
    answer_id SERIAL PRIMARY KEY,
    answer VARCHAR(200),
    question_id INTEGER REFERENCES questions,
    is_correct BOOLEAN
);