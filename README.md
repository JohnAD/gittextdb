# GitTextDB

A document-oriented local filesystem database that supports git tracking.

## Goals

This is a local filesystem database with the following goals:

 - Git trackable: any git commits should show (via diff or other tools)
   precisely and predicably what has changed in the database content.
 - human readable: even if the database files are NOT meant to be human
   writable, they should be reasonably human readable.
 - document-oriented: this is not a relational database and must fully
   support the safety paradigm of such. Schema requirement is okay.
 - no server: the database must operate by convention, not mediating server
   software. For "cached data" a cache-update client is expected to run,
   but that software runs as a periodic client, not the server.
 - safe for simultaneous use: using the atomic operations of modern OS
   filesystem, it should be perfectly safe for multiple unconnected programs
   to read/write to the same database.


Not goals:

 - performance: the database can be slow when large enough.
 - scalability: it is okay if the database cannot be clustered.
 - automatic backups: it is fully expected that off-site "backups" are
   invoked using the standard git commit & push mechanism.
 - security: a side effect of the "no server" goal, the only security
   mechanishm is via the filesystem access rights. That is it.

## Document Database Pricipals

While the relational-database industry has a lot academic and public
understanding, the document-database industry is sadly lacking. So the
following notes describe the philisophical expectations of this
database's use and role.

### Read vs Write Scaling

Fast reads; slow writes. Reads should require O(1) access to underlying
documents. Writes, on the other hand, can have any forms of scaling;
including O(n^2) or worse.

### Collections

Documents should be organized by subject into "Collections". Some document
oriented databases allow for open-ended content in the documents, but
this DB requires that each collection have a fixed-but-human-readable
schema defined for each collection.

### Source-Of-Truth (SOT), Data Duplication, and "Eventually Correct"

Data duplication is embraced and expected. More specifically, each document
contains:

 - truth: data in the document that is the _sole_ _source_ of truth about the
   individual target.
 - reflection: data in the document that is a almost-live copy of the truth in
   other documents.
 - cache: data in the document that is "eventually correct" copies of the
   truth in other documents.

In terms of this database, this means that writing to a single document's
truth will:

1. Update the single document AND the "reflected" documents immediately as a
   single almost-atomic operation.
2. A later "batch" operation will update any documents with the cache copies
   of the truth. This might be a few minutes later; or days later.

Reflections of truth are computationally expensive and should be used with
care.

### repeated "searches" are a sign of design failure

A generic search through all the fields of a Collection is generally
discouraged.

For example, if a PlantHarvest collection includes quantities harvested each 
season. Searching for the "top 100 biggest seasons" would certainly work.
But, if this is a common query, it is much better to build a "TopSeasons"
collection that has little truth, but contains a reflection/cache of the
documents. This increases the "cost" of writing to PlantHarvest, but turns
the later query into a single document read.

