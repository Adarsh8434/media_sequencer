# Multi-Window Media Sequencer with Sync Playback

A full-stack media sequencing application that manages multiple display windows, plays configured media playlists continuously, supports dynamic playlist updates, and provides synchronized playback across all windows.

## Live Demo

- Frontend: `https://media-sequencer-frontend-ig60.onrender.com`
- Backend API: `https://media-sequencer-backend-mwpq.onrender.com`


---

## Project Overview

The Multi-Window Media Sequencer allows multiple display windows to continuously play their own configured media sequences.

Each window has an independent playlist and playback sequence. Administrators can dynamically add media to a window's playlist.

The application also provides a synchronization feature that allows one selected media item to be displayed across all windows simultaneously for a specified duration.

After synchronization ends, each window automatically resumes its own playlist.

---

## Features

### Multi-Window Playback

- Supports multiple display windows.
- Currently configured with 4 windows.
- Each window has its own playlist.
- Each window plays its playlist independently.
- Playlists continuously repeat.

### Media Support

The application supports:

- Images
- Videos
- Blank/fallback display

Each media item contains:

- Media ID
- Name
- Type
- URL
- Duration

### Dynamic Playlist Updates

Media can be added to any window dynamically from the frontend.

When media is added:

1. The frontend sends the request to the backend.
2. The backend stores the playlist entry in PostgreSQL.
3. The frontend reloads the playlist.
4. The new media becomes part of the playback sequence.

### Synchronized Playback

The synchronization feature allows a selected media item to be displayed across all windows simultaneously.

Example:

```text
Window 1 → M1
Window 2 → M5
Window 3 → M2
Window 4 → M6

             ↓ Start Sync

Window 1 → M3
Window 2 → M3
Window 3 → M3
Window 4 → M3

             ↓ Sync duration ends

Window 1 → resumes its playlist
Window 2 → resumes its playlist
Window 3 → resumes its playlist
Window 4 → resumes its playlist
```

The normal playlist configuration is not modified during synchronization.

### 5-Hour Playback Cycle

The assignment treats each window's total playback size as a 5-hour cycle.

The configured playlist continuously repeats within this cycle.

The implementation does not create a 5-hour timer. Instead, each media item's configured duration is used for playback, and the playlist loops back to its first item after the final item.

For example:

```text
M1 → 10 seconds
M2 → 10 seconds
M3 → 10 seconds

After M3:
       ↓
M1 → M2 → M3 → ...
```

This provides continuous playback while preserving the configured sequence.

---

## Technology Stack

### Frontend

- React
- Vite
- JavaScript
- HTML
- CSS

### Backend

- Go
- Gin
- GORM
- PostgreSQL
- REST API

### Database

- PostgreSQL

### Deployment

- GitHub
- Render

---

## Project Structure

```text
media_sequencer/
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── media_handler.go
│   │   │   ├── playlist_handler.go
│   │   │   ├── sync_handler.go
│   │   │   └── window_handler.go
│   │   │
│   │   ├── model/
│   │   │   ├── media.go
│   │   │   ├── playlist.go
│   │   │   └── window.go
│   │   │
│   │   ├── repository/
│   │   │   ├── database.go
│   │   │   ├── media_repository.go
│   │   │   ├── playlist_repository.go
│   │   │   └── window_repository.go
│   │   │
│   │   └── service/
│   │       └── sync_service.go
│   │
│   ├── .env
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── App.jsx
│   │   ├── App.css
│   │   └── index.css
│   │
│   ├── .env
│   ├── package.json
│   └── vite.config.js
│
├── .gitignore
└── README.md
```

> `.env` files contain environment-specific configuration and are not committed to GitHub.

---

## Backend API

Base URL (local):

```text
http://localhost:8080
```

Use the deployed backend URL after deployment.

### Health Check

```http
GET /health
```

Checks whether the backend is running.

**Response**

```json
{
  "status": "ok"
}
```

### Get Windows

```http
GET /windows
```

Returns all configured display windows.

**Example Response**

```json
[
  { "id": 1, "name": "Window 1" },
  { "id": 2, "name": "Window 2" },
  { "id": 3, "name": "Window 3" },
  { "id": 4, "name": "Window 4" }
]
```

### Get Media

```http
GET /media
```

Returns all available media items.

**Example Response**

```json
[
  {
    "id": 1,
    "name": "M1",
    "type": "image",
    "url": "https://picsum.photos/id/1015/800/600",
    "durationSeconds": 10
  }
]
```

### Create Media

```http
POST /media
```

Creates a new media item.

**Request**

```json
{
  "name": "M7",
  "type": "image",
  "url": "https://example.com/image.jpg",
  "durationSeconds": 10
}
```

**Response**

Returns the newly created media item.

### Get Window Playlist

```http
GET /windows/:windowId/playlist
```

Returns the playlist configured for a specific window.

**Example**

```http
GET /windows/1/playlist
```

**Response**

```json
[
  { "id": 1, "windowId": 1, "mediaId": 1, "position": 1 },
  { "id": 2, "windowId": 1, "mediaId": 2, "position": 2 }
]
```

The playlist is ordered by the `position` field.

### Add Media to Playlist

```http
POST /playlist
```

Adds a media item to a window's playlist.

**Request**

```json
{
  "windowId": 1,
  "mediaId": 5,
  "position": 4
}
```

**Response**

Returns the created playlist entry.

---

## Synchronization API

### Start Sync

```http
POST /sync
```

Starts synchronized playback.

**Request**

```json
{
  "mediaId": 3,
  "durationSeconds": 10
}
```

All display windows use the selected media during the synchronization period.

**Response**

```json
{
  "message": "sync started"
}
```

### Get Sync State

```http
GET /sync
```

Returns the current synchronization state.

**Active Sync**

```json
{
  "active": true,
  "mediaId": 3,
  "durationSeconds": 10,
  "startedAt": "2026-09-18T12:00:00Z"
}
```

**No Active Sync**

```json
{
  "active": false
}
```

---

## Database Design

The application uses three main tables.

### `windows`

| Column | Type    | Description  |
| ------ | ------- | ------------ |
| id     | SERIAL  | Primary key  |
| name   | VARCHAR | Window name  |

### `media`

| Column           | Type    | Description       |
| ---------------- | ------- | ----------------- |
| id               | SERIAL  | Primary key       |
| name             | VARCHAR | Media name        |
| type             | VARCHAR | Media type        |
| url              | TEXT    | Media URL         |
| duration_seconds | INTEGER | Playback duration |

### `playlists`

| Column    | Type    | Description        |
| --------- | ------- | ------------------ |
| id        | SERIAL  | Primary key        |
| window_id | INTEGER | Associated window  |
| media_id  | INTEGER | Associated media   |
| position  | INTEGER | Playback order     |

### Relationships

```text
Windows
   │
   │ 1
   │
   │ many
   ▼
Playlists
   │
   │ many
   │
   │ 1
   ▼
Media
```

Foreign keys are used to maintain relationships between windows, playlists, and media.

The `playlists` table also has a unique constraint on:

```text
(window_id, position)
```

to prevent duplicate positions within the same window.

---

## Sync Playback Behavior

The sync state is maintained by the backend sync service.

When synchronization starts:

```text
POST /sync
       ↓
Backend stores sync state
       ↓
Frontend polls /sync
       ↓
All windows detect active sync
       ↓
All windows display selected media
```

When the configured synchronization duration expires:

```text
Sync expires
       ↓
/sync returns active=false
       ↓
Windows stop showing sync media
       ↓
Each window continues its normal playlist
```

The current playlist index is not changed during synchronization, so the window can continue from its normal playback sequence.

---

## Local Setup

### Prerequisites

Install:

- Go
- Node.js
- PostgreSQL
- Git

### Clone Repository

```bash
git clone https://github.com/Adarsh8434/media_sequencer.git
cd media_sequencer
```

### Backend Setup

Go to the backend:

```bash
cd backend
```

Install dependencies:

```bash
go mod tidy
```

Create a `.env` file:

```env
DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASSWORD=YOUR_POSTGRES_PASSWORD
DB_NAME=media_sequencer
DB_PORT=5432
DB_SSLMODE=disable
```

Run the backend:

```bash
go run ./cmd/server
```

The backend runs on `http://localhost:8080`.

Test:

```bash
curl http://localhost:8080/health
```

### Frontend Setup

Open another terminal:

```bash
cd frontend
npm install
```

Create a `.env` file:

```env
VITE_API_URL=http://localhost:8080
```

Start the development server:

```bash
npm run dev
```

The frontend will normally be available at `http://localhost:5173`.

---

## Deployment

The project is deployed using Render.

### Backend

The Go backend is deployed as a Render Web Service.

Build command:

```bash
go build -o server ./cmd/server
```

Start command:

```bash
./server
```

The backend uses the Render-provided `PORT` environment variable.

### PostgreSQL

The application uses Render PostgreSQL for persistent storage.

The deployed backend uses environment variables for the database connection:

```env
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=5432
DB_SSLMODE=require
```

The database schema contains: `windows`, `media`, `playlists`.

### Frontend

The React frontend is deployed as a Render Static Site.

Build command:

```bash
npm install && npm run build
```

Publish directory:

```text
dist
```

The production environment variable is:

```env
VITE_API_URL=(https://media-sequencer-frontend-ig60.onrender.com/)
```

---

## Environment Variables

### Backend

```env
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=
DB_SSLMODE=
PORT=
```

### Frontend

```env
VITE_API_URL= https://media-sequencer-backend-mwpq.onrender.com

---

## Testing

The following functionality was tested during development:

- Backend health endpoint
- Four display windows
- Media retrieval
- Window-specific playlists
- Continuous playlist playback
- Dynamic media addition
- Synchronization across all windows
- Configurable synchronization duration
- Returning to normal playback after synchronization
- PostgreSQL persistence
- Frontend/backend communication
- CORS configuration

---

## Assumptions

- Each media item has a configured playback duration.
- Playlist order is determined by the `position` field.
- A playlist repeats continuously after its final item.
- Synchronization temporarily overrides normal playback.
- Synchronization does not modify a window's stored playlist.
- After synchronization ends, each window resumes normal playback.
- Blank/fallback playback is only displayed when explicitly configured or when no media is available.
- PostgreSQL is used for persistent application data.
- The synchronization state is maintained in backend memory and is not persisted as a database record.
- The application uses polling to obtain the current synchronization state from the backend.

---

## Repository

GitHub: <https://github.com/Adarsh8434/media_sequencer>

## Author

**Adarsh Kumar Choubey**
