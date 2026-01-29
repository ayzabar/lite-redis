# Lite-Redis

[![CI Pipeline](https://github.com/ayzabar/lite-redis/actions/workflows/docker-build.yml/badge.svg)](https://github.com/ayzabar/lite-redis/actions/workflows/docker-build.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/ayzabar/lite-redis?label=Go&color=blue)
![Docker Image Size](https://img.shields.io/badge/docker_image_size-5.04MB-success)
![Repo Size](https://img.shields.io/github/repo-size/ayzabar/lite-redis)
![License](https://img.shields.io/github/license/ayzabar/lite-redis)

**Lite-Redis** is a lightweight, thread-safe, in-memory key-value store written in **Go** from scratch.

I built this project to deep-dive into **TCP socket programming**, **concurrency patterns**, and **system architecture**—basically, the fun stuff that powers high-scale systems like the ones at **Peak Games**. No external dependencies, no bloat. Just pure Go.

---

## Features

- **TCP Server:** Raw socket programming using `net` package.
- **RESP Compatible:** Speaks the Redis Serialization Protocol (works with `redis-cli`).
- **Concurrent & Thread-Safe:** Handles multiple connections via Goroutines, protected by `sync.RWMutex`.
- **TTL (Time To Live):** Supports key expiration (`EX` command).
  - **Lazy Expiration:** Checks validity on read.
  - **Active Expiration (Janitor):** Background process cleans up expired keys every second.
- **Dockerized:** Multi-stage build resulting in a tiny **~5MB** scratch image.
- **CI/CD:** Automated testing and build pipeline via GitHub Actions.

---

## Quick Start

### Option 1: Docker (Recommended)

The image is optimized (Scratch base). It's lighter than a typical MP3 file.

    # 1. Build the image
    docker build -t lite-redis .

    # 2. Run container (background)
    docker run -d -p 6379:6379 --name my-redis lite-redis

### Option 2: Run Manually (Go Required)

    go run main.go

---

## Usage

You can connect using `netcat` or any standard redis client.

**Example Session:**

    $ nc localhost 6379

    PING
    > +PONG

    SET nick archura
    > +OK

    GET nick
    > $7
    > archura

    # Setting a key with 5 seconds TTL
    SET buyokki "just testing" EX 5
    > +OK

    # Wait 6 seconds...
    GET buyokki
    > $-1

---

## Architecture / Under the Hood

How does it handle traffic without crashing?

1.  **The Listener:** Binds to port `6379` and listens for incoming TCP connections.
2.  **The Handler:** Spawns a dedicated **Goroutine** for every new connection (Concurrency).
3.  **The Storage:** Uses a `map[string]Item` struct protected by a **Mutex** (preventing Race Conditions).
4.  **The Janitor:** A background Goroutine that wakes up every second to sweep away expired keys (Active Expiration).

---

## Author

**Bahri "Ayzabar"**
_Computer Engineering Student @ İSTE_

Passionate about System Programming, Backend Development, and DevOps.
I like breaking things to learn how to fix them.

[GitHub](https://github.com/ayzabar)
