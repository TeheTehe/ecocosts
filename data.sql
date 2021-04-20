-- CREATE USER AND GRANT PERMISSIONS
CREATE USER postgres WITH PASSWORD '';
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON client TO postgres;
GRANT ALL PRIVILEGES ON budget TO postgres;
GRANT ALL PRIVILEGES ON transaction TO postgres;
GRANT ALL PRIVILEGES ON category TO postgres;
GRANT ALL PRIVILEGES ON stock TO postgres;

-- PSEUDO CLIENTS
INSERT INTO client VALUES (
	DEFAULT, 'john', '$2a$10$UEBUVEVqVIvmJHTcA.J7gOrtcjPhLIUCZwLfpup4ctnmR6GhLi0tC'
);

-- DEFAULT CATEGORIES
INSERT INTO category VALUES (
	DEFAULT, 'Food'
);

INSERT INTO category VALUES (
	DEFAULT, 'Living'
);

INSERT INTO category VALUES (
	DEFAULT, 'Grocer'
);

INSERT INTO category VALUES (
	DEFAULT, 'Travel'
);

INSERT INTO category VALUES (
	DEFAULT, 'Utility'
);

-- PSEUDO BUDGET
INSERT INTO budget VALUES (
	1, 1, 799.99
);

INSERT INTO budget VALUES (
	1, 2, 200.00
);

-- PSEUDO TRANSACTION
INSERT INTO transaction VALUES (
	DEFAULT, 1, NULL, -12.99, -12.99, 'Amazon: E-Book',
	'2020-12-31 05:15:30'::timestamp
);

INSERT INTO transaction VALUES (
	DEFAULT, 1, NULL, 1300.00, 1287.01, 'Income: Doe LLC',
	'2021-01-05 13:45:30'::timestamp
);

INSERT INTO transaction VALUES (
	DEFAULT, 1, NULL, -4.59, 1282.42, 'McDonalds: Meal',
	'now'::timestamp - '3 day'::interval
);

-- PSEUDO STOCK
INSERT INTO stock VALUES (
	1, 'INTC', 3
);

INSERT INTO stock VALUES (
	1, 'TSLA', 1
);

INSERT INTO stock VALUES (
	1, 'BA', 5
);
