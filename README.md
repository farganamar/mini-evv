# Mini-EVV Service

A lightweight Electronic Visit Verification (EVV) system with appointment management for caregivers.

## Tech Stack

This application uses the following technologies:

- **Go**: Backend language (v1.23)
- **SQLite**: Local database for development and testing
- **Redis**: Cache for user sessions and data
- **Docker**: Containerization for consistent deployment
- **Chi Router**: Lightweight HTTP routing
- **Swagger**: API documentation

## Key Design Decisions

1. **Embedded SQLite Database**: Chosen for simplicity and portability, allowing the application to run without external database dependencies.

2. **Clean Architecture**: Organized in layers (handlers → services → repositories) for clean separation of concerns.

3. **Stateless Authentication**: JWT-based authentication with refresh tokens stored in Redis for security and scalability.

4. **Location Validation**: Appointments include location verification with configurable tolerance to ensure caregivers are at the correct client location.

5. **Concurrent Data Fetching**: Optimized response times for appointment listing by implementing parallel client data retrieval.

6. **Database Schema**: I used 4 tables: 
 - `Appointments` : to capture all appointment for between the caregivers and client (patient)
 - `Users` : all internal user admin/caregivers stored in here separated by column `roles`
 - `Client` :all client information in here includes: `latitude`, `longitude`
 - `Appointments_log`: to capture appointment log such as: `Checkin` , `Checkout`, and `Activity report`

## Brief of assumptions of the feature

As a caregivers you will be able to see all appointments task by the Status, the caregivers will have 4 statuses `SCHEDULED`, `COMPLETED`, `CANCELLED` and `IN_PROGRESS`. 

As a caregivers when you arrived at the client location, you will be able to do check-in and give report Activity before you leave the location (check-out). Every Activity will be recorded on table `appointment_logs` including your device, ip address, latitude, and longitude, this is part of EVV and audit log to make sure everything is provable because we are handling the patient healthcare.

There is some multi layer Validation before you do the checkin:
- make sure you are still in the radius between client location, i put fault tolerance to maximum 600meters
- you are need to put the verification code from the client (you can bypass this to input `0000`)

**Because we are gonna record your realtime location, make sure you're running the seeder appointment data first using this API or click button at the dashboard**

Using the API: 
```curl -X POST http://localhost:3200/v1/evv/seed/appointment \
  -H "Content-Type: application/json" \
  -d '{"latitude": *YOUR_LATITUDE*, "longitude": *YOUR_LONGITUDE*}'
```

Using the button:

<img width="307" alt="image" src="https://github.com/user-attachments/assets/c8950c69-5ed2-4228-b777-125ab808980b" />


## Setup Instructions

Clone this repository
```bash
    git clone https://github.com/farganamar/mini-evv.git
    cd mini-evv
```
Before running the application, you need to set up the environment variables:

```bash
# Copy the example environment file and modify as needed
cp .env.example .env
```
### Initial Configuration
Using Docker (Recommended)

1. Make sure Docker and Docker Compose are installed on your system.
2. Start the application
```
docker-compose up -d
```
3. The application will be available at http://localhost:3200

4. Access the API documentation at http://localhost:3200/swagger/index.html


Using Makefile (Local Development)
1. Make sure Go 1.23+ is installed on your system.

2. Run this command to start the application
```
make dev
```

3. The application will be available at http://localhost:3200

### Additional Features in the future
1. Caregiver: 
- Allows caregivers to decline the appointment, if only the appointment still to far from the scheduled day
- Add attachment to record appointment log, ensure the consistency and secure the data
- Implement EMR (Electronic Medical Record) if possible or if the activities are related with the healthcare
- Add Validation on checkin if only caregivers come early before the scheduled
- Allows caregivers to rate the client
- Add emergency button that already integrated with 911 or related public institution if only something happened to caregivers or client
- integrated with 3rd party calendar (google calendar, apple calender, outlook) so system can make reminder before their scheduled
- Caregiver settings: setup available time, 
- Caregiver milestone: show metrics like review or how much client satisfaction 
- if we develop client side for mobile: we can add automatically checkin if within the client radius

2. Client
- Allows multiple caregivers request at one time book
- Allows to rate caregivers services
- Track caregivers location

3. Tech stack
- migrate to postgresql 
- setup CI/CD to smooth deployment




