CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

insert into 
users(id, name, username, password) 
VALUES(1, "admin", "admin", "$2a$04$DwVMG5iM.ov0xQFEgH3wZu/hGXCgSIPVjWZoGj0DpVDEVaiK9Fw66");