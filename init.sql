CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	patronymic text,
    first_name text,
    last_name text,
    email text,
	password text
);
CREATE TABLE IF NOT EXISTS services (
	id SERIAL PRIMARY KEY,
	service_name text
);
CREATE TABLE IF NOT EXISTS access (
	accessid SERIAL PRIMARY KEY,
    user_id INT,
    service_id INT,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (service_id) REFERENCES Services(id)
);
CREATE TABLE IF NOT EXISTS jwt_blacklist (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL
);
INSERT  INTO users (patronymic, first_name, last_name, email, password) VALUES ('Иванов', 'Иван', 'Иванович', 'ivany@bb.ru', '123superstrong') ;
INSERT  INTO users (patronymic, first_name, last_name, email, password) VALUES ('Петров', 'Алексей', 'Владимирович', 'pav@mm.ru', 'pawpaw123') ; 
INSERT  INTO services (service_name) VALUES ('shop');
INSERT  INTO services (service_name) VALUES ('admin_panel');
INSERT  INTO access (user_id, service_id) VALUES (1, 1);
INSERT  INTO access (user_id, service_id) VALUES (1, 2);
INSERT  INTO access (user_id, service_id) VALUES (2, 1)