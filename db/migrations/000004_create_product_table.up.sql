CREATE TABLE products(
    id int auto_increment primary key,
    category_id int,
    name varchar(255) not null,
    slug varchar(255) not null unique,
    real_price decimal(10, 2),
    discount decimal(10, 2),
    stock int default 0,
    description text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    foreign key(category_id) references categories(id)
) ENGINE = InnoDB;