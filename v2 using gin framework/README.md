# Lib-Managment-System
A basic library managment system created using the gin framework on golang

The sql query for creating the database that is used in this project:
CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255) NOT NULL
);

