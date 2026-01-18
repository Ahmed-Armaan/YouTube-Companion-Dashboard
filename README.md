# YouTube-Companion-Dashboard

A youtube API wrapper which lets user manage their youtube video titles, descriptions, comments along with a notes section to keep work upon your ideas.

## API

### GET /me
Check if the user is authenticated.

### GET /channelId
Get the channel Id of user

**Response**
if successful:
```json
{
    channelId: <ChannelId>
}
```

### GET /channel
Get the list of videos on the channel

**Response**
```json
{
  "nextPageToken": "string",
  "videos": [
    {
      "videoId": "string",
      "title": "string",
      "description": "string",
      "publishedAt": "string (RFC3339 timestamp)",
      "thumbnail": "string (URL)",
      "duration": "string (ISO 8601 duration)",
      "viewCount": "string",
      "likesCount": "string",
      "disLikesCount": "string",
      "embeddedhtml": "string (HTML iframe)"
    }
  ]
}
```

### GET /comments
Get comment thread on the video

**Query parameters**
- videoId
- PageToken

**Response**
```json
{
  "nextPageToken": "string",
  "items": [
    {
      "id": "string",
      "snippet": {
        "channelId": "string",
        "topLevelComment": {
          "id": "string",
          "snippet": {
            "authorDisplayName": "string",
            "authorProfileImageUrl": "string (URL)",
            "authorChannelUrl": "string (URL)",
            "authorChannelId": {
              "value": "string"
            },
            "textOriginal": "string"
          }
        }
      },
      "replies": {
        "comments": [
          {
            "id": "string",
            "snippet": {
              "authorDisplayName": "string",
              "authorProfileImageUrl": "string (URL)",
              "authorChannelUrl": "string (URL)",
              "authorChannelId": {
                "value": "string"
              },
              "textOriginal": "string"
            }
          }
        ]
      }
    }
  ]
}
```

### POST /comment
Upload a comment

### POST /comment/reply
Upload a reply

### DELETE /comment
delete a comeent

### POST /title/ai
Suggest three names for video based on previous title and description (Uses OpenAI API)

**Request**
```json
{
  "title": "string",
  "description": "string"
}
```
**Response**
```json
{
  "suggestions": [
    "string",
    "string",
    "string"
  ]
}
```

### POST /notes
put a note in database

**Request**
```json
{
  "videoId": "string",
  "content": "string",
  "tags": ["string"]
}
```
**Response**
```json
{
  "id": "string",
  "userId": "string (UUID)",
  "videoId": "string",
  "content": "string",
  "tags": ["string"],
  "createdAt": "string (RFC3339 timestamp)"
}
```

### GET /notes
get all notes for a video

**Query Parameters**
- videoId
- limit (how many notes to fetch)
- tags

**Response**
```json
{
  "items": [
    {
      "id": "string",
      "userId": "string (UUID)",
      "videoId": "string",
      "content": "string",
      "tags": ["string"],
      "createdAt": "string (RFC3339 timestamp)"
    }
  ],
  "nextCursor": "string (RFC3339 timestamp) | null"
}
```

### DELETE /note
Deletes a note

# Database Schema

The database is PostgreSQL and managed using GORM. UUIDs are used for primary keys.

## Tables

### `users`
Stores authenticated users mapped to Google OAuth accounts.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_user_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
```
### tokens
Stores encrypted OAuth refresh tokens for YouTube API access.

```sql
CREATE TABLE tokens (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    refresh_token_enc TEXT NOT NULL,
    revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
```
### `Note` Model (GORM)

The `notes` table is managed via the following Go struct. It includes indexing on `UserID` and `VideoID` for optimized querying and uses PostgreSQL's native `TEXT[]` for tagging.

```go
type Note struct {
    ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    uuid.UUID      `gorm:"type:uuid;index"`
    VideoID   string         `gorm:"index"`
    Content   string
    Tags      pq.StringArray `gorm:"type:text[]"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Powered by the YOutube API
