-- gosql Up
CREATE TABLE user (
    user_id INT PRIMARY KEY,
    username VARCHAR(50),
    email VARCHAR(100)
);

-- gosql Down
DROP TABLE "user";
