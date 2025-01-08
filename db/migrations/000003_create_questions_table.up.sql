CREATE TABLE
    questions (
        id int auto_increment primary key,
        question text not null,
        answer text not null,
        count_helpful int default 0,
        count_unhelpful int default 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB;