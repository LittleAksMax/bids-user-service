# User Service

This service exposes user-facing and internal endpoints for Amazon Ads user data, policy attachments, bid history, logs, schedules, and token management.

## Module Summary

- `main.go`: Bootstraps the service.
  - Loads environment configuration.
  - Connects to Postgres and Redis.
  - Runs migrations in development mode.
  - Builds the router and starts the HTTP server.

- `api/`: HTTP transport layer, route registration, middleware wiring, and controller logic.
  - `authController`: Builds Amazon LwA redirect URLs and processes OAuth callbacks.
  - `campaignsController`: Returns seller profiles and campaigns, with Redis-backed response caching.
  - `attachmentController`: Reads, attaches, and detaches policies for profiles; invalidates related campaign cache entries.
  - `bidsController`: Queries bid history and creates bid records.
  - `tokensController`: Returns stored refresh tokens and the set of authenticated regions.
  - `logsController`: Creates user logs and returns paginated log history.
  - `policySchedulesController`: Manages schedule CRUD, prioritisation, due-schedule reads, and internal drive/process actions.
  - `router.go` / `routes.go`: Wires middleware, services, health checks, and endpoint groups together.

- `service/`: Business logic layer and external integration logic.
  - `CampaignsService`: Loads Amazon Ads seller profiles and campaigns across regions.
  - `AuthService`: Exchanges Amazon OAuth codes for refresh tokens and stores them.
  - `LWAStateService`: Encrypts and decrypts OAuth redirect state.
  - `AttachmentService`: Converts attach/detach requests into repository operations.
  - `BidsService`: Applies bid-search filters and persists bid history.
  - `TokensService`: Reads stored region refresh tokens for a user.
  - `LogsService`: Validates log writes and paginates log reads.
  - `PolicySchedulesService`: Creates, deletes, prioritises, drives, and processes schedules.
  - `service_factories.go`: Constructs concrete services from repositories and config.

- `repository/`: Postgres persistence layer.
  - `bids_repository`: Stores and queries bid history.
  - `attached_policies_repository`: Batch-upserts, deletes, and lists attached policies by profile.
  - `user_tokens_repository`: Stores region-specific Amazon refresh tokens.
  - `logs_repository`: Appends user logs and reads paginated log pages.
  - `schedules_repository`: Stores schedules and implements prioritise/drive/process state transitions.

- `cache/`: Redis-backed caching utilities.
  - `RedisRefreshStore`: Owns the shared Redis client created at startup.
  - `RequestCache` / `redisRequestCache`: Caches HTTP-derived data such as profiles and campaigns.
  - `redis_common.go`: Shared Redis read/TTL helper logic.

- `db/`: Database connection and migration helpers.
  - `config.go`: Builds the Postgres DSN.
  - `db.go`: Opens the pooled connection and runs Goose migrations.

- `config/`: Loads typed application config from environment variables for auth, Amazon Ads, Redis, Postgres, CORS, and port selection.

- `contracts/`: Shared request DTOs, response DTOs, token payloads, and domain models passed between layers.

- `health/`: Small health-check abstraction used by `/health`.
  - `DBHealthChecker`: Pings Postgres.
  - `CacheHealthChecker`: Pings Redis.

- `migrations/`: Schema changes for users, bids, attached policies, schedules, and logs.

## `campaigns_service` Semaphore Implementation

`CampaignsService.GetCampaigns` first fetches the campaign list for a profile, filters it down to enabled Sponsored Products campaigns with auto-target creation enabled, and then fans out ad-group loading with one goroutine per campaign.

To avoid unbounded concurrency, it uses a buffered channel as a counting semaphore:

- `sem := make(chan struct{}, maxConcurrentRequests)` creates a pool of `5` available slots.
- Before a goroutine performs an Amazon Ads API call, it writes to `sem`.
- When the call finishes, it reads from `sem`, releasing the slot.
- Because the channel buffer is capped, at most `5` guarded API calls can be in flight at once.

In this implementation, the semaphore is used around:

- The initial per-campaign `GetAdGroups(...)` paginator creation.
- Each paginated `Next(ctx)` call while fetching ad-group pages.

The goroutines send results into a buffered `results` channel along with the original campaign index, and the service rebuilds the final output slice in the original campaign order after the `WaitGroup` completes.

## `LWA` OAuth Flow

The `LWA` flow is intentionally coordinated by the backend rather than the frontend. The User Service creates the Amazon redirect URL itself so it can embed internal state into the OAuth `state` parameter before the browser leaves the application. That state contains the authenticated user ID, the selected Amazon Ads region, and the frontend return URL.

The state is encrypted with a backend-held symmetric key and decrypted again when Amazon redirects back to the callback endpoint. This lets the service carry trusted internal context through the browser round-trip without exposing that state to the client as plain text and without trusting the client to send it back unchanged. The implementation uses AES-GCM. A limitation of this is that the key for the encryption is not rotated.

1. Client starts the flow.
   - An authenticated client calls `GET /lwa/{region}?redirect_uri=...`.
   - The User Service validates the authenticated subject and requested region.
2. Backend creates the Amazon redirect URL.
   - The service packages `user_id`, `region`, and `redirect_url` into a small state payload.
   - `LWAStateService` serialises that payload to JSON and encrypts it with AES-GCM.
   - The User Service returns the Amazon `/ap/oa` URL containing the encrypted `state`.
3. The user authenticates with Amazon.
   - The client redirects the browser to the returned Amazon URL.
   - Amazon authenticates the user and collects consent for Amazon Ads access.
4. Amazon redirects back to the backend.
   - Amazon calls `GET /lwa/process_token?code=...&state=...`.
   - The User Service decrypts the returned `state` and recovers the original internal context.
5. The backend completes the flow.
   - `AuthService` exchanges the authorisation code for a refresh token.
   - The refresh token is stored against the user and region.
   - The User Service redirects the browser back to the original frontend `redirect_uri`.
