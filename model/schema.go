package model

var schema = `
DROP TYPE IF EXISTS request_type CASCADE; 
CREATE TYPE request_type AS ENUM ('GET', 'POST', 'PUT');

DROP TABLE IF EXISTS "user" CASCADE;
CREATE TABLE "user" (
	id 			SERIAL PRIMARY KEY,
	email	 	VARCHAR(30) NOT NULL,
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
	type 			request_type DEFAULT 'GET',
	path 			TEXT NOT NULL,
	destination 	TEXT NOT NULL,
	body			TEXT,
	created_at		TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

DROP TABLE IF EXISTS query CASCADE;
CREATE TABLE IF NOT EXISTS "query" (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name 			VARCHAR(20) NOT NULL,
	value 			VARCHAR(50) NOT NULL
);
`