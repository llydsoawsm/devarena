# DevCode Architecture & Implementation Plan

A comprehensive engineering specification and roadmap for building or recreating **DevCode** (real-time 1v1 and competitive multiplayer coding arena).

---

## 1. Executive Summary & System Vision

**DevCode** is a competitive developer arena platform where engineers enter ranked or casual matchmaking, get paired in real time, and face off head-to-head to solve algorithmic and coding challenges.

### Core Experience
1. **Matchmaking**: Players select a mode (`Ranked`, `Casual`, `Tournament`) and enter a FIFO or rating-based search queue.
2. **Room Transition**: Two players are paired $\rightarrow$ match instance generated with a unique ULID $\rightarrow$ synchronized 10-second countdown.
3. **Live Coding Battle**: Dual-pane IDE powered by Monaco Editor, live synchronized timer, real-time opponent progress indicators, problem statement, and test case execution.
4. **Code Submission & Online Judge**: Untrusted code is isolated and tested against secret test cases in sandboxed ephemeral runners with CPU/memory caps.
5. **Settlement & ELO Rating**: Match finishes upon valid submission, timeout, or forfeit; participants' ELO/MMR ratings update immediately.

---

## 2. High-Level Architecture Topology

```
                   +-------------------------------+
                   |      Client (Web / Mobile)    |
                   |   Vue 3 / Next.js + Monaco    |
                   +---------------+---------------+
                                   |
                     HTTPS (REST)  |  WSS (WebSockets)
                                   v
+--------------------------------------------------------------------+
|                         API Gateway & Real-Time                     |
|  - Reverse Proxy (Nginx / Traefik / Caddy)                         |
|  - Web / HTTP API (Matchmaking, Problems, Auth, History)           |
|  - Real-Time WebSocket Server (Reverb / Socket.io / Go WebSockets) |
+------------------+------------------------------+------------------+
                   |                              |
                   v                              v
      +------------------------+      +------------------------+
      |      PostgreSQL /      |      |   Redis / Dragonfly    |
      |         MySQL          |      |  - Pub/Sub & Sockets   |
      |  - Persistent State    |      |  - Ephemeral Queues    |
      |  - Relational Models   |      |  - Matchmaker Locks    |
      +------------------------+      +-----------+------------+
                                                  |
                                                  v
                                     +---------------------------+
                                     |  Judge & Sandbox Workers  |
                                     |   Docker / gVisor Runners |
                                     |  - C++, Python, TS, Go    |
                                     +---------------------------+
```

---

## 3. Core Subsystems Specification

### 3.1 Matchmaking & Queue Engine

* **Queue Modes**: `ranked`, `casual`, `tournament`.
* **Queue Data Structure**: 
  * In-memory sorted sets (Redis `ZSET` scored by rating or timestamp) or transactional relational database queue entries with row-level locks (`SELECT ... FOR UPDATE`).
* **Matchmaking Algorithm**:
  * **Phase 1 (Immediate)**: Exact MMR window match ($\pm 50$ rating).
  * **Phase 2 (Expanded Search, > 15s in queue)**: Expand window by $+50$ rating every 10 seconds up to max window ($\pm 300$).
  * **Concurrency Protection**: Atomic lock on user pairing to prevent double-matching.

### 3.2 Match State Machine

A finite state machine governs all matches to prevent invalid state progressions:

```
[Created] ---> [Countdown] ---> [Live] ---> [Finished] ---> [Settled]
    |               |             |
    +---------------+-------------+-------> [Cancelled / Forfeit]
```

* **States**:
  * `created`: Match and participant records allocated; room ULID created.
  * `countdown`: 10-second sync window where participants load the room assets.
  * `live`: Timer begins; editor is unlocked; code execution enabled.
  * `finished`: A player passed all tests, timer expired, or player resigned.
  * `settled`: Ratings adjusted, match history updated, room archived.
  * `cancelled`: Match voided due to player failing to connect during countdown.

### 3.3 Real-Time WebSocket Protocol

* **Channels**:
  * Private Room Channel: `match.{public_id}`
  * User Private Channel: `user.{id}`
* **Socket Events**:
  * `match.created`: Emitted when matchmaking pairs two players.
  * `match.countdown_tick`: Broadcasts countdown countdown clock.
  * `match.started`: Live round start; unlocks the Monaco editor.
  * `participant.heartbeat`: Keeps connection alive (10s interval).
  * `participant.state`: Opponent progress (e.g., "Passed 2/5 test cases", keystroke heatmaps).
  * `match.submission_result`: Broadcasts test results without leaking opponent's solution code.
  * `match.finished`: Declares winner, loser, or draw with finish timestamps.

### 3.4 Isolated Code Execution Sandbox (Online Judge)

Untrusted user code must **never** execute on the host machine.

* **Execution Runtime**: Docker containers or lightweight microVMs (e.g., gVisor / Firecracker / nsjail).
* **Execution Parameters**:
  * **Network Isolation**: `--network none` (strictly forbidden from external egress).
  * **CPU Limit**: Maximum 1 vCPU core.
  * **Memory Limit**: 128 MB – 256 MB.
  * **Execution Timeout**: 2.0 – 5.0 seconds per test case.
  * **File System**: Read-only root filesystem with ephemeral tmpfs storage.
* **Grading Verdicts**:
  * `AC` (Accepted): All test cases match expected outputs.
  * `WA` (Wrong Answer): Output mismatch on one or more test cases.
  * `TLE` (Time Limit Exceeded): Execution exceeded timeout.
  * `MLE` (Memory Limit Exceeded): Process exceeded memory limits.
  * `CE` (Compilation Error): Code failed compilation step.
  * `RE` (Runtime Error): Unhandled exception / Segfault.

### 3.5 ELO Rating System

* Standard logistic ELO formulation:
  $$E_A = \frac{1}{1 + 10^{(R_B - R_A) / 400}}$$
  $$R'_A = R_A + K \cdot (S_A - E_A)$$
* **Parameters**:
  * $K = 32$ for provisional players ($< 30$ matches), $K = 16$ for seasoned players.
  * $S_A = 1$ (Win), $0.5$ (Draw), $0$ (Loss / Forfeit).

---

## 4. Recommended Tech Stacks

### Option A: Go (Golang) + Next.js (Maximum Performance & Scalability)
* **Backend**: Go (Fiber or Gin) + Gorilla WebSocket
* **Frontend**: Next.js 15 (React 19) + Tailwind CSS + `@monaco-editor/react`
* **Database & Cache**: PostgreSQL + Redis
* **Judge**: Go microservice wrapping Docker Engine SDK
* **Pros**: Native concurrency handling with goroutines, ultra-low WebSocket overhead, effortless Docker API integration.

### Option B: Full-Stack TypeScript (Bun / Fastify + Vue 3 or React)
* **Backend**: Fastify or NestJS on Bun / Node.js + Socket.io / `ws`
* **Frontend**: Vue 3 (Vite) or Next.js
* **Queue**: BullMQ on Redis
* **Judge**: Integration with self-hosted Judge0 or Piston
* **Pros**: 100% shared types between client, server, and events; maximum development speed.

### Option C: Laravel 12 + Laravel Reverb + Vue 3 (Inertia.js) — *Direct Evolution of Current Repo*
* **Backend**: Laravel 12 + PHP 8.2+
* **Real-time Server**: **Laravel Reverb** (first-party WebSocket server)
* **Frontend**: Inertia.js + Vue 3 + Tailwind CSS
* **Judge Runner**: Asynchronous worker consuming from Redis queue, dispatching to Docker container runner.
* **Pros**: Directly leverages the existing codebase and data models in `c:/laragon/www/devarena`.

---

## 5. Database Schema Blueprint

```sql
-- Users
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    rating INT DEFAULT 1200 NOT NULL,
    matches_played INT DEFAULT 0 NOT NULL,
    wins INT DEFAULT 0 NOT NULL,
    losses INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP NULL
);

-- Problems
CREATE TABLE problems (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(128) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    difficulty ENUM('easy', 'medium', 'hard') NOT NULL,
    description TEXT NOT NULL,
    starter_code JSON NOT NULL, -- { "python": "...", "javascript": "...", "cpp": "..." }
    time_limit_ms INT DEFAULT 2000 NOT NULL,
    memory_limit_kb INT DEFAULT 131072 NOT NULL, -- 128 MB
    created_at TIMESTAMP NULL
);

-- Test Cases
CREATE TABLE test_cases (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    problem_id BIGINT UNSIGNED NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input MEDIUMTEXT NOT NULL,
    expected_output MEDIUMTEXT NOT NULL,
    is_hidden BOOLEAN DEFAULT TRUE NOT NULL,
    sequence_order INT DEFAULT 1 NOT NULL
);

-- Matches
CREATE TABLE matches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    public_id VARCHAR(32) UNIQUE NOT NULL, -- ULID
    mode ENUM('ranked', 'casual', 'tournament') NOT NULL,
    status ENUM('created', 'countdown', 'live', 'finished', 'settled', 'cancelled') DEFAULT 'created' NOT NULL,
    problem_id BIGINT UNSIGNED NULL REFERENCES problems(id),
    max_duration_seconds INT DEFAULT 1800 NOT NULL,
    countdown_seconds SMALLINT DEFAULT 10 NOT NULL,
    started_at TIMESTAMP NULL,
    ended_at TIMESTAMP NULL,
    settled_at TIMESTAMP NULL,
    version INT DEFAULT 1 NOT NULL,
    metadata JSON NULL,
    created_at TIMESTAMP NULL
);

-- Match Participants
CREATE TABLE match_participants (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    match_id BIGINT UNSIGNED NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id BIGINT UNSIGNED NOT NULL REFERENCES users(id),
    seat TINYINT NOT NULL, -- 1 or 2
    result ENUM('pending', 'won', 'lost', 'draw', 'forfeit') DEFAULT 'pending' NOT NULL,
    rating_before INT NULL,
    rating_after INT NULL,
    rating_delta INT NULL,
    submission_count INT DEFAULT 0 NOT NULL,
    is_connected BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP NULL,
    UNIQUE KEY (match_id, seat),
    UNIQUE KEY (match_id, user_id)
);

-- Submissions
CREATE TABLE submissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    match_id BIGINT UNSIGNED NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id BIGINT UNSIGNED NOT NULL REFERENCES users(id),
    problem_id BIGINT UNSIGNED NOT NULL REFERENCES problems(id),
    language VARCHAR(32) NOT NULL,
    code MEDIUMTEXT NOT NULL,
    verdict ENUM('AC', 'WA', 'TLE', 'MLE', 'CE', 'RE', 'PENDING') NOT NULL,
    passed_test_cases INT DEFAULT 0 NOT NULL,
    total_test_cases INT DEFAULT 0 NOT NULL,
    runtime_ms INT NULL,
    memory_kb INT NULL,
    submitted_at TIMESTAMP NOT NULL
);

-- Match Events (Audit Log)
CREATE TABLE match_events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    match_id BIGINT UNSIGNED NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    participant_id BIGINT UNSIGNED NULL REFERENCES match_participants(id) ON DELETE SET NULL,
    event_type VARCHAR(64) NOT NULL,
    event_sequence BIGINT UNSIGNED NOT NULL,
    payload JSON NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    UNIQUE KEY (match_id, event_sequence)
);
```

---

## 6. Phased Implementation Roadmap

### Phase 1: Database Foundation & Core Auth
* [ ] Finalize migrations for `problems`, `test_cases`, `matches`, `match_participants`, and `submissions`.
* [ ] Seed sample algorithmic problems (Easy, Medium, Hard) with inputs, outputs, and starter code templates.
* [ ] Set up user registration, login, and profile stats tracking.

### Phase 2: Matchmaking & Queue Engine
* [ ] Implement atomic enqueue and leave actions with cancellation safety.
* [ ] Build pairing worker that sweeps active queue entries and generates matches.
* [ ] Implement room assignment and seat reservations for matched participants.

### Phase 3: Real-Time WebSockets & Presence
* [ ] Configure WebSocket infrastructure (Laravel Reverb / Socket.io / Go WebSockets).
* [ ] Connect client-side listeners for `match.{public_id}` channels.
* [ ] Synchronize 10-second pre-match countdown and live match timer clock.
* [ ] Implement participant presence detection, disconnection timeouts, and forfeit triggers.

### Phase 4: IDE & Match Room UI
* [ ] Mount Monaco Editor in `Room.vue` / room component.
* [ ] Add multi-language selector (Python, JavaScript/TypeScript, C++, Go).
* [ ] Build test runner interface allowing users to run sample test cases before full submission.
* [ ] Display real-time opponent progress bar (e.g. "Passed 3/6 test cases").

### Phase 5: Code Execution Sandbox (Online Judge)
* [ ] Set up container-based execution worker (Docker / gVisor).
* [ ] Implement compile, run, and output assertion pipeline.
* [ ] Enforce memory limits, CPU caps, and strict timeout terminations.
* [ ] Connect submission queue to judge worker with asynchronous status updates.

### Phase 6: ELO Rating, Match Settlement & Polish
* [ ] Implement ELO rating adjustment calculations upon match victory/defeat.
* [ ] Create post-match results screen displaying rating deltas, execution benchmarks, and full solution diffs.
* [ ] Build user profile match history, win rates, and global leaderboard.
* [ ] Comprehensive test suite covering concurrent queue edge cases, judge timeouts, and disconnects.
