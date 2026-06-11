# Catalog icons

One `<slug>.svg` per server, referenced from each manifest as
`icon: icons/<name>.svg`. These are presentation-only, repo-hosted (so they ride
the signed-index provenance), and **excluded from `RuntimeHash`** — adding or
changing an icon never diverges a baked image from the signed index, forces no
rebuild, and triggers no re-consent.

Each icon is a 64×64 rounded tile:

- **Brand glyph** — where [simple-icons](https://github.com/simple-icons/simple-icons)
  (pinned) ships the brand, its glyph is placed (white or black for contrast) on
  a tile filled with the brand's official hex. simple-icons icon files are
  CC0-1.0; the brand marks themselves remain the trademarks of their owners.
- **Monogram fallback** — for servers simple-icons doesn't carry (long-tail SaaS,
  plus flagships simple-icons has removed on trademark request), a self-authored
  monogram tile: a deterministic brand-ish background (hashed from the slug) with
  the displayName's initials.

The linter (`registryctl lint`) enforces that every `icon:` ref equals
`icons/<name>.svg` and that the file exists, so a ref can never dangle.
