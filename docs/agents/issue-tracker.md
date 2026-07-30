# Issue Tracker

- Backend: local Markdown
- Feature root pattern: `.scratch/<feature-slug>/`
- Specification: `<feature-root>/spec.md`
- Tickets: `<feature-root>/issues/<NN>-<slug>.md`, numbered from `01`
- Coordination: record status and explicit `blocks` / `blocked by` edges in each ticket

Resolve `<feature-root>` in this order:

1. Use an explicit feature root or slug from the current request.
2. Reuse the feature root of a spec or ticket referenced in the current context.
3. Use the single relevant existing feature root when repository evidence makes it unique.
4. For a new feature, derive a lowercase kebab-case slug from the agreed feature name.

Ask the user when multiple existing roots are plausible or a derived slug collides with a different feature. `<feature-slug>` is a pattern placeholder and is never created literally.
