INSERT INTO categories (category_id, name, description, photo_url, background_url)
VALUES  
(2,'ВЕБ-ФРЕЙМВОРКИ', 'изучение различных направлений веб-разработки', 'https://frankfurt.apollo.olxcdn.com/v1/files/yotk6yerz8x63-UZ/image;s=1080x735', 'https://i.pinimg.com/564x/59/76/ec/5976ecb9a136c0b40967e9fb8ee48e21.jpg');

INSERT INTO courses  
(name, description, min_rating, optimal_rating, category_id, photo_url) 
VALUES
('Изучение backend-разработки Python', 'Стань профессиональным разработчиком распределённых систем на Python', 0, 0, 2, 'https://bilgi.uz/upload/resize_cache/iblock/5c7/6n0932l2qee6wh5guz0n50ftv72yhcfe/0_350_2/2.png');

INSERT INTO questions (question, course_id)
values
('Что такое backend на Python?', 6),
('Какие основные технологии используются для создания backend на Python?', 6),
('Что такое Django?', 6),
('Что такое Flask?', 6),
('Что такое Tornado?', 6),
('Что такое RESTful API?', 6),
('Что такое ORM (Object Relational Mapping)?', 6),
('Что такое SQLAlchemy?', 6),
('Что такое Celery?', 6),
('Что такое асинхронное программирование?', 6),
('Что такое MVC (Model-View-Controller)?', 6),
('Что такое middleware?', 6),
('Что такое асинхронность?', 6),
('Что такое асинхронный запрос?', 6);

INSERT INTO answers (answer, question_id, is_correct) values
--1
('Серверная часть веб-приложения.', 11, TRUE),
('Клиентская часть веб-приложения.', 11, FALSE),
('База данных.', 11, FALSE),
('Веб-фреймворк.', 11, FALSE),
--2
('Django, Flask, Tornado.', 12, TRUE),
('React, Angular, Vue.js.', 12, FALSE),
('MySQL, PostgreSQL, MongoDB.', 12, FALSE),
('Python, HTML, CSS.', 12, FALSE),
--3
('Фреймворк для создания backend на Python.', 13, TRUE),
('Библиотека для работы с базами данных.', 13, FALSE),
('Инструмент для тестирования кода.', 13, FALSE),
('Система контроля версий.', 13, FALSE),
--4
('Фреймворк для создания backend на Python.', 14, TRUE),
('Библиотека для работы с базами данных.', 14, FALSE),
('Инструмент для тестирования кода.', 14, FALSE),
('Система контроля версий.', 14, FALSE),
--5
('Фреймворк для создания backend на Python.', 15, TRUE),
('Библиотека для работы с базами данных.', 15, FALSE),
('Инструмент для тестирования кода.', 15, FALSE),
('Система контроля версий.', 15, FALSE),
--6
('Протокол взаимодействия backend и frontend.', 16, TRUE),
('Технология создания веб-приложений.', 16, FALSE),
('Набор инструментов для разработки мобильных приложений.', 16, FALSE),
('Метод оптимизации производительности сайта.', 16, FALSE),
--7
('Технология для работы с базами данных.', 17, FALSE),
('Метод оптимизации производительности сайта.', 17, FALSE),
('Фреймворк для создания backend на Python.', 17, FALSE),
('Библиотека для работы с базами данных.', 17, TRUE),
--8
('Фреймворк для создания backend на Python.', 18, FALSE),
('Библиотека для работы с базами данных.', 18, TRUE),
('Инструмент для тестирования кода.', 18, FALSE),
('Система контроля версий.', 18, FALSE),
--9
('Фреймворк для создания backend на Python.', 19, TRUE),
('Библиотека для работы с базами данных.', 19, FALSE),
('Инструмент для тестирования кода.', 19, FALSE),
('Система контроля версий.', 19, FALSE),
--10
('Метод оптимизации производительности сайта.', 20, FALSE),
('Технология создания веб-приложений.', 20, FALSE),
('Способ обработки большого количества запросов от пользователей.', 20, TRUE),
('Способ обеспечения безопасности данных.', 20, FALSE),
--11
('Архитектура backend на Python.', 21, TRUE),
('Архитектура frontend на JavaScript.', 21, FALSE),
('Архитектура базы данных.', 21, FALSE),
('Архитектура системы контроля версий.', 21, FALSE),
--12
('Промежуточное ПО для обработки запросов к базе данных.', 22, TRUE),
('Промежуточное ПО для аутентификации пользователей.', 22, FALSE),
('Промежуточное ПО для кеширования данных.', 22, FALSE),
('Промежуточное ПО для балансировки нагрузки.', 22, FALSE),
--13
('Способ обработки большого количества запросов от пользователей.', 23, TRUE),
('Метод оптимизации производительности сайта.', 23, FALSE),
('Способ обеспечения безопасности данных.', 23, FALSE),
('Способ обработки транзакций в базе данных.', 23, FALSE),
--14
('Запрос, который обрабатывается параллельно с другими запросами.', 24, TRUE),
('Запрос, который выполняется без блокировки основного потока.', 24, FALSE),
('Запрос, который выполняется быстрее синхронного запроса.', 24, FALSE),
('Запрос, который не блокирует основной поток выполнения программы.', 24, FALSE);


INSERT INTO categories (category_id, name, description, photo_url, background_url)
VALUES  
(3, 'МОБИЛЬНАЯ РАЗРАБОТКА', 'изучение разработки приложений для мобильных платформ', 'https://i.pinimg.com/736x/28/ec/e5/28ece574570a14034716faba99e1e92c.jpg', 'https://i.pinimg.com/564x/05/e6/cf/05e6cfb35ef994b22bde2512fe3d3049.jpg');

INSERT INTO courses  
(name, description, min_rating, optimal_rating, category_id, photo_url) 
VALUES
('Разработка мобильных приложений на Flutter', 'Flutter для iOS и Android', 0, 0, 3, 'https://docs.flutter.dev/assets/images/flutter-logo-sharing.png');

INSERT INTO questions (question, course_id)
VALUES
('Что такое Flutter?', 7),
('Что такое Dart?', 7),
('Какие основные компоненты Flutter?', 7),
('Что такое hot reload в Flutter?', 7),
('Какие основные платформы поддерживает Flutter?', 7),
('Что такое state management в Flutter?', 7),
('Что такое виджеты в Flutter?', 7),
('Что такое Material Design в контексте Flutter?', 7),
('Что такое Cupertino в Flutter?', 7),
('Что такое Flutter SDK?', 7),
('Как работает Flutter с платформенными API?', 7),
('Что такое pubspec.yaml?', 7),
('Что такое StatefulWidget и StatelessWidget?', 7),
('Как управлять состоянием в Flutter?', 7);

INSERT INTO answers (answer, question_id, is_correct) VALUES
-- Вопрос 1
('Фреймворк для создания мобильных приложений.', 25, TRUE),
('Язык программирования.', 25, FALSE),
('Инструмент для работы с базами данных.', 25, FALSE),
('Редактор кода.', 25, FALSE),
-- Вопрос 2
('Язык программирования для Flutter.', 26, TRUE),
('Фреймворк для создания веб-приложений.', 26, FALSE),
('Технология работы с базами данных.', 26, FALSE),
('Система контроля версий.', 26, FALSE),
-- Вопрос 3
('Виджеты, state management, платформа.', 27, TRUE),
('Компоненты базы данных.', 27, FALSE),
('Функции для работы с сетью.', 27, FALSE),
('API для работы с камерой.', 27, FALSE),
-- Вопрос 4
('Перезагрузка кода без перезапуска приложения.', 28, TRUE),
('Функция для работы с базами данных.', 28, FALSE),
('Способ ускорить тестирование.', 28, FALSE),
('Инструмент для работы с графикой.', 28, FALSE),
-- Вопрос 5
('iOS и Android.', 29, TRUE),
('Только Android.', 29, FALSE),
('Только iOS.', 29, FALSE),
('Windows и macOS.', 29, FALSE),
-- Вопрос 6
('Механизм управления состоянием виджетов.', 30, TRUE),
('Фреймворк для управления запросами.', 30, FALSE),
('Язык программирования.', 30, FALSE),
('Метод обработки данных.', 30, FALSE),
-- Вопрос 7
('Основные строительные блоки интерфейса.', 31, TRUE),
('Компоненты для работы с сетью.', 31, FALSE),
('Элементы для работы с базами данных.', 31, FALSE),
('Инструменты для управления памятью.', 31, FALSE),
-- Вопрос 8
('Гайдлайны по дизайну от Google.', 32, TRUE),
('Язык программирования.', 32, FALSE),
('Средство для работы с анимацией.', 32, FALSE),
('Фреймворк для тестирования.', 32, FALSE),
-- Вопрос 9
('Набор стилей от Apple для Flutter.', 33, TRUE),
('Набор библиотек для Flutter.', 33, FALSE),
('Язык программирования.', 33, FALSE),
('Инструмент для работы с графикой.', 33, FALSE),
-- Вопрос 10
('Набор инструментов для создания приложений на Flutter.', 34, TRUE),
('Фреймворк для работы с сетью.', 34, FALSE),
('Средство для тестирования приложений.', 34, FALSE),
('Язык программирования для Flutter.', 34, FALSE),
-- Вопрос 11
('Через вызов нативных API через платформенные каналы.', 35, TRUE),
('Использует базовые HTTP-запросы.', 35, FALSE),
('Посредством баз данных.', 35, FALSE),
('Через удалённые API.', 35, FALSE),
-- Вопрос 12
('Конфигурационный файл, содержащий зависимости проекта.', 36, TRUE),
('Файл для работы с базами данных.', 36, FALSE),
('Файл для описания архитектуры приложения.', 36, FALSE),
('Файл для описания структуры UI.', 36, FALSE),
-- Вопрос 13
('StatefulWidget изменяется со временем, StatelessWidget - нет.', 37, TRUE),
('StatefulWidget создаёт виджеты, StatelessWidget - базу данных.', 37, FALSE),
('StatefulWidget используется для работы с сетью.', 37, FALSE),
('StatelessWidget обрабатывает анимации.', 37, FALSE),
-- Вопрос 14
('Использовать state management.', 38, TRUE),
('Использовать систему баз данных.', 38, FALSE),
('Использовать языки программирования.', 38, FALSE),
('Использовать сторонние API.', 38, FALSE);

INSERT INTO courses  
(name, description, min_rating, optimal_rating, category_id, photo_url) 
VALUES
('Основы веб-разработки с использованием React и Redux.', 'front-end разработка', 100, 300, 2, 'https://i.pinimg.com/736x/cb/07/cf/cb07cfa10b0fad38e078dee3784d5e1a.jpg'),
('Веб-разработка на Go: создание эффективных и масштабируемых веб-приложений', 'back-end разработка', 100, 300, 2, 'https://i.pinimg.com/736x/cb/07/cf/cb07cfa10b0fad38e078dee3784d5e1a.jpg');

ALTER SEQUENCE categories_category_id_seq RESTART WITH 4;
ALTER SEQUENCE answers_answer_id_seq RESTART WITH 153;
ALTER SEQUENCE questions_question_id_seq RESTART WITH 39;
ALTER SEQUENCE courses_course_id_seq RESTART WITH 10;