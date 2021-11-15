package model

import "fmt"

var schema = fmt.Sprintf(`
DROP TYPE IF EXISTS request_type CASCADE; 
CREATE TYPE request_type AS ENUM ('GET', 'POST', 'PUT', 'DELETE');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS "user" CASCADE;
CREATE TABLE "user" (
	id 			SERIAL PRIMARY KEY,
	email	 	VARCHAR(30) NOT NULL UNIQUE,
	created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW() 
);

DROP TABLE IF EXISTS project CASCADE;
CREATE TABLE project (
	id  		SERIAL PRIMARY KEY,
	user_id 	INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name 		TEXT NOT NULL UNIQUE,
	created_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

DROP TABLE IF EXISTS route CASCADE;
CREATE TABLE route (
	id  			SERIAL PRIMARY KEY,
	project_id		INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	method 			request_type DEFAULT 'GET',
	path 			TEXT NOT NULL,
	description 	TEXT NOT NULL DEFAULT '',
	destination 	TEXT NOT NULL,
	body			TEXT DEFAULT '',
	created_at		TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

DROP TABLE IF EXISTS query CASCADE;
CREATE TABLE IF NOT EXISTS "query" (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name 			TEXT NOT NULL,
	value 			TEXT NOT NULL
);

DROP TABLE IF EXISTS header CASCADE;
CREATE TABLE IF NOT EXISTS header (
	id 				SERIAL PRIMARY KEY,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	route_id	 	INTEGER REFERENCES route(id) ON DELETE CASCADE,
	name 			TEXT NOT NULL,
	value 			TEXT NOT NULL
);

DROP TABLE IF EXISTS perm_origin CASCADE;
CREATE TABLE IF NOT EXISTS perm_origin (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	url 			TEXT NOT NULL
);

INSERT INTO "user"(email) VALUES ('ikeohachidi@gmail.com');
INSERT INTO project(user_id, name) VALUES (1, 'foo');
INSERT INTO route(project_id, user_id, method, path, description, destination) VALUES (1, 1, 'GET', '/users', 'None existent for now', 'http://localhost:5001');
INSERT INTO query (user_id, route_id, name, value) VALUES (1, 1, pgp_sym_encrypt('private-key-1', '%[1]v'), pgp_sym_encrypt('value-pair', '%[1]v'));
INSERT INTO perm_origin (route_id, url) VALUES (1, 'http://localhost:8080');
`, PG_CRYPT_KEY)
