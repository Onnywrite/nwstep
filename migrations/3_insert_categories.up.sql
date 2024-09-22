INSERT INTO categories (category_id, name, description, photo_url, background_url)
VALUES (1, 'Языки программирования', 'Различные языки программирования', 'https://www.theschoolrun.com/sites/theschoolrun.com/files/article_images/what_is_a_programming_language.jpg', 'https://i.pinimg.com/736x/db/71/5f/db715f3f37f54dbd809186aeb307c6dd.jpg');

INSERT INTO courses (name, description, min_rating, optimal_rating, category_id, photo_url) VALUES
('Основы JavaScript', 'Пришло время начать ваш путь с бравого языка..', 0, 0, 1, 'https://upload.wikimedia.org/wikipedia/commons/6/6a/JavaScript-logo.png'),
('ReactJS', 'Вы уже достаточно продвинулись чтобы изучить React', 50, 52, 1, 'https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg'),
('Python', 'Изучите Python и станьте настоящим мастером', 100, 102, 1, 'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Python-logo-notext.svg/1200px-Python-logo-notext.svg.png'),
('Golang', 'Golang - это язык для настоящих мужчин', 150, 152, 1, 'https://upload.wikimedia.org/wikipedia/commons/thumb/2/27/Gopher_mascot.svg/1200px-Gopher_mascot.svg.png'),
('C++', 'Освойте C++ и станьте гигачадом', 200, 202, 1, 'https://upload.wikimedia.org/wikipedia/commons/thumb/1/18/ISO_C%2B%2B_Logo.svg/1200px-ISO_C%2B%2B_Logo.svg.png');