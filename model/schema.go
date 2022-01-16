package model

var schema = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS "user" (
	id 			SERIAL PRIMARY KEY,
	email	 	VARCHAR(30) NOT NULL UNIQUE,
	created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW() 
);

CREATE TABLE IF NOT EXISTS project (
	id  		SERIAL PRIMARY KEY,
	user_id 	INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name 		TEXT NOT NULL UNIQUE,
	created_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS route (
	id  			SERIAL PRIMARY KEY,
	project_id		INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	method 			TEXT NOT NULL DEFAULT 'GET',
	path 			TEXT NOT NULL,
	description 	TEXT NOT NULL DEFAULT '',
	destination 	TEXT NOT NULL,
	body			TEXT DEFAULT '',
	created_at		TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "query" (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	name 			TEXT NOT NULL,
	value 			TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS header (
	id 				SERIAL PRIMARY KEY,
	user_id 		INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
	route_id	 	INTEGER REFERENCES route(id) ON DELETE CASCADE,
	name 			TEXT NOT NULL,
	value 			TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS perm_origin (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	url 			TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS request_config (
	id 				SERIAL PRIMARY KEY,
	route_id		INTEGER REFERENCES route(id) ON DELETE CASCADE,
	merge_header	BOOL NOT NULL DEFAULT false,
	merge_body		BOOL NOT NULL DEFAULT false
)
`
