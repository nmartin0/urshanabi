# AGENTS

Instructions for coding agents working in this repository.

1. Read `RULES.md` before changing anything. Its hard rules have no
   exceptions without the owner's recorded authorization, and its
   engineering standards (part three) are enforced in CI.
2. `ROADMAP.md` is the only list of work. Do not keep another.
3. Never name a company, product, project or programming language
   outside a dependency manifest or tool configuration (`RULES.md`
   H1).
4. Every component is its own project. Work inside one component at a
   time, and reach other components only through `contracts/`.
5. Run `script/test` from the repository root before proposing a
   change, and `script/test-integration` when the change touches a
   dependency. Both must pass.
6. When a choice is uncertain, research current practice before
   making it (H7); when the options are roughly balanced, stop and
   ask the owner (H8).
