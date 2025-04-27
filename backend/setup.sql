-- User stuff

CREATE TABLE IF NOT EXISTS UserTable(
    uuid VARCHAR(36) PRIMARY KEY,
    username VARCHAR(60) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_pass VARCHAR(60) NOT NULL,
    salt VARCHAR(50) NOT NULL,
    role ENUM ('Student', 'Recruiter', 'Admin')
);

INSERT INTO UserTable(uuid, username, email, hashed_pass, salt, role) VALUES ('f52a683e-1896-489b-8c78-d8f03b735e7a', 'admin', 'admin@gmail.com', '$2a$10$8uzZdQR66ZEoWa.iPF7n7ektlX2E.0LG8Fo.QFGz2cNXQ6NZ.L8WK', 'SUkgJGkyqbnHDJV3pwWiTg', 'Admin');

CREATE TABLE IF NOT EXISTS TagsTable (
     id INT AUTO_INCREMENT PRIMARY KEY,
     tagName VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS UserTagsLiked(
    uuid VARCHAR(36),
    tagID int,
    PRIMARY KEY (uuid, tagID),
    FOREIGN KEY (uuid) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (tagID) REFERENCES TagsTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS UserTagsDisLiked(
    uuid VARCHAR(36),
    tagID int,
    PRIMARY KEY (uuid, tagID),
    FOREIGN KEY (uuid) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (tagID) REFERENCES TagsTable(id) ON DELETE CASCADE
);


-- Recruiter

CREATE TABLE IF NOT EXISTS RecruiterTable(
    uuid VARCHAR(36) PRIMARY KEY,
    organisationName VARCHAR (100),
    applicationStatus BOOL DEFAULT FALSE,
    FOREIGN KEY (uuid) REFERENCES UserTable (uuid) ON DELETE CASCADE
);

-- Points stuff
CREATE TABLE IF NOT EXISTS StudentTable (uuid VARCHAR(36) PRIMARY KEY, points int NOT NULL);

-- Student table
CREATE TABLE IF NOT EXISTS StudentInfoTable(
    uuid VARCHAR(36) PRIMARY KEY,
    description TEXT,
    profile TEXT,
    FOREIGN KEY (uuid) REFERENCES UserTable(uuid) ON DELETE CASCADE
);


-- Opportunity stuff
CREATE TABLE IF NOT EXISTS OpportunitiesTable (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    points INT NOT NULL,
    location VARCHAR(100),
    opportunityType ENUM('event', 'volunteer', 'job', 'issue') NOT NULL,
    postedByUUID VARCHAR(36) NOT NULL,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    approved BOOL DEFAULT FALSE,
    FOREIGN KEY (postedByUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS OpportunityLikesTable (
    userUUID VARCHAR(36),
    opportunityUUID VARCHAR(36),
    likedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userUUID, opportunityUUID),
    FOREIGN KEY (userUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS OpportunityDislikesTable (
                                                        userUUID VARCHAR(36),
    opportunityUUID VARCHAR(36),
    dislikedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (userUUID, opportunityUUID),
    FOREIGN KEY (userUUID) REFERENCES UserTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS OpportunityTagsTable (
    opportunityUUID VARCHAR(36),
    tagID INT,
    PRIMARY KEY (opportunityUUID, tagID),
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE,
    FOREIGN KEY (tagID) REFERENCES TagsTable(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS OpportunityMediaTable (
    id INT AUTO_INCREMENT PRIMARY KEY,
    opportunityUUID VARCHAR(36),
    mediaURL TEXT NOT NULL,
    mediaType VARCHAR(50) NOT NULL,
    FOREIGN KEY (opportunityUUID) REFERENCES OpportunitiesTable(uuid) ON DELETE CASCADE
);