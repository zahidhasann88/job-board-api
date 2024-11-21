# Job Board API

The **Job Board API** is a backend service that facilitates job seekers, recruiters, and administrators to manage and interact with job postings, user accounts, and job applications.

## Features
- **User Management**: Register and authenticate users (Job Seekers, Recruiters, and Admins).
- **Job Management**: Create, update, delete, and list job postings.
- **Application Management**: Apply for jobs, view applications, and filter them.
- **Role-Based Access Control**: Protect routes based on user roles.
- **Analytics and Insights**: Gain insights into job applications and recommendations.

---

## Requirements

- **Go**: `>=1.18`
- **Database**: PostgreSQL
- **Environment Variables**: Defined in `.env` file.

---

## Installation and Setup

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd job-board-api
   ```

2. **Set Up Environment Variables**
   Create a `.env` file in the root directory and provide the following values:
   ```env
   DATABASE_URL=postgres://<username>:<password>@localhost:5432/job_board?sslmode=disable
   PORT=8080
   JWT_SECRET=your-secret-key
   FILE_STORAGE_PATH=./uploads
   ```

3. **Run Database Migrations**
   Ensure the following tables are set up in your PostgreSQL database:
   - `users`
   - `jobs`
   - `applications`

   Use the provided SQL scripts in the `migrations/` folder.

4. **Start the Server**
   ```bash
   go run cmd/api/main.go
   ```

---

## API Endpoints

### Public Endpoints
| Method | Endpoint                  | Description                      |
|--------|---------------------------|----------------------------------|
| POST   | `/api/v1/users/register`  | Register a new user.            |
| POST   | `/api/v1/users/login`     | Authenticate and get a JWT.     |
| GET    | `/api/v1/jobs`            | List available jobs.            |
| GET    | `/api/v1/jobs/:id`        | Get details of a specific job.  |

### Protected Endpoints (Requires Authentication)
#### Recruiter
| Method | Endpoint                             | Description                                 |
|--------|--------------------------------------|---------------------------------------------|
| POST   | `/api/v1/jobs`                       | Create a new job.                          |
| PUT    | `/api/v1/jobs/:id`                   | Update an existing job.                    |
| PATCH  | `/api/v1/jobs/:id/status`            | Change the status of a job.                |
| DELETE | `/api/v1/jobs/:id`                   | Delete a job.                              |
| GET    | `/api/v1/jobs/analytics`             | Get analytics for jobs.                    |
| POST   | `/api/v1/jobs/bulk`                  | Bulk create job postings.                  |
| GET    | `/api/v1/jobs/:id/application-insights` | View insights for job applications.     |
| GET    | `/api/v1/jobs/:id/recommended-candidates` | View recommended candidates.          |

#### Job Seeker
| Method | Endpoint                | Description               |
|--------|-------------------------|---------------------------|
| POST   | `/api/v1/applications`  | Submit a job application. |
| GET    | `/api/v1/applications`  | List job applications.    |

#### Common
| Method | Endpoint                  | Description            |
|--------|---------------------------|------------------------|
| PUT    | `/api/v1/users/profile`   | Update user profile.   |

---

## Configuration

- **Configuration File**: `config/config.go`
- Defaults:
  - Database URL: `postgres://username:password@localhost:5432/job_board?sslmode=disable`
  - Port: `8080`
  - JWT Secret: `your-secret-key`
  - File Storage Path: `./uploads`

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    resume_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Jobs Table
```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    company_id UUID NOT NULL,
    location VARCHAR(255) NOT NULL,
    salary_range VARCHAR(255),
    job_type VARCHAR(50) NOT NULL,
    experience_level VARCHAR(50) NOT NULL,
    skills TEXT[] NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Applications Table
```sql
CREATE TABLE applications (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL,
    applicant_id UUID NOT NULL,
    cover_letter TEXT,
    resume_url VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Testing
- Use **Postman** or **cURL** to test API endpoints.
- JWT tokens are required for protected routes.

---

## Contribution
1. Fork the repository.
2. Create a new branch.
3. Submit a pull request.

---

## License
MIT License.