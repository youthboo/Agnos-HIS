# Agnos-HIS

## Installation & Setup
1. **Clone the repository**
   - `git clone https://github.com/yourgithub/Agnos-HIS.git`
   - `cd Agnos-HIS`

2. **Setup Environment Variables**  
   Edit the `.env` file in the root directory and modify the following values as needed:

    - ``` DB_HOST=db DB_USER=admin DB_PASSWORD=password DB_NAME=his_db JWT_SECRET=your_secret_key ```

3. **Start the project using Docker**
- `docker-compose up --build`
- 📌 *Note: Ensure that Docker is installed and running.*

## Running Unit Tests
- `go test -v ./tests/`

## Contributors
**Palita Lertsaksrisakul**






