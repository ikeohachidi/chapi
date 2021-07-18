package models

var schema = `
CREATE TABLE project (
	id  		SERIAL PRIMARY KEY ON DELETE CASCADE,
	name 		TEXT NOT NULL UNIQUE,
	user_id 	INTEGER NOT NULL,
	created_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)

CREATE TABLE routes (
	id  			SERIAL PRIMARY KEY,
	path 			TEXT NOT NULL,
	project_id		INTEGER NOT NULL,
	request_body	TEXT,
	created_at	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
) 
`
