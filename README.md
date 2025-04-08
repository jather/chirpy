# Chirpy

Chirpy is a http server that handles requests for users to be able to create accounts and chirps (short messages), as well as read, update and delete them. 

Chirpy is made using go (http package from standard lib), with postgres as database. It has authentification with JWT, and supports queries to filter results. 

## Installation
### Prerequisites
- Go
- Postgresql
### 1. Clone the repository
### 2. Set up the database
```
createdb chirpy
```
### 3. Configure .env file

DB_URL="your database connection string"

PLATFORM="dev"

SECRET="your secret"

POLKA_KEY="polka key"

### 4. Build and run
from root directory of project:
```
go build
./chirpy
```
for windows:
```
go build
chirpy.exe
```
## API Endpoints

### Authentication
- `POST /api/login` - Authenticates a user and returns a JWT (1 hour)
- `POST /api/refresh` - Refreshes an authentication token

### Users
- `POST /api/users` - Creates a new user
- `PUT /api/users` - Updates the current user's details

### Chirps
- `GET /api/chirps` - Retrieves all chirps (can filter by creater of chirp using`?author_id=<id>` and sort the results by most recent to oldest using `?sort=desc`. Default is asc)
- `GET /api/chirps/{id}` - Retrieves a specific chirp
- `POST /api/chirps` - Creates a new chirp
- `DELETE /api/chirps/{id}` - Deletes a specific chirp