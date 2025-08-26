# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.4.6] - 2025-08-27 Wed

- Add: IDs for flat classification files.
- Add: Portugeese nomenclatural statuses.
- Add: update modules.
- Fix [#16]: process synonyms correctly if in DwCA input there are
  acceptedNameUsageID and taxonomicStatus fields together.
- Fix: remove unused error handlings.

## [v0.4.5] - 2025-08-21 Thu

- Add: back `with quotes` option for tsv/pipe files.

## [v0.4.4] - 2025-08-21 Thu

Add: improve nomenclatural status.

## [v0.4.3] - 2025-05-27 Tue

Add: additional converters for nomenclatural status.
Add: vernaculars to coldp.Data.

## [v0.4.2] - 2025-05-21 Wed

Add: read flat classification IDs

## [v0.4.1] - 2025-05-21 Wed

Add: flat classification IDs.

## [v0.4.0] - 2025-05-15 Thu

Add: missing data to SFGA names using parser.

## [v0.3.9] - 2025-05-13 Tue

Add: ability to harvest extensions with multiple files.

## [v0.3.8] - 2025-05-12 Mon

Add [#7]: import DwCA.

## [v0.3.7] - 2025-04-29 Tue

Add: 'department' to actors.

## [v0.3.6] - 2025-04-15 Tue

Add [#13]: update old SFGA to current schema.

## [v0.3.5] - 2025-04-10 Thu

Add: update schema to v0.3.32.

## [v0.3.4] - 2025-04-09 Wed

Add [#12]: writer to CoLDP.
Add [#11]: reader from SFGA.

## [v0.3.3] - 2025-04-06 Sun

Remove: dependency on sqlite3 executable.
Fix: geo-temporal IDs, other small fixes.

## [v0.3.2] - 2025-03-28 Fri

Fix: environment data in coldp.Taxon

## [v0.3.1] - 2025-03-26 Wed

Add [#9]: add XSV to library.

## [v0.3.0] - 2025-03-25 Tue

Add [#8]: add text to library.
Add [#6]: migrate CoLDP to new structure.
Add [#5]: migrate SFGA library to new structure.
Add [#4]: refactor code to accomodate CoLDP and DwCA archives.

## [v0.2.9] - 2025-03-11 Tue

Add: move to coldp v0.3.15: add city, state, country to actor

## [v0.2.8] - 2025-03-10 Mon

Add: add more GN fields to coldp inserts.

## [v0.2.7] - 2025-03-08 Sat

Fix: sfga.Close implementation does not break if sfga database connection
is null.

## [v0.2.6] - 2025-03-05 Wed

Add: move to sfga v0.3.27: name_match table.

## [v0.2.5] - 2025-03-04 Tue

Add: move to sfga v0.3.26: canonical forms in name table.

## [v0.2.4] - 2025-03-01 Sat

Add: move to sfga v0.3.25
Add: tests for coldp insert methods.
Add: update modules.

## [v0.2.3] - 2025-02-20 Thu

Add: coldp insert methods.

## [v0.2.2] - 2025-02-20 Thu

Add: incorporate sfga schema version with the lib.

## [v0.2.1] - 2025-02-20 Thu

Add: improve docs, update modules.

## [v0.2.0] - 2025-02-20 Thu

Add: major refactoring.

## [v0.1.8] - 2025-02-12 Wed

Add: method IsCompatible to check if an app can use the archive.

## [v0.1.7] - 2025-01-23 Thu

Add: fix logs.

## [v0.1.6] - 2025-01-19 Sun

Add: more logs.

## [v0.1.5] - 2024-09-09 Mon

Add [#3]: Improve error handling.

## [v0.1.4] - 2024-09-06 Fri

Add: from-dwca integration.

## [v0.1.3] - 2024-08-24

Add: Decouple database and archive.

## [v0.1.2] - 2024-08-23 Fri

Add [#2]: Connect to SQLite database.

## [v0.1.1] - 2024-08-22 Thu

Add [#1]: Prepare SFGArchive file for use.

## [v0.1.0] - 2024-08-21 Wed

Add: Use cached schema when possible

## [v0.0.1] - 2024-08-20 Tue

Add: Fetching repo works.

## Footnotes

This document follows [changelog guidelines]

[v0.3.6]: https://github.com/sfborg/sflib/compare/v0.3.5...v0.3.6
[v0.3.5]: https://github.com/sfborg/sflib/compare/v0.3.4...v0.3.5
[v0.3.4]: https://github.com/sfborg/sflib/compare/v0.3.3...v0.3.4
[v0.3.3]: https://github.com/sfborg/sflib/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/sfborg/sflib/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/sfborg/sflib/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/sfborg/sflib/compare/v0.2.9...v0.3.0
[v0.2.9]: https://github.com/sfborg/sflib/compare/v0.2.8...v0.2.9
[v0.2.8]: https://github.com/sfborg/sflib/compare/v0.2.7...v0.2.8
[v0.2.7]: https://github.com/sfborg/sflib/compare/v0.2.6...v0.2.7
[v0.2.6]: https://github.com/sfborg/sflib/compare/v0.2.5...v0.2.6
[v0.2.5]: https://github.com/sfborg/sflib/compare/v0.2.4...v0.2.5
[v0.2.4]: https://github.com/sfborg/sflib/compare/v0.2.3...v0.2.4
[v0.2.3]: https://github.com/sfborg/sflib/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/sfborg/sflib/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/sfborg/sflib/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/sfborg/sflib/compare/v0.1.8...v0.2.0
[v0.1.8]: https://github.com/sfborg/sflib/compare/v0.1.7...v0.1.8
[v0.1.7]: https://github.com/sfborg/sflib/compare/v0.1.6...v0.1.7
[v0.1.6]: https://github.com/sfborg/sflib/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/sfborg/sflib/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/sfborg/sflib/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/sfborg/sflib/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/sfborg/sflib/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/sfborg/sflib/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/sfborg/sflib/compare/v0.0.1...v0.1.0
[v0.0.1]: https://github.com/sfborg/sflib/tree/v0.0.1
[#20]: https://github.com/sfborg/sflib/issues/20
[#19]: https://github.com/sfborg/sflib/issues/19
[#18]: https://github.com/sfborg/sflib/issues/18
[#17]: https://github.com/sfborg/sflib/issues/17
[#16]: https://github.com/sfborg/sflib/issues/16
[#15]: https://github.com/sfborg/sflib/issues/15
[#14]: https://github.com/sfborg/sflib/issues/14
[#13]: https://github.com/sfborg/sflib/issues/13
[#12]: https://github.com/sfborg/sflib/issues/12
[#11]: https://github.com/sfborg/sflib/issues/11
[#10]: https://github.com/sfborg/sflib/issues/10
[#9]: https://github.com/sfborg/sflib/issues/9
[#8]: https://github.com/sfborg/sflib/issues/8
[#7]: https://github.com/sfborg/sflib/issues/7
[#6]: https://github.com/sfborg/sflib/issues/6
[#5]: https://github.com/sfborg/sflib/issues/5
[#4]: https://github.com/sfborg/sflib/issues/4
[#3]: https://github.com/sfborg/sflib/issues/3
[#2]: https://github.com/sfborg/sflib/issues/2
[#1]: https://github.com/sfborg/sflib/issues/1
[changelog guidelines]: https://keepachangelog.com/en/1.0.0/
