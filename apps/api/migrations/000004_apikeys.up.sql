ALTER TABLE
    apikeys
ADD
    COLUMN apikey_id UUID NOT NULL;

ALTER TABLE
    apikeys DROP CONSTRAINT pk_apikeys;

ALTER TABLE
    apikeys
ADD
    CONSTRAINT pk_apikeys PRIMARY KEY (apikey_id);

ALTER TABLE
    apikeys
ADD
    CONSTRAINT uq_apikeys_key_hash UNIQUE (key_hash);
