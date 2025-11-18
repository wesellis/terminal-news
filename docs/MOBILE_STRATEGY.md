# Mobile App Strategy (Future Phase)

## Core Principle: Terminal Authenticity on Mobile

**The mobile app must feel like a real terminal, even on a touchscreen.**

---

## Design Philosophy

### Non-Negotiables:
- **Monospace font** - Always
- **Terminal aesthetic** - Green-on-black, amber, or classic terminal themes
- **Text-first** - No modern app bloat
- **Keyboard-inspired** - Swipe gestures that mimic terminal commands
- **Authentic feel** - Should feel like SSH-ing from your phone

### What This Looks Like:

**Portrait Mode (Vertical):**
```
┌─────────────────────────┐
│ $ terminal-news         │
│ ========================│
│                         │
│ [Hot] Controversial Ris│
│                         │
│ ▲142 ▼12  [1.2k] 2h    │
│ New AI Model Breaks...  │
│ TechCrunch             │
│ [O][L][D][C]           │
│                         │
│ ▲98 ▼3  [890] 4h       │
│ Major Climate Agree... │
│ Reuters                │
│ [O][L][D][C]           │
│                         │
│ ▲76 ▼45  [1.8k] 1h 🔥  │
│ Controversial Tech CEO │
│ The Guardian           │
│ [O][L][D][C]           │
│                         │
│ ════════════════════════│
│ [↑][↓][Tab][?][Q]      │
└─────────────────────────┘
```

**Key Features:**
- Compact article cards
- Terminal-style buttons (keyboard letters)
- Monospace everywhere
- Minimal graphics
- Fast scrolling
- Gesture support (but looks like commands)

---

## Interaction Design

### Gestures That Feel Like Terminal:

**Swipe Right:** Like (L key)
**Swipe Left:** Dislike (D key)
**Tap:** Select/Open (Enter key)
**Long Press:** Options menu (like right-click)
**Two-finger swipe up:** Refresh (R key)
**Shake:** Help menu (? key)

**Visual Feedback:**
When you swipe right, show: `> like` in terminal text
When you swipe left, show: `> dislike`
When you tap: `> open article`

---

## UI Elements (Terminal-Authentic)

### Buttons:
```
Not this: [Like 👍]
But this: [L]ike  or  > l
```

### Navigation:
```
Not this: Bottom nav with icons
But this:
┌─────────────────────────┐
│ :tabs                   │
│ 1. hot                  │
│ 2. controversial        │
│ 3. rising              │
│ 4. profile             │
│ 5. weather             │
└─────────────────────────┘
```

### Typing/Search:
```
Actual command prompt feel:
$ search: bitcoin
```

---

## Technical Approach

### Stack Options:

**Option 1: React Native + Terminal UI Library**
- React Native for cross-platform
- Custom terminal-themed components
- Monospace font package
- Haptic feedback for "typing" feel

**Option 2: Flutter + Custom Terminal Theme**
- Flutter for smooth performance
- Custom terminal widgets
- Great for complex UI while staying lightweight

**Option 3: Native (Swift/Kotlin) + Terminal Emulator Libraries**
- Most authentic terminal feel
- Use actual terminal rendering libraries
- More work, but best result

**Recommendation:** Flutter + custom terminal theme (best balance)

---

## Features for Mobile

### Must Have (Day 1):
- Browse news (Hot, Controversial, Rising)
- Vote (like/dislike)
- Read articles (in-app browser styled like terminal)
- Comments (read and post)
- Weather widget
- Profile view

### Nice to Have (Later):
- Classifieds browsing
- Post classified (mobile-optimized form)
- Notifications (terminal-style alerts)
- Offline mode
- Dark/light terminal themes

### Mobile-Specific:
- Swipe gestures
- Haptic feedback
- Share sheet integration
- Widget for home screen (terminal-styled)

---

## Themes

### Classic Terminal Themes:

**Green on Black:**
```
Background: #000000
Text: #00FF00
Accent: #00AA00
```

**Amber:**
```
Background: #000000
Text: #FFAA00
Accent: #CC8800
```

**White on Blue (IBM):**
```
Background: #0000AA
Text: #FFFFFF
Accent: #AAAAAA
```

**Retro (C64):**
```
Background: #40318D
Text: #7869C4
Accent: #A3A3FF
```

Users can switch themes like terminal color schemes.

---

## Typography

**Font:** Monospace REQUIRED

**Options:**
- JetBrains Mono
- Fira Code
- Source Code Pro
- SF Mono (iOS)
- Roboto Mono (Android)

**Size:**
- Articles: 13-14pt
- Headlines: 16pt (bold)
- Buttons: 14pt

**Line Height:** 1.4 (terminal standard)

---

## Animations

### Terminal-Authentic Animations:

**NOT:** Smooth fades, bouncy transitions
**YES:**
- Blink cursor when typing
- Text appearing character-by-character (fast)
- Terminal "clear screen" effect
- Scanline effect (subtle)

**Example:**
When opening article:
```
$ open article_142
> loading...
> ████████████ 100%
> done.
[Article appears]
```

---

## Notifications

### Style:
```
┌─────────────────────────┐
│ terminal-news           │
├─────────────────────────┤
│ > new reply             │
│ @techie_sam replied:    │
│ "Great point!"          │
│                         │
│ [V]iew  [D]ismiss       │
└─────────────────────────┘
```

Not a normal phone notification - looks like a terminal alert.

---

## Launch Strategy

### Phase 1: Mobile Web (PWA)
**Months 6-12**
- Responsive web version
- Install as PWA
- Terminal styling on mobile web
- No app store needed yet
- Test mobile UX

### Phase 2: Native App
**Year 2**
- iOS app (TestFlight beta)
- Android app (beta)
- Refine based on PWA learnings
- Terminal aesthetic fully realized

### Phase 3: App Store Launch
**Year 2+**
- Polish and submit
- App Store / Play Store
- Marketing push

---

## Differentiation

### Why Mobile Matters:

**Most terminal users also have phones.**

They'll want:
- Check news on commute
- Quick classified browsing
- Respond to comments
- But they still want the terminal vibe

**Mobile app = convenience**
**Terminal app = power user mode**

Both should feel like the same product, just different contexts.

---

## Monetization (Mobile)

**Same as desktop:**
- Free to use
- Premium classifieds
- Boosts
- Sponsorships (show in mobile too)

**Mobile-specific:**
- In-app purchases for premium (Apple/Google take 30% - factor in)
- Or: Direct people to web for payment (avoid fee)

---

## Development Timeline

**Don't build this yet.**

**Priority:**
1. Desktop terminal app (Months 1-6)
2. Web version/API (Months 4-9)
3. Mobile-responsive web (Months 9-12)
4. Native mobile app (Year 2)

Mobile should come AFTER:
- 10k+ users on desktop
- Proven business model
- Users requesting mobile

---

## Technical Notes

### Terminal Rendering on Mobile:

**Challenges:**
- Small screen (need compact layout)
- Touch vs keyboard (gestures must replace shortcuts)
- Variable screen sizes (iPhone Mini vs iPad)
- Performance (terminal rendering can be heavy)

**Solutions:**
- Simplified terminal layout for mobile
- Smart defaults (no 80-column restrictions)
- Efficient rendering (React/Flutter optimization)
- Adaptive layout (works on all sizes)

---

## Design Mockup (Conceptual)

### Home Screen Widget:
```
┌─────────────────┐
│ terminal-news   │
├─────────────────┤
│ Top Story:      │
│ ▲142 New AI... │
│ TechCrunch      │
│                 │
│ $ tap to open   │
└─────────────────┘
```

Looks like a terminal window on your home screen.

---

## The Important Part

**The mobile app should feel like you're SSH-ing into Terminal News from your phone.**

It's not a "mobile-first modern app."
It's a **terminal emulator for news.**

If someone sees it and thinks "that looks like code" - we've succeeded.

---

## Launch Checklist (When Ready)

- [ ] Desktop app has 10k+ users
- [ ] API is stable and documented
- [ ] Mobile-responsive web works
- [ ] Terminal aesthetic translated to mobile
- [ ] Gestures feel natural
- [ ] Performance is smooth (60fps)
- [ ] Beta testers approve
- [ ] App Store guidelines met
- [ ] Marketing materials ready

**Timeline: Year 2+**

For now: Focus on desktop. Mobile will come when it's ready.
