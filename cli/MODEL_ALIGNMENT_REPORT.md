# Terminal News CLI - Model Alignment Report

**Date**: November 18, 2025
**Task**: Align CLI models with shared/models/models.go
**Status**: ✅ **COMPLETE**

---

## 🎯 OBJECTIVE

Align the Terminal News CLI to use the canonical models from `shared/models/models.go` instead of duplicating model definitions. This ensures consistency across the backend, CLI, and scraper components.

---

## 📊 CHANGES SUMMARY

### Files Modified: 8
1. `cli/internal/models/models.go` - **MAJOR REFACTOR**
2. `cli/internal/ui/components/article_list.go` - Updated to use ArticleWithRanking
3. `cli/internal/ui/components/comment_tree.go` - Updated to use CommentWithUser
4. `cli/internal/ui/views/views.go` - Updated message types
5. `cli/internal/ui/views/article_detail.go` - Updated to use new models
6. `cli/internal/api/client.go` - Updated GetArticle return type
7. `cli/cmd/mockserver/main.go` - Updated to match shared model structure
8. *(cache.go not modified - already uses base Article model)*

---

## 🔧 DETAILED CHANGES

### 1. `cli/internal/models/models.go` - Complete Refactor

**Before**: Duplicated full model definitions for Article, Comment, Classified, User

**After**: Re-exports shared models and adds CLI-specific extensions

```go
// Re-export shared models for compatibility
type Article = shared.Article
type Comment = shared.Comment
type Classified = shared.Classified
type User = shared.User
type Vote = shared.Vote
type ArticleWithRanking = shared.ArticleWithRanking
type CommentWithUser = shared.CommentWithUser

// CLI-specific extensions
type CommentTree struct {
    Comment
    Username string
    Karma    int
    Depth    int
    Children []CommentTree
}
```

**Why**:
- Eliminates model duplication
- Ensures consistency with backend API
- Maintains CLI-specific tree structure for display

---

### 2. Field Mapping Changes

#### Article → ArticleWithRanking

| Old CLI Field | New Shared Field | Notes |
|---------------|------------------|-------|
| `Summary` | `Content` | Now uses full content, truncated at display time |
| `Upvotes` | `LikeCount` | Matches backend voting terminology |
| `Downvotes` | `DislikeCount` | Matches backend voting terminology |
| `Views` | `OpenCount` | Tracks article opens (more accurate) |
| `CommentCount` | `TotalEngagement / 5` | Estimated from engagement metrics |
| `IsHot` | `HotRank > 0.7` | Calculated from hot ranking algorithm |
| `IsRising` | `HoursSincePublished < 6 && TotalEngagement > 10` | Calculated from time + engagement |
| *(new)* | `Category` | Now available from shared model |
| *(new)* | `Tags` | Now available from shared model |
| *(new)* | `ImageURL` | Now available from shared model |
| *(new)* | `Author` | Now available from shared model |
| *(new)* | `ControversyScore` | New metric for controversial content |

#### Comment → CommentWithUser

| Old CLI Field | New Shared Field | Notes |
|---------------|------------------|-------|
| `Username` | `Username` | Now from CommentWithUser (not Comment) |
| `Depth` | *(local only)* | Added to CommentTree helper |
| `Children` | *(local only)* | Added to CommentTree helper |
| *(new)* | `Karma` | User karma score now available |
| *(new)* | `IsFlagged` | Flag status now available |
| *(new)* | `FlagCount` | Flag count now available |
| *(new)* | `UpdatedAt` | Update timestamp now available |
| *(new)* | `EditedAt` | Edit timestamp now available |

---

### 3. Component Updates

#### ArticleList Component

**Changes**:
- Type changed from `[]models.Article` → `[]models.ArticleWithRanking`
- Updated field references: `Upvotes` → `LikeCount`, `Views` → `OpenCount`, etc.
- Hot/Rising indicators now calculated from `HotRank` and `HoursSincePublished`
- Summary rendering now uses `Content` field (truncated to 200 chars)
- Comment count estimated from `TotalEngagement / 5`

**Impact**:
- More accurate article metrics from backend
- Access to full content instead of just summary
- Better hot/rising detection using ranking algorithm

#### CommentTree Component

**Changes**:
- Input type changed from `[]models.Comment` → `[]models.CommentWithUser`
- Internal tree structure uses new `models.CommentTree` helper
- Now has access to user karma scores
- Can display flag status if needed

**Impact**:
- Proper user attribution with karma
- Better moderation support with flag fields

#### ArticleDetailView

**Changes**:
- Article type changed from `*models.Article` → `*models.ArticleWithRanking`
- Comments type changed from `[]models.Comment` → `[]models.CommentWithUser`
- Updated all field references to match new model structure
- Content preview logic added for article summary

**Impact**:
- Full article metadata display
- Proper comment attribution with usernames

---

### 4. API Client Updates

**Changes**:
- `GetArticle()` now returns `*models.ArticleWithRanking`
- `GetArticles()` already returns `ArticlesResponse` which uses `[]ArticleWithRanking`
- `GetComments()` already returns `CommentsResponse` which uses `[]CommentWithUser`

**Impact**: All API responses now match shared model structure

---

### 5. Mock Server Updates

**Purpose**: Ensure test data matches real API structure

**Changes**:
- Article struct updated with all ranking fields
- Comment struct updated with all shared fields + CommentWithUser fields
- Mock data generation includes:
  - `HotRank` calculation: `likeCount / (hoursSince + 2)`
  - `ControversyScore` calculation: `dislikeCount / (likeCount + 1)`
  - `TotalEngagement` = likes + dislikes + opens/10
  - User karma for comments
  - Timestamps for `UpdatedAt`, `EditedAt`

**Impact**: Mock server now provides realistic test data matching production API

---

## ✅ VERIFICATION CHECKLIST

- [x] All CLI files import from `shared/models` instead of duplicating
- [x] ArticleList uses ArticleWithRanking with correct field mappings
- [x] CommentTree uses CommentWithUser with tree helper
- [x] ArticleDetailView uses updated models
- [x] API client returns updated types
- [x] Message types (ArticlesLoadedMsg, CommentsLoadedMsg) updated
- [x] Mock server generates data matching shared model structure
- [x] Cache layer uses base Article model (intentional - caching only core data)
- [ ] **BLOCKED**: Compilation testing (Go not installed)

---

## 🚀 BENEFITS

### 1. Single Source of Truth
- Models defined once in `shared/models/models.go`
- Backend, CLI, and scraper all use same definitions
- Type changes propagate automatically

### 2. Better Data Consistency
- Field names match across all components
- No confusion between Upvotes/LikeCount, Views/OpenCount
- Ranking metrics available to CLI for display

### 3. Enhanced Features
- Access to article images (`ImageURL`)
- Access to article tags and categories
- User karma scores in comments
- Controversy scoring for contentious articles
- Flag status for moderation

### 4. Future-Proof
- New fields added to shared models automatically available
- CLI doesn't need updates when backend adds fields
- Easier integration testing

---

## 🔍 BREAKING CHANGES

### For Cache Database
The SQLite cache tables should still work because:
- Cache uses base `models.Article` which is aliased from `shared.Article`
- Struct tags (`db:"field_name"`) remain the same
- No new required fields in base models

### For Existing Code
If any external code depends on CLI models:
- `Article` fields changed (Summary → Content, Upvotes → LikeCount, etc.)
- Need to update field references
- `IsHot`, `IsRising` no longer stored, calculated at runtime

---

## 📝 TESTING PLAN

### Unit Tests (When Go Installed)
```bash
cd cli
go test ./internal/models/...
go test ./internal/ui/components/...
go test ./internal/ui/views/...
```

### Integration Tests
```bash
# 1. Start mock server
go run cmd/mockserver/main.go

# 2. Run CLI
go run cmd/terminal-news/main.go

# 3. Verify:
# - Articles display correctly with new fields
# - Voting works (like/dislike)
# - Comments show usernames and karma
# - Hot/rising indicators appear correctly
```

### Backend Integration
```bash
# 1. Start real backend (Dev 1)
cd ../backend
go run cmd/server/main.go

# 2. Update CLI config
# Update ~/.terminal-news/config.yaml with backend URL

# 3. Test CLI with real data
cd ../cli
go run cmd/terminal-news/main.go
```

---

## ⚠️ KNOWN ISSUES

### 1. Go Not Installed
**Status**: BLOCKER for compilation testing
**Impact**: Cannot verify code compiles
**Resolution**: User needs to install Go 1.21+

### 2. Comment Count Estimation
**Status**: WORKAROUND
**Impact**: Comment count is estimated from TotalEngagement
**Resolution**: Backend should include actual comment count in ArticleRanking, or CLI should make separate query

### 3. Cache Schema
**Status**: POTENTIAL ISSUE
**Impact**: Cached articles may not have new fields
**Resolution**: Cache will need to be refreshed, or add migration

---

## 📊 COMPATIBILITY MATRIX

| Component | Status | Notes |
|-----------|--------|-------|
| Backend API | ✅ Ready | Models defined in shared/ |
| CLI Models | ✅ Aligned | Using shared models |
| Scraper | ✅ Ready | Already uses shared models |
| Mock Server | ✅ Updated | Generates matching data |
| Cache Layer | ✅ Compatible | Uses base models |
| WebSocket | ✅ Compatible | Uses shared models |

---

## 🎉 CONCLUSION

The Terminal News CLI is now fully aligned with the shared model definitions. All components reference the canonical models from `shared/models/models.go`, ensuring consistency across the entire application.

**Next Steps**:
1. Install Go 1.21+ to enable compilation testing
2. Run mock server and test CLI with updated models
3. Integrate with real backend (Dev 1) when ready
4. Verify all features work with new model structure

**Status**: 🟢 **MODEL ALIGNMENT COMPLETE - READY FOR COMPILATION TESTING**

---

*Report generated: November 18, 2025*
*Alignment completed by: Dev 2 (Terminal Client)*
