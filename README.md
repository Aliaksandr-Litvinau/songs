
# Songs Library API

This project built using **Gin Gonic** (a Go web framework) and **Gorm** (a Go ORM). The API allows you to manage a library of songs, including adding songs, retrieving song data, and fetching song lyrics with pagination by verses.

## Features

- **Add Songs**: Create a new song in the database.
- **Get Songs**: Retrieve all songs with options for filtering and pagination.
- **Get Lyrics by Verses**: Fetch song lyrics with pagination, splitting the song text into verses.

## Technologies

- **Go (Golang)**: The core programming language.
- **Gin Gonic**: Web framework for building APIs.
- **Gorm**: ORM for working with the PostgreSQL database.
- **PostgreSQL**: The database used for storing song information.
- **Docker**: For containerizing the application.
- **Logrus**: Logging library used for detailed logs.

## Getting Started

### Prerequisites

- Docker
- Go (if running outside Docker)

### Installation

1. **Clone the repository:**

```bash
git clone git@github.com:Aliaksandr-Litvinau/songs.git
cd songs
```

2. **Set up environment variables:**

Create a `.env` file in the root directory of the project, with the following contents:

```bash
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/music_library
```

3. **Run the application with Docker:**

This project uses Docker for containerization. The `docker-compose.yml` file provided will set up the database and the API.

```bash
docker-compose up --build
```

> After the containers are built and running, the API and Swagger documentation will be accessible.
Swagger Documentation:
Visit ```/swagger/index.html``` to explore and interact with the API.

4. **Database Migration:**

The project uses Gorm’s `AutoMigrate` feature to automatically create the `songs` table in the PostgreSQL database. The migration will be applied automatically when the application is started.

### 🚨Why Clean Architecture Wasn't Used from the Start
> **At the start of this project, the priority was rapid development.** Implementing Clean Architecture from day one would have introduced unnecessary complexity and slowed down the initial build. When you’re trying to hit deadlines, especially for a smaller project with no immediate plans for scaling, adding layers of abstraction can be overkill. Clean Architecture requires more planning, design, and development time, which can be counterproductive when the focus is on getting a working product out the door quickly.

> That being said, as the project grows and its complexity increases, transitioning to Clean Architecture will be a natural progression. This will allow the codebase to handle scaling, new features, and long-term maintainability more effectively.

> In summary, while we started with a pragmatic, simpler approach to get results fast, the project is structured in a way that makes it possible to evolve into Clean Architecture when the time comes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
