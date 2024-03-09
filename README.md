# Contact Service

The contact service is responsible for managing the contacts of the users. It allows the users to add, remove, and
update their contacts. The service also provides the ability to search for contacts.

- Notes:
  - The service supports the following operations:
    - Create a contact for an existing user
    - Update a contact of an existing user
    - Get a contact of an existing user
    - Search contacts of an existing user
    - Delete a contact of an existing user
  - The service supports logging by stdout. In production, the logs would be sent to a log aggregator like Datadog.
  - The service uses the `myerror` package to wrap errors and handle different HTTP status codes.
  - Due to the lack of cloud resources, the service uses only in-memory resources. This means that the data will be lost
    when the service is restarted.
  - In order to handle high scale, the service would use a document-based database like MongoDB and a distributed cache like Redis.   
  - User management is out of the scope of the service.
  - The service currently does not support user authentication and authorization.


- ⭐ Bonuses 
  - As a bonus feature, the service provides a search endpoint that allows the user to search for contacts by their first name,
    last name, address, and phone number. The search endpoint also supports pagination.
  - Another bonus feature that wasn't implemented is the ability to upload a contact's photo, and get all photos of a contact that were uploaded by other users. This can be implemented by adding another microservice (image-service) to the system, and using a message broker like Kafka to communicate between the services. In addition, the service would use an object-store like AWS S3 to store the images.
  
---

## Running the service

1. Clone the repository:

```bash
git clone https://github.com/talgro/phone-book.git
```

2. Change directory to the service
3. Build the docker image:

```bash
docker build -t contact-service -f <path-to-Dockerfile> .
```

4. Run the docker image:

```bash
docker run -p 8080:8080 --name contact-service contact-service
```

5. The service will be available at `http://localhost:8080`

---

## API Requests and Responses

Requests and responses are in JSON format.
All responses are in the following format:

_Success Response 200_

```json
{
  "data": {
    // Response data
  }
}
```

_Error Responses_

```json
{
  "error": "<Error message>"
}
```

---

## API Endpoints

### Create a contact for an existing user

```http
POST /users/:userID/contacts
```

#### Request Body

| Field     | Type   | Comment                |
|-----------|--------|------------------------|
| firstName | string | mandatory              |
| lastName  | string | mandatory              |
| address   | string | mandatory              |
| phone     | string | mandatory, digits only |

###### Example Request

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "address": "123 Main St, Springfield, IL 62701",
  "phone": "5555555555"
}
```

#### Response

Success Response 200

| Field | Type   |
|-------|--------|
| id    | string |

###### Example

```json
{
  "id": "a1b2c3d4"
}
```

-----

### Update a contact of an existing user

```http
PUT /users/:userID/contacts/:contactID
```

#### Request Body

| Field     | Type   | Comment                                                                                |
|-----------|--------|----------------------------------------------------------------------------------------|
| firstName | string | optional                                                                               |
| lastName  | string | optional                                                                               |
| address   | string | optional                                                                               |
| phone     | string | optional, digits only                                                                  |
| updatedAt | string | mandatory, used to validate that the contact hasn't been updated since the client read |

###### Example Request

```json
{
  "firstName": "John",
  "lastName": "Doe",
  "address": "123 Main St, Springfield, IL 62701",
  "phone": "5555555555",
  "updatedAt": "2024-03-09T18:56:05.814330211Z"
}
```

#### Response

Success Response 200 - No content

---

### Get a contact of an existing user

```http
GET /users/:userID/contacts/:contactID
```

#### Response

Success Response 200

| Field     | Type   |
|-----------|--------|
| id        | string |
| firstName | string |
| lastName  | string |
| address   | string |
| phone     | string |
| updatedAt | string |
| createdAt | string |

###### Example

```json
{
  "data": {
    "id": "a1b2c3d4",
    "firstName": "John",
    "lastName": "Doe",
    "address": "123 Main St, Springfield, IL 62701",
    "phone": "5555555555",
    "updatedAt": "2024-03-09T18:56:05.814330211Z",
    "createdAt": "2024-03-09T18:56:05.814330211Z"
  }
}
```

---

### Search contacts of an existing user

```http
GET /users/:userID/contacts
```

#### Query Parameters

| Field     | Type                   |
|-----------|------------------------|
| firstName | string                 |
| lastName  | string                 |
| address   | string                 |
| phone     | string                 |
| limit     | integer between [0,10] |
| offset    | non-negative integer   |

#### Response

Success Response 200

| Field                 | Type           | Comment                  |
|-----------------------|----------------|--------------------------|
| contacts              | list of object |                          |
| contacts[i].id        | string         |                          |
| contacts[i].firstName | string         |                          |
| contacts[i].lastName  | string         |                          |
| contacts[i].address   | string         |                          |
| contacts[i].phone     | string         |                          |
| contacts[i].updatedAt | string         |                          |
| contacts[i].createdAt | string         |                          |
| pagination            | object         |                          |
| pagination.previous   | string         | URL of the previous page |
| pagination.next       | string         | URL of the next page     |

###### Example

```json
{
  "data": {
    "contacts": [
      {
        "id": "a1b2c3d4",
        "firstName": "John",
        "lastName": "Doe",
        "address": "123 Main St, Springfield, IL 62701",
        "phone": "5555555555",
        "updatedAt": "2024-03-09T18:56:05.814330211Z",
        "createdAt": "2024-03-09T18:56:05.814330211Z"
      }
    ],
    "pagination": {
      "previous": "",
      "next": "users/203012323/contacts?phone=&firstName=John&lastName=Doe&limit=2&offset=2"
    }
  }
}
```

---

### Delete a contact of an existing user

```http
DELETE /users/:userID/contacts/:contactID
```

#### Response

Success Response 200 - No content

---
