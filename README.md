# Appointment Scheduling

## Setup

1. To load the database with the initial data, run the following command:

```bash
go run utils/load-database.go
```

2. To run the server, run the following command:
```bash
go run main.go
```

## Available endpoints

### Get available appointments
```bash
curl -X GET 'http://localhost:8080/appointments/available?trainer_id=1&start_date=2025-01-01&end_date=2025-01-07'
```

### Create an appointment
```bash
curl -X POST http://localhost:8080/appointments -H "Content-Type: application/json" -d '{"trainer_id": 1, "user_id": 1, "starts_at": "2019-01-25T09:00:00-08:00", "ends_at": "2019-01-25T09:30:00-08:00"}'
```


### Get scheduled appointments
```bash
curl -X GET 'http://localhost:8080/appointments/scheduled?trainer_id=1'
```

## What I wish I would have had time to do

- Testing
- Add Multiple Timezone Support
  - Here I would make sure that the data in the database is stored in UTC or even unix to ensure consistency.
- Improve database interactions
  - Ensure CRUD properties
  - If two insert calls where made at the same time, there could be a race condition with appointment ID currently because of the way I was 
    incrementing an id
  - Since id’s were incrementing in sample data I followed suit but would be nice to have uuid support instead, or use auto increment sql column 
    feature
- Add documentation for the API
  - Looking back I should have started with this to lay the groundwork for the project
  - Would be nice to provide the user clean docs on the use of each endpoint



