```markdown
# Daisi REST Postgres API Documentation

## Authentication

All endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <token>
```

---

## Standard Response Format

All API responses are wrapped in this structure:

```json
{
  "success": true,
  "data": ...,
  "error": "...",    // only present if success is false
  "total": 123       // only present for paginated endpoints
}
```

---

## Endpoints

### Agents

#### List Agents

- **GET** `/api/v1/agents`
- **Query:** `agentids` (optional, comma-separated)

**Response:**
```json
{
  "success": true,
  "data": [ { ...Agent }, ... ]
}
```

#### Get Agent

- **GET** `/api/v1/agents/:agent_id`

**Response:**
```json
{
  "success": true,
  "data": { ...Agent }
}
```

#### Create Agent

- **POST** `/api/v1/agents`
- **Body:** Agent object

**Response:**
```json
{
  "success": true,
  "data": { ...Agent }
}
```

#### Update Agent Name

- **PATCH** `/api/v1/agents/:id`
- **Body:** `{ "agent_name": "New Name" }`

**Response:**
```json
{
  "success": true,
  "data": { ...Agent }
}
```

#### Delete Agent

- **DELETE** `/api/v1/agents/:id`

**Response:** HTTP 204 No Content

---

### Chats

#### List Chats

- **GET** `/api/v1/chats?limit=20&offset=0&...filters`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Chat }, ... ],
  "total": 123
}
```

#### Range Chats

- **GET** `/api/v1/chats/range?start=0&end=9&...filters`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Chat }, ... ]
}
```

#### Search Chats

- **GET** `/api/v1/chats/search?q=term`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Chat }, ... ],
  "total": 10
}
```

---

### Messages

#### List Messages by Chat

- **GET** `/api/v1/messages?agent_id=...&chat_id=...&limit=20&offset=0`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Message }, ... ],
  "total": 123
}
```

#### Range Messages by Chat

- **GET** `/api/v1/messages/range?agent_id=...&chat_id=...&start=0&end=9`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Message }, ... ]
}
```

---

### Contacts

#### List Contacts

- **GET** `/api/v1/contacts?limit=20&offset=0&...filters`

**Response:**
```json
{
  "success": true,
  "data": [ { ...Contact }, ... ],
  "total": 123
}
```

#### Get Contact

- **GET** `/api/v1/contacts/:id`

**Response:**
```json
{
  "success": true,
  "data": { ...Contact }
}
```

#### Update Contact

- **PATCH** `/api/v1/contacts/:id`
- **Body:** `{ "custom_name": "...", "assigned_to": "...", "tags": "..." }`

**Response:**
```json
{
  "success": true,
  "data": { ...Contact }
}
```

---

## Model Examples

### Agent

```json
{
  "agent_id": "string",
  "qr_code": "string",
  "status": "string",
  "agent_name": "string",
  "host_name": "string",
  "version": "string",
  "company_id": "string",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

### Chat

```json
{
  "id": "string",
  "jid": "string",
  "push_name": "string",
  "is_group": true,
  "group_name": "string",
  "unread_count": 0,
  "last_message": { ... }, // object
  "conversation_timestamp": 0,
  "not_spam": false,
  "agent_id": "string",
  "company_id": "string",
  "phone_number": "string",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

### Message

```json
{
  "id": "string",
  "from_user": "string",
  "to_user": "string",
  "chat_id": "string",
  "jid": "string",
  "flow": "string",
  "type": "string",
  "agent_id": "string",
  "company_id": "string",
  "message_obj": { ... }, // object
  "edited_message_obj": { ... }, // object
  "key": { ... }, // object
  "status": "string",
  "is_deleted": false,
  "message_timestamp": 0,
  "message_date": "2024-06-01T12:00:00Z",
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

### Contact

```json
{
  "id": "string",
  "phone_number": "string",
  "agent_id": "string",
  "type": "string",
  "custom_name": "string",
  "notes": "string",
  "tags": "string",
  "company_id": "string",
  "avatar": "string",
  "assigned_to": "string",
  "pob": "string",
  "dob": "2024-06-01T12:00:00Z",
  "gender": "string",
  "origin": "string",
  "push_name": "string",
  "status": "string",
  "first_message_id": "string",
  "first_message_timestamp": 0,
  "created_at": "2024-06-01T12:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

---

## Error Response Example

```json
{
  "success": false,
  "error": "agent not found"
}
```
```