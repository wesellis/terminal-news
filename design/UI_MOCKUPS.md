# Terminal News - UI Mockups

## Design Principles

- **Monospace aesthetic** - Embrace the terminal
- **ASCII art borders** - Retro newspaper feel
- **Information density** - Maximize content, minimize chrome
- **Keyboard-first** - Every action has a shortcut
- **Responsive layouts** - Adapt to terminal size

---

## Main View: Hot News Feed

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ TERMINAL NEWS                    [Hot] Controversial Rising Profile Weather │
├─────────────────────────────────────────────────────────────────────────────┤
│ San Francisco, CA  │  72°F Partly Cloudy  │  Updated: 2:34 PM              │
│ Weather powered by Local Coffee Co.                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ▲ 142 ▼ 12  [1.2k views]  2h ago  │  TechCrunch                            │
│ New AI Model Breaks Reasoning Benchmarks                                    │
│ Researchers at Stanford unveil GPT-5 with improved capabilities...          │
│ [O]pen [L]ike [D]islike [C]omment (45)                                      │
│                                                                              │
│ ▲ 98 ▼ 3  [890 views]  4h ago  │  Reuters                                  │
│ Major Climate Agreement Reached at UN Summit                                │
│ World leaders commit to carbon neutrality by 2050...                        │
│ [O]pen [L]ike [D]islike [C]omment (23)                                      │
│                                                                              │
│ ▲ 76 ▼ 45  [1.8k views]  1h ago  │  The Guardian                           │
│ Controversial Tech CEO Steps Down Amid Scandal                              │
│ Silicon Valley titan faces allegations of workplace misconduct...           │
│ [O]pen [L]ike [D]islike [C]omment (156) 🔥                                  │
│                                                                              │
│ ▲ 54 ▼ 2  [420 views]  30m ago  │  BBC                                     │
│ Space Mission Discovers Water on Distant Moon                               │
│ NASA's latest probe sends back promising data...                            │
│ [O]pen [L]ike [D]islike [C]omment (8) ⚡                                    │
│                                                                              │
│ ▲ 32 ▼ 1  [245 views]  15m ago  │  NPR                                     │
│ Local Communities Embrace Renewable Energy                                  │
│ Small towns leading the charge in sustainable power...                      │
│ [O]pen [L]ike [D]islike [C]omment (4)                                       │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [Tab] Switch Tab  [R]efresh  [/] Search  [?] Help  [Q] Quit│
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Tab navigation at top (active tab highlighted)
- Weather widget always visible
- Vote counts (upvotes/downvotes)
- View count and age
- Source attribution
- Quick action shortcuts
- Fire emoji for controversial, lightning for rising
- Status bar with keyboard shortcuts

---

## Controversial Feed

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ TERMINAL NEWS                    Hot [Controversial] Rising Profile Weather │
├─────────────────────────────────────────────────────────────────────────────┤
│ Showing articles with high engagement and mixed sentiment                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ▲ 234 ▼ 218  [3.2k views]  5h ago  │  Politico                             │
│ New Policy Divides Nation on Immigration Reform                             │
│ Heated debate emerges over proposed border security measures...             │
│ [O]pen [L]ike [D]islike [C]omment (342) 🔥🔥🔥                              │
│ Controversy Score: 93/100                                                   │
│                                                                              │
│ ▲ 156 ▼ 142  [2.1k views]  3h ago  │  Tech Insider                         │
│ Apple's Latest iPhone Sparks Heated Reviews                                 │
│ Users split on whether new features justify the price...                    │
│ [O]pen [L]ike [D]islike [C]omment (189) 🔥🔥                                │
│ Controversy Score: 89/100                                                   │
│                                                                              │
│ ▲ 92 ▼ 87  [1.5k views]  2h ago  │  Sports Network                         │
│ Star Athlete's Comments on League Rules Draw Backlash                       │
│ Social media erupts following controversial interview...                    │
│ [O]pen [L]ike [D]islike [C]omment (98) 🔥                                   │
│ Controversy Score: 84/100                                                   │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [Tab] Switch Tab  [R]efresh  [?] Help  [Q] Quit            │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Controversy score visible
- Fire emoji intensity based on score
- Mixed vote counts highlighted

---

## Rising Feed

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ TERMINAL NEWS                    Hot Controversial [Rising] Profile Weather │
├─────────────────────────────────────────────────────────────────────────────┤
│ Articles rapidly gaining traction in the last hour                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ▲ 45 ▼ 1  [320 views]  12m ago  │  Breaking News Wire                     │
│ BREAKING: Major Earthquake Strikes Pacific Region                           │
│ Initial reports suggest magnitude 7.2, tsunami warnings issued...           │
│ [O]pen [L]ike [D]islike [C]omment (12) ⚡⚡⚡ +38 in last hour               │
│                                                                              │
│ ▲ 32 ▼ 0  [180 views]  8m ago  │  Crypto News                              │
│ Bitcoin Surges 15% Following Regulatory Announcement                        │
│ Cryptocurrency markets rally after surprise policy shift...                 │
│ [O]pen [L]ike [D]islike [C]omment (7) ⚡⚡ +32 in last hour                  │
│                                                                              │
│ ▲ 28 ▼ 2  [210 views]  18m ago  │  Science Daily                           │
│ Researchers Announce Breakthrough in Fusion Energy                          │
│ New reactor design shows promise for sustainable power...                   │
│ [O]pen [L]ike [D]islike [C]omment (14) ⚡⚡ +24 in last hour                 │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [Tab] Switch Tab  [R]efresh  [?] Help  [Q] Quit            │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Lightning bolt intensity for velocity
- Vote gain in last hour shown
- Recent timestamps emphasized

---

## Article Detail + Comments

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ ← Back to Feed                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ New AI Model Breaks Reasoning Benchmarks                                    │
│ TechCrunch  │  2 hours ago  │  ▲ 142 ▼ 12  │  1.2k views                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Researchers at Stanford University have unveiled a new artificial           │
│ intelligence model that significantly outperforms existing systems on       │
│ complex reasoning tasks. The model, dubbed "ReasonGPT," demonstrates        │
│ unprecedented capabilities in mathematical problem-solving and logical      │
│ deduction...                                                                │
│                                                                              │
│ [Full article: https://techcrunch.com/2024/...]                            │
│                                                                              │
│ [O]pen in Browser  [L]ike  [D]islike  [Shift+C] Comment                    │
│                                                                              │
├─── Comments (45) ────────────────────────────────────────────────────────────┤
│                                                                              │
│ @techie_sam  •  1h ago  •  ▲ 24                                            │
│ This is genuinely impressive. I tested it on some Leetcode hard problems   │
│ and it got 9/10 correct on first try. Game changer for code assistance.    │
│   [R]eply  [L]ike                                                           │
│                                                                              │
│   └─ @dev_alice  •  45m ago  •  ▲ 12                                       │
│      Did you try the edge cases? I found it struggles with unusual inputs. │
│      [R]eply  [L]ike                                                        │
│                                                                              │
│      └─ @techie_sam  •  30m ago  •  ▲ 8                                    │
│         Good point! I should test more thoroughly. Will report back.       │
│         [R]eply  [L]ike                                                     │
│                                                                              │
│ @skeptical_dev  •  2h ago  •  ▲ 6 ▼ 2                                      │
│ Benchmarks are one thing, but real-world performance is what matters.      │
│ Let's see how it handles production code reviews.                          │
│   [R]eply  [L]ike                                                           │
│                                                                              │
│ @ai_researcher  •  1h ago  •  ▲ 18                                         │
│ I work adjacent to this team. The architecture is fascinating - they're    │
│ using a novel attention mechanism that I think will be widely adopted.     │
│   [R]eply  [L]ike                                                           │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [C] New Comment  [R] Reply  [Esc] Back  [Q] Quit           │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Article summary/excerpt
- Link to full article
- Threaded comments with indentation
- Vote counts on comments
- Quick reply functionality
- Timestamp on comments

---

## Profile / My Activity

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ TERMINAL NEWS                    Hot Controversial Rising [Profile] Weather │
├─────────────────────────────────────────────────────────────────────────────┤
│ @wesley_dev  •  Member since Jan 2025  •  San Francisco, CA                │
│ ▲ 1,234 karma  •  156 comments  •  89 posts                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ┌─ Recent Activity ──────────────────────────────────────────────────────┐  │
│ │                                                                         │  │
│ │ Commented 2h ago on:                                                   │  │
│ │ "New AI Model Breaks Reasoning Benchmarks"                             │  │
│ │ "This is a game changer for developers..."  ▲ 12                      │  │
│ │                                                                         │  │
│ │ Liked 3h ago:                                                          │  │
│ │ "Major Climate Agreement Reached at UN Summit"                         │  │
│ │                                                                         │  │
│ │ Posted Classified 1d ago:                                              │  │
│ │ "Selling: MacBook Pro M3 - $2000 OBO"                                  │  │
│ │ Views: 45  •  Responses: 3                                             │  │
│ │                                                                         │  │
│ │ Commented 1d ago on:                                                   │  │
│ │ "Controversial Tech CEO Steps Down"                                    │  │
│ │ "About time. The evidence was clear."  ▲ 34 ▼ 12                      │  │
│ │                                                                         │  │
│ └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│ ┌─ Your Classifieds ─────────────────────────────────────────────────────┐  │
│ │                                                                         │  │
│ │ MacBook Pro M3 - $2000 OBO  •  [ACTIVE]                                │  │
│ │ Electronics > Computers  •  Posted 1 day ago  •  45 views              │  │
│ │ [E]dit  [B]oost ($3)  [D]elete                                         │  │
│ │                                                                         │  │
│ └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [E] Edit Profile  [S] Settings  [Esc] Back  [Q] Quit       │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- User stats at top
- Activity timeline
- Your classified listings
- Quick actions on classifieds

---

## Weather Widget (Expanded)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ TERMINAL NEWS                    Hot Controversial Rising Profile [Weather] │
├─────────────────────────────────────────────────────────────────────────────┤
│                      WEATHER - SAN FRANCISCO, CA                            │
│                     Powered by Local Coffee Co.                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Current Conditions  │  2:34 PM PST  │  Updated 5 min ago                  │
│                                                                              │
│        .--.           72°F  Partly Cloudy                                   │
│     .-(    ).         Feels like: 70°F                                      │
│    (___.__)__)        Humidity: 65%                                         │
│                       Wind: 8 mph NW                                        │
│                       Pressure: 30.12 in                                    │
│                                                                              │
│ ┌─ 5-Day Forecast ─────────────────────────────────────────────────────┐    │
│ │                                                                       │    │
│ │ Wed      Thu      Fri      Sat      Sun                              │    │
│ │  ☀       ⛅       ☁       🌧      ⛅                                  │    │
│ │ 75°/62°  72°/60°  68°/58°  65°/56°  70°/59°                          │    │
│ │                                                                       │    │
│ └───────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│ ┌─ Alerts ──────────────────────────────────────────────────────────────┐    │
│ │ No active weather alerts for your area                                │    │
│ └───────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│ Data from NOAA National Weather Service                                     │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [L] Change Location  [Esc] Back  [Q] Quit                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- ASCII art weather icons
- Current conditions detail
- 5-day forecast
- Weather alerts
- Sponsor attribution

---

## Classifieds Browse

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASSIFIEDS  │  [All] Jobs Housing For Sale Services Events                │
├─────────────────────────────────────────────────────────────────────────────┤
│ San Francisco Bay Area  •  Sort: Recent  │  [/] Search  [F] Filter  [P]ost │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ⭐ PREMIUM  •  15m ago  •  San Francisco                                    │
│ Senior React Developer - Remote OK - $150k-200k                             │
│ Fast-growing startup seeking experienced frontend engineer...               │
│ Jobs > Software Engineering  │  @tech_recruiter                             │
│ [O]pen  [S]ave  [R]eport                                                    │
│                                                                              │
│ 2h ago  •  Oakland                                                          │
│ 2BR Apartment in Lake Merritt - $2800/mo                                    │
│ Beautiful renovated unit with lake views, parking included...               │
│ Housing > Apartments  │  @landlord_jane                                     │
│ [O]pen  [S]ave  [R]eport                                                    │
│                                                                              │
│ ⭐ PREMIUM  •  4h ago  •  Berkeley                                          │
│ MacBook Pro M3 Max 16" - Like New - $2500                                   │
│ Barely used, still under warranty. Original box and accessories...          │
│ For Sale > Electronics  │  @seller_bob                                      │
│ [O]pen  [S]ave  [R]eport                                                    │
│                                                                              │
│ 1d ago  •  San Jose                                                         │
│ Guitar Lessons - Beginner to Advanced - $50/hr                              │
│ Professional musician with 15 years teaching experience...                  │
│ Services > Music  │  @guitar_pro                                            │
│ [O]pen  [S]ave  [R]eport                                                    │
│                                                                              │
│ 1d ago  •  San Francisco                                                    │
│ Tech Meetup - AI & Machine Learning - Thu 7pm                               │
│ Monthly gathering of AI enthusiasts. Free pizza and drinks...               │
│ Events > Technology  │  @ai_meetup_sf                                       │
│ [O]pen  [S]ave  [R]eport                                                    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [Tab] Category  [P] Post New  [?] Help  [Q] Quit           │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Category tabs
- Premium listings highlighted (⭐)
- Location and timestamp
- Category breadcrumb
- User attribution
- Quick actions

---

## Post Classified Form

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ POST CLASSIFIED                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Category: [For Sale ▾]                                                      │
│           > Electronics                                                     │
│                                                                              │
│ Title:                                                                      │
│ ┌─────────────────────────────────────────────────────────────────────────┐ │
│ │ MacBook Pro M3 - Excellent Condition                                    │ │
│ └─────────────────────────────────────────────────────────────────────────┘ │
│ (50 characters remaining)                                                   │
│                                                                              │
│ Description:                                                                │
│ ┌─────────────────────────────────────────────────────────────────────────┐ │
│ │ 2023 MacBook Pro M3 14" in excellent condition.                         │ │
│ │                                                                          │ │
│ │ Specs:                                                                   │ │
│ │ - M3 chip (8-core CPU, 10-core GPU)                                     │ │
│ │ - 16GB RAM                                                               │ │
│ │ - 512GB SSD                                                              │ │
│ │ - AppleCare+ until 2026                                                 │ │
│ │                                                                          │ │
│ │ Includes original box, charger, and USB-C cable. Minor wear on bottom.  │ │
│ │_                                                                         │ │
│ └─────────────────────────────────────────────────────────────────────────┘ │
│ (380 characters remaining)                                                  │
│                                                                              │
│ Price: [$] [ 2000    ]                                                      │
│                                                                              │
│ Location: [ San Francisco, CA          ]                                    │
│                                                                              │
│ Contact: [✓] Email (wes***@***.com)                                         │
│          [ ] Phone                                                          │
│          [ ] Direct message                                                 │
│                                                                              │
│ ┌─ Premium Listing ($10) ────────────────────────────────────────────────┐  │
│ │ [ ] Make this a premium listing                                        │  │
│ │     • Highlighted with ⭐ in feed                                       │  │
│ │     • 3x visibility                                                     │  │
│ │     • 60-day duration (vs 30)                                           │  │
│ └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
│ [Tab] Next Field  [Shift+Tab] Previous  [Ctrl+S] Preview  [Ctrl+P] Post    │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [Ctrl+P] Post  [Ctrl+S] Save Draft  [Esc] Cancel                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Category dropdown
- Character limits shown
- Multiple contact methods
- Premium upgrade option
- Draft saving
- Preview before posting

---

## Help / Keyboard Shortcuts

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ HELP & KEYBOARD SHORTCUTS                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ┌─ Navigation ────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  ↑/↓ or J/K     Navigate items                                          │ │
│ │  Tab            Switch tabs/sections                                    │ │
│ │  Shift+Tab      Previous tab/section                                    │ │
│ │  Enter          Open selected item                                      │ │
│ │  Esc            Go back / Close                                         │ │
│ │  Q              Quit application                                        │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ News & Voting ─────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  O              Open article in browser                                 │ │
│ │  L              Like article                                            │ │
│ │  D              Dislike article                                         │ │
│ │  C              View/add comments                                       │ │
│ │  R              Refresh feed                                            │ │
│ │  /              Search articles                                         │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ Classifieds ───────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  P              Post new classified                                     │ │
│ │  F              Filter classifieds                                      │ │
│ │  S              Save/bookmark classified                                │ │
│ │  B              Boost your classified                                   │ │
│ │  E              Edit your classified                                    │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ Other ─────────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  ?              Show this help                                          │ │
│ │  Ctrl+S         Settings                                                │ │
│ │  Ctrl+R         Force refresh                                           │ │
│ │  Ctrl+C         Copy current URL                                        │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [Esc] Close Help                                                            │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Settings Panel

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ SETTINGS                                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ ┌─ Account ───────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  Username:        wesley_dev                                            │ │
│ │  Email:           wes***@***.com                       [Change]         │ │
│ │  Password:        ****************                     [Change]         │ │
│ │  Location:        San Francisco, CA                    [Change]         │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ Display ───────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  Theme:           [•] Classic  [ ] Dark  [ ] Light                      │ │
│ │  Compact Mode:    [✓] Enable                                            │ │
│ │  Show Emojis:     [✓] Enable (🔥⚡)                                     │ │
│ │  Font Size:       [─────•────] Medium                                   │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ Content ───────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  Auto-refresh:    [ ] Enabled (every ____ minutes)                      │ │
│ │  NSFW Content:    [ ] Show  [•] Hide  [ ] Blur                          │ │
│ │  Default Tab:     [Hot ▾]                                               │ │
│ │  Articles/Page:   [50 ▾]                                                │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│ ┌─ Privacy ───────────────────────────────────────────────────────────────┐ │
│ │                                                                          │ │
│ │  Analytics:       [ ] Allow anonymous usage data                        │ │
│ │  Read History:    [✓] Save locally                                      │ │
│ │  Location Data:   [•] Local only  [ ] Share for better results          │ │
│ │                                                                          │ │
│ └──────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [S] Save  [R] Reset to Defaults  [Esc] Cancel                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Login Screen

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│                                                                              │
│                        ████████╗███╗   ██╗                                  │
│                        ╚══██╔══╝████╗  ██║                                  │
│                           ██║   ██╔██╗ ██║                                  │
│                           ██║   ██║╚██╗██║                                  │
│                           ██║   ██║ ╚████║                                  │
│                           ╚═╝   ╚═╝  ╚═══╝                                  │
│                                                                              │
│                          TERMINAL NEWS                                      │
│                   AM Radio for the Information Age                          │
│                                                                              │
│                                                                              │
│                    ┌─────────────────────────────┐                          │
│                    │ Username:                   │                          │
│                    │ ┌─────────────────────────┐ │                          │
│                    │ │ wesley_dev_             │ │                          │
│                    │ └─────────────────────────┘ │                          │
│                    │                             │                          │
│                    │ Password:                   │                          │
│                    │ ┌─────────────────────────┐ │                          │
│                    │ │ ******************      │ │                          │
│                    │ └─────────────────────────┘ │                          │
│                    │                             │                          │
│                    │        [L] Login            │                          │
│                    │     [R] Register            │                          │
│                    │  [F] Forgot Password        │                          │
│                    │                             │                          │
│                    └─────────────────────────────┘                          │
│                                                                              │
│                                                                              │
│                         Version 1.0.0                                       │
│                                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ [Q] Quit                                                                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Design Notes

**Colors (Terminal):**
- Headings: Bold white
- Normal text: Gray/white
- Highlights: Bright white
- Links: Cyan
- Success: Green
- Warning: Yellow
- Error: Red
- Premium: Gold/yellow

**Responsive Behavior:**
- Min width: 80 columns
- Ideal width: 100-120 columns
- Graceful degradation for smaller terminals
- Horizontal scrolling if needed

**Accessibility:**
- All features keyboard accessible
- Clear focus indicators
- High contrast mode option
- Screen reader friendly (plain text)

**Performance:**
- Lazy load comments
- Virtual scrolling for long lists
- Debounced input
- Cached rendering
