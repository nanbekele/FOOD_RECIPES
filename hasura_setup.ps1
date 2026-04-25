# === Hasura Permissions and Metadata Setup ===
$headers = @{
  "x-hasura-admin-secret" = "myhasurasecret"
  "Content-Type"          = "application/json"
}

function Send-Hasura($json) {
  try {
    Invoke-RestMethod -Method Post -Uri http://localhost:8082/v1/metadata -Headers $headers -Body $json | Out-Null
    Write-Host "✅ OK -> $($json.Substring(0, [Math]::Min(60, $json.Length)))..."
  } catch {
    Write-Host "⚠️ WARN: $($_.ErrorDetails.Message)"
  }
}

$payloads = @(
# --- Track tables ---
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "recipe_ratings"} } }
"@,
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "favorites"} } }
"@,

@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "recipes"} } }
"@,
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "categories"} } }
"@,
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "recipe_ingredients"} } }
"@,
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "recipe_steps"} } }
"@,
@"
{ "type": "pg_track_table", "args": { "source": "default", "table": {"schema": "public", "name": "recipe_media"} } }
"@,

# --- Relationships ---
@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipe_ratings"}, "name":"recipe", "using":{"foreign_key_constraint_on":"recipe_id"} } }
"@,
@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipe_ratings"}, "name":"user", "using":{"foreign_key_constraint_on":"user_id"} } }
"@,
@"
{ "type":"pg_create_array_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"recipe_ratings", "using": { "foreign_key_constraint_on": { "table": {"schema":"public","name":"recipe_ratings"}, "column": "recipe_id" } } } }
"@,
@"
{ "type":"pg_create_array_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"favorites", "using": { "foreign_key_constraint_on": { "table": {"schema":"public","name":"favorites"}, "column": "recipe_id" } } } }
"@,

@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipe_ingredients"}, "name":"recipe", "using":{"foreign_key_constraint_on":"recipe_id"} } }
"@,
@"
{ "type":"pg_create_array_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"recipe_ingredients", "using": { "foreign_key_constraint_on": { "table": {"schema":"public","name":"recipe_ingredients"}, "column": "recipe_id" } } } }
"@,
@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipe_steps"}, "name":"recipe", "using":{"foreign_key_constraint_on":"recipe_id"} } }
"@,
@"
{ "type":"pg_create_array_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"steps", "using": { "foreign_key_constraint_on": { "table": {"schema":"public","name":"recipe_steps"}, "column": "recipe_id" } } } }
"@,
@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipe_media"}, "name":"recipe", "using":{"foreign_key_constraint_on":"recipe_id"} } }
"@,
@"
{ "type":"pg_create_array_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"media", "using": { "foreign_key_constraint_on": { "table": {"schema":"public","name":"recipe_media"}, "column": "recipe_id" } } } }
"@,

@"
{ "type":"pg_create_object_relationship", "args": { "source":"default", "table":{"schema":"public","name":"recipes"}, "name":"category", "using": { "foreign_key_constraint_on": "category_id" } } }
"@,

@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"public" } }
"@,
@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"user" } }
"@,
@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"public" } }
"@,
@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"user" } }
"@,
@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"public" } }
"@,
@"
{ "type":"pg_drop_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"user" } }
"@,

# --- recipe_ratings permissions ---
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"public","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"user","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"user","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id","rating"] } } }
"@,
@"
{ "type":"pg_create_update_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"user","permission":{"columns":["rating"],"filter":{"user_id":{"_eq":"X-Hasura-User-Id"}}} } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"guest","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"chef","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"guest","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id","rating"] } } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"chef","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id","rating"] } } }
"@,
@"
{ "type":"pg_create_update_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"guest","permission":{"columns":["rating"],"filter":{"user_id":{"_eq":"X-Hasura-User-Id"}}} } }
"@,
@"
{ "type":"pg_create_update_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ratings"},"role":"chef","permission":{"columns":["rating"],"filter":{"user_id":{"_eq":"X-Hasura-User-Id"}}} } }
"@,

# --- favorites permissions ---
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"public","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"user","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"user","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id"] } } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"user","permission":{ "filter":{"user_id":{"_eq":"X-Hasura-User-Id"}} } } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"guest","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"chef","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"guest","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id"] } } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"chef","permission":{ "check":{"user_id":{"_eq":"X-Hasura-User-Id"}}, "set":{"user_id":"X-Hasura-User-Id"}, "columns":["user_id","recipe_id"] } } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"guest","permission":{ "filter":{"user_id":{"_eq":"X-Hasura-User-Id"}} } } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"favorites"},"role":"chef","permission":{ "filter":{"user_id":{"_eq":"X-Hasura-User-Id"}} } } }
"@,

# --- public select permissions for recipes and components ---
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"public","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"user","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"user","permission":{
  "check": {"user_id":{"_eq":"X-Hasura-User-Id"}},
  "set": {"user_id":"X-Hasura-User-Id"},
  "columns": ["user_id","category_id","title","description","featured_image","prep_time_minutes"]
} } }
"@,
@"
{ "type":"pg_create_update_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"user","permission":{
  "columns": ["category_id","title","description","featured_image","prep_time_minutes"],
  "filter": {"user_id":{"_eq":"X-Hasura-User-Id"}}
} } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"guest","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"chef","permission":{"columns":"*","filter":{},"allow_aggregations":true} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"chef","permission":{
  "check": {"user_id":{"_eq":"X-Hasura-User-Id"}},
  "set": {"user_id":"X-Hasura-User-Id"},
  "columns": ["user_id","category_id","title","description","featured_image","prep_time_minutes"]
} } }
"@,
@"
{ "type":"pg_create_update_permission", "args": { "source":"default","table":{"schema":"public","name":"recipes"},"role":"chef","permission":{
  "columns": ["category_id","title","description","featured_image","prep_time_minutes"],
  "filter": {"user_id":{"_eq":"X-Hasura-User-Id"}}
} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"public","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"user","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","name","quantity"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"user","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"guest","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"chef","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"chef","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","name","quantity"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_ingredients"},"role":"chef","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"public","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"user","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","position","description"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"user","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"guest","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"chef","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"chef","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","position","description"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_steps"},"role":"chef","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"public","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"user","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","url","media_type","position"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"user","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,

@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"guest","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"chef","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_insert_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"chef","permission":{
  "check": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}},
  "set": {},
  "columns": ["recipe_id","url","media_type","position"]
} } }
"@,
@"
{ "type":"pg_create_delete_permission", "args": { "source":"default","table":{"schema":"public","name":"recipe_media"},"role":"chef","permission":{
  "filter": {"recipe": {"user_id": {"_eq": "X-Hasura-User-Id"}}}
} } }
"@,

# --- categories permissions ---
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"categories"},"role":"public","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"categories"},"role":"user","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"categories"},"role":"guest","permission":{"columns":"*","filter":{}} } }
"@,
@"
{ "type":"pg_create_select_permission", "args": { "source":"default","table":{"schema":"public","name":"categories"},"role":"chef","permission":{"columns":"*","filter":{}} } }
"@,

# --- Reload metadata ---
@"
{ "type":"reload_metadata", "args": {} }
"@
)

Write-Host "=== Applying Hasura metadata and permissions ==="
foreach ($p in $payloads) { Send-Hasura $p }
Write-Host "✅ All Hasura metadata setup complete!"