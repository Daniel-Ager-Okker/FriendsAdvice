CREATE TYPE SUBJECT AS ENUM
    ('Russian', 'Mathematics', 'Physics', 'Literature', 'English', 'History', 'Technology');

CREATE TYPE TYPETEST AS ENUM
    ('Annual test', 'Quarter test', 'Test', 'Independent work', 'Answer to the board');

--Table for currencies
CREATE TABLE Grades (
    id              SERIAL PRIMARY KEY,
    pupil           VARCHAR NOT NULL,
    establishment   VARCHAR NOT NULL,
	subj            SUBJECT NOT NULL,
	knowlegde_check TYPETEST NOT NULL,
	grade           NUMERIC(5) NOT NULL
);



--Simulating some data in the database
INSERT INTO Grades (id, pupil, establishment, subj, knowlegde_check, grade)
VALUES
    (1, 'Daniel', 'AMTEK', 'Physics', 'Test', 5),
    (2, 'Maria', 'School 26', 'English', 'Independent work', 4),
    (3, 'Zima', 'School 17', 'Mathematics', 'Quarter test', 2);
----------------------