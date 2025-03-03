CREATE TABLE IF NOT EXISTS movies (
    id bigint NOT NULL AUTO_INCREMENT,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title varchar(255) NOT NULL,
    `year` int NOT NULL,
    PRIMARY KEY (id)
);
