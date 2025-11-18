# Terminal News

> AM Radio for the Information Age

A native terminal-based news aggregator that combines community-driven content curation, local classifieds, and real-time weather. Think: CNN meets Reddit meets Craigslist, all running in your command line.

## Vision

Terminal News reimagines the newspaper experience for developers and terminal enthusiasts. Fast, focused, and community-driven—no ads, no tracking, no infinite scroll. Just the information you need, ranked by what the community values.

## Core Features

### News Aggregation & Curation
- **Hot** - Most liked articles
- **Controversial** - High engagement, mixed sentiment
- **Rising** - Rapidly trending stories
- **My Activity** - Your comments and interactions

### Community Engagement
- **Vote system**: Opens (+1), Likes (+2), Dislikes (-1)
- **Comment threads** on articles
- **Real-time ranking** based on community engagement

### Local Integration
- **NOAA Weather** - Always-visible weather for your location
- **Classifieds** - Local buy/sell/trade marketplace
- **Geographic feeds** - News and classifieds relevant to your area

## Design Philosophy

**Speed First**
- Terminal-native performance
- No images, no bloat
- Information at the speed of thought

**Community Driven**
- User votes determine visibility
- Not algorithmic manipulation
- Organic content discovery

**Local Focus**
- Weather from NOAA
- Classifieds for your community
- Regional news prioritization

**Retro Aesthetic**
- Ancient terminal look
- Monospace everything
- Like reading a newspaper from 1985, but smart

## Why Terminal News?

- **For developers**: Lives where you already work
- **For privacy advocates**: No tracking, minimal data collection
- **For focus seekers**: Distraction-free news consumption
- **For communities**: Local classifieds and connections

## Technical Overview

**Native CLI Application**
- Cross-platform terminal UI (Rust/Go)
- Connects to backend API for shared data
- Local caching for offline reading

**Backend Services**
- RESTful API for news, votes, comments
- Real-time WebSocket updates
- PostgreSQL for persistence
- Redis for ranking calculations

**Data Sources**
- RSS/Atom feeds from major news outlets
- News APIs (NewsAPI, Guardian, etc.)
- NOAA weather API
- User-submitted classifieds

## Project Status

**Current Phase**: Planning & Architecture

See `/docs` folder for detailed design documents.

## Getting Started

*Coming soon - MVP in development*

## Contributing

*Contribution guidelines coming soon*

## License

*TBD*

---

**"News at the speed of terminal"**
