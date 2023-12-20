BEGIN TRANSACTION;

CREATE TABLE gauge_metrics(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(200) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    updated_at BIGINT NOT NULL
);

CREATE TABLE counter_metrics(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(200) NOT NULL,
    value BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

COMMIT;