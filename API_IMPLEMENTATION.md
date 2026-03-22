# API Endpoints Implementation - Simplified Architecture

This document describes the scaffolded endpoints and their simplified structure.

## Endpoints

### 1. GET `/users/bids/{campaignID}`
**Description:** Get all bids for a specific campaign  
**Authentication:** Required (access token)  
**Response:**
```json
{
  "success": true,
  "data": [
    {
      "user_id": 123,
      "campaign_id": 456,
      "adgroup_id": 789,
      "policy_id": 1,
      "from_bid": 1.50,
      "to_bid": 2.00,
      "change_date": "2026-02-26T00:00:00Z"
    }
  ]
}
```

### 2. POST `/users/bids`
**Description:** Create a new bid  
**Authentication:** Required (access token)  
**Request Body:**
```json
{
  "campaign_id": 456,
  "adgroup_id": 789,
  "policy_id": 1,
  "from_bid": 1.50,
  "to_bid": 2.00
}
```

### 3. GET `/user/tokens`
**Description:** Get user's Amazon Ads refresh tokens  
**Authentication:** Required (access token)  
**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": 123,
    "refresh_token_eu": "token_eu_here",
    "refresh_token_us": null,
    "refresh_token_fe": "token_fe_here"
  }
}
```

### 4. POST `/user/tokens`
**Description:** Set a refresh token for a specific region  
**Authentication:** Required (access token)  
**Request Body:**
```json
{
  "region": "eu",
  "token": "new_refresh_token_here"
}
```

### 5. GET `/campaigns`
**Description:** Get all campaigns with nested structure (sellers > profiles > campaigns > adgroups > ads)  
**Authentication:** Required (access token)  
**Response:**
```json
{
  "success": true,
  "data": {
    "sellers": [
      {
        "id": "seller_123",
        "profiles": [
          {
            "id": 123456,
            "name": "Profile Name",
            "campaigns": [
              {
                "id": 789,
                "name": "Campaign Name",
                "adgroups": [
                  {
                    "id": 101,
                    "name": "AdGroup Name",
                    "ads": [
                      {
                        "id": 202,
                        "name": "Ad Name"
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### 6. GET `/process_token`
**Description:** Amazon LwA OAuth callback endpoint  
**Authentication:** None (public endpoint)  
**Query Parameters:**
- `code`: Authorization code from Amazon
- `state`: State parameter for CSRF protection and user/region identification

## Simplified Architecture

### Controllers Layer (`api/`)
Each controller is independent and receives its specific service via dependency injection:

- **`bids_controller.go`** - Injects `BidsService`
- **`tokens_controller.go`** - Injects `TokensService`
- **`campaigns_controller.go`** - Injects `CampaignsService`
- **`auth_controller.go`** - Injects `AuthService`
- **`user_controller.go`** - Empty scaffold for future extensions
- **`routes.go`** - Route registration
- **`router.go`** - Dependency injection and wiring

### Service Layer (`service/`)
Services are **completely independent** - no nested service pattern:

**Interfaces:**
- `BidsService` - Bid operations
- `TokensService` - Token operations  
- `CampaignsService` - Campaign retrieval from Amazon Ads API
- `AuthService` - Amazon LwA authentication

**Constructors (exported):**
- `NewBidsService(db *sql.DB)` - Creates BidsService
- `NewTokensService(db *sql.DB)` - Creates TokensService
- `NewCampaignsService(db *sql.DB)` - Creates CampaignsService
- `NewAuthService(db *sql.DB)` - Creates AuthService

**Implementations (private):**
- `bids_service.impl.go` - Private `bidsService` struct
- `tokens_service.impl.go` - Private `tokensService` struct
- `campaigns_service.impl.go` - Private `campaignsService` struct (with TODOs for Amazon SDK)
- `auth_service.impl.go` - Private `authService` struct (with TODOs for OAuth)

### Repository Layer (`repository/`)
- `bids_repository.go` - Database operations for bids
- `user_tokens_repository.go` - Database operations for tokens

### Contracts Layer (`contracts/`)
- `models.go` - All request/response DTOs

## Dependency Injection Flow

```
main.go
  └─> api.NewRouter(pool, cache, cfg)
       └─> service.NewBidsService(pool)      ──> bidsController
       └─> service.NewTokensService(pool)    ──> tokensController
       └─> service.NewCampaignsService(pool) ──> campaignsController
       └─> service.NewAuthService(pool)      ──> authController
```

**Key Points:**
- ✅ Services are injected separately at the router level
- ✅ No nested service dependencies (services don't contain other services)
- ✅ Each controller only gets the service it needs
- ✅ Simple, flat architecture
- ✅ Easy to test - mock individual services

## TODO Items

1. **Amazon Ads API Integration** (`service/campaigns_service.impl.go`)
   - Import your `amazon-ads-api-sdk-go` package
   - Implement profile fetching
   - Implement campaign, adgroup, and ad fetching
   - Implement seller ID binning logic
   - Build nested response structure

2. **Amazon LwA OAuth** (`service/auth_service.impl.go`)
   - Implement OAuth configuration
   - Implement code-to-token exchange
   - Implement state validation (CSRF protection)
   - Parse state to extract userID and region
   - Store refresh tokens appropriately

3. **Configuration**
   - Add Amazon Ads API credentials to config
   - Add Amazon LwA OAuth credentials to config
   - Add redirect URI configuration

4. **Error Handling**
   - Add more specific error types
   - Add validation for request payloads
   - Add logging throughout

5. **Testing**
   - Add unit tests for services
   - Add integration tests for repositories
   - Add API endpoint tests

## Key Design Decisions

### Why Independent Services?
- **Simplicity:** No nested service pattern to understand
- **Testability:** Each service can be mocked independently
- **Maintainability:** Clear dependencies - each controller knows exactly what it needs
- **Flexibility:** Easy to swap service implementations

### Why Private Service Implementations?
- Only the interface is public (`BidsService`, etc.)
- Implementation details (`bidsService` struct) are hidden
- Consumers depend on interfaces, not concrete types
- Follows Go best practices

## Notes

- All endpoints except `/process_token` require authentication via access token
- Bids are immutable (can only be created, not updated or deleted)
- Tokens are stored per region (EU, US, FE)
- The campaigns endpoint returns a nested structure binned by seller ID
- User ID is extracted from the JWT access token by the auth middleware
- Services are completely independent - no service contains another service

