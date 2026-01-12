CREATE TABLE Users (
		username TEXT NOT NULL PRIMARY KEY,
		photo_id INTEGER DEFAULT NULL
	);
CREATE TABLE Users_photos (
		photo_id INTEGER PRIMARY KEY AUTOINCREMENT,
		photo_data BLOB NOT NULL
	);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE Groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		conversation_id INTEGER UNIQUE,
		photo_id INTEGER DEFAULT NULL);
CREATE TABLE Groups_memberships (
		group_id INTEGER NOT NULL,
		member_username TEXT NOT NULL,
		PRIMARY KEY (group_id, member_username));
CREATE TABLE Conversations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT CHECK( type IN ('private','group') ) NOT NULL
	);
CREATE TABLE Private_conversations_memberships (
		conversation_id INTEGER NOT NULL,
		member_username TEXT NOT NULL,
		PRIMARY KEY (conversation_id, member_username));
CREATE TABLE Messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		conversation_id INTEGER NOT NULL,
		sender_username TEXT NOT NULL,
		content TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		attachment BLOB DEFAULT NULL,
		reacted_to_message_id INTEGER DEFAULT NULL,
		status TEXT CHECK( status IN ('sent','delivered','received','hidden') ), 
		CHECK (content IS NOT NULL OR attachment IS NOT NULL)
	);
