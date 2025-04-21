# GreenUni

GreenUni is structured into two main parts: the frontend and the backend.


## Prerequisites 

-- Mysql server
-- Maybe smt else?


## Frontend

Built with React Native using Expo.

### Getting Started

To run the app locally:

1. Open your terminal.  
2. Navigate to the frontend directory:
    ```bash
    cd Frontend
    ```
3. Start the Expo development server:
    ```bash
    npx expo start
    ```
4. A QR code will appear in your terminal or browser. You can:
   - Scan it with the Expo Go app on your phone (recommended).  
   - Or run it in your browser (not recommended).

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