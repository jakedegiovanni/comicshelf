# Comicshelf

...

## Todo

- view this weeks releases (complete)
- follow certain series (complete)
- view when issues in series will be released (complete)
    - view all comics within a series (complete)
    - show marvel unlimited release date on comic card (complete)
- series page doesn't show all items in series
    - add page size (complete)
    - support pagination on the pages
- be able to see which series you follow
- be notified when issues in a series are released
- ignore results.json, better caching of results (complete)
- cache limit and eviction
- change weekly (complete)
- better error handling
- better logging
    - use slog everywhere (complete)
    - no more os.exit from non root paths
- htmx to enable better html structure (in progress)
    - comic-card : follow / unfollow (complete)
    - navbar
    - don't use cdn
    - accessibility
- middleware for enforcing date query parameter on marvel endpoints
- in-mem db persists beyond restarts (complete)
- real db? object storage sufficient? something on filesystem enough?
- efficient network usage, lots of network requests happening with html setup as it is
- makefile supports build for different platforms
- deploy to aws
- support more than just marvel unlimited
- Middleware should have custom futures to try and reduce and cloning and split box pins into methods