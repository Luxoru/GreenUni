# GreenUni

GreenUni is structured into two main parts: the frontend and the backend.


## Prerequisites

-- Mysql server
-- neo4j server
-- Expo go app for android/ios (tested on android only)
-- Postman recommended

# Setup

### Backend

1. Make sure you have a mysql server and a neo4j server running
2. Change directory `backend/`
3. Run ` mysql -u root -p mydatabase < setup.sql` to setup all the tables etc
4. Open `cmd/main.go`
5. Change the authenticationConfig to what is needed to connect to your databases
6. In the terminal run `go run main.go`. This will run the backend server on port 8080
7. Get the network ip for your pc. Windows is `ipconfig`

### Frontend

1. change directory to `frontend/utils`
2. open config.js and change the `apiURL` to your network ip for your pc with the port still being 8080
3. open the terminal and run `npx expo start`
4. scan the qr code using expo go app to launch app


## Important Notes

- When a recruiter (not student) has created their account its in a pending state. To verify the account follow these steps:
   1. http://{yourIP}/api/v1/auth/login
      `
       Method: "POST",
       body: {
          "username": "admin",
          "password": "password"
      }`
   3. This will return info.token save this somwhere
   4. Grab the recruiter uuid by doing `http://localhost:8080/api/v1/user/username/{recruiterName}`
   5. Grab opportunity uuid by doing `http://localhost:8080/api/v1/opportunities/author/{recruiterUUID}`
   6. Then u can update the application status by doing `http://localhost:8080/api/v1/opportunities/status?uuid={opportunityUUID}&status=True` Settings the Authorisation header to `Bearer {adminToken}`



## Backend

Written in Go, the backend follows a three-layer architecture:

- Controller (Routing & HTTP)
- Service (Business Logic)
- Repository (Database Access)

### Tech Stack

- Database: MySQL (active)
- Caching: Redis (setup complete, caching logic coming soon)
- Authentication: JWT tokens (token generation implemented, endpoints auth protection planned)



### Repositories

- Represent data stored
- SQL repositories need to implement methods found in repositories/base.go
- RN just putting all queries in code can move later if needed


### Build errors

- If you get table creation errors you might need to just rerun a few times/change the ordering of the table creations.
- This just happens once when the tables are not on your db


## TODO

- Add endpoint/service for saving tag preferences (db already done: UserTagsLiked/UserTagsDisliked (/backend/internal/repositories/user.go))
- Add endpoint/service for saving recruiters. Basically register/update register status (db already done: RecruiterTable (/backend/internal/repositories/user.go))
- Add recommendation system endpoint
- Create a 