CREATE TABLE store(
    id int auto_increment primary key,
    name varchar(255) not null,
    description text,
    location varchar(255) not null,
    latitude decimal(10, 8),
    longtitude decimal(10, 2),
    phone_number varchar(255) not null,
    email varchar(255) not null,
    image_url varchar(255),
    whatsapp_link varchar(255) not null,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;