CREATE TYPE LETTER AS ENUM
    ('A', 'B', 'V');

--Table for currencies
CREATE TABLE Pupils (
    id            SERIAL PRIMARY KEY,
    pupil         VARCHAR NOT NULL,
    establishment VARCHAR NOT NULL,
    class_num     NUMERIC(11) NOT NULL,
	letter        LETTER NOT NULL
);


--Simulating some data in the database
INSERT INTO Pupils (id, pupil, establishment, class_num, letter)
VALUES
    (1, 'Daniel', 'AMTEK', 11, 'A'),
    (2, 'Maria', 'School 26', 10, 'A'),
    (3, 'Zima', 'School 17', 10, 'V');
----------------------