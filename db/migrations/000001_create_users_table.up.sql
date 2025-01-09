CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

