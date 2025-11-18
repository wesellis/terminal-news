# Business Model

## Core Principle: "AM Radio for the Information Age"

Like AM radio stations, Terminal News is **free for listeners** but **monetized through local business relationships**. Users never pay for access; revenue comes from those who want to reach the community.

## Revenue Streams

### 1. Classified Ads (Primary Revenue)

**Free Tier:**
- Basic text listings
- 30-day duration
- Standard placement in feed
- No limit on number of posts

**Premium Listings: $5-25 per listing**
- Highlighted/featured placement
- Longer duration (60-90 days)
- Image support (rendered as ASCII art in terminal)
- Priority in search results
- "Bump to top" refresh option

**Categories:**
- Jobs
- Housing (rent, sale, roommates)
- Services (local contractors, tutors, etc.)
- For Sale (electronics, furniture, vehicles)
- Events (concerts, meetups, workshops)
- Gigs (short-term work)

**Revenue Projection:**
```
Conservative (Year 1):
- 5,000 active users
- 10% post classifieds = 500 listings/month
- 15% upgrade to premium = 75 × $10 avg = $750/month

Growth (Year 2):
- 25,000 active users
- 10% post = 2,500 listings/month
- 20% premium = 500 × $12 avg = $6,000/month

Mature (Year 3):
- 100,000 active users
- 15% post = 15,000 listings/month
- 25% premium = 3,750 × $15 avg = $56,250/month
```

### 2. Local Business Sponsorships

**City/Region Sponsors: $50-500/month**

**Sponsorship Placements:**
- Weather widget: "Weather powered by [Local Coffee Shop]"
- Section headers: "Hot News - Sponsored by [Bookstore]"
- Daily tip/quote: "Today's wisdom brought to you by [Yoga Studio]"

**Tiers:**
- Small town (<50k pop): $50-100/month
- Mid-size city (50k-500k): $150-300/month
- Major metro (500k+): $300-500/month

**Value Proposition:**
- Reach tech-savvy, employed demographic
- Local businesses supporting local community tool
- Non-intrusive, tasteful integration
- Geographic targeting

**Revenue Projection:**
```
Year 1: 20 cities × 2 sponsors avg × $100 = $4,000/month
Year 2: 100 cities × 3 sponsors avg × $150 = $45,000/month
Year 3: 500 cities × 5 sponsors avg × $200 = $500,000/month
```

### 3. Classified "Boost" Feature

**Pay-per-boost: $2-5**
- Bumps listing to top of feed
- Lasts 24 hours
- Can be used multiple times
- Perfect for time-sensitive listings

**Use Cases:**
- Urgent job posting
- Event happening soon
- Flash sale
- Quick sale needed

**Revenue Projection:**
```
Year 1: 100 boosts/month × $3 = $300/month
Year 2: 1,000 boosts/month × $3 = $3,000/month
Year 3: 5,000 boosts/month × $4 = $20,000/month
```

### 4. API Access (Future)

**Free Tier:**
- 100 requests/day
- Personal use
- Rate limited

**Business Tier: $50-200/month**
- Higher rate limits
- Bulk posting classifieds
- Analytics access
- Priority support

**Enterprise Tier: Custom pricing**
- White-label options
- Custom integrations
- Dedicated support

### 5. Premium User Features (Optional, Future)

**$3-5/month per user**
- Ad-free experience (no sponsor messages)
- Advanced filters
- Saved searches with alerts
- Custom news sources
- Priority customer support

**Note:** Only implement if needed. Core should remain free.

## Total Revenue Projections

### Year 1 (Conservative)
```
Classifieds:     $750/month
Sponsorships:    $4,000/month
Boosts:          $300/month
-------------------------
Total:           $5,050/month ($60,600/year)

Costs:           $500/month (server, APIs, domain)
Net:             $4,550/month ($54,600/year)
```

### Year 2 (Growth)
```
Classifieds:     $6,000/month
Sponsorships:    $45,000/month
Boosts:          $3,000/month
API Access:      $2,000/month
-------------------------
Total:           $56,000/month ($672,000/year)

Costs:           $2,000/month (scaled infrastructure, 1 part-time mod)
Net:             $54,000/month ($648,000/year)
```

### Year 3 (Mature)
```
Classifieds:     $56,250/month
Sponsorships:    $500,000/month
Boosts:          $20,000/month
API Access:      $10,000/month
Premium Users:   $15,000/month
-------------------------
Total:           $601,250/month ($7.2M/year)

Costs:           $50,000/month (team of 3-5, infrastructure, support)
Net:             $551,250/month ($6.6M/year)
```

## Cost Structure

### Year 1
- **Server/Infrastructure**: $200/month (Digital Ocean)
- **APIs**: $100/month (NewsAPI, weather)
- **Domain/SSL**: $20/month
- **Monitoring**: $50/month
- **Payment Processing**: 3% of revenue
- **Miscellaneous**: $130/month
- **Total**: ~$500/month

### Year 2
- Infrastructure: $1,000/month (scaled)
- APIs: $300/month
- Part-time moderator: $500/month
- Payment processing: 3%
- Tools & services: $200/month
- **Total**: ~$2,000/month

### Year 3
- Infrastructure: $5,000/month
- Team (3-5 people): $35,000/month
- APIs & services: $2,000/month
- Marketing: $3,000/month
- Legal/accounting: $2,000/month
- Miscellaneous: $3,000/month
- **Total**: ~$50,000/month

## Competitive Analysis

### vs. Craigslist
**Craigslist Strengths:**
- Established brand (28 years)
- Massive user base
- Simple, fast

**Terminal News Advantages:**
- Modern tech stack
- Community curation (voting)
- Integrated news + classifieds
- Better UX for developers
- Real-time updates

**Craigslist Revenue:** ~$1B/year (mostly job postings)
**Our Target:** 1% of Craigslist's market in tech communities

### vs. Reddit
**Reddit Strengths:**
- Massive community
- Rich features
- Well-known

**Terminal News Advantages:**
- Terminal-native (faster for CLI users)
- No infinite scroll
- Local focus
- Cleaner signal-to-noise
- Actual classifieds marketplace

### vs. Hacker News
**HN Strengths:**
- Tech credibility
- Smart community
- No monetization pressure

**Terminal News Advantages:**
- Terminal interface
- Multi-category (not just tech)
- Classifieds
- Local weather/info
- Vote weighting system

## Market Opportunity

**Target Addressable Market:**
- Developers: 27M worldwide
- Terminal users: ~5M active daily
- Tech-savvy CLI enthusiasts: 1M+

**Initial Target (Year 1):**
- 5,000 daily active users
- 50 cities with critical mass

**Growth Target (Year 3):**
- 100,000 daily active users
- 500+ cities

**Market Size:**
- Online classifieds market: $20B/year globally
- News aggregation: $5B/year
- Terminal tool market: Growing (see: Warp, Fig funding)

## Go-to-Market Strategy

### Phase 1: Launch (Months 0-3)
**Target:** Tech early adopters

**Channels:**
1. **Hacker News** - Launch post + Show HN
2. **Reddit** - r/commandline, r/linux, r/programming
3. **Product Hunt** - Dev tools category
4. **Twitter/X** - Developer influencers
5. **Dev.to** - Blog post about building it

**Goal:** 1,000 users, 10 cities

### Phase 2: Growth (Months 4-12)
**Target:** Broader developer community

**Channels:**
1. Word of mouth (referral system)
2. GitHub stars/trending
3. Dev YouTube channels
4. Tech podcasts (Changelog, etc.)
5. Conference talks/demos

**Partnerships:**
- Terminal emulator developers (Warp, Hyper)
- Developer tool companies
- Local tech meetups

**Goal:** 10,000 users, 50 cities

### Phase 3: Scale (Year 2+)
**Target:** Mainstream tech users

**Channels:**
1. Local business partnerships
2. University CS programs
3. Tech company internal tools
4. Press coverage (TechCrunch, Ars Technica)

**Goal:** 100,000 users, 500 cities

## Key Success Metrics

**User Engagement:**
- Daily active users (DAU)
- Session length
- Articles read per session
- Vote rate
- Comment rate

**Marketplace Health:**
- Classified posts per day
- Premium conversion rate
- Time to first response (classifieds)
- Successful transactions

**Revenue Metrics:**
- Monthly recurring revenue (MRR)
- Customer acquisition cost (CAC)
- Lifetime value (LTV)
- Premium conversion rate
- Sponsor retention rate

**Community Health:**
- User retention (30-day, 90-day)
- Geographic distribution
- Content quality (moderation queue size)
- Response time on classifieds

## Risk Analysis

### Risks & Mitigations

**1. Chicken-and-egg (classifieds need users, users need classifieds)**
- Mitigation: Launch with news first, add classifieds at 500+ users per city
- Pre-seed classifieds with job postings from APIs

**2. Moderation burden (spam, scams, abuse)**
- Mitigation: Start with user reports + basic filters
- Community moderation (trust system)
- Hire moderator at $5k/month revenue

**3. News API costs/limits**
- Mitigation: Heavy use of free RSS feeds
- Cache aggressively
- Community-submitted links

**4. Competition from incumbents**
- Mitigation: Stay focused on terminal niche
- Move fast on features
- Build loyal community

**5. Low adoption (terminal users are niche)**
- Mitigation: Realistic expectations
- Sustainable indie project mindset
- Can still be profitable at 10k users

## Exit Strategy (Optional)

**Potential Acquirers:**
- Terminal emulator companies (Warp, Hyper)
- Developer tool companies (JetBrains, GitHub)
- News aggregators
- Local classifieds platforms

**Alternative:** Sustainable indie business
- No VC needed
- Stay small, profitable
- Serve community indefinitely

## Summary

Terminal News has a **proven business model** (Craigslist's classified revenue) applied to an **underserved niche** (terminal users). Revenue comes from those who want to reach the community (businesses, classified posters), while users get a free, fast, focused news experience.

The path to profitability is clear:
1. Build user base with free news aggregation
2. Activate classifieds when cities hit critical mass
3. Add local business sponsorships
4. Scale geographically

**Conservative estimate:** Profitable within 6 months at just 5,000 active users.

**This is not a unicorn startup.** It's a sustainable, community-focused business that can support a small team while providing real value to terminal users worldwide.
