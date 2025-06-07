## 1. Install Docker

Make sure Docker is installed on your system. Follow the installation instructions for your operating system:

- [Install Docker](https://docs.docker.com/get-docker/)
- [Install Docker Compose](https://docs.docker.com/compose/install/)

## **How the Application Works**

### **Overview**
This application is a Go-based backend that interacts with a MySQL database. It is designed with a **multi-service architecture** that includes:
1. A **ProjectService** to manage project data.
2. A **TaskService** to handle tasks within those projects.

The application uses **Docker** for containerization, **MySQL** as the database, and **Alpine** images to ensure a lightweight and efficient deployment. The Go application runs on port **8000**.

### **Architecture**
1. **Backend**:
   - Built using **Go** (version 1.23) and follows a **layered architecture**.
   - The **ProjectService** handles business logic for projects, including creating, reading, updating, and deleting projects.
   - The **TaskService** handles similar logic for tasks within those projects.
   - The services communicate with the database using **MySQL** and GORM ORM for database operations.

2. **Database**:
   - **MySQL** is used for storing project and task data.
   - The **MySQL container** is managed by Docker and configured through environment variables in the `.env` file.
   - To add new tables or remove tables, ensure you add schema's for these on the `./schema` folder.
   - To start of the MySQL service seperately"
   ```bash
      docker compose up MySQL
   ```
   - To exec into the MySQL containers
   ```bash
         # Open a shell inside the container
      docker exec -it hearx-mysql-1 bash  # or 'sh' if bash not available

      # From within container, launch MySQL client
      mysql -u user -puserpassword project_db
   ```
   - Cleaning up the containers
   ```bash
      # Stop containers (keep data)
      docker compose down 

      # Stop & remove containers + volumes:
      docker compose down -v
   ```

3. **Containerization**:
   - The **Dockerfile** defines a **multi-stage build**:
     - **Stage 1**: Builds the Go application.
     - **Stage 2**: Uses a minimal `Alpine` image to run the Go application in a production-ready environment, exposing it on port `8000`.

4. **Migration**:
   - **Database migrations** can be handled by the Go application or manually using tools like `golang-migrate`. Migrations are applied on startup via the `entrypoint.sh` script before the Go service is started.

### **Running the Application**

1. **Docker Setup**:
   - The application uses **Docker Compose** to spin up both the backend service and MySQL container, ensuring easy management of dependencies.
   - **Build and start the services** with:

     ```bash
     docker-compose up --build
     ```

2. **Endpoints**:
   - The application exposes several API endpoints to manage projects and tasks:
     - `GET /projects`: Retrieves all projects.
     - `POST /projects`: Creates a new project.
     - `GET /projects/{id}`: Retrieves a specific project by ID.
     - `PUT /projects/{id}`: Updates a specific project.
     - `DELETE /projects/{id}`: Deletes a specific project.
     - Similar endpoints exist for tasks.

3. **Testing**:
   - You can test the application using tools like `Postman` or `curl` by making HTTP requests to `http://localhost:8000` for interacting with the project and task management API.
      - Example
      ```bash
         curl -X POST http://localhost:8000/projects \
         -H "Content-Type: application/json" \
         -d '{
         "title": "New Project",
         "description": "This is a description for the new project",
         "deadline": "2025-12-31"
         }'

         {
            "ID": 1,
            "Title": "New Project",
            "Description": "This is a description for the new project",
            "Deadline": "2025-12-31"
         }
      ```
   - To access the swagger page with the endpoints used on the app simply visit `http://localhost:8000/swagger/index.html`.
   - Accessing the database
   ```bash
   docker exec -it mysql_container mysql -u user -p
   ```

4. **Stopping the Application**:
   - To stop the containers, run:

     ```bash
     docker-compose down
     ```

