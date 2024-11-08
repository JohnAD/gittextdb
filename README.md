# GitTextDB

A document-oriented local filesystem database that supports git tracking.

## Goals

This is a local filesystem database with the following goals:

 - Git trackable: any git commits should show (via diff or other tools)
   precisely and predictably what has changed in the database content.
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

And, added as of Oct 2024:

 - supports existing industry standard for data storage.

Not goals:

 - performance: the database can be slow when large enough.
 - scalability: it is okay if the database cannot be clustered.
 - automatic backups: it is fully expected that off-site "backups" are
   invoked using the standard git commit & push mechanism.
 - security: a side effect of the "no server" goal, the only security
   mechanism is via the filesystem access rights. That is it.

## Document Database Principals

While the relational-database industry has a lot academic and public
understanding, the document-database industry is sadly lacking. So the
following notes describe the philosophical expectations of this
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
 - reflection: data in the document that is a live copy of the truth in
   other documents.
 - cache: data in the document that is "eventually correct" copies of the
   truth in other documents.

Neither "reflection" or "cache" need to be all of the truth. They can be
limited subsets of it per the schema.

In terms of this database, this means that writing to a single document's
truth will:

1. Update the single document AND the "reflected" documents immediately as a
   single almost-atomic operation.
2. A later "batch" operation will update any documents with the cache copies
   of the truth. This might be a few seconds later; or days later.
3. A document not only must store any reflections it needs, but must also store
   the document ids that reflect it's own truth. 

In general "reflections" of truth are computationally expensive and should be used
with care.

### Repeated "searches" are a sign of design failure

A generic search through all the documents of a Collection is generally
discouraged; except for occasional diagnostic queries.

For an example, if a Plant_Harvest collection includes quantities harvested each 
season. Searching for the "top 100 biggest seasons" would certainly work.
But, if this is a common query, it is much better to build a "Top_Seasons"
collection that contains a sorted list with reflection/cache of the
documents. This increases the "cost" of writing to Plant_Harvest, but turns
the later query into a single document read.

## Example Document

The following is an example of what a document could look like:

`{dbdir}/Plant/01HVYS0MR2CZJZ1JNRP5NSD1SQ.JSON`

```json
{
  "_id": "Plant/01HVYS0MR2CZJZ1JNRP5NSD1SQ",
  "_ver": "2024-10-22/01HVYS0MR2CZJZ1JNRP5NSD1JJ",
  "current_price": "mirror/Pricing/01HVYSSY6W062GQEBXVZ3P4XBW",
  "edible": true,
  "name": "Yellow Corn",
  "output_name": "bushels",
  "regions_grown": [
    "cache/Growth_Regions/01HVYSA9Y5X6J0AH2Q8ZVHDHC5",
    "cache/Growth_Regions/01HVYSDYHZVRSTW33G2V7N31Q6"
  ],
  "species": "Zea mays"
}
```

`{dbdir}/Plant/__cache/Growth_Regions/summary_01HVYSA9Y5X6J0AH2Q8ZVHDHC5.JSON`
```json
{
   "_ver": "2024-10-22/01HVYS0MR2CZJZ1JNRP5NS9999",
  "name": "United States Midwest",
  "zone": "temperate"
}
```

And, defining one of the above objects, the `schema` for the `Plant` collection:

`{dir}/Plant.schema.json`

```json
{
   "title": "Plant",
   "description": "Each plant used in the agricultural catalog",
   "type": "object",
   "properties": {
      "_id": "string",
      "_ver": "string",
      "current_price": {
         "description": "daily market price reference",
         "type": "string",
         "ref_type": "mirror"
      },
      "edible": {
         "description": "daily market price reference",
         "type": "boolean"
      },
      "name": {
         "description": "common name",
         "type": "string"
      },
      "output_name": {
         "description": "qty unit of sale",
         "type": "string"
      },
      "regions_grown": {
         "description": "list of geographic region references of where the plant is grown",
         "type": "array",
         "order": "sorted",
         "items": {
            "type": "string",
            "ref_type": "reflect"
         }
      },
      "species": {
         "description": "the scientific species name",
         "type": "string"
      }
   },
   "required": [
      "edible",
      "name",
      "species"
   ]
}
```

This schema follows [https://json-schema.org/](https://json-schema.org/) as a custom schema.
