package migrations

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/market/modules/util"
)

func initialSchema(tx *gorm.DB) error {

	group_uuid := util.NewUuid()

	sql := fmt.Sprintf(`
	
CREATE TABLE users
(
    id          TEXT	NOT NULL UNIQUE,
    nid         SERIAL  NOT NULL UNIQUE,
    name        TEXT    NOT NULL,
    lower_name  TEXT    NOT NULL UNIQUE,
    email       TEXT    NOT NULL UNIQUE,
    password    TEXT    NOT NULL,
    last_seen   TIMESTAMP WITH TIME ZONE,
    is_active   BOOLEAN NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT user_pkey PRIMARY KEY (id)
);


CREATE TABLE groups 
(
    nid         SERIAL  NOT NULL,
    id          TEXT    NOT NULL, 
    name        TEXT    NOT NULL,
    lower_name  TEXT    NOT NULL UNIQUE,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT groups_pkey PRIMARY KEY (id)
);
	
CREATE TABLE user_groups 
(
	user_nid INTEGER,
	group_nid INTEGER,

	CONSTRAINT users_groups_pkey PRIMARY KEY (user_nid, group_nid)
);


CREATE TABLE documents
(
    nid         SERIAL  NOT NULL,
    id          TEXT    NOT NULL,
    name        TEXT    NOT NULL,
    lower_name  TEXT    NOT NULL,
	description TEXT,
    mimetype    TEXT    NOT NULL,
    size        INTEGER NOT NULL,
    created_by  INTEGER 
        REFERENCES users(nid) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT documents_pkey PRIMARY KEY (nid)
);


CREATE TABLE metadata_keys
(
    key     TEXT,
    multiple BOOLEAN NOT NULL,
    comment TEXT,
    active  BOOLEAN NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT keys_pkey PRIMARY KEY (key)
);


CREATE TABLE metadatas
(
    id  SERIAL  NOT NULL,  
    key TEXT
        REFERENCES metadata_keys(key) NOT NULL,
    value TEXT  NOT NULL,
    document INTEGER
        REFERENCES documents(nid) NOT NULL,

    CONSTRAINT metadata_pkey PRIMARY KEY (id)
);


CREATE TABLE role_levels 
(
    level       INTEGER NOT NULL,
    name        TEXT    NOT NULL,
    comment     TEXT,
    group_nid   INTEGER,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE,

    CONSTRAINT role_levels_pkey PRIMARY KEY (level),
    CONSTRAINT group_unique_roles UNIQUE (level, group_nid),
    CONSTRAINT group_unique_names UNIQUE (name, group_nid)
);

CREATE TABLE roles 
(
    id          SERIAL  NOT NULL,
    user_nid     INTEGER,
    group_nid    INTEGER,
    role_level   INTEGER 
        REFERENCES role_levels(level) NOT NULL,
    metadata_key INTEGER,
    
    CONSTRAINT roles_pkey PRIMARY KEY (id),
    CONSTRAINT user_or_group UNIQUE (user_nid, group_nid)
);

CREATE TABLE actions
(
    id          SERIAL  NOT NULL,
    user_nid    INTEGER
        REFERENCES users(nid) NOT NULL,
    document    INTEGER
        REFERENCES documents(nid) NOT NULL,
    operation   INTEGER NOT NULL,
    comment     TEXT,
	comment_lower TEXT,
    ts          TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT actions_pkey PRIMARY KEY (id)
);
	
CREATE TABLE sessions
(
	id		SERIAL NOT NULL,
	user_nid INTEGER
		REFERENCES users(nid) NOT NULL,
	nonce	TEXT	NOT NULL,
	user_agent TEXT NOT NULL,
	issued_at	TIMESTAMP WITH TIME ZONE NOT NULL,

	CONSTRAINT sessions_pkey PRIMARY KEY (id),
	CONSTRAINT user_nonce_unique UNIQUE (user_nid, nonce)
);
	
INSERT INTO groups VALUES (1, '%s', 'default group', 'default group', TIMENOW(), TIMENOW());
INSERT INTO role_levels VALUES (0, 'none', 'No permission', 1, TIMENOW(), TIMENOW());
INSERT INTO role_levels VALUES (100, 'admin', 'Site administrator', 1, TIMENOW(), TIMENOW());
`, group_uuid)
	return tx.Exec(sql).Error
}
