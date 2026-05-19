# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — 2026-05-19

### Added

- User registration and sign-in by phone number and password.
- User profile: view and update fields (`full_name`, `phone`, `email`, `birth_date`).
- Match listing with opponent info, score, status and derby flag.
- Match statistics: wins and losses broken down by finish type (regular time, overtime, penalty shootout).
- Stadium layout: sectors with price coefficients and color codes.
- Seats per sector with availability and wheelchair-accessible flags.
- Ticket purchase: seat reservation, QR hash generation, buyer data storage.
- User ticket history with match and seat details.
- Admin endpoints for sales statistics.
- 8 database migrations: schema, opponents, matches, sectors, seats, tickets and demo data.
- Multi-stage `Dockerfile`; image includes `migrate` CLI and migration files.
- CI: GitHub Actions PR validation (formatting, `go vet`, build, Docker image).
- CI: image publish to GHCR on `v*.*.*` tag push.
- CI: automatic server deploy on GitHub Release publish.
