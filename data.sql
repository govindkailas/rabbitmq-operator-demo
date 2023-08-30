CREATE TABLE consumer_data(
  	unique_id UUID DEFAULT gen_random_uuid (),
    "data" character varying(255),
  	consumer_name VARCHAR NOT NULL,
  	PRIMARY KEY (unique_id) 
);